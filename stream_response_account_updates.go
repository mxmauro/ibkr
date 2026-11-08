package ibkr

import (
	"sync"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type accountUpdatesStreamResponse struct {
	resp       *models.AccountUpdatesResponse
	closedCh   chan struct{}
	closeOnce  sync.Once
	publishMtx sync.RWMutex
}

// -----------------------------------------------------------------------------

func newAccountUpdatesStreamResponse() *accountUpdatesStreamResponse {
	s := &accountUpdatesStreamResponse{
		resp: &models.AccountUpdatesResponse{
			Channel: make(chan models.AccountUpdateEvent, 4),
		},
		closedCh:   make(chan struct{}),
		closeOnce:  sync.Once{},
		publishMtx: sync.RWMutex{},
	}

	return s
}

func (s *accountUpdatesStreamResponse) close() {
	s.closeOnce.Do(func() {
		s.publishMtx.Lock()
		defer s.publishMtx.Unlock()

		close(s.closedCh)
		close(s.resp.Channel)
	})
}

func (s *accountUpdatesStreamResponse) publish(data models.AccountUpdateEvent) {
	s.publishMtx.RLock()
	defer s.publishMtx.RUnlock()

	select {
	case <-s.closedCh:
		return
	default:
	}

	for {
		select {
		case s.resp.Channel <- data:
			return
		default:
			select {
			case <-s.resp.Channel:
			default:
			}
		}
	}
}
