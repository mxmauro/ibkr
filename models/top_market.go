package models

import (
	"strings"
	"time"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type TopMarketData interface {
	TickType() TickType // The type of the price being received (i.e., ask price).
	String() string
}

type TopMarketDataPrice struct {
	tickType       TickType
	Price          float64 // The actual price.
	CanAutoExecute bool    // Specifies whether the price tick is available for automatic execution
	PastLimit      bool    // Indicates if the bid price is lower than the day's lowest value or the ask price is higher than the highest ask.
	PreOpen        bool    // Indicates whether the bid/ask price tick is from a pre-open session.
}

type TopMarketDataSize struct {
	tickType TickType
	Size     Decimal // The actual size. US stocks have a multiplier of 100.
}

type TopMarketDataGeneric struct {
	tickType TickType
	Value    float64
}

type TopMarketDataString struct {
	tickType TickType
	Value    string
}

type TopMarketDataEFP struct {
	tickType                 TickType
	BasisPoints              float64 // Annualized basis points, which is representative of the financing rate that can be directly compared to broker rates.
	FormattedBasisPoints     string  // Annualized basis points as a formatted string that depicts them in percentage form.
	TotalDividends           float64 // The implied Futures price.
	HoldDays                 int     // The number of hold days until the lastTradeDate of the EFP.
	FutureLastTradeDate      string  // The expiration date of the single stock future.
	DividendImpact           float64 // The dividend impact upon the annualized basis points interest rate.
	DividendsToLastTradeDate float64 // The dividends expected until the expiration of the single stock future.
}

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

type TopMarketDataTimestamp struct {
	tickType  TickType
	Timestamp time.Time
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
	_, _ = sb.WriteString(utils.FloatString(t.Price))
	_, _ = sb.WriteString(", CanAutoExecute=")
	_, _ = sb.WriteString(utils.BoolString(t.CanAutoExecute))
	_, _ = sb.WriteString(", PastLimit=")
	_, _ = sb.WriteString(utils.BoolString(t.PastLimit))
	_, _ = sb.WriteString(", PreOpen=")
	_, _ = sb.WriteString(utils.BoolString(t.PreOpen))
	return sb.String()
}

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
	_, _ = sb.WriteString(utils.FloatString(t.Value))
	return sb.String()
}

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

func NewTopMarketDataEFP(tickType TickType) *TopMarketDataEFP {
	return &TopMarketDataEFP{
		tickType: tickType,
	}
}

func (t *TopMarketDataEFP) TickType() TickType {
	return t.tickType
}

func (t *TopMarketDataEFP) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString(", BasisPoints=")
	_, _ = sb.WriteString(utils.FloatString(t.BasisPoints))
	_, _ = sb.WriteString(", FormattedBasisPoints=\"")
	_, _ = sb.WriteString(t.FormattedBasisPoints)
	_, _ = sb.WriteString("\", TotalDividends=")
	_, _ = sb.WriteString(utils.FloatString(t.TotalDividends))
	_, _ = sb.WriteString(", HoldDays=")
	_, _ = sb.WriteString(utils.IntString(int64(t.HoldDays)))
	_, _ = sb.WriteString(", FutureLastTradeDate=\"")
	_, _ = sb.WriteString(t.FutureLastTradeDate)
	_, _ = sb.WriteString("\", DividendsToLastTradeDate=")
	_, _ = sb.WriteString(utils.FloatString(t.DividendsToLastTradeDate))
	return sb.String()
}

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
	_, _ = sb.WriteString(utils.BoolString(t.IsPriceBased))
	if t.Price != nil {
		_, _ = sb.WriteString(", Price=")
		_, _ = sb.WriteString(utils.FloatString(*t.Price))
	}
	if t.ImpliedVolatility != nil {
		_, _ = sb.WriteString(", ImpliedVolatility=")
		_, _ = sb.WriteString(utils.FloatString(*t.ImpliedVolatility))
	}
	if t.Delta != nil {
		_, _ = sb.WriteString(", Delta=")
		_, _ = sb.WriteString(utils.FloatString(*t.Delta))
	}
	if t.PvDividend != nil {
		_, _ = sb.WriteString(", PvDividend=")
		_, _ = sb.WriteString(utils.FloatString(*t.PvDividend))
	}
	if t.Gamma != nil {
		_, _ = sb.WriteString(", Gamma=")
		_, _ = sb.WriteString(utils.FloatString(*t.Gamma))
	}
	if t.Vega != nil {
		_, _ = sb.WriteString(", Vega=")
		_, _ = sb.WriteString(utils.FloatString(*t.Vega))
	}
	if t.Theta != nil {
		_, _ = sb.WriteString(", Theta=")
		_, _ = sb.WriteString(utils.FloatString(*t.Theta))
	}
	if t.UnderlyingPrice != nil {
		_, _ = sb.WriteString(", UnderlyingPrice=")
		_, _ = sb.WriteString(utils.FloatString(*t.UnderlyingPrice))
	}
	return sb.String()
}

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
