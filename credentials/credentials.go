/*
Plug in credentials to dynago from aws-sdk-go

This allows dynago to use credential management that comes from AWS-SDK-Go
such as environment, file, and most importantly EC2 role credentials providers.

This uses internal packages in dynago that are in the middle of refactoring,
so consider this in flux.

Usage:
	import "github.com/aws/aws-sdk-go/aws"
	import "github.com/crast/dynatools/credentials"

	executor := dynago.NewAwsExecutor(...)
	credentials.New("us-east-1", aws.DefaultChainCredentials).Install(executor)
	client := dynago.NewClient(executor)
*/
package credentials

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"gopkg.in/underarmour/dynago.v1"
	dynaws "gopkg.in/underarmour/dynago.v1/internal/aws"
)

/*
Cred is a credentials proxy for dynago to adapt the aws-sdk-go creds.
*/
type Cred struct {
	credentials *credentials.Credentials
	stop        chan none
	signer      *dynaws.AwsSigner
	tick        *time.Ticker
}

// Runs in a goroutine, updates creds
func (c *Cred) run() {
	select {
	case <-c.stop:
		c.stop = nil
		c.tick.Stop()
		return
	case <-c.tick.C:
		if c.credentials.IsExpired() {
			c.updateCreds()
		}
	}
}

func (c *Cred) updateCreds() {
	// TODO actual error handling.
	if val, err := c.credentials.Get(); err == nil {
		c.signer.AccessKey = val.AccessKeyID
		c.signer.SecretKey = val.SecretAccessKey
	}
}

// Implement the AwsSigner interface
func (c *Cred) SignRequest(request *http.Request, bodyBytes []byte) {
	c.signer.SignRequest(request, bodyBytes)
}

/*
Create a new credentials proxy.

Will block until the credentials can be retrieved (might take time on the
network) and also start a goroutine in the background that will retrieve
credentials periodically from the provided credentials.
*/
func New(region string, credentials *credentials.Credentials) *Cred {
	c := &Cred{
		credentials: credentials,
		stop:        make(chan none),
		tick:        time.NewTicker(10 * time.Second),
		signer:      &dynaws.AwsSigner{Region: region, Service: "dynamodb"},
	}
	go c.run()
	c.updateCreds()
	return c
}

// Install to the executor. Temporary hack that will go away soon.
func (c *Cred) Install(executor *dynago.AwsExecutor) {
	requester := executor.Requester.(*dynaws.RequestMaker)
	if s, ok := requester.Signer.(*dynaws.AwsSigner); ok {
		c.signer.Service = s.Service
		c.signer.SecretKey = s.SecretKey
	}
	requester.Signer = c
}

type none struct{}
