/*
Plug in credentials to dynago from aws-sdk-go

This allows dynago to use credential management that comes from AWS-SDK-Go
such as environment, file, and most importantly EC2 role credentials providers.

This uses internal packages in dynago that are in the middle of refactoring,
so consider this in flux.
*/
package credentials

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/underarmour/dynago"
	dynaws "github.com/underarmour/dynago/internal/aws"
)

type RoleCred struct {
	credentials *credentials.Credentials
	stop        chan none
	signer      *dynaws.AwsSigner
	tick        *time.Ticker
}

// Runs in a goroutine, updates creds
func (c *RoleCred) run() {
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

func (c *RoleCred) updateCreds() {
	// TODO actual error handling.
	if val, err := c.credentials.Get(); err == nil {
		c.signer.AccessKey = val.AccessKeyID
		c.signer.SecretKey = val.SecretAccessKey
	}
}

func (c *RoleCred) SignRequest(request *http.Request, bodyBytes []byte) {
	c.signer.SignRequest(request, bodyBytes)
}

func New(region string, credentials *credentials.Credentials) *RoleCred {
	c := &RoleCred{
		credentials: credentials,
		stop:        make(chan none),
		tick:        time.NewTicker(1 * time.Second),
		signer:      &dynaws.AwsSigner{Region: region, Service: "dynamodb"},
	}
	go c.run()
	c.updateCreds()
	return c
}

// Install to the executor. Temporary hack that will go away soon.
func (c *RoleCred) Install(executor *dynago.AwsExecutor) {
	requester := executor.Requester.(*dynaws.RequestMaker)
	if s, ok := requester.Signer.(*dynaws.AwsSigner); ok {
		c.signer.Service = s.Service
		c.signer.SecretKey = s.SecretKey
	}
	requester.Signer = c
}

type none struct{}
