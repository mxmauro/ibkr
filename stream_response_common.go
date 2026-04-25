package ibkr

import "sync"

// -----------------------------------------------------------------------------

func newRequestCancelFunc(cancel func()) func() {
	var once sync.Once

	return func() {
		once.Do(cancel)
	}
}
