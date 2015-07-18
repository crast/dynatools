package streams_test

import (
	"log"

	"github.com/crast/dynatools/streams"
	"github.com/underarmour/dynago"
)

func ExampleStreamer(executor *dynago.AwsExecutor) {
	client := dynago.NewClient(executor)
	result, err := client.DescribeTable("mytable")
	if err != nil {
		return
	}
	config := streams.NewConfig().WithExecutor(executor).WithArn(result.Table.LatestStreamArn)
	streamer := streams.NewStreamer(config)

	// This is actually the mainloop of the application. It doesn't need to do anything else.
	for shard := range streamer.ShardUpdater() {
		log.Printf("Got shard with ID %s", shard.Id())
		go worker(shard)
	}
	log.Printf("Streamer exiting.")
	streamer.Close()
}

// Each instance of worker runs in its own goroutine, consuming a shard of the stream.
func worker(shard *streams.StreamerShard) {
	for packet := range shard.Consume() {
		for _, record := range packet.Records {
			change := record.Dynamodb
			log.Printf("Got record on %s: action=%s, sequence=%s || Id=%s",
				shard.Id(), record.EventName, change.SequenceNumber, change.Keys["Id"],
			)
		}
	}
	log.Printf("Work complete, shard %s", shard.Id())
}
