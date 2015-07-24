package streams

import (
	"gopkg.in/underarmour/dynago.v1"
	"gopkg.in/underarmour/dynago.v1/schema"
)

type Stream struct {
	StreamArn   string
	StreamLabel string
	TableName   string
}

type StreamDescription struct {
	Stream
	KeySchema      []schema.KeySchema
	Shards         []Shard
	StreamStatus   string
	StreamViewType string

	CreationRequestDateTime float64
	LastEvaluatedShardId    string
}

type Shard struct {
	ParentShardId       string
	SequenceNumberRange SequenceNumberRange
	ShardId             string
}

type SequenceNumberRange struct {
	EndingSequenceNumber   string
	StartingSequenceNumber string
}

type StreamRecord struct {
	Keys           dynago.Document
	OldImage       dynago.Document
	NewImage       dynago.Document
	SequenceNumber string
	SizeBytes      uint64
	StreamViewType string
}

type Record struct {
	Dynamodb     StreamRecord
	AwsRegion    string
	EventId      string
	EventName    string
	EventSource  string
	EventVersion string
}
