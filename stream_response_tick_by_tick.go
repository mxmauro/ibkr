package ibkr

import "github.com/mxmauro/ibkr/models"

// -----------------------------------------------------------------------------

type tickByTickStreamResponse struct {
	resp     *models.TickByTickResponse
	queueCh  chan models.TickByTickData
	closedCh chan struct{}
}

// -----------------------------------------------------------------------------

func newTickByTickStreamResponse() *tickByTickStreamResponse {
	s := &tickByTickStreamResponse{
		resp: &models.TickByTickResponse{
			Channel: make(chan models.TickByTickData, 4),
		},
		queueCh:  make(chan models.TickByTickData, streamDispatchQueueSize),
		closedCh: make(chan struct{}),
	}

	go func() {
		defer close(s.resp.Channel)

		for {
			select {
			case <-s.closedCh:
				return

			case data := <-s.queueCh:
				select {
				case <-s.closedCh:
					return
				case s.resp.Channel <- data:
				}
			}
		}
	}()

	return s
}

func (s *tickByTickStreamResponse) close() {
	close(s.closedCh)
}

func (s *tickByTickStreamResponse) enqueue(data models.TickByTickData) {
	enqueueDropOldest(s.queueCh, data, s.closedCh)
}
