package models

import (
	"fmt"
	"math"
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type Order struct {
	// order identifier
	OrderID  int32
	ClientID int32
	PermID   int64
	ParentID int32 // Parent order id, to associate Auto STP or TRAIL orders with the original order.

	// primary attributes
	Action        string
	TotalQuantity *Decimal
	DisplaySize   int32
	OrderType     string
	LmtPrice      *float64
	AuxPrice      *float64
	TIF           TimeInForce

	// clearing info
	Account         string // IB account
	SettlingFirm    string
	ClearingAccount string // True beneficiary of the order
	ClearingIntent  string // "" (Default), "IB", "Away", "PTA" (PostTrade)

	// secondary attributes
	AllOrNone       bool
	BlockOrder      bool
	Hidden          bool
	OutsideRTH      bool
	SweepToFill     bool
	PercentOffset   *float64 // REL orders only
	TrailingPercent *float64
	TrailStopPrice  *float64 // TRAILLIMIT orders only
	MinQty          *int32
	GoodAfterTime   string // Format: 20060505 08:00:00 {time zone}
	GoodTillDate    string // Format: 20060505 08:00:00 {time zone}
	OCAGroup        string // one cancels all group name
	OrderRef        string // order reference
	Rule80A         Rule80A
	OCAType         Oca
	TriggerMethod   TriggerMethod

	// extended order fields
	ActiveStartTime string // for GTC orders
	ActiveStopTime  string // for GTC orders

	// advisor allocation orders
	FAGroup      string
	FAMethod     string
	FAPercentage string

	// volatility orders
	Volatility                     *float64
	VolatilityType                 Volatility
	ContinuousUpdate               bool
	ReferencePriceType             int32 // 1=Average, 2 = BidOrAsk
	DeltaNeutralOrderType          string
	DeltaNeutralAuxPrice           *float64
	DeltaNeutralConID              int32
	DeltaNeutralOpenClose          string
	DeltaNeutralShortSale          bool
	DeltaNeutralShortSaleSlot      int32
	DeltaNeutralDesignatedLocation string

	// scale orders
	ScaleInitLevelSize       *int32
	ScaleSubsLevelSize       *int32
	ScalePriceIncrement      *float64
	ScalePriceAdjustValue    *float64
	ScalePriceAdjustInterval *int32
	ScaleProfitOffset        *float64
	ScaleAutoReset           bool
	ScaleInitPosition        *int32
	ScaleInitFillQty         *int32
	ScaleRandomPercent       bool
	ScaleTable               string

	// hedge orders
	HedgeType  string // 'D' - delta, 'B' - beta, 'F' - FX, 'P' - pair
	HedgeParam string // 'beta=X' value for beta hedge, 'ratio=Y' for pair hedge

	// algo orders
	AlgoStrategy string
	AlgoParams   []TagValue
	AlgoID       string

	// combo orders
	SmartComboRoutingParams []TagValue
	OrderComboLegs          []OrderComboLeg

	// processing control
	WhatIf                        bool
	Transmit                      bool // if false, order will be created but not transmitted
	OverridePercentageConstraints bool

	// institutional orders (ie non-cleared)
	OpenClose                   string // O=Open, C=Close
	Origin                      int32  // 0=Customer, 1=Firm
	ShortSaleSlot               int32  // 1 if you hold the shares, 2 if they will be delivered from elsewhere.  Only for Action=SSHORT
	DesignatedLocation          string // set when slot=2 only.
	ExemptCode                  int32  `default:"-1"`
	DeltaNeutralSettlingFirm    string
	DeltaNeutralClearingAccount string
	DeltaNeutralClearingIntent  string

	// SMART routing only
	DiscretionaryAmt   float64
	OptOutSmartRouting bool

	// box or volatility orders only
	AuctionStrategy AuctionStrategy // AUCTION_UNSET, AUCTION_MATCH, AUCTION_IMPROVEMENT, AUCTION_TRANSPARENT

	// box orders only
	StartingPrice *float64
	StockRefPrice *float64
	Delta         *float64

	// pegged to stock and volatility orders only
	StockRangeLower *float64
	StockRangeUpper *float64

	// combo orders only
	BasisPoints     *float64 // EFP orders only
	BasisPointsType *int32   // EFP orders only

	// not held
	NotHeld bool

	// order misc options
	OrderMiscOptions []TagValue

	//order algo id
	Solicited bool

	RandomizeSize  bool
	RandomizePrice bool

	//VER PEG2BENCH fields:
	ReferenceContractID          int32
	PeggedChangeAmount           float64
	IsPeggedChangeAmountDecrease bool
	ReferenceChangeAmount        float64
	ReferenceExchangeID          string
	AdjustedOrderType            string
	TriggerPrice                 *float64
	AdjustedStopPrice            *float64
	AdjustedStopLimitPrice       *float64
	AdjustedTrailingAmount       *float64
	AdjustableTrailingUnit       int32
	LmtPriceOffset               *float64

	Conditions            []OrderCondition
	ConditionsCancelOrder bool
	ConditionsIgnoreRth   bool

	// models
	ModelCode string

	ExtOperator    string
	SoftDollarTier SoftDollarTier

	// native cash quantity
	CashQty *float64

	Mifid2DecisionMaker   string
	Mifid2DecisionAlgo    string
	Mifid2ExecutionTrader string
	Mifid2ExecutionAlgo   string

	// don't use auto price for hedge
	DontUseAutoPriceForHedge bool

	IsOmsContainer              bool
	DiscretionaryUpToLimitPrice bool

	AutoCancelDate       string
	FilledQuantity       *Decimal
	RefFuturesConID      int32
	AutoCancelParent     bool
	Shareholder          string
	ImbalanceOnly        bool
	RouteMarketableToBbo bool
	ParentPermID         int64

	UsePriceMgmtAlgo         bool
	Duration                 *int32
	PostToAts                *int32
	AdvancedErrorOverride    string
	ManualOrderTime          string
	MinTradeQty              *int32
	MinCompeteSize           *int32
	CompeteAgainstBestOffset *float64
	MidOffsetAtWhole         *float64
	MidOffsetAtHalf          *float64
	CustomerAccount          string
	ProfessionalCustomer     bool
	BondAccruedInterest      string
	IncludeOvernight         bool
	ManualOrderIndicator     *int32
	Submitter                string
}

var (
	CompeteAgainstBestOffsetUpToMid = math.Inf(1)
)

// -----------------------------------------------------------------------------

// NewOrder creates a default Order.
func NewOrder() *Order {
	order := &Order{
		ExemptCode:      -1,
		AuctionStrategy: AuctionStrategyUnset,
	}

	return order
}

func (o *Order) HasSameID(other *Order) bool {
	if o.PermID != 0 && other.PermID != 0 {
		return o.PermID == other.PermID
	}
	return o.OrderID == other.OrderID && o.ClientID == other.ClientID
}

func (o *Order) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString(fmt.Sprintf("%s, %s, %s: %s %s %s@%s %s",
		formatter.Int32String(o.OrderID),
		formatter.Int32String(o.ClientID),
		formatter.Int64String(o.PermID),
		o.OrderType,
		o.Action,
		o.TotalQuantity.StringMax(),
		formatter.FloatMaxString(o.LmtPrice),
		o.TIF,
	))

	if len(o.OrderComboLegs) > 0 {
		_, _ = sb.WriteString(" CMB(")
		for idx, leg := range o.OrderComboLegs {
			if idx > 0 {
				_, _ = sb.WriteRune(',')
			}
			_, _ = sb.WriteString(leg.String())
		}
		_, _ = sb.WriteRune(')')
	}

	if len(o.Conditions) > 0 {
		_, _ = sb.WriteString(" COND(")
		for idx, cond := range o.Conditions {
			if idx > 0 {
				_, _ = sb.WriteRune(',')
			}
			_, _ = sb.WriteString(cond.String())
		}
		_, _ = sb.WriteRune(')')
	}

	return sb.String()
}
