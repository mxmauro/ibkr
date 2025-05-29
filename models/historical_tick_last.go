package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

// HistoricalTickLast is the historical last tick's description.
// Used when requesting historical tick data with whatToShow = TRADES.
type HistoricalTickLast struct {
	Time              time.Time // Epoch in seconds
	TickAttribLast    TickAttribLast
	Price             float64
	Size              Decimal
	Exchange          string
	SpecialConditions string
}

// -----------------------------------------------------------------------------

func NewHistoricalTickLast() HistoricalTickLast {
	htl := HistoricalTickLast{
		Size: UNSET_DECIMAL,
	}
	return htl
}

func (h HistoricalTickLast) String() string {
	return fmt.Sprintf("Time: %s, TickAttribLast: %s, Price: %f, Size: %s, Exchange: %s, SpecialConditions: %s",
		h.Time.Format("2006/01/02 15:04:05"), h.TickAttribLast, h.Price, h.Size.StringMax(), h.Exchange, h.SpecialConditions)
}
