package streams

import (
	"github.com/underarmour/dynago"
	"github.com/underarmour/dynago/schema"
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

///"eventVersion":"1.0", "awsRegion":"us-east-1", "dynamodb":map[string]interface {}{"StreamViewType":"NEW_AND_OLD_IMAGES", "Keys":map[string]interface {}{"Id":map[string]interface {}{"N":"1"}}, "NewImage":map[string]interface {}{"Id":map[string]interface {}{"N":"1"}, "someval":map[string]interface {}{"N":"1437178404"}}, "OldImage":map[string]interface {}{"Id":map[string]interface {}{"N":"1"}, "someval":map[string]interface {}{"N":"1437178402"}}, "SequenceNumber":"247230900000000000003420708", "SizeBytes":38
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
