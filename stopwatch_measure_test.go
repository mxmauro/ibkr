package ibkr_test

import (
	"fmt"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

type StopWatchMeasure struct {
	t         *testing.T
	sw        *stopwatchImpl
	startTime time.Duration
}

// -----------------------------------------------------------------------------

func newStopWatchMeasure(t *testing.T, sw *stopwatchImpl) *StopWatchMeasure {
	t.Helper()
	swm := &StopWatchMeasure{
		t:         t,
		sw:        sw,
		startTime: sw.Elapsed(),
	}
	t.Log("Start time: " + swm.toHuman(swm.startTime))
	return swm
}

func (swm *StopWatchMeasure) End() {
	swm.t.Helper()
	stopTime := swm.sw.Elapsed()
	swm.t.Log("Stop time: " + swm.toHuman(stopTime) + " / Elapsed: " + swm.toHuman(stopTime-swm.startTime))
}

func (swm *StopWatchMeasure) toHuman(duration time.Duration) string {
	us := duration.Microseconds()
	secs := us / 1e6
	us = us % 1e6
	return fmt.Sprintf("%d.%06ds", secs, us)
}
