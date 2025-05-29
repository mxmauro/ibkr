package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

// TickType identifiers for all available tick types in the Interactive Brokers TWS API.
type TickType int

const (
	TickTypeUnset                         TickType = -1
	TickTypeBidSize                       TickType = 0
	TickTypeBid                           TickType = 1
	TickTypeAsk                           TickType = 2
	TickTypeAskSize                       TickType = 3
	TickTypeLast                          TickType = 4
	TickTypeLastSize                      TickType = 5
	TickTypeHigh                          TickType = 6
	TickTypeLow                           TickType = 7
	TickTypeVolume                        TickType = 8
	TickTypeClose                         TickType = 9
	TickTypeBidOptionComputation          TickType = 10
	TickTypeAskOptionComputation          TickType = 11
	TickTypeLastOptionComputation         TickType = 12
	TickTypeModelOptionComputation        TickType = 13
	TickTypeOpen                          TickType = 14
	TickTypeLow13Weeks                    TickType = 15
	TickTypeHigh13Weeks                   TickType = 16
	TickTypeLow26Weeks                    TickType = 17
	TickTypeHigh26Weeks                   TickType = 18
	TickTypeLow52Weeks                    TickType = 19
	TickTypeHigh52Weeks                   TickType = 20
	TickTypeAverageVolume                 TickType = 21
	TickTypeOpenInterest                  TickType = 22
	TickTypeOptionHistoricalVolatility    TickType = 23
	TickTypeOptionImpliedVolatility       TickType = 24
	TickTypeOptionBidExchange             TickType = 25
	TickTypeOptionAskExchange             TickType = 26
	TickTypeOptionCallOpenInterest        TickType = 27
	TickTypeOptionPutOpenInterest         TickType = 28
	TickTypeOptionCallVolume              TickType = 29
	TickTypeOptionPutVolume               TickType = 30
	TickTypeIndexFuturePremium            TickType = 31
	TickTypeBidExchange                   TickType = 32
	TickTypeAskExchange                   TickType = 33
	TickTypeAuctionVolume                 TickType = 34
	TickTypeAuctionPrice                  TickType = 35
	TickTypeAuctionImbalance              TickType = 36
	TickTypeMarkPrice                     TickType = 37
	TickTypeBidEFPComputation             TickType = 38
	TickTypeAskEFPComputation             TickType = 39
	TickTypeLastEFPComputation            TickType = 40
	TickTypeOpenEFPComputation            TickType = 41
	TickTypeHighEFPComputation            TickType = 42
	TickTypeLowEFPComputation             TickType = 43
	TickTypeCloseEFPComputation           TickType = 44
	TickTypeLastTimestamp                 TickType = 45
	TickTypeShortable                     TickType = 46
	TickTypeRTVolume                      TickType = 48
	TickTypeHalted                        TickType = 49
	TickTypeBidYield                      TickType = 50
	TickTypeAskYield                      TickType = 51
	TickTypeLastYield                     TickType = 52
	TickTypeCustomOptionComputation       TickType = 53
	TickTypeTradeCount                    TickType = 54
	TickTypeTradeRate                     TickType = 55
	TickTypeVolumeRate                    TickType = 56
	TickTypeLastRTHTrade                  TickType = 57
	TickTypeRTHistoricalVolatility        TickType = 58
	TickTypeIBDividends                   TickType = 59
	TickTypeBondFactorMultiplier          TickType = 60
	TickTypeRegulatoryImbalance           TickType = 61
	TickTypeNews                          TickType = 62
	TickTypeShortTermVolume3Min           TickType = 63
	TickTypeShortTermVolume5Min           TickType = 64
	TickTypeShortTermVolume10Min          TickType = 65
	TickTypeDelayedBid                    TickType = 66
	TickTypeDelayedAsk                    TickType = 67
	TickTypeDelayedLast                   TickType = 68
	TickTypeDelayedBidSize                TickType = 69
	TickTypeDelayedAskSize                TickType = 70
	TickTypeDelayedLastSize               TickType = 71
	TickTypeDelayedHigh                   TickType = 72
	TickTypeDelayedLow                    TickType = 73
	TickTypeDelayedVolume                 TickType = 74
	TickTypeDelayedClose                  TickType = 75
	TickTypeDelayedOpen                   TickType = 76
	TickTypeRTTradeVolume                 TickType = 77
	TickTypeCreditManagerMarkPrice        TickType = 78
	TickTypeCreditManagerSlowMarkPrice    TickType = 79
	TickTypeDelayedBidOption              TickType = 80
	TickTypeDelayedAskOption              TickType = 81
	TickTypeDelayedLastOption             TickType = 82
	TickTypeDelayedModelOptionComputation TickType = 83
	TickTypeLastExchange                  TickType = 84
	TickTypeLastRegulatoryTime            TickType = 85
	TickTypeFuturesOpenInterest           TickType = 86
	TickTypeAvgOptVolume                  TickType = 87
	TickTypeDelayedLastTimestamp          TickType = 88
	TickTypeShortableShares               TickType = 89
	TickTypeDelayedHalted                 TickType = 90
	TickTypeReutersToMutualFunds          TickType = 91
	TickTypeETFNavClose                   TickType = 92
	TickTypeETFNavPriorClose              TickType = 93
	TickTypeETFNavBid                     TickType = 94
	TickTypeETFNavAsk                     TickType = 95
	TickTypeETFNavLast                    TickType = 96
	TickTypeETFNavFrozenLast              TickType = 97
	TickTypeETFNavHigh                    TickType = 98
	TickTypeETFNavLow                     TickType = 99
	TickTypeEstimatedIPOMidpoint          TickType = 101
	TickTypeFinalIPOPrice                 TickType = 102
	TickTypeDelayedYieldBid               TickType = 103
	TickTypeDelayedYieldAsk               TickType = 104
)

