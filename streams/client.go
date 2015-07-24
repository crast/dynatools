package streams

import (
	"gopkg.in/underarmour/dynago.v1"
)

type MakeRequester interface {
	MakeRequestUnmarshal(method string, document interface{}, dest interface{}) (err error)
}

const TargetPrefix = "DynamoDBStreams_20120810." // This is the Dynamo API version we support

type Config struct {
	arn      string
	executor MakeRequester
}

func NewConfig() *Config {
	return &Config{}
}

func (c Config) WithExecutor(e *dynago.AwsExecutor) *Config {
	c.executor = e
	return &c
}

func (c Config) WithArn(arn string) *Config {
	c.arn = arn
	return &c
}

func NewClient(config *Config) *Client {
	return &Client{config.executor}
}

/*
Client is the low-level interface to the streams API.
Here, all the individual Streams API requests are available.

The API can be quite clunky, and for many use cases it's
easier to use the Streamer to build applications.
*/
type Client struct {
	caller MakeRequester
}

func (s *Client) DescribeStream(request *DescribeStreamRequest) (dest *DescribeStreamResponse, err error) {
	err = s.caller.MakeRequestUnmarshal(TargetPrefix+"DescribeStream", request, &dest)
	return
}

func (s *Client) GetShardIterator(req *GetIteratorRequest) (dest *GetIteratorResult, err error) {
	err = s.caller.MakeRequestUnmarshal(TargetPrefix+"GetShardIterator", req, &dest)
	return
}

func (s *Client) GetRecords(iterator string) (result *GetRecordsResponse, err error) {
	req := GetRecordsRequest{iterator, 100}
	err = s.caller.MakeRequestUnmarshal(TargetPrefix+"GetRecords", req, &result)
	return
}
