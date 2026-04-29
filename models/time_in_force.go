package models

// -----------------------------------------------------------------------------

// TimeInForce describes how long an order should remain active.
type TimeInForce string

const (
	TimeInForceNone              TimeInForce = ""
	TimeInForceDay               TimeInForce = "DAY" // Day order
	TimeInForceGoodTillCanceled  TimeInForce = "GTC" // Good-Til-Canceled
	TimeInForceOnOpen            TimeInForce = "OPG" // On-Open (MOO/LOO)
	TimeInForceImmediateOrCancel TimeInForce = "IOC" // Immediate-Or-Cancel
	TimeInForceGoodTillDate      TimeInForce = "GTD" // Good-Til-Date
	TimeInForceGoodTillTime      TimeInForce = "GTT" // Good-Til-Time
	TimeInForceAuction           TimeInForce = "AUC" // Auction
	TimeInForceFillOrKill        TimeInForce = "FOK" // Fill-Or-Kill
	TimeInForceGoodTillCrossing  TimeInForce = "GTX" // Good-Til-Canceled (variant)
	TimeInForceDayTillCanceled   TimeInForce = "DTC" // Day-Till-Canceled
	TimeInForceMinutes           TimeInForce = "Minutes"
)
