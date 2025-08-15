package models

import (
	"strings"
)

// -----------------------------------------------------------------------------

type TopMarketDataSize struct {
	tickType TickType
	Size     Decimal // The actual size. US stocks have a multiplier of 100.
}

// -----------------------------------------------------------------------------

func NewTopMarketDataSize(tickType TickType) *TopMarketDataSize {
	return &TopMarketDataSize{
		tickType: tickType,
	}
}

func (t *TopMarketDataSize) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataSize) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", Size=")
	_, _ = sb.WriteString(t.Size.String())
	return sb.String()
}
