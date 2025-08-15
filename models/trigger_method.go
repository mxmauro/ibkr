package models

// -----------------------------------------------------------------------------

type TriggerMethod int32

const (
	TriggerMethodDefault      TriggerMethod = 0
	TriggerMethodDoubleBidAsk TriggerMethod = 1
	TriggerMethodLast         TriggerMethod = 2
	TriggerMethodDoubleLast   TriggerMethod = 3
	TriggerMethodBidAsk       TriggerMethod = 4
	TriggerMethodLastBidAsk   TriggerMethod = 7
	TriggerMethodMidPoint     TriggerMethod = 8
)

// -----------------------------------------------------------------------------

func (tm TriggerMethod) String() string {
	switch tm {
	case TriggerMethodDefault:
		return "Default"
	case TriggerMethodDoubleBidAsk:
		return "Double Bid Ask"
	case TriggerMethodLast:
		return "Last"
	case TriggerMethodDoubleLast:
		return "Double Last"
	case TriggerMethodBidAsk:
		return "Bid Ask"
	case TriggerMethodLastBidAsk:
		return "Last Bid Ask"
	case TriggerMethodMidPoint:
		return "Mid Point"
	}
	return ""
}
