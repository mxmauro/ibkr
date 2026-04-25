package ibkr

import (
	"sync/atomic"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type topMarketDataStreamResponse struct {
	resp     *models.TopMarketDataResponse
	queueCh  chan models.TopMarketData
	closedCh chan struct{}
	dataType int32
}

// -----------------------------------------------------------------------------

func newTopMarketDataStreamResponse(dataType models.MarketDataType) *topMarketDataStreamResponse {
	s := &topMarketDataStreamResponse{
		resp: &models.TopMarketDataResponse{
			Channel: make(chan models.TopMarketData, 4),
		},
		queueCh:  make(chan models.TopMarketData, streamDispatchQueueSize),
		closedCh: make(chan struct{}),
		dataType: int32(dataType),
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

func (s *topMarketDataStreamResponse) close() {
	close(s.closedCh)
}

func (s *topMarketDataStreamResponse) enqueue(data models.TopMarketData) {
	data.SetMarketDataType(s.marketDataTypeForTickType(data.TickType()))
	enqueueDropOldest(s.queueCh, data, s.closedCh)
}

func (s *topMarketDataStreamResponse) setMarketDataType(dataType models.MarketDataType) {
	atomic.StoreInt32(&s.dataType, int32(dataType))
}

func (s *topMarketDataStreamResponse) marketDataTypeForTickType(tickType models.TickType) models.MarketDataType {
	dataType := models.MarketDataType(atomic.LoadInt32(&s.dataType))
	if tickType.IsDelayed() {
		if dataType == models.MarketDataTypeFrozen || dataType == models.MarketDataTypeDelayedFrozen {
			return models.MarketDataTypeDelayedFrozen
		}
		return models.MarketDataTypeDelayed
	}
	return dataType
}
