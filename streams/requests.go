package streams

type DescribeStreamRequest struct {
	ExclusiveStartShardId string `json:",omitempty"`
	Limit                 uint   `json:",omitempty"`
	StreamArn             string
}

type DescribeStreamResponse struct {
	StreamDescription StreamDescription
}

type GetIteratorRequest struct {
	StreamArn         string
	ShardId           string
	ShardIteratorType IteratorType
	SequenceNumber    string `json:",omitempty"`
}

type GetIteratorResult struct {
	ShardIterator string
}

type GetRecordsRequest struct {
	ShardIterator string
	Limit         uint `json:",omitempty"`
}

type GetRecordsResponse struct {
	NextShardIterator string
	Records           []Record
}
