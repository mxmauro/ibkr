package utils

import (
	"math"
	"time"
)

// -----------------------------------------------------------------------------

var (
	// EpochDate is the Unix epoch in UTC and is used as a lower bound for IB timestamps.
	EpochDate = time.Unix(0, 0).UTC()
)

// -----------------------------------------------------------------------------

// IsPrintableAsciiString reports whether s only contains printable ASCII characters.
func IsPrintableAsciiString(s string) bool {
	for _, r := range s {
		if (r < 32 || r > 126) && r != 9 && r != 10 && r != 13 {
			return false
		}
	}
	return true
}

// EqualFloat reports whether two float64 values are equal within a small tolerance.
func EqualFloat(a float64, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}
