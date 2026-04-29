package models

import (
	"time"
)

// -----------------------------------------------------------------------------

// CurrentTimeResponse holds the current server time returned by IB.
type CurrentTimeResponse struct {
	CurrentTime time.Time
}

// ManagedAccountsResponse holds the list of managed accounts returned by IB.
type ManagedAccountsResponse struct {
	Accounts []string
}

// HistoricalDataRequestOptions configures a historical bar request.
type HistoricalDataRequestOptions struct {
	Contract                *Contract
	EndDate                 time.Time
	Duration                int
	DurationUnit            DurationUnit
	BarSize                 BarSize
	WhatToShow              WhatToShow
	OnlyRegularTradingHours bool
	KeepUpToDate            bool
	ChartOptions            TagValueList
}

// HistoricalDataResponse contains the initial bar snapshot and optional update stream.
type HistoricalDataResponse struct {
	Bars    []HistoricalDataBar
	Channel chan HistoricalDataBar
	Cancel  CancelFunc
	Err     ErrFunc
}

// HistoricalTicksRequestOptions configures a historical tick request.
type HistoricalTicksRequestOptions struct {
	Contract                *Contract
	StartDate               time.Time
	EndDate                 time.Time
	NumberOfTicks           int
	WhatToShow              WhatToShow
	OnlyRegularTradingHours bool
	IgnoreSize              bool
}

// HistoricalTicksResponse contains historical ticks in the shape requested by WhatToShow.
type HistoricalTicksResponse struct {
	Ticks       []HistoricalTick
	TicksBidAsk []HistoricalTickBidAsk
	TicksLast   []HistoricalTickLast
}

// ContractDetailsRequestOptions configures a contract details lookup.
type ContractDetailsRequestOptions struct {
	Contract *Contract
}

// ContractDetailsResponse contains all matching contract details.
type ContractDetailsResponse struct {
	ContractDetails []*ContractDetails
}

// MatchingSymbolsRequestOptions configures a matching symbols lookup.
type MatchingSymbolsRequestOptions struct {
	Pattern string
}

// MatchingSymbolsResponse contains all matching symbol descriptions.
type MatchingSymbolsResponse struct {
	ContractDescriptions []*ContractDescription
}

// MarketDataTypeRequestOptions configures the preferred market data mode.
type MarketDataTypeRequestOptions struct {
	Type MarketDataType
}

// TopMarketDataRequestOptions configures a top-of-book market data subscription.
type TopMarketDataRequestOptions struct {
	Contract               *Contract
	AdditionalGenericTicks []GenericTick
	Snapshot               bool
	RegulatorySnapshot     bool
}

// TopMarketDataResponse exposes the streaming top market data channel and controls.
type TopMarketDataResponse struct {
	Channel chan TopMarketData
	Cancel  CancelFunc
	Err     ErrFunc
}

// MarketDepthDataRequestOptions configures a market depth subscription.
type MarketDepthDataRequestOptions struct {
	Contract   *Contract
	RowsCount  int
	SmartDepth bool
}

// MarketDepthDataResponse exposes the streaming market depth channel and current book snapshot.
type MarketDepthDataResponse struct {
	Channel chan MarketDepthData
	Book    *MarketDepthBook
	Cancel  CancelFunc
	Err     ErrFunc
}

// RealTimeBarsRequestOptions configures a real-time bars subscription.
type RealTimeBarsRequestOptions struct {
	Contract                *Contract
	BarSize                 int
	WhatToShow              WhatToShow
	OnlyRegularTradingHours bool
	Options                 TagValueList
}

// RealTimeBarsResponse exposes the streaming real-time bars channel and controls.
type RealTimeBarsResponse struct {
	Channel chan RealTimeBar
	Cancel  CancelFunc
	Err     ErrFunc
}

// TickByTickRequestOptions configures a tick-by-tick market data subscription.
type TickByTickRequestOptions struct {
	Contract      *Contract
	TickType      TickByTickType
	NumberOfTicks int
	IgnoreSize    bool
}

// TickByTickResponse exposes the streaming tick-by-tick data channel and controls.
type TickByTickResponse struct {
	Channel chan TickByTickData
	Cancel  CancelFunc
	Err     ErrFunc
}

// SecurityDefinitionOptionParametersRequestOptions configures an option chain metadata request.
type SecurityDefinitionOptionParametersRequestOptions struct {
	UnderlyingSymbol       string
	FutFopExchange         string
	UnderlyingSecurityType SecurityType
	UnderlyingContractID   int32
}

// SecurityDefinitionOptionParametersResponse contains the returned option chain parameters.
type SecurityDefinitionOptionParametersResponse struct {
	Parameters []SecurityDefinitionOptionParameter
}

// HeadTimestampRequestOptions configures a head timestamp request.
type HeadTimestampRequestOptions struct {
}

// HeadTimestampResponse contains the oldest timestamp available for a data request.
type HeadTimestampResponse struct {
	Timestamp time.Time
}

// CancelFunc cancels an active streaming request.
type CancelFunc func()

// ErrFunc returns the terminal error associated with a request or stream.
type ErrFunc func() error
