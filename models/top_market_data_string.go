package models

import (
	"strings"
)

// -----------------------------------------------------------------------------

type TopMarketDataString struct {
	tickType TickType
	Value    string
}

// -----------------------------------------------------------------------------

func NewTopMarketDataString(tickType TickType) *TopMarketDataString {
	return &TopMarketDataString{
		tickType: tickType,
	}
}

func (t *TopMarketDataString) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataString) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", Value=\"")
	_, _ = sb.WriteString(t.Value)
	_, _ = sb.WriteString("\"")
	return sb.String()
}
