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
}

type HistoricalDataResponse struct {
	Bars []Bar
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
	ContractDetails []ContractDetails
}

type MatchingSymbolsRequestOptions struct {
	Pattern string
}

type MatchingSymbolsResponse struct {
	ContractDescriptions []ContractDescription
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
	Book    MarketDepthBook
	Cancel  CancelFunc
	Err     ErrFunc
}

type CancelFunc func()

type ErrFunc func() error
