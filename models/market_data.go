package models

// -----------------------------------------------------------------------------

type MarketDataType int

const (
	MarketDataTypeRealtime      MarketDataType = 1
	MarketDataTypeFrozen        MarketDataType = 2
	MarketDataTypeDelayed       MarketDataType = 3
	MarketDataTypeDelayedFrozen MarketDataType = 4
)

// -----------------------------------------------------------------------------

func (mdt MarketDataType) String() string {
	switch mdt {
	case MarketDataTypeRealtime:
		return "Real time"
	case MarketDataTypeFrozen:
		return "Frozen"
	case MarketDataTypeDelayed:
		return "Delayed"
	case MarketDataTypeDelayedFrozen:
		return "Delayed & Frozen"
	}
	return ""
}
