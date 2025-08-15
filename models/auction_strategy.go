package models

// -----------------------------------------------------------------------------

type AuctionStrategy int32

const (
	AuctionStrategyUnset       AuctionStrategy = 0
	AuctionStrategyMatch       AuctionStrategy = 1
	AuctionStrategyImprovement AuctionStrategy = 2
	AuctionStrategyTransparent AuctionStrategy = 3
)

// -----------------------------------------------------------------------------

func (as AuctionStrategy) String() string {
	switch as {
	case AuctionStrategyMatch:
		return "Match"
	case AuctionStrategyImprovement:
		return "Improvement"
	case AuctionStrategyTransparent:
		return "Transparent"
	}
	return ""
}
