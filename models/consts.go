package models

// -----------------------------------------------------------------------------

type WhatToShow int

const (
	WhatToShowTrades WhatToShow = iota + 1
	WhatToShowMidPoint
	WhatToShowBid
	WhatToShowAsk
	WhatToShowBidAsk
	WhatToShowHistoricalVolatility
	WhatToShowOptionImpliedVolatility
	WhatToShowSchedule
)

func (wts WhatToShow) String() string {
	switch wts {
	case WhatToShowTrades:
		return "TRADES"
	case WhatToShowMidPoint:
		return "MIDPOINT"
	case WhatToShowBid:
		return "BID"
	case WhatToShowAsk:
		return "ASK"
	case WhatToShowBidAsk:
		return "BID_ASK"
	case WhatToShowHistoricalVolatility:
		return "HISTORICAL_VOLATILITY"
	case WhatToShowOptionImpliedVolatility:
		return "OPTION_IMPLIED_VOLATILITY"
	case WhatToShowSchedule:
		return "SCHEDULE"
	}
	return ""
}

// -----------------------------------------------------------------------------

type BarSize int

const (
	BarSizeOneSecond BarSize = iota + 1
	BarSizeFiveSeconds
	BarSizeTenSeconds
	BarSizeFifteenSeconds
	BarSizeThirtySeconds
	BarSizeOneMinute
	BarSizeTwoMinutes
	BarSizeThreeMinutes
	BarSizeFiveMinutes
	BarSizeTenMinutes
	BarSizeFifteenMinutes
	BarSizeThirtyMinutes
	BarSizeOneHour
	BarSizeTwoHours
	BarSizeFourHours
	BarSizeOneDay
	BarSizeOneWeek
	BarSizeOneMonth
)

func (bs BarSize) String() string {
	switch bs {
	case BarSizeOneSecond:
		return "1 secs"
	case BarSizeFiveSeconds:
		return "5 secs"
	case BarSizeTenSeconds:
		return "10 secs"
	case BarSizeFifteenSeconds:
		return "15 secs"
	case BarSizeThirtySeconds:
		return "30 secs"
	case BarSizeOneMinute:
		return "1 min"
	case BarSizeTwoMinutes:
		return "2 mins"
	case BarSizeThreeMinutes:
		return "3 mins"
	case BarSizeFiveMinutes:
		return "5 mins"
	case BarSizeTenMinutes:
		return "10 mins"
	case BarSizeFifteenMinutes:
		return "15 mins"
	case BarSizeThirtyMinutes:
		return "30 mins"
	case BarSizeOneHour:
		return "1 hour"
	case BarSizeTwoHours:
		return "2 hours"
	case BarSizeFourHours:
		return "4 hours"
	case BarSizeOneDay:
		return "1 day"
	case BarSizeOneWeek:
		return "1 week"
	case BarSizeOneMonth:
		return "1 month"
	}
	return ""
}

func (bs BarSize) Seconds() int64 {
	switch bs {
	case BarSizeOneSecond:
		return 1
	case BarSizeFiveSeconds:
		return 5
	case BarSizeTenSeconds:
		return 10
	case BarSizeFifteenSeconds:
		return 15
	case BarSizeThirtySeconds:
		return 30
	case BarSizeOneMinute:
		return 60
	case BarSizeTwoMinutes:
		return 120
	case BarSizeThreeMinutes:
		return 180
	case BarSizeFiveMinutes:
		return 300
	case BarSizeTenMinutes:
		return 600
	case BarSizeFifteenMinutes:
		return 900
	case BarSizeThirtyMinutes:
		return 1800
	case BarSizeOneHour:
		return 3600
	case BarSizeTwoHours:
		return 7200
	case BarSizeFourHours:
		return 14400
	case BarSizeOneDay:
		return 86400
	case BarSizeOneWeek:
		return 604800
	case BarSizeOneMonth:
		return 2592000 // Assuming equal months of 30 days
	}
	return 0
}

// -----------------------------------------------------------------------------

type DurationUnit int

const (
	DurationUnitSeconds DurationUnit = iota + 1
	DurationUnitDays
	DurationUnitWeeks
	DurationUnitMonths
	DurationUnitYears
)

func (du DurationUnit) String() string {
	switch du {
	case DurationUnitSeconds:
		return "S"
	case DurationUnitDays:
		return "D"
	case DurationUnitWeeks:
		return "W"
	case DurationUnitMonths:
		return "M"
	case DurationUnitYears:
		return "Y"
	}
	return ""
}
