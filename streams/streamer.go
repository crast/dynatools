package streams

import (
	"github.com/underarmour/dynago"
	"time"
)

type ShardWorker func(chan<- Record) error

func NewStreamer(config *Config) *Streamer {
	return &Streamer{
		arn:      config.arn,
		client:   NewClient(config),
		shutdown: make(chan struct{}),
	}
}

/*
Streamer is the high-level interface to streaming.

Streamer allows you to build an application which consumes from a
dynamo table which handles a lot of the ugliness of dealing with
dynamo DB streams.

It provides:
  * Automated polling for changes in the shard topology
  * Easily handle creating new shard consumers
  * Self-tuning read frequency based on data rate in each shard.
  * Clean shutdown of shards which have reached completion
  * Ability to write an application as a simple set of channel range loops
  * timeouts/backoff and retry

Check example_test.go for an example of an application using streamer.
*/
type Streamer struct {
	arn      string
	client   *Client
	shutdown chan struct{}
}

// Describes the stream, making multiple requests if needed to list all shards.
func (s *Streamer) Describe() (*StreamDescription, error) {
	req := &DescribeStreamRequest{
		Limit:     100,
		StreamArn: s.arn,
	}
	var desc *StreamDescription
	for {
		resp, err := s.client.DescribeStream(req)
		if err != nil {
			return nil, err
		}
		if desc == nil {
			desc = &resp.StreamDescription
		} else {
			desc.Shards = append(desc.Shards, resp.StreamDescription.Shards...)
		}
		req.ExclusiveStartShardId = resp.StreamDescription.LastEvaluatedShardId
		if req.ExclusiveStartShardId == "" {
			break
		}
	}
	return desc, nil
}

/*
Generate unique shards in this stream over a channel.

This will start a goroutine which will periodically run Describe on the stream;
remembering all the shards it knows about. If there are any new shards in the
stream on subsequent checks, it will yield those as well.

If the stream is no longer valid, then the channel will be closed.
*/
func (s *Streamer) ShardUpdater() <-chan *StreamerShard {
	c := make(chan *StreamerShard)
	closer := func() { close(c) }
	shards := map[string]*StreamerShard{}

	s.withTimeout(shardUpdaterAdjust, closer, func(adjust time.Duration) time.Duration {
		adjust = adjustSlower
		desc, _ := s.Describe() // TODO error handling
		for _, shard := range desc.Shards {
			if shards[shard.ShardId] == nil {
				ss := &StreamerShard{
					id:       shard.ShardId,
					streamer: s,
				}
				shards[shard.ShardId] = ss
				c <- ss
				adjust = adjustFaster
			}
		}
		return adjust
	})
	return c
}

// Helper to enable the timeout mechanism.
func (s *Streamer) withTimeout(conf adjustConfig, closer func(), callback func(time.Duration) time.Duration) {
	go func() {
		timeout := conf.start
		callback(0)
		for {
			select {
			case <-s.shutdown:
				closer()
				return
			case <-time.After(timeout):
				adjust := callback(timeout)
				switch adjust {
				case adjustFaster:
					if timeout > conf.min {
						timeout -= conf.decreaseBy
					}
				case adjustSlower:
					if timeout < conf.max {
						timeout += conf.increaseBy
					}
				case noAdjust:
					// nothing
				case stopLoop:
					closer()
					return
				default:
					timeout = adjust
				}
			}
		}
	}()
}

type StreamerShard struct {
	id       string
	streamer *Streamer
}

func (s *StreamerShard) Id() string {
	return s.id
}

func (s *StreamerShard) Consume() <-chan Packet {
	ch := make(chan Packet)
	closer := func() { close(ch) }
	iterator, _ := s.myIterator()
	s.streamer.withTimeout(consumeAdjust, closer, func(adjust time.Duration) time.Duration {
		result, err := s.streamer.client.GetRecords(iterator)
		if err == nil {
			iterator = result.NextShardIterator
			ch <- Packet{
				Timeout: adjust,
				Records: result.Records,
			}
			if iterator == "" {
				return stopLoop
			} else if len(result.Records) == 0 {
				return adjustSlower
			} else {
				return adjustFaster
			}
		} else if e, ok := err.(*dynago.Error); ok {
			panic(e)
			_ = e // TODO even more error handling
		}
		return adjust
	})
	return ch
}

func (s *StreamerShard) myIterator() (iterator string, err error) {
	result, err := s.streamer.client.GetShardIterator(&GetIteratorRequest{
		StreamArn:         s.streamer.arn,
		ShardId:           s.id,
		ShardIteratorType: "LATEST", // TODO variables
	})
	if err == nil {
		iterator = result.ShardIterator
	}
	return
}

type Packet struct {
	Timeout time.Duration
	Records []Record
}
