package models

// -----------------------------------------------------------------------------

type Liquidities int32

const (
	LiquiditiesNone       Liquidities = iota
	LiquiditiesAdded      Liquidities = iota
	LiquiditiesRemoved    Liquidities = iota
	LiquiditiesRoundedOut Liquidities = iota
)

// -----------------------------------------------------------------------------

func (l Liquidities) String() string {
	switch l {
	case LiquiditiesNone:
		return "NONE"
	case LiquiditiesAdded:
		return "ADDED"
	case LiquiditiesRemoved:
		return "REMOVED"
	case LiquiditiesRoundedOut:
		return "ROUNDED OUT"
	}
	return ""
}
