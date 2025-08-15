package models

// -----------------------------------------------------------------------------

type TopMarketData interface {
	TickType() TickType // The type of the price being received (i.e., ask price).
	String() string
}
