package models

import (
	"time"
)

// -----------------------------------------------------------------------------

type CurrentTimeResponse struct {
	CurrentTime time.Time
}

type ManagedAccountsResponse struct {
	Accounts []string
}

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

type HistoricalDataResponse struct {
	Bars    []HistoricalDataBar
	Channel chan HistoricalDataBar
	Cancel  CancelFunc
	Err     ErrFunc
}

type HistoricalTicksRequestOptions struct {
	Contract                *Contract
	StartDate               time.Time
	EndDate                 time.Time
	NumberOfTicks           int
	WhatToShow              WhatToShow
	OnlyRegularTradingHours bool
	IgnoreSize              bool
}

type HistoricalTicksResponse struct {
	Ticks       []HistoricalTick
	TicksBidAsk []HistoricalTickBidAsk
	TicksLast   []HistoricalTickLast
}

type ContractDetailsRequestOptions struct {
	Contract *Contract
}

type ContractDetailsResponse struct {
	ContractDetails []*ContractDetails
}

type MatchingSymbolsRequestOptions struct {
	Pattern string
}

type MatchingSymbolsResponse struct {
	ContractDescriptions []*ContractDescription
}

type MarketDataTypeRequestOptions struct {
	Type MarketDataType
}

type TopMarketDataRequestOptions struct {
	Contract               *Contract
	AdditionalGenericTicks []GenericTick
	Snapshot               bool
	RegulatorySnapshot     bool
}

type TopMarketDataResponse struct {
	Channel chan TopMarketData
	Cancel  CancelFunc
	Err     ErrFunc
}

type MarketDepthDataRequestOptions struct {
	Contract   *Contract
	RowsCount  int
	SmartDepth bool
}

type MarketDepthDataResponse struct {
	Channel chan MarketDepthData
	Book    *MarketDepthBook
	Cancel  CancelFunc
	Err     ErrFunc
}

type RealTimeBarsRequestOptions struct {
	Contract                *Contract
	BarSize                 int
	WhatToShow              WhatToShow
	OnlyRegularTradingHours bool
	Options                 TagValueList
}

type RealTimeBarsResponse struct {
	Channel chan RealTimeBar
	Cancel  CancelFunc
	Err     ErrFunc
}

type TickByTickRequestOptions struct {
	Contract      *Contract
	TickType      TickByTickType
	NumberOfTicks int
	IgnoreSize    bool
}

type TickByTickResponse struct {
	Channel chan TickByTickData
	Cancel  CancelFunc
	Err     ErrFunc
}

type SecurityDefinitionOptionParametersRequestOptions struct {
	UnderlyingSymbol       string
	FutFopExchange         string
	UnderlyingSecurityType SecurityType
	UnderlyingContractID   int32
}

type SecurityDefinitionOptionParametersResponse struct {
	Parameters []SecurityDefinitionOptionParameter
}

type HeadTimestampRequestOptions struct {
}

type HeadTimestampResponse struct {
	Timestamp time.Time
}

type CancelFunc func()

type ErrFunc func() error
