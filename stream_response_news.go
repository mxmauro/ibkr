package ibkr

import (
	"sync"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type newsStreamResponse struct {
	client   *Client
	resp     *models.NewsResponse
	queueCh  chan models.NewsHeadline
	closedCh chan struct{}

	mtx         sync.Mutex
	requests    []*Request
	terminalErr error
	closeOnce   sync.Once
}

type newsTopicStreamResponse struct {
	stream *newsStreamResponse
	label  string
}

// -----------------------------------------------------------------------------

func newNewsStreamResponse(client *Client) *newsStreamResponse {
	s := &newsStreamResponse{
		client:   client,
		resp:     &models.NewsResponse{Channel: make(chan models.NewsHeadline, 4)},
		queueCh:  make(chan models.NewsHeadline, streamDispatchQueueSize),
		closedCh: make(chan struct{}),
	}

	go func() {
		defer close(s.resp.Channel)

		for {
			select {
			case <-s.closedCh:
				return
			case headline := <-s.queueCh:
				select {
				case <-s.closedCh:
					return
				case s.resp.Channel <- headline:
				}
			}
		}
	}()

	return s
}

func (s *newsStreamResponse) addRequest(req *Request) {
	s.mtx.Lock()
	closed := s.isClosedLocked()
	if !closed {
		s.requests = append(s.requests, req)
	}
	s.mtx.Unlock()
	if closed {
		s.client.cancelTopMarketData(req)
	}
}

func (s *newsStreamResponse) enqueue(headline models.NewsHeadline) {
	enqueueDropOldest(s.queueCh, headline, s.closedCh)
}

func (s *newsStreamResponse) close(err error) {
	s.closeOnce.Do(func() {
		s.mtx.Lock()
		s.terminalErr = err
		requests := append([]*Request(nil), s.requests...)
		close(s.closedCh)
		s.mtx.Unlock()

		for _, req := range requests {
			s.client.cancelTopMarketData(req)
		}
	})
}

func (s *newsStreamResponse) err() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.terminalErr
}

func (s *newsStreamResponse) isClosedLocked() bool {
	select {
	case <-s.closedCh:
		return true
	default:
		return false
	}
}
