package models

// -----------------------------------------------------------------------------

type Volatility int32

const (
	VolatilityNone   Volatility = 0
	VolatilityDaily  Volatility = 1
	VolatilityAnnual Volatility = 2
)

// -----------------------------------------------------------------------------

func (v Volatility) String() string {
	switch v {
	case VolatilityNone:
		return "None"
	case VolatilityDaily:
		return "Daily"
	case VolatilityAnnual:
		return "Annual"
	}
	return ""
}
