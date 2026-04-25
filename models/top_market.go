package models

// -----------------------------------------------------------------------------

type TopMarketData interface {
	TickType() TickType // The type of the price being received (i.e., ask price).
	MarketDataType() MarketDataType
	Delayed() bool
	Frozen() bool
	SetMarketDataType(MarketDataType)
	String() string
}

// -----------------------------------------------------------------------------

type TopMarketDataState struct {
	dataType  MarketDataType
	isDelayed bool
	isFrozen  bool
}

// -----------------------------------------------------------------------------

func (s *TopMarketDataState) MarketDataType() MarketDataType {
	return s.dataType
}

func (s *TopMarketDataState) Delayed() bool {
	return s.isDelayed
}

func (s *TopMarketDataState) Frozen() bool {
	return s.isFrozen
}

func (s *TopMarketDataState) SetMarketDataType(dataType MarketDataType) {
	s.dataType = dataType
	s.isDelayed = dataType == MarketDataTypeDelayed || dataType == MarketDataTypeDelayedFrozen
	s.isFrozen = dataType == MarketDataTypeFrozen || dataType == MarketDataTypeDelayedFrozen
}
