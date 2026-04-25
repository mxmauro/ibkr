package ibkr

import (
	"context"
)

// -----------------------------------------------------------------------------

const streamDispatchQueueSize = 64

func enqueueDropOldest[T any](queueCh chan T, value T, closedCh <-chan struct{}) {
	for {
		select {
		case <-closedCh:
			return

		case queueCh <- value:
			return

		default:
			select {
			case <-closedCh:
				return
			case <-queueCh:
			default:
			}
		}
	}
}

func bindStreamRequestContext(ctx context.Context, req *Request, cancel func()) {
	if ctx == nil {
		return
	}

	go func() {
		select {
		case <-ctx.Done():
			cancel()

		case <-req.CompleteCh():
		}
	}()
}
