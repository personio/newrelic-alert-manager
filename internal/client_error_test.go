package internal_test

import (
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	"testing"
)

func TestIsRetryableError_NullError(t *testing.T) {
	var err error

	isRetryable := internal.IsClientError(err)
	if isRetryable {
		t.Error("Nil error should not be a client error")
	}
}


func TestIsRetryableError_GenericError(t *testing.T) {
	err := fmt.Errorf("something went wrong")

	isRetryable := internal.IsClientError(err)
	if isRetryable {
		t.Error("Generic error should not be a client error")
	}
}


func TestIsRetryableError_RetryableError(t *testing.T) {
	err := internal.NewClientError("something went wrong")

	isRetryable := internal.IsClientError(err)
	if !isRetryable {
		t.Error("Generic error should not be a client error")
	}
}