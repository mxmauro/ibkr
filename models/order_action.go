package models

// -----------------------------------------------------------------------------

// OrderAction identifies the side or short-sale action used when submitting an order.
type OrderAction = string

const (
	OrderActionNone      OrderAction = ""
	OrderActionBuy       OrderAction = "BUY"
	OrderActionSell      OrderAction = "SELL"
	OrderActionSellShort OrderAction = "SSHORT"
	OrderActionSellLong  OrderAction = "SLONG"
)
