package connection

import (
	"errors"
	"syscall"
)

// -----------------------------------------------------------------------------

func IsConnectionDropError(err error) bool {
	var errno syscall.Errno

	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.WSAECONNRESET) || errors.Is(errno, syscall.WSAECONNABORTED)
	}
	return false
}
