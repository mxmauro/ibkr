package models

import (
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type TopMarketDataOptionComputation struct {
	tickType          TickType
	IsPriceBased      bool     // If false, then it is return-based.
	ImpliedVolatility *float64 // The implied volatility calculated by the TWS option modeler, using the specified tick type value.
	Delta             *float64 // The option delta value.
	Price             *float64 // The option price.
	PvDividend        *float64 // The present value of dividends expected on the option's underlying.
	Gamma             *float64 // The option gamma value.
	Vega              *float64 // The option vega value.
	Theta             *float64 // The option theta value.
	UnderlyingPrice   *float64 // The price of the underlying.
}

// -----------------------------------------------------------------------------

func NewTopMarketDataOptionComputation(tickType TickType) *TopMarketDataOptionComputation {
	return &TopMarketDataOptionComputation{
		tickType: tickType,
	}
}

func (t *TopMarketDataOptionComputation) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataOptionComputation) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Type=")
	_, _ = sb.WriteString(t.tickType.String())
	_, _ = sb.WriteString(", IsPriceBased=")
	_, _ = sb.WriteString(formatter.BoolString(t.IsPriceBased))
	if t.Price != nil {
		_, _ = sb.WriteString(", Price=")
		_, _ = sb.WriteString(formatter.FloatString(*t.Price))
	}
	if t.ImpliedVolatility != nil {
		_, _ = sb.WriteString(", ImpliedVolatility=")
		_, _ = sb.WriteString(formatter.FloatString(*t.ImpliedVolatility))
	}
	if t.Delta != nil {
		_, _ = sb.WriteString(", Delta=")
		_, _ = sb.WriteString(formatter.FloatString(*t.Delta))
	}
	if t.PvDividend != nil {
		_, _ = sb.WriteString(", PvDividend=")
		_, _ = sb.WriteString(formatter.FloatString(*t.PvDividend))
	}
	if t.Gamma != nil {
		_, _ = sb.WriteString(", Gamma=")
		_, _ = sb.WriteString(formatter.FloatString(*t.Gamma))
	}
	if t.Vega != nil {
		_, _ = sb.WriteString(", Vega=")
		_, _ = sb.WriteString(formatter.FloatString(*t.Vega))
	}
	if t.Theta != nil {
		_, _ = sb.WriteString(", Theta=")
		_, _ = sb.WriteString(formatter.FloatString(*t.Theta))
	}
	if t.UnderlyingPrice != nil {
		_, _ = sb.WriteString(", UnderlyingPrice=")
		_, _ = sb.WriteString(formatter.FloatString(*t.UnderlyingPrice))
	}
	return sb.String()
}
