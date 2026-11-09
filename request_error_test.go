package ibkr

import (
	"errors"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

func TestRequestErrorSubscriptionRequired(t *testing.T) {
	for _, code := range []int{10089, 10189, 2188} {
		if !errors.Is(newRequestError(time.Time{}, code, "missing subscription", ""), ErrSubscriptionRequired) {
			t.Fatalf("error code %d was not classified as subscription required", code)
		}
	}
}

func TestRequestErrorOtherCodesAreNotSubscriptionRequired(t *testing.T) {
	if errors.Is(newRequestError(time.Time{}, 300, "cancelled", ""), ErrSubscriptionRequired) {
		t.Fatal("error code 300 was incorrectly classified as subscription required")
	}
}
