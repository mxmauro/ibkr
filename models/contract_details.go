package models

import (
	"fmt"
	"strings"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

//goland:noinspection SpellCheckingInspection
type ContractDetails struct {
	Contract               *Contract
	MarketName             string
	MinTick                *float64
	OrderTypes             string
	ValidExchanges         string
	PriceMagnifier         *int32
	UnderConID             *int32
	LongName               string
	ContractMonth          string
	Industry               string
	Category               string
	Subcategory            string
	TimeZoneID             string
	TradingHours           string
	LiquidHours            string
	EVRule                 string
	EVMultiplier           float64
	AggGroup               int32
	UnderSymbol            string
	UnderSecType           SecurityType
	MarketRuleIDs          string
	RealExpirationDate     string
	LastTradeTime          string
	StockType              string
	MinSize                *Decimal
	SizeIncrement          *Decimal
	SuggestedSizeIncrement *Decimal

	SecIDList TagValueList

	// BOND values
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string //XXX
	CouponType        string
	Callable          bool
	Puttable          bool
	Coupon            float64
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	BondNotes         string

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

func NewContractDetails() *ContractDetails {
	cd := ContractDetails{
		Contract: NewContract(),
	}
	return &cd
}

func NewContractDetailsFromProtobufDecoder(
	msgDec *protofmt.Decoder, pb *protobuf.ContractDetails, pbContract *protobuf.Contract,
) *ContractDetails {
	cd := NewContractDetails()
	if pb == nil {
		return cd
	}
	cd.Contract = NewContractFromProtobufDecoder(msgDec, pbContract)
	cd.MarketName = msgDec.String(pb.MarketName)
	cd.MinTick = msgDec.FloatMaxFromString(pb.MinTick)
	cd.PriceMagnifier = msgDec.Int32Max(pb.PriceMagnifier)
	cd.OrderTypes = msgDec.String(pb.OrderTypes)
	cd.ValidExchanges = msgDec.String(pb.ValidExchanges)
	cd.UnderConID = msgDec.Int32Max(pb.UnderConId)
	cd.LongName = msgDec.String(pb.LongName)
	cd.ContractMonth = msgDec.String(pb.ContractMonth)
	cd.Industry = msgDec.String(pb.Industry)
	cd.Category = msgDec.String(pb.Category)
	cd.Subcategory = msgDec.String(pb.Subcategory)
	cd.TimeZoneID = msgDec.String(pb.TimeZoneId)
	cd.TradingHours = msgDec.String(pb.TradingHours)
	cd.LiquidHours = msgDec.String(pb.LiquidHours)
	cd.EVRule = msgDec.String(pb.EvRule)
	cd.EVMultiplier = msgDec.Float(pb.EvMultiplier)
	cd.SecIDList = make([]TagValue, 0, len(pb.SecIdList))
	for k, v := range pb.SecIdList {
		cd.SecIDList = append(cd.SecIDList, TagValue{
			Tag:   k,
			Value: v,
		})
	}
	cd.AggGroup = msgDec.Int32(pb.AggGroup)
	cd.UnderSymbol = msgDec.String(pb.UnderSymbol)
	cd.UnderSecType = SecurityType(msgDec.String(pb.UnderSecType))
	cd.MarketRuleIDs = msgDec.String(pb.MarketRuleIds)
	cd.RealExpirationDate = msgDec.String(pb.RealExpirationDate)
	cd.StockType = msgDec.String(pb.StockType)
	cd.MinSize = NewDecimalMaxFromProtobufDecoder(msgDec, pb.MinSize)
	cd.SizeIncrement = NewDecimalMaxFromProtobufDecoder(msgDec, pb.SizeIncrement)
	cd.SuggestedSizeIncrement = NewDecimalMaxFromProtobufDecoder(msgDec, pb.SuggestedSizeIncrement)
	// bond fields
	cd.Cusip = msgDec.String(pb.Cusip)
	cd.Ratings = msgDec.String(pb.Ratings)
	cd.DescAppend = msgDec.String(pb.DescAppend)
	cd.BondType = msgDec.String(pb.BondType)
	cd.Coupon = msgDec.Float(pb.Coupon)
	cd.CouponType = msgDec.String(pb.CouponType)
	cd.Callable = msgDec.Bool(pb.Callable)
	cd.Puttable = msgDec.Bool(pb.Puttable)
	cd.Convertible = msgDec.Bool(pb.Convertible)
	cd.IssueDate = msgDec.String(pb.IssueDate)
	cd.NextOptionDate = msgDec.String(pb.NextOptionDate)
	cd.NextOptionType = msgDec.String(pb.NextOptionType)
	cd.NextOptionPartial = msgDec.Bool(pb.NextOptionPartial)
	cd.BondNotes = msgDec.String(pb.BondNotes)
	// fund	fields
	cd.FundName = msgDec.String(pb.FundName)
	cd.FundFamily = msgDec.String(pb.FundFamily)
	cd.FundType = msgDec.String(pb.FundType)
	cd.FundFrontLoad = msgDec.String(pb.FundFrontLoad)
	cd.FundBackLoad = msgDec.String(pb.FundBackLoad)
	cd.FundBackLoadTimeInterval = msgDec.String(pb.FundBackLoadTimeInterval)
	cd.FundManagementFee = msgDec.String(pb.FundManagementFee)
	cd.FundClosed = msgDec.Bool(pb.FundClosed)
	cd.FundClosedForNewInvestors = msgDec.Bool(pb.FundClosedForNewInvestors)
	cd.FundClosedForNewMoney = msgDec.Bool(pb.FundClosedForNewMoney)
	cd.FundNotifyAmount = msgDec.String(pb.FundNotifyAmount)
	cd.FundMinimumInitialPurchase = msgDec.String(pb.FundMinimumInitialPurchase)
	cd.FundSubsequentMinimumPurchase = msgDec.String(pb.FundMinimumSubsequentPurchase)
	cd.FundBlueSkyStates = msgDec.String(pb.FundBlueSkyStates)
	cd.FundBlueSkyTerritories = msgDec.String(pb.FundBlueSkyTerritories)
	cd.FundDistributionPolicyIndicator = FundDistributionPolicyIndicator(msgDec.String(pb.FundDistributionPolicyIndicator))
	cd.FundAssetType = FundAsset(msgDec.String(pb.FundAssetType))
	// ineligibility reason
	for _, ir := range pb.IneligibilityReasonList {
		cd.IneligibilityReasonList = append(cd.IneligibilityReasonList, IneligibilityReason{
			ID:          msgDec.String(ir.Id),
			Description: msgDec.String(ir.Description),
		})
	}
	return cd
}

func (cd *ContractDetails) OnContractLastTradeDateOrContractMonthUpdated(isBond bool) {
	//goland:noinspection
	lastTradeDateOrContractMonth := cd.Contract.LastTradeDateOrContractMonth // YYYYMM or YYYYMMDD
	if len(lastTradeDateOrContractMonth) > 0 {
		var split []string

		if strings.Contains(lastTradeDateOrContractMonth, "-") {
			split = strings.Split(lastTradeDateOrContractMonth, "-")
		} else {
			split = strings.Split(lastTradeDateOrContractMonth, " ")

			n := 0
			for _, s := range split {
				if len(s) > 0 {
					split[n] = s
					n++
				}
			}
			split = split[:n]
		}
		if len(split) > 0 {
			if isBond {
				cd.Maturity = split[0]
			} else {
				cd.Contract.LastTradeDateOrContractMonth = split[0]
			}
		}
		if len(split) > 1 {
			cd.LastTradeTime = split[1]
		}
		if isBond && len(split) > 2 {
			cd.TimeZoneID = split[2]
		}
	}
}

func (cd *ContractDetails) String() string {
	return fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %t, %f, %t, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		cd.Contract,
		cd.MarketName,
		formatter.FloatMaxString(cd.MinTick),
		cd.OrderTypes,
		cd.ValidExchanges,
		formatter.Int32MaxString(cd.PriceMagnifier),
		formatter.Int32MaxString(cd.UnderConID),
		cd.LongName,
		cd.ContractMonth,
		cd.Industry,
		cd.Category,
		cd.Subcategory,
		cd.TimeZoneID,
		cd.TradingHours,
		cd.LiquidHours,
		cd.EVRule,
		formatter.FloatString(cd.EVMultiplier),
		cd.UnderSymbol,
		cd.UnderSecType,
		cd.MarketRuleIDs,
		formatter.Int32String(cd.AggGroup),
		cd.SecIDList,
		cd.RealExpirationDate,
		cd.StockType,
		// Bond
		cd.Cusip,
		cd.Ratings,
		cd.DescAppend,
		cd.BondType,
		cd.CouponType,
		cd.Callable,
		cd.Puttable,
		cd.Coupon,
		cd.Convertible,
		cd.Maturity,
		cd.IssueDate,
		cd.NextOptionDate,
		cd.NextOptionType,
		cd.NextOptionPartial,
		cd.BondNotes,
		cd.MinSize.StringMax(),
		cd.SizeIncrement.StringMax(),
		cd.SuggestedSizeIncrement.StringMax(),
	)
}
