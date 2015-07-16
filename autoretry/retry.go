// Dynago has no retry logic currently, this is a way to test out auto-retry logic.
package autoretry

import "time"
import "github.com/underarmour/dynago"

type AutoRetry struct {
	Requester  dynago.AwsRequester
	MaxRetries uint // Maximum number of retries.
}

func (r *AutoRetry) MakeRequest(target string, body []byte) (response []byte, err error) {
	backoff := 10 * time.Millisecond
	for i := r.MaxRetries; i > 0; i-- {
		response, err = r.Requester.MakeRequest(target, body)
		if err == nil {
			return
		} else if e, ok := err.(*dynago.Error); !ok {
			break
		} else {
			switch e.Type {
			case dynago.ErrorThroughputExceeded, dynago.ErrorThrottling:
				backoff *= 3
			case dynago.ErrorServiceUnavailable, dynago.ErrorInternalFailure:
				backoff *= 2
			default:
				return
			}
		}
	}
	return
}

func (r *AutoRetry) Install(e Executor) {
	var requester dynago.AwsRequester = r
	awsExec := e.(*dynago.AwsExecutor)
	if awsExec.Requester != requester {
		r.Requester = awsExec.Requester
		awsExec.Requester = r
	}
}
