//go:build windows

package ibkr_test

import (
	"math/bits"
	"syscall"
	"time"
	"unsafe"
)

// -----------------------------------------------------------------------------

type stopwatchImpl struct {
	now int64
}

// -----------------------------------------------------------------------------

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procQueryPerformanceCounter = kernel32.NewProc("QueryPerformanceCounter")
	procQueryPerformanceFreq    = kernel32.NewProc("QueryPerformanceFrequency")

	queryPerformanceFreq int64
)

// -----------------------------------------------------------------------------

func newStopWatch() *stopwatchImpl {
	return &stopwatchImpl{
		now: queryPerformanceCounter(),
	}
}

func (t *stopwatchImpl) Elapsed() time.Duration {
	now := queryPerformanceCounter()
	freq := uint64(queryPerformanceFrequency())
	hi, lo := bits.Mul64(uint64(now-t.now), uint64(time.Second)/uint64(time.Nanosecond))
	quo, _ := bits.Div64(hi, lo, freq)
	return time.Duration(quo)
}

func queryPerformanceCounter() int64 {
	var counter int64

	_, _, _ = procQueryPerformanceCounter.Call(uintptr(unsafe.Pointer(&counter)))
	return counter
}

func queryPerformanceFrequency() int64 {
	if queryPerformanceFreq == 0 {
		_, _, _ = procQueryPerformanceFreq.Call(uintptr(unsafe.Pointer(&queryPerformanceFreq)))
	}
	return queryPerformanceFreq
}
