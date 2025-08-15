package models

import (
	"errors"

	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderDecoder struct {
	order         *Order
	contract      *Contract
	orderState    *OrderState
	version       int
	serverVersion int
}

// -----------------------------------------------------------------------------

func (d *OrderDecoder) decodeOrderId(msgDec *message.Decoder) {
	d.order.OrderID = msgDec.Int32()
}

func (d *OrderDecoder) decodeContractFields(msgDec *message.Decoder) {
	d.contract.ConID = msgDec.Int32()
	d.contract.Symbol = msgDec.String()
	d.contract.SecType = NewSecurityTypeFromString(msgDec.String())
	d.contract.LastTradeDateOrContractMonth = msgDec.String()
	d.contract.Strike = msgDec.FloatMax()
	d.contract.Right = msgDec.String()
	d.contract.Multiplier = msgDec.FloatMax()
	d.contract.Exchange = msgDec.String()
	d.contract.Currency = msgDec.String()
	d.contract.LocalSymbol = msgDec.String()
	d.contract.TradingClass = msgDec.String()
}

func (d *OrderDecoder) decodeAction(msgDec *message.Decoder) {
	d.order.Action = msgDec.String()
}

func (d *OrderDecoder) decodeTotalQuantity(msgDec *message.Decoder) {
	d.order.TotalQuantity = NewDecimalMaxFromMessageDecoder(msgDec)
}

func (d *OrderDecoder) decodeOrderType(msgDec *message.Decoder) {
	d.order.OrderType = msgDec.String()
}

func (d *OrderDecoder) decodeLmtPrice(msgDec *message.Decoder) {
	d.order.LmtPrice = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeAuxPrice(msgDec *message.Decoder) {
	d.order.AuxPrice = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeTIF(msgDec *message.Decoder) {
	d.order.TIF = NewTimeInForceFromString(msgDec.String())
}

func (d *OrderDecoder) decodeOcaGroup(msgDec *message.Decoder) {
	d.order.OCAGroup = msgDec.String()
}

func (d *OrderDecoder) decodeAccount(msgDec *message.Decoder) {
	d.order.Account = msgDec.String()
}

func (d *OrderDecoder) decodeOpenClose(msgDec *message.Decoder) {
	d.order.OpenClose = msgDec.String()
}

func (d *OrderDecoder) decodeOrigin(msgDec *message.Decoder) {
	d.order.Origin = msgDec.Int32()
}

func (d *OrderDecoder) decodeOrderRef(msgDec *message.Decoder) {
	d.order.OrderRef = msgDec.String()
}

func (d *OrderDecoder) decodeClientId(msgDec *message.Decoder) {
	d.order.ClientID = msgDec.Int32()
}

func (d *OrderDecoder) decodePermId(msgDec *message.Decoder) {
	d.order.PermID = msgDec.Int64()
}

func (d *OrderDecoder) decodeOutsideRth(msgDec *message.Decoder) {
	d.order.OutsideRTH = msgDec.Bool()
}

func (d *OrderDecoder) decodeHidden(msgDec *message.Decoder) {
	d.order.Hidden = msgDec.Bool()
}

func (d *OrderDecoder) decodeDiscretionaryAmount(msgDec *message.Decoder) {
	d.order.DiscretionaryAmt = msgDec.Float()
}

func (d *OrderDecoder) decodeGoodAfterTime(msgDec *message.Decoder) {
	d.order.GoodAfterTime = msgDec.String()
}

func (d *OrderDecoder) skipSharesAllocation(msgDec *message.Decoder) {
	msgDec.Skip()
}

func (d *OrderDecoder) decodeFAParams(msgDec *message.Decoder) {
	d.order.FAGroup = msgDec.String()
	d.order.FAMethod = msgDec.String()
	d.order.FAPercentage = msgDec.String()
	msgDec.Skip() // skip deprecated FAProfile
}

func (d *OrderDecoder) decodeModelCode(msgDec *message.Decoder) {
	d.order.ModelCode = msgDec.String()
}

func (d *OrderDecoder) decodeGoodTillDate(msgDec *message.Decoder) {
	d.order.GoodTillDate = msgDec.String()
}

func (d *OrderDecoder) decodeRule80A(msgDec *message.Decoder) {
	d.order.Rule80A = NewRule80aFromString(msgDec.String())
}

func (d *OrderDecoder) decodePercentOffset(msgDec *message.Decoder) {
	d.order.PercentOffset = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeSettlingFirm(msgDec *message.Decoder) {
	d.order.SettlingFirm = msgDec.String()
}

func (d *OrderDecoder) decodeShortSaleParams(msgDec *message.Decoder) {
	d.order.ShortSaleSlot = msgDec.Int32()
	d.order.DesignatedLocation = msgDec.String()
	d.order.ExemptCode = msgDec.Int32()
}

func (d *OrderDecoder) decodeAuctionStrategy(msgDec *message.Decoder) {
	d.order.AuctionStrategy = AuctionStrategy(msgDec.Int32())
}

func (d *OrderDecoder) decodeBoxOrderParams(msgDec *message.Decoder) {
	d.order.StartingPrice = msgDec.FloatMax()
	d.order.StockRefPrice = msgDec.FloatMax()
	d.order.Delta = msgDec.FloatMax()
}

func (d *OrderDecoder) decodePegToStkOrVolOrderParams(msgDec *message.Decoder) {
	d.order.StockRangeLower = msgDec.FloatMax()
	d.order.StockRangeUpper = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeDisplaySize(msgDec *message.Decoder) {
	d.order.DisplaySize = msgDec.Int32()
}

func (d *OrderDecoder) decodeBlockOrder(msgDec *message.Decoder) {
	d.order.BlockOrder = msgDec.Bool()
}

func (d *OrderDecoder) decodeSweepToFill(msgDec *message.Decoder) {
	d.order.SweepToFill = msgDec.Bool()
}

func (d *OrderDecoder) decodeAllOrNone(msgDec *message.Decoder) {
	d.order.AllOrNone = msgDec.Bool()
}

func (d *OrderDecoder) decodeMinQty(msgDec *message.Decoder) {
	d.order.MinQty = msgDec.Int32Max()
}

func (d *OrderDecoder) decodeOcaType(msgDec *message.Decoder) {
	d.order.OCAType = Oca(msgDec.Int32())
}

func (d *OrderDecoder) skipETradeOnly(msgDec *message.Decoder) {
	msgDec.Skip() // deprecated order.ETradeOnly
}

func (d *OrderDecoder) skipFirmQuoteOnly(msgDec *message.Decoder) {
	msgDec.Skip() // deprecated order.FirmQuoteOnly
}

func (d *OrderDecoder) skipNbboPriceCap(msgDec *message.Decoder) {
	msgDec.Skip() // deprecated order.NBBOPriceCap
}

func (d *OrderDecoder) decodeParentId(msgDec *message.Decoder) {
	d.order.ParentID = msgDec.Int32()
}

func (d *OrderDecoder) decodeTriggerMethod(msgDec *message.Decoder) {
	d.order.TriggerMethod = TriggerMethod(msgDec.Int32())
}

func (d *OrderDecoder) decodeVolOrderParams(msgDec *message.Decoder, readOpenOrderAttribs bool) {
	d.order.Volatility = msgDec.FloatMax()
	d.order.VolatilityType = Volatility(msgDec.Int32())
	d.order.DeltaNeutralOrderType = msgDec.String()
	d.order.DeltaNeutralAuxPrice = msgDec.FloatMax()

	if len(d.order.DeltaNeutralOrderType) > 0 {
		d.order.DeltaNeutralConID = msgDec.Int32()
		if readOpenOrderAttribs {
			d.order.DeltaNeutralSettlingFirm = msgDec.String()
			d.order.DeltaNeutralClearingAccount = msgDec.String()
			d.order.DeltaNeutralClearingIntent = msgDec.String()
			d.order.DeltaNeutralOpenClose = msgDec.String()
		}

		d.order.DeltaNeutralShortSale = msgDec.Bool()
		d.order.DeltaNeutralShortSaleSlot = msgDec.Int32()
		d.order.DeltaNeutralDesignatedLocation = msgDec.String()
	}

	d.order.ContinuousUpdate = msgDec.Bool()
	d.order.ReferencePriceType = msgDec.Int32()
}

func (d *OrderDecoder) decodeTrailParams(msgDec *message.Decoder) {
	d.order.TrailStopPrice = msgDec.FloatMax()
	d.order.TrailingPercent = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeBasisPoints(msgDec *message.Decoder) {
	d.order.BasisPoints = msgDec.FloatMax()
	d.order.BasisPointsType = msgDec.Int32Max()
}

func (d *OrderDecoder) decodeComboLegs(msgDec *message.Decoder) {
	d.contract.ComboLegsDescription = msgDec.String()
	comboLegsCount := int(msgDec.Int32())
	if comboLegsCount < 0 {
		msgDec.SetErr(errors.New("invalid combolegs count"))
		return
	}
	d.contract.ComboLegs = make([]*ComboLeg, 0, comboLegsCount)
	for i := 0; i < comboLegsCount; i++ {
		comboleg := NewComboLeg()
		comboleg.ConID = msgDec.Int32()
		comboleg.Ratio = msgDec.Int32()
		comboleg.Action = msgDec.String()
		comboleg.Exchange = msgDec.String()
		comboleg.OpenClose = msgDec.Int32()
		comboleg.ShortSalesSlot = msgDec.Int32()
		comboleg.DesignatedLocation = msgDec.String()
		comboleg.ExemptCode = msgDec.Int32()

		d.contract.ComboLegs = append(d.contract.ComboLegs, comboleg)
	}

	orderComboLegsCount := int(msgDec.Int32())
	if orderComboLegsCount < 0 {
		msgDec.SetErr(errors.New("invalid order combolegs count"))
		return
	}
	d.order.OrderComboLegs = make([]OrderComboLeg, 0, orderComboLegsCount)
	for i := 0; i < orderComboLegsCount; i++ {
		orderComboLeg := OrderComboLeg{}
		orderComboLeg.Price = msgDec.FloatMax()

		d.order.OrderComboLegs = append(d.order.OrderComboLegs, orderComboLeg)
	}
}

func (d *OrderDecoder) decodeSmartComboRoutingParams(msgDec *message.Decoder) {
	smartComboRoutingParamsCount := int(msgDec.Int32())
	if smartComboRoutingParamsCount < 0 {
		msgDec.SetErr(errors.New("invalid smart combo routing parameters count"))
		return
	}
	d.order.SmartComboRoutingParams = make([]TagValue, 0, smartComboRoutingParamsCount)
	for i := 0; i < smartComboRoutingParamsCount; i++ {
		tagValue := TagValue{}
		tagValue.Tag = msgDec.String()
		tagValue.Value = msgDec.String()

		d.order.SmartComboRoutingParams = append(d.order.SmartComboRoutingParams, tagValue)
	}
}

func (d *OrderDecoder) decodeScaleOrderParams(msgDec *message.Decoder) {
	d.order.ScaleInitLevelSize = msgDec.Int32Max()
	d.order.ScaleSubsLevelSize = msgDec.Int32Max()
	d.order.ScalePriceIncrement = msgDec.FloatMax()
	if d.order.ScalePriceIncrement != nil && *d.order.ScalePriceIncrement > 0.0 {
		d.order.ScalePriceAdjustValue = msgDec.FloatMax()
		d.order.ScalePriceAdjustInterval = msgDec.Int32Max()
		d.order.ScaleProfitOffset = msgDec.FloatMax()
		d.order.ScaleAutoReset = msgDec.Bool()
		d.order.ScaleInitPosition = msgDec.Int32Max()
		d.order.ScaleInitFillQty = msgDec.Int32Max()
		d.order.ScaleRandomPercent = msgDec.Bool()
	}
}

func (d *OrderDecoder) decodeHedgeParams(msgDec *message.Decoder) {
	d.order.HedgeType = msgDec.String()
	if len(d.order.HedgeType) > 0 {
		d.order.HedgeParam = msgDec.String()
	}
}

func (d *OrderDecoder) decodeOptOutSmartRouting(msgDec *message.Decoder) {
	d.order.OptOutSmartRouting = msgDec.Bool()
}

func (d *OrderDecoder) decodeClearingParams(msgDec *message.Decoder) {
	d.order.ClearingAccount = msgDec.String()
	d.order.ClearingIntent = msgDec.String()
}

func (d *OrderDecoder) decodeNotHeld(msgDec *message.Decoder) {
	d.order.NotHeld = msgDec.Bool()
}

func (d *OrderDecoder) decodeDeltaNeutral(msgDec *message.Decoder) {
	deltaNeutralContractPresent := msgDec.Bool()
	if deltaNeutralContractPresent {
		d.contract.DeltaNeutralContract = &DeltaNeutralContract{}
		d.contract.DeltaNeutralContract.ConID = msgDec.Int32()
		d.contract.DeltaNeutralContract.Delta = msgDec.Float()
		d.contract.DeltaNeutralContract.Price = msgDec.Float()
	}
}

func (d *OrderDecoder) decodeAlgoParams(msgDec *message.Decoder) {
	d.order.AlgoStrategy = msgDec.String()
	if len(d.order.AlgoStrategy) > 0 {
		algoParamsCount := int(msgDec.Int32())
		if algoParamsCount < 0 {
			msgDec.SetErr(errors.New("invalid algo parameters count"))
			return
		}
		d.order.AlgoParams = make([]TagValue, 0, algoParamsCount)
		for i := 0; i < algoParamsCount; i++ {
			tagValue := TagValue{}
			tagValue.Tag = msgDec.String()
			tagValue.Value = msgDec.String()

			d.order.AlgoParams = append(d.order.AlgoParams, tagValue)
		}
	}
}

func (d *OrderDecoder) decodeSolicited(msgDec *message.Decoder) {
	d.order.Solicited = msgDec.Bool()
}

func (d *OrderDecoder) decodeWhatIfInfoAndCommissionAndFees(msgDec *message.Decoder) {
	d.order.WhatIf = msgDec.Bool()

	d.decodeOrderStatus(msgDec)

	d.orderState.InitMarginBefore = msgDec.String()
	d.orderState.MaintMarginBefore = msgDec.String()
	d.orderState.EquityWithLoanBefore = msgDec.String()
	d.orderState.InitMarginChange = msgDec.String()
	d.orderState.MaintMarginChange = msgDec.String()
	d.orderState.EquityWithLoanChange = msgDec.String()

	d.orderState.InitMarginAfter = msgDec.String()
	d.orderState.MaintMarginAfter = msgDec.String()
	d.orderState.EquityWithLoanAfter = msgDec.String()

	d.orderState.CommissionAndFees = msgDec.FloatMax()
	d.orderState.MinCommissionAndFees = msgDec.FloatMax()
	d.orderState.MaxCommissionAndFees = msgDec.FloatMax()
	d.orderState.CommissionAndFeesCurrency = msgDec.String()

	d.orderState.MarginCurrency = msgDec.String()
	d.orderState.InitMarginBeforeOutsideRTH = msgDec.FloatMax()
	d.orderState.MaintMarginBeforeOutsideRTH = msgDec.FloatMax()
	d.orderState.EquityWithLoanBeforeOutsideRTH = msgDec.FloatMax()
	d.orderState.InitMarginChangeOutsideRTH = msgDec.FloatMax()
	d.orderState.MaintMarginChangeOutsideRTH = msgDec.FloatMax()
	d.orderState.EquityWithLoanChangeOutsideRTH = msgDec.FloatMax()
	d.orderState.InitMarginAfterOutsideRTH = msgDec.FloatMax()
	d.orderState.MaintMarginAfterOutsideRTH = msgDec.FloatMax()
	d.orderState.EquityWithLoanAfterOutsideRTH = msgDec.FloatMax()
	d.orderState.SuggestedSize = NewDecimalMaxFromMessageDecoder(msgDec)
	d.orderState.RejectReason = msgDec.String()

	accountsCount := int(msgDec.Int32())
	if accountsCount < 0 {
		msgDec.SetErr(errors.New("invalid accounts count"))
		return
	}
	d.orderState.OrderAllocations = make([]*OrderAllocation, 0, accountsCount)
	for i := 0; i < accountsCount; i++ {
		oa := NewOrderAllocation()
		oa.Account = msgDec.String()
		oa.Position = NewDecimalMaxFromMessageDecoder(msgDec)
		oa.PositionDesired = NewDecimalMaxFromMessageDecoder(msgDec)
		oa.PositionAfter = NewDecimalMaxFromMessageDecoder(msgDec)
		oa.DesiredAllocQty = NewDecimalMaxFromMessageDecoder(msgDec)
		oa.AllowedAllocQty = NewDecimalMaxFromMessageDecoder(msgDec)
		oa.IsMonetary = msgDec.Bool()

		d.orderState.OrderAllocations = append(d.orderState.OrderAllocations, oa)
	}

	d.orderState.WarningText = msgDec.String()
}

func (d *OrderDecoder) decodeOrderStatus(msgDec *message.Decoder) {
	d.orderState.Status = msgDec.String()
}

func (d *OrderDecoder) decodeVolRandomizeFlags(msgDec *message.Decoder) {
	d.order.RandomizeSize = msgDec.Bool()
	d.order.RandomizePrice = msgDec.Bool()
}

func (d *OrderDecoder) decodePegBenchParams(msgDec *message.Decoder) {
	if d.order.OrderType == "PEG BENCH" {
		d.order.ReferenceContractID = msgDec.Int32()
		d.order.IsPeggedChangeAmountDecrease = msgDec.Bool()
		d.order.PeggedChangeAmount = msgDec.Float()
		d.order.ReferenceChangeAmount = msgDec.Float()
		d.order.ReferenceExchangeID = msgDec.String()
	}
}

func (d *OrderDecoder) decodeConditions(msgDec *message.Decoder) {
	conditionsSize := int(msgDec.Int32())
	if conditionsSize < 0 {
		msgDec.SetErr(errors.New("invalid conditions count"))
		return
	}
	d.order.Conditions = make([]OrderCondition, 0, conditionsSize)
	if conditionsSize > 0 {
		for i := 0; i < conditionsSize; i++ {
			cond := NewOrderConditionFromMessageDecoder(msgDec)
			if cond != nil {
				d.order.Conditions = append(d.order.Conditions, cond)
			}
		}

		d.order.ConditionsIgnoreRth = msgDec.Bool()
		d.order.ConditionsCancelOrder = msgDec.Bool()
	}
}

func (d *OrderDecoder) decodeAdjustedOrderParams(msgDec *message.Decoder) {
	d.order.AdjustedOrderType = msgDec.String()
	d.order.TriggerPrice = msgDec.FloatMax()
	d.decodeStopPriceAndLmtPriceOffset(msgDec)
	d.order.AdjustedStopPrice = msgDec.FloatMax()
	d.order.AdjustedStopLimitPrice = msgDec.FloatMax()
	d.order.AdjustedTrailingAmount = msgDec.FloatMax()
	d.order.AdjustableTrailingUnit = msgDec.Int32()
}

func (d *OrderDecoder) decodeStopPriceAndLmtPriceOffset(msgDec *message.Decoder) {
	d.order.TrailStopPrice = msgDec.FloatMax()
	d.order.LmtPriceOffset = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeSoftDollarTier(msgDec *message.Decoder) {
	d.order.SoftDollarTier = SoftDollarTier{
		Name:        msgDec.String(),
		Value:       msgDec.String(),
		DisplayName: msgDec.String(),
	}
}

func (d *OrderDecoder) decodeCashQty(msgDec *message.Decoder) {
	d.order.CashQty = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeDontUseAutoPriceForHedge(msgDec *message.Decoder) {
	d.order.DontUseAutoPriceForHedge = msgDec.Bool()
}

func (d *OrderDecoder) decodeIsOmsContainer(msgDec *message.Decoder) {
	d.order.IsOmsContainer = msgDec.Bool()
}

func (d *OrderDecoder) decodeDiscretionaryUpToLimitPrice(msgDec *message.Decoder) {
	d.order.DiscretionaryUpToLimitPrice = msgDec.Bool()
}

func (d *OrderDecoder) decodeAutoCancelDate(msgDec *message.Decoder) {
	d.order.AutoCancelDate = msgDec.String()
}

func (d *OrderDecoder) decodeFilledQuantity(msgDec *message.Decoder) {
	d.order.FilledQuantity = NewDecimalMaxFromMessageDecoder(msgDec)
}

func (d *OrderDecoder) decodeRefFuturesConId(msgDec *message.Decoder) {
	d.order.RefFuturesConID = msgDec.Int32()
}

func (d *OrderDecoder) decodeAutoCancelParent(msgDec *message.Decoder) {
	d.order.AutoCancelParent = msgDec.Bool()
}

func (d *OrderDecoder) decodeShareholder(msgDec *message.Decoder) {
	d.order.Shareholder = msgDec.String()
}

func (d *OrderDecoder) decodeImbalanceOnly(msgDec *message.Decoder) {
	d.order.ImbalanceOnly = msgDec.Bool()
}

func (d *OrderDecoder) decodeRouteMarketableToBbo(msgDec *message.Decoder) {
	d.order.RouteMarketableToBbo = msgDec.Bool()
}

func (d *OrderDecoder) decodeParentPermId(msgDec *message.Decoder) {
	d.order.ParentPermID = msgDec.Int64()
}

func (d *OrderDecoder) decodeCompletedTime(msgDec *message.Decoder) {
	d.orderState.CompletedTime = msgDec.String()
}

func (d *OrderDecoder) decodeCompletedStatus(msgDec *message.Decoder) {
	d.orderState.CompletedStatus = msgDec.String()
}

func (d *OrderDecoder) decodeUsePriceMgmtAlgo(msgDec *message.Decoder) {
	d.order.UsePriceMgmtAlgo = msgDec.Bool()
}

func (d *OrderDecoder) decodeDuration(msgDec *message.Decoder) {
	d.order.Duration = msgDec.Int32Max()
}

func (d *OrderDecoder) decodePostToAts(msgDec *message.Decoder) {
	d.order.PostToAts = msgDec.Int32Max()
}

func (d *OrderDecoder) decodePegBestPegMidOrderAttributes(msgDec *message.Decoder) {
	d.order.MinTradeQty = msgDec.Int32Max()
	d.order.MinCompeteSize = msgDec.Int32Max()
	d.order.CompeteAgainstBestOffset = msgDec.FloatMax()
	d.order.MidOffsetAtWhole = msgDec.FloatMax()
	d.order.MidOffsetAtHalf = msgDec.FloatMax()
}

func (d *OrderDecoder) decodeCustomerAccount(msgDec *message.Decoder) {
	d.order.CustomerAccount = msgDec.String()
}

func (d *OrderDecoder) decodeProfessionalCustomer(msgDec *message.Decoder) {
	d.order.ProfessionalCustomer = msgDec.Bool()
}

func (d *OrderDecoder) decodeBondAccruedInterest(msgDec *message.Decoder) {
	d.order.BondAccruedInterest = msgDec.String()
}

func (d *OrderDecoder) decodeIncludeOvernight(msgDec *message.Decoder) {
	d.order.IncludeOvernight = msgDec.Bool()
}

func (d *OrderDecoder) decodeCMETaggingFields(msgDec *message.Decoder) {
	d.order.ExtOperator = msgDec.String()
	d.order.ManualOrderIndicator = msgDec.Int32Max()
}

func (d *OrderDecoder) decodeSubmitter(msgDec *message.Decoder) {
	d.order.Submitter = msgDec.String()
}
