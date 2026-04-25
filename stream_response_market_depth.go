package ibkr

import "github.com/mxmauro/ibkr/models"

// -----------------------------------------------------------------------------

type marketDepthDataStreamResponse struct {
	resp     *models.MarketDepthDataResponse
	queueCh  chan models.MarketDepthData
	closedCh chan struct{}
}

// -----------------------------------------------------------------------------

func newMarketDepthDataStreamResponse(maxDepth int) *marketDepthDataStreamResponse {
	s := &marketDepthDataStreamResponse{
		resp: &models.MarketDepthDataResponse{
			Book:    models.NewMarketDepthBook(maxDepth),
			Channel: make(chan models.MarketDepthData, 4),
		},
		queueCh:  make(chan models.MarketDepthData, streamDispatchQueueSize),
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

func (s *marketDepthDataStreamResponse) close() {
	close(s.closedCh)
}

func (s *marketDepthDataStreamResponse) enqueue(data models.MarketDepthData) {
	s.resp.Book.UpdateBook(data)
	enqueueDropOldest(s.queueCh, data, s.closedCh)
}
