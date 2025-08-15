package models

import (
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type TopMarketDataGeneric struct {
	tickType TickType
	Value    float64
}

// -----------------------------------------------------------------------------

func NewTopMarketDataGeneric(tickType TickType) *TopMarketDataGeneric {
	return &TopMarketDataGeneric{
		tickType: tickType,
	}
}

func (t *TopMarketDataGeneric) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataGeneric) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", Value=")
	_, _ = sb.WriteString(formatter.FloatString(t.Value))
	return sb.String()
}
