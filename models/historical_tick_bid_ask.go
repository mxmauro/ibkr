package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

// HistoricalTickBidAsk is the historical tick's description.
// Used when requesting historical tick data with whatToShow = BID_ASK.
type HistoricalTickBidAsk struct {
	Time             time.Time // Epoch in seconds
	TickAttribBidAsk TickAttribBidAsk
	PriceBid         float64
	PriceAsk         float64
	SizeBid          Decimal
	SizeAsk          Decimal
}

// -----------------------------------------------------------------------------

func NewHistoricalTickBidAsk() HistoricalTickBidAsk {
	htba := HistoricalTickBidAsk{
		SizeBid: UNSET_DECIMAL,
		SizeAsk: UNSET_DECIMAL,
	}
	return htba
}

func (h HistoricalTickBidAsk) String() string {
	return fmt.Sprintf("Time: %s, TickAttriBidAsk: %s, PriceBid: %f, PriceAsk: %f, SizeBid: %s, SizeAsk: %s",
		h.Time.Format("2006/01/02 15:04:05"), h.TickAttribBidAsk, h.PriceBid, h.PriceAsk, h.SizeBid.StringMax(),
		h.SizeAsk.StringMax())
}
