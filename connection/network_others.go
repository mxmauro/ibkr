//go:build !windows

package connection

import (
	"errors"
	"syscall"
)

// -----------------------------------------------------------------------------

func IsConnectionDropError(err error) bool {
	var errno syscall.Errno

	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.ECONNRESET) || errors.Is(errno, syscall.ECONNABORTED)
	}
	return false
}