func (tt TickType) IsPrice() bool {
	switch tt {
	case TickTypeBid:
	case TickTypeAsk:
	case TickTypeLast:
	case TickTypeDelayedBid:
	case TickTypeDelayedAsk:
	case TickTypeDelayedLast:
		return true
	}
	return false
}

func (tt TickType) String() string {
	switch tt {
	case TickTypeUnset:
		return ""
	case TickTypeBidSize:
		return "Bid Size"
	case TickTypeBid:
		return "Bid"
	case TickTypeAsk:
		return "Ask"
	case TickTypeAskSize:
		return "Ask Size"
	case TickTypeLast:
		return "Last"
	case TickTypeLastSize:
		return "Last Size"
	case TickTypeHigh:
		return "High"
	case TickTypeLow:
		return "Low"
	case TickTypeVolume:
		return "Volume"
	case TickTypeClose:
		return "Close"
	case TickTypeBidOptionComputation:
		return "Bid Option Computation"
	case TickTypeAskOptionComputation:
		return "Ask Option Computation"
	case TickTypeLastOptionComputation:
		return "Last Option Computation"
	case TickTypeModelOptionComputation:
		return "Model Option Computation"
	case TickTypeOpen:
		return "Open"
	case TickTypeLow13Weeks:
		return "Low 13 Weeks"
	case TickTypeHigh13Weeks:
		return "High 13 Weeks"
	case TickTypeLow26Weeks:
		return "Low 26 Weeks"
	case TickTypeHigh26Weeks:
		return "High 26 Weeks"
	case TickTypeLow52Weeks:
		return "Low 52 Weeks"
	case TickTypeHigh52Weeks:
		return "High 52 Weeks"
	case TickTypeAverageVolume:
		return "Average Volume"
	case TickTypeOpenInterest:
		return "Open Interest"
	case TickTypeOptionHistoricalVolatility:
		return "Option Historical Volatility"
	case TickTypeOptionImpliedVolatility:
		return "Option Implied Volatility"
	case TickTypeOptionBidExchange:
		return "Option Bid Exchange"
	case TickTypeOptionAskExchange:
		return "Option Ask Exchange"
	case TickTypeOptionCallOpenInterest:
		return "Option Call Open Interest"
	case TickTypeOptionPutOpenInterest:
		return "Option Put Open Interest"
	case TickTypeOptionCallVolume:
		return "Option Call Volume"
	case TickTypeOptionPutVolume:
		return "Option Put Volume"
	case TickTypeIndexFuturePremium:
		return "Index Future Premium"
	case TickTypeBidExchange:
		return "Bid Exchange"
	case TickTypeAskExchange:
		return "Ask Exchange"
	case TickTypeAuctionVolume:
		return "Auction Volume"
	case TickTypeAuctionPrice:
		return "Auction Price"
	case TickTypeAuctionImbalance:
		return "Auction Imbalance"
	case TickTypeMarkPrice:
		return "Mark Price"
	case TickTypeBidEFPComputation:
		return "Bid EFP Computation"
	case TickTypeAskEFPComputation:
		return "Ask EFP Computation"
	case TickTypeLastEFPComputation:
		return "Last EFP Computation"
	case TickTypeOpenEFPComputation:
		return "Open EFP Computation"
	case TickTypeHighEFPComputation:
		return "High EFP Computation"
	case TickTypeLowEFPComputation:
		return "Low EFP Computation"
	case TickTypeCloseEFPComputation:
		return "Close EFP Computation"
	case TickTypeLastTimestamp:
		return "Last Timestamp"
	case TickTypeShortable:
		return "Shortable"
	case TickTypeRTVolume:
		return "RT Volume"
	case TickTypeHalted:
		return "Halted"
	case TickTypeBidYield:
		return "Bid Yield"
	case TickTypeAskYield:
		return "Ask Yield"
	case TickTypeLastYield:
		return "Last Yield"
	case TickTypeCustomOptionComputation:
		return "Custom Option Computation"
	case TickTypeTradeCount:
		return "Trade Count"
	case TickTypeTradeRate:
		return "Trade Rate"
	case TickTypeVolumeRate:
		return "Volume Rate"
	case TickTypeLastRTHTrade:
		return "Last RTH Trade"
	case TickTypeRTHistoricalVolatility:
		return "RT Historical Volatility"
	case TickTypeIBDividends:
		return "IB Dividends"
	case TickTypeBondFactorMultiplier:
		return "Bond Factor Multiplier"
	case TickTypeRegulatoryImbalance:
		return "Regulatory Imbalance"
	case TickTypeNews:
		return "News"
	case TickTypeShortTermVolume3Min:
		return "Short Term Volume 3 Min"
	case TickTypeShortTermVolume5Min:
		return "Short Term Volume 5 Min"
	case TickTypeShortTermVolume10Min:
		return "Short Term Volume 10 Min"
	case TickTypeDelayedBid:
		return "Delayed Bid"
	case TickTypeDelayedAsk:
		return "Delayed Ask"
	case TickTypeDelayedLast:
		return "Delayed Last"
	case TickTypeDelayedBidSize:
		return "Delayed Bid Size"
	case TickTypeDelayedAskSize:
		return "Delayed Ask Size"
	case TickTypeDelayedLastSize:
		return "Delayed Last Size"
	case TickTypeDelayedHigh:
		return "Delayed High"
	case TickTypeDelayedLow:
		return "Delayed Low"
	case TickTypeDelayedVolume:
		return "Delayed Volume"
	case TickTypeDelayedClose:
		return "Delayed Close"
	case TickTypeDelayedOpen:
		return "Delayed Open"
	case TickTypeRTTradeVolume:
		return "RT Trade Volume"
	case TickTypeCreditManagerMarkPrice:
		return "Credit-manager Mark Price"
	case TickTypeCreditManagerSlowMarkPrice:
		return "Credit-manager Slow Mark Price"
	case TickTypeDelayedBidOption:
		return "Delayed Bid Option"
	case TickTypeDelayedAskOption:
		return "Delayed Ask Option"
	case TickTypeDelayedLastOption:
		return "Delayed Last Option"
	case TickTypeDelayedModelOptionComputation:
		return "Delayed Model Option Computation"
	case TickTypeLastExchange:
		return "Last Exchange"
	case TickTypeLastRegulatoryTime:
		return "Last Regulatory Time"
	case TickTypeFuturesOpenInterest:
		return "Futures Open Interest"
	case TickTypeAvgOptVolume:
		return "Average Option Volume"
	case TickTypeDelayedLastTimestamp:
		return "Delayed Last Timestamp"
	case TickTypeShortableShares:
		return "Shortable Shares"
	case TickTypeDelayedHalted:
		return "Delayed Halted"
	case TickTypeReutersToMutualFunds:
		return "Reuters To Mutual Funds"
	case TickTypeETFNavClose:
		return "ETF NAV Close"
	case TickTypeETFNavPriorClose:
		return "ETF NAV Prior Close"
	case TickTypeETFNavBid:
		return "ETF NAV Bid"
	case TickTypeETFNavAsk:
		return "ETF NAV Ask"
	case TickTypeETFNavLast:
		return "ETF NAV Last"
	case TickTypeETFNavFrozenLast:
		return "ETF NAV Frozen Last"
	case TickTypeETFNavHigh:
		return "ETF NAV High"
	case TickTypeETFNavLow:
		return "ETF NAV Low"
	case TickTypeEstimatedIPOMidpoint:
		return "Estimated IPO Midpoint"
	case TickTypeFinalIPOPrice:
		return "Final IPO Price"
	case TickTypeDelayedYieldBid:
		return "Delayed Yield Bid"
	case TickTypeDelayedYieldAsk:
		return "Delayed Yield Ask"
	default:
		return fmt.Sprintf("TickType(%d)", tt)
	}

}
