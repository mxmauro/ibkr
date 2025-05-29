package models

import (
	"fmt"
	"math"
	"strings"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type AuctionStrategy = int64

const (
	AuctionStrategyUnset       AuctionStrategy = 0
	AuctionStrategyMatch       AuctionStrategy = 1
	AuctionStrategyImprovement AuctionStrategy = 2
	AuctionStrategyTransparent AuctionStrategy = 3
)

type Order struct {
	// order identifier
	OrderID  int64
	ClientID int64
	PermID   int64

	// main order fields
	Action        string
	TotalQuantity Decimal `default:"UNSET_DECIMAL"`
	OrderType     string
	LmtPrice      float64 `default:"UNSET_FLOAT"`
	AuxPrice      float64 `default:"UNSET_FLOAT"`

	// extended order fields
	TIF                           string // "Time in Force" - DAY, GTC, etc.
	ActiveStartTime               string // for GTC orders
	ActiveStopTime                string // for GTC orders
	OCAGroup                      string // one cancels all group name
	OCAType                       int64  // 1 = CANCEL_WITH_BLOCK, 2 = REDUCE_WITH_BLOCK, 3 = REDUCE_NON_BLOCK
	OrderRef                      string // order reference
	Transmit                      bool   `default:"true"` // if false, order will be created but not transmitted
	ParentID                      int64  // Parent order Id, to associate Auto STP or TRAIL orders with the original order.
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64 // 0=Default, 1=Double_Bid_Ask, 2=Last, 3=Double_Last, 4=Bid_Ask, 7=Last_or_Bid_Ask, 8=Mid-point
	OutsideRTH                    bool
	Hidden                        bool
	GoodAfterTime                 string // Format: 20060505 08:00:00 {time zone}
	GoodTillDate                  string // Format: 20060505 08:00:00 {time zone}
	Rule80A                       string // Individual = 'I', Agency = 'A', AgentOtherMember = 'W', IndividualPTIA = 'J', AgencyPTIA = 'U', AgentOtherMemberPTIA = 'M', IndividualPT = 'K', AgencyPT = 'Y', AgentOtherMemberPT = 'N'
	AllOrNone                     bool
	MinQty                        int64   `default:"UNSET_INT"`
	PercentOffset                 float64 `default:"UNSET_FLOAT"` // REL orders only
	OverridePercentageConstraints bool
	TrailStopPrice                float64 `default:"UNSET_FLOAT"` // TRAILLIMIT orders only
	TrailingPercent               float64 `default:"UNSET_FLOAT"`

	// financial advisors only
	FAGroup      string
	FAMethod     string
	FAPercentage string

	// institutional (ie non-cleared) only
	OpenClose          string // O=Open, C=Close
	Origin             int64  // 0=Customer, 1=Firm
	ShortSaleSlot      int64  // 1 if you hold the shares, 2 if they will be delivered from elsewhere.  Only for Action=SSHORT
	DesignatedLocation string // set when slot=2 only.
	ExemptCode         int64  `default:"-1"`

	// SMART routing only
	DiscretionaryAmt   float64
	OptOutSmartRouting bool

	// BOX exchange orders only
	AuctionStrategy AuctionStrategy // AUCTION_UNSET, AUCTION_MATCH, AUCTION_IMPROVEMENT, AUCTION_TRANSPARENT
	StartingPrice   float64         `default:"UNSET_FLOAT"`
	StockRefPrice   float64         `default:"UNSET_FLOAT"`
	Delta           float64         `default:"UNSET_FLOAT"`

	// pegged to stock and VOL orders only
	StockRangeLower float64 `default:"UNSET_FLOAT"`
	StockRangeUpper float64 `default:"UNSET_FLOAT"`

	RandomizeSize  bool
	RandomizePrice bool

	// VOLATILITY ORDERS ONLY
	Volatility                     float64 `default:"UNSET_FLOAT"`
	VolatilityType                 int64   `default:"UNSET_INT"`
	DeltaNeutralOrderType          string
	DeltaNeutralAuxPrice           float64 `default:"UNSET_FLOAT"`
	DeltaNeutralConID              int64
	DeltaNeutralSettlingFirm       string
	DeltaNeutralClearingAccount    string
	DeltaNeutralClearingIntent     string
	DeltaNeutralOpenClose          string
	DeltaNeutralShortSale          bool
	DeltaNeutralShortSaleSlot      int64
	DeltaNeutralDesignatedLocation string
	ContinuousUpdate               bool
	ReferencePriceType             int64 `default:"UNSET_INT"` // 1=Average, 2 = BidOrAsk

	// COMBO ORDERS ONLY
	BasisPoints     float64 `default:"UNSET_FLOAT"` // EFP orders only
	BasisPointsType int64   `default:"UNSET_INT"`   // EFP orders only

	// SCALE ORDERS ONLY
	ScaleInitLevelSize       int64   `default:"UNSET_INT"`
	ScaleSubsLevelSize       int64   `default:"UNSET_INT"`
	ScalePriceIncrement      float64 `default:"UNSET_FLOAT"`
	ScalePriceAdjustValue    float64 `default:"UNSET_FLOAT"`
	ScalePriceAdjustInterval int64   `default:"UNSET_INT"`
	ScaleProfitOffset        float64 `default:"UNSET_FLOAT"`
	ScaleAutoReset           bool
	ScaleInitPosition        int64 `default:"UNSET_INT"`
	ScaleInitFillQty         int64 `default:"UNSET_INT"`
	ScaleRandomPercent       bool
	ScaleTable               string

	// HEDGE ORDERS
	HedgeType  string // 'D' - delta, 'B' - beta, 'F' - FX, 'P' - pair
	HedgeParam string // 'beta=X' value for beta hedge, 'ratio=Y' for pair hedge

	// Clearing info
	Account         string // IB account
	SettlingFirm    string
	ClearingAccount string // True beneficiary of the order
	ClearingIntent  string // "" (Default), "IB", "Away", "PTA" (PostTrade)

	// ALGO ORDERS ONLY
	AlgoStrategy string

	AlgoParams              []TagValue
	SmartComboRoutingParams []TagValue

	AlgoID string

	// What-if
	WhatIf bool

	// Not Held
	NotHeld   bool
	Solictied bool

	// models
	ModelCode string

	// order combo legs
	OrderComboLegs   []OrderComboLeg
	OrderMiscOptions []TagValue

	//VER PEG2BENCH fields:
	ReferenceContractID          int64
	PeggedChangeAmount           float64
	IsPeggedChangeAmountDecrease bool
	ReferenceChangeAmount        float64
	ReferenceExchangeID          string
	AdjustedOrderType            string
	TriggerPrice                 float64 `default:"UNSET_FLOAT"`
	AdjustedStopPrice            float64 `default:"UNSET_FLOAT"`
	AdjustedStopLimitPrice       float64 `default:"UNSET_FLOAT"`
	AdjustedTrailingAmount       float64 `default:"UNSET_FLOAT"`
	AdjustableTrailingUnit       int64
	LmtPriceOffset               float64 `default:"UNSET_FLOAT"`

	Conditions            []OrderCondition
	ConditionsCancelOrder bool
	ConditionsIgnoreRth   bool

	// ext operator
	ExtOperator string

	SoftDollarTier SoftDollarTier

	// native cash quantity
	CashQty float64 `default:"UNSET_FLOAT"`

	Mifid2DecisionMaker   string
	Mifid2DecisionAlgo    string
	Mifid2ExecutionTrader string
	Mifid2ExecutionAlgo   string

	// don't use auto price for hedge
	DontUseAutoPriceForHedge bool

	IsOmsContainer bool

	DiscretionaryUpToLimitPrice bool

	AutoCancelDate       string
	FilledQuantity       Decimal `default:"UNSET_DECIMAL"`
	RefFuturesConID      int64
	AutoCancelParent     bool
	Shareholder          string
	ImbalanceOnly        bool
	RouteMarketableToBbo bool
	ParentPermID         int64

	UsePriceMgmtAlgo         bool
	Duration                 int64 `default:"UNSET_INT"`
	PostToAts                int64 `default:"UNSET_INT"`
	AdvancedErrorOverride    string
	ManualOrderTime          string
	MinTradeQty              int64   `default:"UNSET_INT"`
	MinCompeteSize           int64   `default:"UNSET_INT"`
	CompeteAgainstBestOffset float64 `default:"UNSET_FLOAT"`
	MidOffsetAtWhole         float64 `default:"UNSET_FLOAT"`
	MidOffsetAtHalf          float64 `default:"UNSET_FLOAT"`
	CustomerAccount          string
	ProfessionalCustomer     bool
	BondAccruedInterest      string
	IncludeOvernight         bool
	ManualOrderIndicator     int64 `default:"UNSET_INT"`
	Submitter                string
}

var (
	COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID = math.Inf(1)
)

// -----------------------------------------------------------------------------

// NewOrder creates a default Order.
func NewOrder() *Order {
	order := &Order{
		TotalQuantity: UNSET_DECIMAL,
		LmtPrice:      common.UNSET_FLOAT,
		AuxPrice:      common.UNSET_FLOAT,

		Transmit:        true,
		MinQty:          common.UNSET_INT,
		PercentOffset:   common.UNSET_FLOAT,
		TrailStopPrice:  common.UNSET_FLOAT,
		TrailingPercent: common.UNSET_FLOAT,

		ExemptCode: -1,

		AuctionStrategy: AuctionStrategyUnset,
		StartingPrice:   common.UNSET_FLOAT,
		StockRefPrice:   common.UNSET_FLOAT,
		Delta:           common.UNSET_FLOAT,

		StockRangeLower: common.UNSET_FLOAT,
		StockRangeUpper: common.UNSET_FLOAT,

		Volatility:           common.UNSET_FLOAT,
		VolatilityType:       common.UNSET_INT,
		DeltaNeutralAuxPrice: common.UNSET_FLOAT,
		ReferencePriceType:   common.UNSET_INT,

		BasisPoints:     common.UNSET_FLOAT,
		BasisPointsType: common.UNSET_INT,

		ScaleInitLevelSize:       common.UNSET_INT,
		ScaleSubsLevelSize:       common.UNSET_INT,
		ScalePriceIncrement:      common.UNSET_FLOAT,
		ScalePriceAdjustValue:    common.UNSET_FLOAT,
		ScalePriceAdjustInterval: common.UNSET_INT,
		ScaleProfitOffset:        common.UNSET_FLOAT,
		ScaleInitPosition:        common.UNSET_INT,
		ScaleInitFillQty:         common.UNSET_INT,

		TriggerPrice:           common.UNSET_FLOAT,
		AdjustedStopPrice:      common.UNSET_FLOAT,
		AdjustedStopLimitPrice: common.UNSET_FLOAT,
		AdjustedTrailingAmount: common.UNSET_FLOAT,
		LmtPriceOffset:         common.UNSET_FLOAT,

		CashQty: common.UNSET_FLOAT,

		FilledQuantity: UNSET_DECIMAL,

		Duration:                 common.UNSET_INT,
		PostToAts:                common.UNSET_INT,
		MinTradeQty:              common.UNSET_INT,
		MinCompeteSize:           common.UNSET_INT,
		CompeteAgainstBestOffset: common.UNSET_FLOAT,
		MidOffsetAtWhole:         common.UNSET_FLOAT,
		MidOffsetAtHalf:          common.UNSET_FLOAT,
		ManualOrderIndicator:     common.UNSET_INT,
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
		utils.IntMaxString(o.OrderID),
		utils.IntMaxString(o.ClientID),
		utils.IntMaxString(o.PermID),
		o.OrderType,
		o.Action,
		o.TotalQuantity.StringMax(),
		utils.FloatMaxString(o.LmtPrice),
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
