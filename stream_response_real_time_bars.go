package ibkr

import "github.com/mxmauro/ibkr/models"

// -----------------------------------------------------------------------------

type realTimeBarsStreamResponse struct {
	resp            *models.RealTimeBarsResponse
	queueCh         chan models.RealTimeBar
	closedCh        chan struct{}
	barDurationSecs int64
}

// -----------------------------------------------------------------------------

func newRealTimeBarsStreamResponse(barDurationSecs int64) *realTimeBarsStreamResponse {
	s := &realTimeBarsStreamResponse{
		resp: &models.RealTimeBarsResponse{
			Channel: make(chan models.RealTimeBar, 4),
		},
		queueCh:         make(chan models.RealTimeBar, streamDispatchQueueSize),
		closedCh:        make(chan struct{}),
		barDurationSecs: barDurationSecs,
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

func (s *realTimeBarsStreamResponse) close() {
	close(s.closedCh)
}

func (s *realTimeBarsStreamResponse) enqueue(data models.RealTimeBar) {
	enqueueDropOldest(s.queueCh, data, s.closedCh)
}
