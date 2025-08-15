package models

import (
	"strings"
	"time"
)

// -----------------------------------------------------------------------------

type TopMarketDataTimestamp struct {
	tickType  TickType
	Timestamp time.Time
}

// -----------------------------------------------------------------------------

func NewTopMarketDataTimestamp(tickType TickType) *TopMarketDataTimestamp {
	return &TopMarketDataTimestamp{
		tickType: tickType,
	}
}

func (t *TopMarketDataTimestamp) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataTimestamp) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", Timestamp=")
	_, _ = sb.WriteString(t.Timestamp.Format("2006/01/02 15:04:05"))
	return sb.String()
}
