package internal

import "gopkg.in/underarmour/dynago.v0"

// Assert an error as a dynago error
func AssertError(err error) *dynago.Error {
	if e, ok := err.(*dynago.Error); ok {
		return e
	}
	return nil
}
