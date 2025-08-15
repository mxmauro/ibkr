package models

import (
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type TopMarketDataEFP struct {
	tickType                 TickType
	BasisPoints              float64 // Annualized basis points, which is representative of the financing rate that can be directly compared to broker rates.
	FormattedBasisPoints     string  // Annualized basis points as a formatted string that depicts them in percentage form.
	ImpliedFuturesPrice      float64 // The implied Futures price.
	HoldDays                 int32   // The number of hold days until the lastTradeDate of the EFP.
	FutureLastTradeDate      string  // The expiration date of the single stock future.
	DividendImpact           float64 // The dividend impact upon the annualized basis points interest rate.
	DividendsToLastTradeDate float64 // The dividends expected until the expiration of the single stock future.
}

// -----------------------------------------------------------------------------

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
	_, _ = sb.WriteString(formatter.FloatString(t.BasisPoints))
	_, _ = sb.WriteString(", FormattedBasisPoints=\"")
	_, _ = sb.WriteString(t.FormattedBasisPoints)
	_, _ = sb.WriteString("\", ImpliedFuturesPrice=")
	_, _ = sb.WriteString(formatter.FloatString(t.ImpliedFuturesPrice))
	_, _ = sb.WriteString(", HoldDays=")
	_, _ = sb.WriteString(formatter.Int32String(t.HoldDays))
	_, _ = sb.WriteString(", FutureLastTradeDate=\"")
	_, _ = sb.WriteString(t.FutureLastTradeDate)
	_, _ = sb.WriteString("\", DividendsToLastTradeDate=")
	_, _ = sb.WriteString(formatter.FloatString(t.DividendsToLastTradeDate))
	return sb.String()
}
