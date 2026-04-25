package ibkr

import "github.com/mxmauro/ibkr/models"

// -----------------------------------------------------------------------------

type historicalDataStreamResponse struct {
	resp            *models.HistoricalDataResponse
	queueCh         chan models.HistoricalDataBar
	closedCh        chan struct{}
	keepUpToDate    bool
	barDurationSecs int64
}

// -----------------------------------------------------------------------------

func newHistoricalDataStreamResponse(keepUpToDate bool, barDurationSecs int64) *historicalDataStreamResponse {
	s := &historicalDataStreamResponse{
		resp: &models.HistoricalDataResponse{
			Bars: make([]models.HistoricalDataBar, 0),
		},
		keepUpToDate:    keepUpToDate,
		barDurationSecs: barDurationSecs,
	}
	if keepUpToDate {
		s.resp.Channel = make(chan models.HistoricalDataBar, 4)
		s.queueCh = make(chan models.HistoricalDataBar, streamDispatchQueueSize)
		s.closedCh = make(chan struct{})

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
	}

	return s
}

func (s *historicalDataStreamResponse) close() {
	if s.closedCh != nil {
		close(s.closedCh)
	}
}

func (s *historicalDataStreamResponse) appendBars(bars []models.HistoricalDataBar) {
	s.resp.Bars = append(s.resp.Bars, bars...)
	if s.keepUpToDate {
		for _, bar := range bars {
			enqueueDropOldest(s.queueCh, bar, s.closedCh)
		}
	}
}

func (s *historicalDataStreamResponse) applyUpdate(bar models.HistoricalDataBar) {
	if barsCount := len(s.resp.Bars); barsCount > 0 && s.resp.Bars[barsCount-1].Date.Equal(bar.Date) {
		s.resp.Bars[barsCount-1] = bar
	} else {
		s.resp.Bars = append(s.resp.Bars, bar)
	}
	if s.keepUpToDate {
		enqueueDropOldest(s.queueCh, bar, s.closedCh)
	}
}
