package models

// -----------------------------------------------------------------------------

// OrderType identifies the execution style IB expects for an order.
type OrderType = string

const (
	OrderTypeNone                 OrderType = ""
	OrderTypeMarket               OrderType = "MKT"
	OrderTypeLimit                OrderType = "LMT"
	OrderTypeMarketToLimit        OrderType = "MTL"
	OrderTypeMarketIfTouched      OrderType = "MIT"
	OrderTypeMarketOnClose        OrderType = "MOC"
	OrderTypeMidPrice             OrderType = "MIDPRICE"
	OrderTypePeggedToMarket       OrderType = "PEG MKT"
	OrderTypePeggedToStock        OrderType = "PEG STK"
	OrderTypeRelative             OrderType = "REL"
	OrderTypeBoxTop               OrderType = "BOX TOP"
	OrderTypeLimitIfTouched       OrderType = "LIT"
	OrderTypeLimitOnClose         OrderType = "LOC"
	OrderTypePassiveRelative      OrderType = "PASSV REL"
	OrderTypePeggedToMidpoint     OrderType = "PEG MID"
	OrderTypeStop                 OrderType = "STP"
	OrderTypeStopLimit            OrderType = "STP LMT"
	OrderTypeStopWithProtection   OrderType = "STP PRT"
	OrderTypeTrailingStop         OrderType = "TRAIL"
	OrderTypeTrailingStopLimit    OrderType = "TRAIL LIMIT"
	OrderTypeMarketWithProtection OrderType = "MKT PRT"
	OrderTypeRelativePlusLimit    OrderType = "REL + LMT"
	OrderTypeRelativePlusMarket   OrderType = "REL + MKT"
	OrderTypeVolatility           OrderType = "VOL"
	OrderTypePeggedToBenchmark    OrderType = "PEG BENCH"
	OrderTypePeggedToBest         OrderType = "PEG BEST"
)
