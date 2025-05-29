package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

// HistoricalTick is the historical tick's description.
// Used when requesting historical tick data with whatToShow = MIDPOINT.
type HistoricalTick struct {
	Time  time.Time
	Price float64
	Size  Decimal
}

// -----------------------------------------------------------------------------

func NewHistoricalTick() HistoricalTick {
	ht := HistoricalTick{
		Size: UNSET_DECIMAL,
	}
	return ht
}

func (h HistoricalTick) String() string {
	return fmt.Sprintf("Time: %s, Price: %f, Size: %s", h.Time.Format("2006/01/02 15:04:05"), h.Price,
		h.Size.StringMax())
}
