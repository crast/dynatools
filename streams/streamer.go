package streams

import (
	"gopkg.in/underarmour/dynago.v1"
	"log"
	"sync"
	"time"
)

type ShardWorker func(chan<- Record) error

func NewStreamer(config *Config) *Streamer {
	return &Streamer{
		arn:      config.arn,
		client:   NewClient(config),
		shutdown: make(chan none),
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
	shutdown chan none
	wakeUp   chan none
	wg       sync.WaitGroup
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

// Close causes the streamer to shut down all goroutines and end.
// After close is called, this streamer is no longer valid.
func (s *Streamer) Close() error {
	close(s.shutdown)
	s.wg.Wait()
	return nil
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

	s.wakeUp = s.withTimeout(shardUpdaterAdjust, closer, func(adjust time.Duration) time.Duration {
		log.Printf("Updating shard list")
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
func (s *Streamer) withTimeout(conf adjustConfig, closer func(), callback func(time.Duration) time.Duration) chan none {
	s.wg.Add(1)
	wake := make(chan none)
	go func() {
		defer s.wg.Done()
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
			case <-wake:
				callback(-1)
			}
		}
	}()
	return wake
}

// Individual shard within a stream
type StreamerShard struct {
	id       string
	streamer *Streamer
	iterator string
}

// Get the amazon AWS shard unique ID
func (s *StreamerShard) Id() string {
	return s.id
}

// Denote that we want to start at the latest event in this shard.
func (s *StreamerShard) AtLatest() {
	s.myIterator(IteratorLatest, "")
}

// Denote we want to start at the oldest event in this shard.
func (s *StreamerShard) AtTrimHorizon() {
	s.myIterator(IteratorTrimHorizon, "")
}

// Start at the given sequencenumber
func (s *StreamerShard) AtSequenceNum(seq string) {
	s.myIterator(IteratorAfterSequence, seq)
}

/*
Begin consuming this shard, yielding the results on a channel.
This consumer will self-adjust how fast it's asking for requests
based on the rate of data received on the channel.

If the consumer reaches the end of a shard's data stream (such as this shard
is no longer actively updating) or if our Streamer is closed, then the
goroutine will end and the channel will be closed.
*/
func (s *StreamerShard) Consume() <-chan Update {
	ch := make(chan Update)
	closer := func() { close(ch) }
	iterator := s.iterator
	if iterator == "" {
		panic("Iterator must be set using one of the At functions first.")
	}
	nothingTimes := 10
	s.streamer.withTimeout(consumeAdjust, closer, func(timeout time.Duration) time.Duration {
		result, err := s.streamer.client.GetRecords(iterator)
		if err == nil {
			ch <- Update{
				Timeout: timeout,
				Records: result.Records,
			}
			if result.NextShardIterator == "" {
				nothingTimes--
				log.Printf("Getting to end of line  %#v %d", result, nothingTimes)
				if nothingTimes <= 0 {
					s.streamer.wakeUp <- none{}
					return stopLoop
				}
			} else {
				iterator = result.NextShardIterator
			}
			if len(result.Records) == 0 {
				timeout = adjustSlower
			} else {
				timeout = adjustFaster
			}
		} else {
			if e, ok := err.(*dynago.Error); ok {
				switch e.Type {
				case dynago.ErrorThrottling, dynago.ErrorThroughputExceeded, dynago.ErrorInternalFailure:
					return timeout + time.Second
				case dynago.ErrorExpiredIterator, dynago.ErrorTrimmedData:
					// TODO determine what we do on an expired iterator
				}
			}
			ch <- Update{
				Timeout: timeout,
				Error:   err,
			}
			return stopLoop
		}
		return timeout
	})
	return ch
}

func (s *StreamerShard) myIterator(iType IteratorType, sequenceNumber string) {
	result, err := s.streamer.client.GetShardIterator(&GetIteratorRequest{
		StreamArn:         s.streamer.arn,
		ShardId:           s.id,
		ShardIteratorType: iType,
		SequenceNumber:    sequenceNumber,
	})
	if err == nil {
		s.iterator = result.ShardIterator
	}
	return
}

/*
Update is received from Consume() every time it gets a response of any sort.

It's not necessarily an error if there are no records; this can happen even
on a full shard if you're paging past a portion of the shard where there is
no data (empty segments, old segments, etc).
*/
type Update struct {
	Timeout time.Duration // How long we waited for this update
	Records []Record      // The records we received for this update.
	Error   error         // Any error we received from the API
}

type none struct{}
