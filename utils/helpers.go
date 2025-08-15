package utils

import (
	"math"
	"time"
)

// -----------------------------------------------------------------------------

var (
	EpochDate = time.Unix(0, 0).UTC()
)

// -----------------------------------------------------------------------------

func IsPrintableAsciiString(s string) bool {
	for _, r := range s {
		if (r < 32 || r > 126) && r != 9 && r != 10 && r != 13 {
			return false
		}
	}
	return true
}

func EqualFloat(a float64, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}
