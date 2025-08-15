//go:build !windows

package ibkr_test

import (
	"time"
)

// -----------------------------------------------------------------------------

type stopwatchImpl struct {
	now time.Time
}

// -----------------------------------------------------------------------------

func newStopWatch() *stopwatchImpl {
	return &stopwatchImpl{
		now: time.Now(),
	}
}

func (t *stopwatchImpl) Elapsed() time.Duration {
	return time.Since(t.now)
}
