package utils

import (
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
