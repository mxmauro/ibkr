package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type ContractDetails struct {
	Contract               Contract
	MarketName             string
	MinTick                float64
	OrderTypes             string
	ValidExchanges         string
	PriceMagnifier         int64
	UnderConID             int64
	LongName               string
	ContractMonth          string
	Industry               string
	Category               string
	Subcategory            string
	TimeZoneID             string
	TradingHours           string
	LiquidHours            string
	EVRule                 string
	EVMultiplier           int64
	AggGroup               int64
	UnderSymbol            string
	UnderSecType           SecurityType
	MarketRuleIDs          string
	RealExpirationDate     string
	LastTradeTime          string
	StockType              string
	MinSize                Decimal
	SizeIncrement          Decimal
	SuggestedSizeIncrement Decimal

	SecIDList []TagValue

	// BOND values
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string
	CouponType        string
	Callable          bool
	Putable           bool
	Coupon            float64
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	Notes             string

	// FUND values
	FundName                        string
	FundFamily                      string
	FundType                        string
	FundFrontLoad                   string
	FundBackLoad                    string
	FundBackLoadTimeInterval        string
	FundManagementFee               string
	FundClosed                      bool
	FundClosedForNewInvestors       bool
	FundClosedForNewMoney           bool
	FundNotifyAmount                string
	FundMinimumInitialPurchase      string
	FundSubsequentMinimumPurchase   string
	FundBlueSkyStates               string
	FundBlueSkyTerritories          string
	FundDistributionPolicyIndicator FundDistributionPolicyIndicator
	FundAssetType                   FundAsset
	IneligibilityReasonList         []IneligibilityReason
}

// -----------------------------------------------------------------------------

func NewContractDetails() ContractDetails {
	cd := ContractDetails{
		Contract:               *NewContract(),
		MinSize:                UNSET_DECIMAL,
		SizeIncrement:          UNSET_DECIMAL,
		SuggestedSizeIncrement: UNSET_DECIMAL,
	}
	return cd
}

func (c ContractDetails) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %t, %f, %t, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		&c.Contract,
		c.MarketName,
		utils.FloatMaxString(c.MinTick),
		c.OrderTypes,
		c.ValidExchanges,
		utils.IntMaxString(c.PriceMagnifier),
		utils.IntMaxString(c.UnderConID),
		c.LongName,
		c.ContractMonth,
		c.Industry,
		c.Category,
		c.Subcategory,
		c.TimeZoneID,
		c.TradingHours,
		c.LiquidHours,
		c.EVRule,
		utils.IntMaxString(c.EVMultiplier),
		c.UnderSymbol,
		c.UnderSecType,
		c.MarketRuleIDs,
		utils.IntMaxString(c.AggGroup),
		c.SecIDList,
		c.RealExpirationDate,
		c.StockType,
		// Bond
		c.Cusip,
		c.Ratings,
		c.DescAppend,
		c.BondType,
		c.CouponType,
		c.Callable,
		c.Putable,
		c.Coupon,
		c.Convertible,
		c.Maturity,
		c.IssueDate,
		c.NextOptionDate,
		c.NextOptionType,
		c.NextOptionPartial,
		c.Notes,
		c.MinSize.StringMax(),
		c.SizeIncrement.StringMax(),
		c.SuggestedSizeIncrement.StringMax(),
	)
}
