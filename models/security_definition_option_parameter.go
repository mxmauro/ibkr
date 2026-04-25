package models

// -----------------------------------------------------------------------------

type SecurityDefinitionOptionParameter struct {
	Exchange             string
	UnderlyingContractID int32
	TradingClass         string
	Multiplier           string
	Expirations          []string
	Strikes              []float64
}
