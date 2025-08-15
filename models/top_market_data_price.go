package models

import (
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type TopMarketDataPrice struct {
	tickType       TickType
	Price          float64 // The actual price.
	CanAutoExecute bool    // Specifies whether the price tick is available for automatic execution
	PastLimit      bool    // Indicates if the bid price is lower than the day's lowest value or the ask price is higher than the highest ask.
	PreOpen        bool    // Indicates whether the bid/ask price tick is from a pre-open session.
}

// -----------------------------------------------------------------------------

func NewTopMarketDataPrice(tickType TickType) *TopMarketDataPrice {
	return &TopMarketDataPrice{
		tickType: tickType,
	}
}

func (t *TopMarketDataPrice) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataPrice) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", Price=")
	_, _ = sb.WriteString(formatter.FloatString(t.Price))
	_, _ = sb.WriteString(", CanAutoExecute=")
	_, _ = sb.WriteString(formatter.BoolString(t.CanAutoExecute))
	_, _ = sb.WriteString(", PastLimit=")
	_, _ = sb.WriteString(formatter.BoolString(t.PastLimit))
	_, _ = sb.WriteString(", PreOpen=")
	_, _ = sb.WriteString(formatter.BoolString(t.PreOpen))
	return sb.String()
}
