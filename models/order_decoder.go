package models

import (
	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
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

func (d *OrderDecoder) decodeOrderId(msgDec *utils.MessageDecoder) {
	d.order.OrderID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeContractFields(msgDec *utils.MessageDecoder) {
	d.contract.ConID = msgDec.Int64(false)
	d.contract.Symbol = msgDec.String(false)
	d.contract.SecType = NewSecurityTypeFromString(msgDec.String(false))
	d.contract.LastTradeDateOrContractMonth = msgDec.String(false)
	d.contract.Strike = msgDec.Float64(false)
	d.contract.Right = msgDec.String(false)
	d.contract.Multiplier = msgDec.String(false)
	d.contract.Exchange = msgDec.String(false)
	d.contract.Currency = msgDec.String(false)
	d.contract.LocalSymbol = msgDec.String(false)
	d.contract.TradingClass = msgDec.String(false)
}

func (d *OrderDecoder) decodeAction(msgDec *utils.MessageDecoder) {
	d.order.Action = msgDec.String(false)
}

func (d *OrderDecoder) decodeTotalQuantity(msgDec *utils.MessageDecoder) {
	d.order.TotalQuantity = NewDecimalFromMessageDecoder(msgDec, false)
}

func (d *OrderDecoder) decodeOrderType(msgDec *utils.MessageDecoder) {
	d.order.OrderType = msgDec.String(false)
}

func (d *OrderDecoder) decodeLmtPrice(msgDec *utils.MessageDecoder) {
	d.order.LmtPrice = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeAuxPrice(msgDec *utils.MessageDecoder) {
	d.order.AuxPrice = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeTIF(msgDec *utils.MessageDecoder) {
	d.order.TIF = msgDec.String(false)
}

func (d *OrderDecoder) decodeOcaGroup(msgDec *utils.MessageDecoder) {
	d.order.OCAGroup = msgDec.String(false)
}

func (d *OrderDecoder) decodeAccount(msgDec *utils.MessageDecoder) {
	d.order.Account = msgDec.String(false)
}

func (d *OrderDecoder) decodeOpenClose(msgDec *utils.MessageDecoder) {
	d.order.OpenClose = msgDec.String(false)
}

func (d *OrderDecoder) decodeOrigin(msgDec *utils.MessageDecoder) {
	d.order.Origin = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeOrderRef(msgDec *utils.MessageDecoder) {
	d.order.OrderRef = msgDec.String(false)
}

func (d *OrderDecoder) decodeClientId(msgDec *utils.MessageDecoder) {
	d.order.ClientID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodePermId(msgDec *utils.MessageDecoder) {
	d.order.PermID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeOutsideRth(msgDec *utils.MessageDecoder) {
	d.order.OutsideRTH = msgDec.Bool()
}

func (d *OrderDecoder) decodeHidden(msgDec *utils.MessageDecoder) {
	d.order.Hidden = msgDec.Bool()
}

func (d *OrderDecoder) decodeDiscretionaryAmount(msgDec *utils.MessageDecoder) {
	d.order.DiscretionaryAmt = msgDec.Float64(false)
}

func (d *OrderDecoder) decodeGoodAfterTime(msgDec *utils.MessageDecoder) {
	d.order.GoodAfterTime = msgDec.String(false)
}

func (d *OrderDecoder) skipSharesAllocation(msgDec *utils.MessageDecoder) {
	msgDec.Skip()
}

func (d *OrderDecoder) decodeFAParams(msgDec *utils.MessageDecoder) {
	d.order.FAGroup = msgDec.String(false)
	d.order.FAMethod = msgDec.String(false)
	d.order.FAPercentage = msgDec.String(false)
	msgDec.Skip() // skip deprecated FAProfile
}

func (d *OrderDecoder) decodeModelCode(msgDec *utils.MessageDecoder) {
	d.order.ModelCode = msgDec.String(false)
}

func (d *OrderDecoder) decodeGoodTillDate(msgDec *utils.MessageDecoder) {
	d.order.GoodTillDate = msgDec.String(false)
}

func (d *OrderDecoder) decodeRule80A(msgDec *utils.MessageDecoder) {
	d.order.Rule80A = msgDec.String(false)
}

func (d *OrderDecoder) decodePercentOffset(msgDec *utils.MessageDecoder) {
	d.order.PercentOffset = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeSettlingFirm(msgDec *utils.MessageDecoder) {
	d.order.SettlingFirm = msgDec.String(false)
}

func (d *OrderDecoder) decodeShortSaleParams(msgDec *utils.MessageDecoder) {
	d.order.ShortSaleSlot = msgDec.Int64(false)
	d.order.DesignatedLocation = msgDec.String(false)
	d.order.ExemptCode = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeAuctionStrategy(msgDec *utils.MessageDecoder) {
	d.order.AuctionStrategy = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeBoxOrderParams(msgDec *utils.MessageDecoder) {
	d.order.StartingPrice = msgDec.Float64(true)
	d.order.StockRefPrice = msgDec.Float64(true)
	d.order.Delta = msgDec.Float64(true)
}

func (d *OrderDecoder) decodePegToStkOrVolOrderParams(msgDec *utils.MessageDecoder) {
	d.order.StockRangeLower = msgDec.Float64(true)
	d.order.StockRangeUpper = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeDisplaySize(msgDec *utils.MessageDecoder) {
	d.order.DisplaySize = msgDec.Int64(true)
}

func (d *OrderDecoder) decodeBlockOrder(msgDec *utils.MessageDecoder) {
	d.order.BlockOrder = msgDec.Bool()
}

func (d *OrderDecoder) decodeSweepToFill(msgDec *utils.MessageDecoder) {
	d.order.SweepToFill = msgDec.Bool()
}

func (d *OrderDecoder) decodeAllOrNone(msgDec *utils.MessageDecoder) {
	d.order.AllOrNone = msgDec.Bool()
}

func (d *OrderDecoder) decodeMinQty(msgDec *utils.MessageDecoder) {
	d.order.MinQty = msgDec.Int64(true)
}

func (d *OrderDecoder) decodeOcaType(msgDec *utils.MessageDecoder) {
	d.order.OCAType = msgDec.Int64(false)
}

func (d *OrderDecoder) skipETradeOnly(msgDec *utils.MessageDecoder) {
	msgDec.Skip() // deprecated order.ETradeOnly
}

func (d *OrderDecoder) skipFirmQuoteOnly(msgDec *utils.MessageDecoder) {
	msgDec.Skip() // deprecated order.FirmQuoteOnly
}

func (d *OrderDecoder) skipNbboPriceCap(msgDec *utils.MessageDecoder) {
	msgDec.Skip() // deprecated order.NBBOPriceCap
}

func (d *OrderDecoder) decodeParentId(msgDec *utils.MessageDecoder) {
	d.order.ParentID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeTriggerMethod(msgDec *utils.MessageDecoder) {
	d.order.TriggerMethod = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeVolOrderParams(msgDec *utils.MessageDecoder, readOpenOrderAttribs bool) {
	d.order.Volatility = msgDec.Float64(true)
	d.order.VolatilityType = msgDec.Int64(false)
	d.order.DeltaNeutralOrderType = msgDec.String(false)
	d.order.DeltaNeutralAuxPrice = msgDec.Float64(true)

	if len(d.order.DeltaNeutralOrderType) > 0 {
		d.order.DeltaNeutralConID = msgDec.Int64(false)
		if readOpenOrderAttribs {
			d.order.DeltaNeutralSettlingFirm = msgDec.String(false)
			d.order.DeltaNeutralClearingAccount = msgDec.String(false)
			d.order.DeltaNeutralClearingIntent = msgDec.String(false)
			d.order.DeltaNeutralOpenClose = msgDec.String(false)
		}

		d.order.DeltaNeutralShortSale = msgDec.Bool()
		d.order.DeltaNeutralShortSaleSlot = msgDec.Int64(false)
		d.order.DeltaNeutralDesignatedLocation = msgDec.String(false)
	}

	d.order.ContinuousUpdate = msgDec.Bool()
	d.order.ReferencePriceType = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeTrailParams(msgDec *utils.MessageDecoder) {
	d.order.TrailStopPrice = msgDec.Float64(true)
	d.order.TrailingPercent = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeBasisPoints(msgDec *utils.MessageDecoder) {
	d.order.BasisPoints = msgDec.Float64(true)
	d.order.BasisPointsType = msgDec.Int64(true)
}

func (d *OrderDecoder) decodeComboLegs(msgDec *utils.MessageDecoder) {
	d.contract.ComboLegsDescription = msgDec.String(false)

	comboLegsCount := msgDec.Int64(false)
	d.contract.ComboLegs = make([]ComboLeg, 0, comboLegsCount)
	for i := int64(0); i < comboLegsCount; i++ {
		comboleg := ComboLeg{}
		comboleg.ConID = msgDec.Int64(false)
		comboleg.Ratio = msgDec.Int64(false)
		comboleg.Action = msgDec.String(false)
		comboleg.Exchange = msgDec.String(false)
		comboleg.OpenClose = msgDec.Int64(false)
		comboleg.ShortSaleSlot = msgDec.Int64(false)
		comboleg.DesignatedLocation = msgDec.String(false)
		comboleg.ExemptCode = msgDec.Int64(false)

		d.contract.ComboLegs = append(d.contract.ComboLegs, comboleg)
	}

	orderComboLegsCount := msgDec.Int64(false)
	d.order.OrderComboLegs = make([]OrderComboLeg, 0, orderComboLegsCount)
	for i := int64(0); i < orderComboLegsCount; i++ {
		orderComboLeg := OrderComboLeg{}
		orderComboLeg.Price = msgDec.Float64(true)

		d.order.OrderComboLegs = append(d.order.OrderComboLegs, orderComboLeg)
	}
}

func (d *OrderDecoder) decodeSmartComboRoutingParams(msgDec *utils.MessageDecoder) {
	smartComboRoutingParamsCount := msgDec.Int64(false)
	d.order.SmartComboRoutingParams = make([]TagValue, 0, smartComboRoutingParamsCount)
	for i := int64(0); i < smartComboRoutingParamsCount; i++ {
		tagValue := TagValue{}
		tagValue.Tag = msgDec.String(false)
		tagValue.Value = msgDec.String(false)

		d.order.SmartComboRoutingParams = append(d.order.SmartComboRoutingParams, tagValue)
	}
}

func (d *OrderDecoder) decodeScaleOrderParams(msgDec *utils.MessageDecoder) {
	d.order.ScaleInitLevelSize = msgDec.Int64(true)
	d.order.ScaleSubsLevelSize = msgDec.Int64(true)
	d.order.ScalePriceIncrement = msgDec.Float64(true)
	if d.order.ScalePriceIncrement != common.UNSET_FLOAT && d.order.ScalePriceIncrement > 0.0 {
		d.order.ScalePriceAdjustValue = msgDec.Float64(true)
		d.order.ScalePriceAdjustInterval = msgDec.Int64(true)
		d.order.ScaleProfitOffset = msgDec.Float64(true)
		d.order.ScaleAutoReset = msgDec.Bool()
		d.order.ScaleInitPosition = msgDec.Int64(true)
		d.order.ScaleInitFillQty = msgDec.Int64(true)
		d.order.ScaleRandomPercent = msgDec.Bool()
	}
}

func (d *OrderDecoder) decodeHedgeParams(msgDec *utils.MessageDecoder) {
	d.order.HedgeType = msgDec.String(false)
	if len(d.order.HedgeType) > 0 {
		d.order.HedgeParam = msgDec.String(false)
	}
}

func (d *OrderDecoder) decodeOptOutSmartRouting(msgDec *utils.MessageDecoder) {
	d.order.OptOutSmartRouting = msgDec.Bool()
}

func (d *OrderDecoder) decodeClearingParams(msgDec *utils.MessageDecoder) {
	d.order.ClearingAccount = msgDec.String(false)
	d.order.ClearingIntent = msgDec.String(false)
}

func (d *OrderDecoder) decodeNotHeld(msgDec *utils.MessageDecoder) {
	d.order.NotHeld = msgDec.Bool()
}

func (d *OrderDecoder) decodeDeltaNeutral(msgDec *utils.MessageDecoder) {
	deltaNeutralContractPresent := msgDec.Bool()
	if deltaNeutralContractPresent {
		d.contract.DeltaNeutralContract = &DeltaNeutralContract{}
		d.contract.DeltaNeutralContract.ConID = msgDec.Int64(false)
		d.contract.DeltaNeutralContract.Delta = msgDec.Float64(false)
		d.contract.DeltaNeutralContract.Price = msgDec.Float64(false)
	}
}

func (d *OrderDecoder) decodeAlgoParams(msgDec *utils.MessageDecoder) {
	d.order.AlgoStrategy = msgDec.String(false)
	if len(d.order.AlgoStrategy) > 0 {
		algoParamsCount := msgDec.Int64(false)
		d.order.AlgoParams = make([]TagValue, 0, algoParamsCount)
		for i := int64(0); i < algoParamsCount; i++ {
			tagValue := TagValue{}
			tagValue.Tag = msgDec.String(false)
			tagValue.Value = msgDec.String(false)

			d.order.AlgoParams = append(d.order.AlgoParams, tagValue)
		}
	}
}

func (d *OrderDecoder) decodeSolicited(msgDec *utils.MessageDecoder) {
	d.order.Solictied = msgDec.Bool()
}

func (d *OrderDecoder) decodeWhatIfInfoAndCommissionAndFees(msgDec *utils.MessageDecoder) {
	d.order.WhatIf = msgDec.Bool()

	d.decodeOrderStatus(msgDec)

	d.orderState.InitMarginBefore = msgDec.String(false)
	d.orderState.MaintMarginBefore = msgDec.String(false)
	d.orderState.EquityWithLoanBefore = msgDec.String(false)
	d.orderState.InitMarginChange = msgDec.String(false)
	d.orderState.MaintMarginChange = msgDec.String(false)
	d.orderState.EquityWithLoanChange = msgDec.String(false)

	d.orderState.InitMarginAfter = msgDec.String(false)
	d.orderState.MaintMarginAfter = msgDec.String(false)
	d.orderState.EquityWithLoanAfter = msgDec.String(false)

	d.orderState.CommissionAndFees = msgDec.Float64(true)
	d.orderState.MinCommissionAndFees = msgDec.Float64(true)
	d.orderState.MaxCommissionAndFees = msgDec.Float64(true)
	d.orderState.CommissionAndFeesCurrency = msgDec.String(false)

	d.orderState.MarginCurrency = msgDec.String(false)
	d.orderState.InitMarginBeforeOutsideRTH = msgDec.Float64(true)
	d.orderState.MaintMarginBeforeOutsideRTH = msgDec.Float64(true)
	d.orderState.EquityWithLoanBeforeOutsideRTH = msgDec.Float64(true)
	d.orderState.InitMarginChangeOutsideRTH = msgDec.Float64(true)
	d.orderState.MaintMarginChangeOutsideRTH = msgDec.Float64(true)
	d.orderState.EquityWithLoanChangeOutsideRTH = msgDec.Float64(true)
	d.orderState.InitMarginAfterOutsideRTH = msgDec.Float64(true)
	d.orderState.MaintMarginAfterOutsideRTH = msgDec.Float64(true)
	d.orderState.EquityWithLoanAfterOutsideRTH = msgDec.Float64(true)
	d.orderState.SuggestedSize = NewDecimalFromMessageDecoder(msgDec, false)
	d.orderState.RejectReason = msgDec.String(false)

	accountsCount := msgDec.Int64(false)
	d.orderState.OrderAllocations = make([]*OrderAllocation, 0, accountsCount)
	for i := int64(0); i < accountsCount; i++ {
		oa := NewOrderAllocation()
		oa.Account = msgDec.String(false)
		oa.Position = NewDecimalFromMessageDecoder(msgDec, false)
		oa.PositionDesired = NewDecimalFromMessageDecoder(msgDec, false)
		oa.PositionAfter = NewDecimalFromMessageDecoder(msgDec, false)
		oa.DesiredAllocQty = NewDecimalFromMessageDecoder(msgDec, false)
		oa.AllowedAllocQty = NewDecimalFromMessageDecoder(msgDec, false)
		oa.IsMonetary = msgDec.Bool()

		d.orderState.OrderAllocations = append(d.orderState.OrderAllocations, oa)
	}

	d.orderState.WarningText = msgDec.String(false)
}

func (d *OrderDecoder) decodeOrderStatus(msgDec *utils.MessageDecoder) {
	d.orderState.Status = msgDec.String(false)
}

func (d *OrderDecoder) decodeVolRandomizeFlags(msgDec *utils.MessageDecoder) {
	d.order.RandomizeSize = msgDec.Bool()
	d.order.RandomizePrice = msgDec.Bool()
}

func (d *OrderDecoder) decodePegBenchParams(msgDec *utils.MessageDecoder) {
	if d.order.OrderType == "PEG BENCH" {
		d.order.ReferenceContractID = msgDec.Int64(false)
		d.order.IsPeggedChangeAmountDecrease = msgDec.Bool()
		d.order.PeggedChangeAmount = msgDec.Float64(false)
		d.order.ReferenceChangeAmount = msgDec.Float64(false)
		d.order.ReferenceExchangeID = msgDec.String(false)
	}
}

func (d *OrderDecoder) decodeConditions(msgDec *utils.MessageDecoder) {
	conditionsSize := msgDec.Int64(false)
	d.order.Conditions = make([]OrderCondition, 0, conditionsSize)
	if conditionsSize > 0 {
		for i := int64(0); i < conditionsSize; i++ {
			cond := NewOrderConditionFromMessageDecoder(msgDec)
			if cond != nil {
				d.order.Conditions = append(d.order.Conditions, cond)
			}
		}

		d.order.ConditionsIgnoreRth = msgDec.Bool()
		d.order.ConditionsCancelOrder = msgDec.Bool()
	}
}

func (d *OrderDecoder) decodeAdjustedOrderParams(msgDec *utils.MessageDecoder) {
	d.order.AdjustedOrderType = msgDec.String(false)
	d.order.TriggerPrice = msgDec.Float64(false)
	d.decodeStopPriceAndLmtPriceOffset(msgDec)
	d.order.AdjustedStopPrice = msgDec.Float64(false)
	d.order.AdjustedStopLimitPrice = msgDec.Float64(false)
	d.order.AdjustedTrailingAmount = msgDec.Float64(false)
	d.order.AdjustableTrailingUnit = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeStopPriceAndLmtPriceOffset(msgDec *utils.MessageDecoder) {
	d.order.TrailStopPrice = msgDec.Float64(false)
	d.order.LmtPriceOffset = msgDec.Float64(false)
}

func (d *OrderDecoder) decodeSoftDollarTier(msgDec *utils.MessageDecoder) {
	d.order.SoftDollarTier = SoftDollarTier{
		Name:        msgDec.String(false),
		Value:       msgDec.String(false),
		DisplayName: msgDec.String(false),
	}
}

func (d *OrderDecoder) decodeCashQty(msgDec *utils.MessageDecoder) {
	d.order.CashQty = msgDec.Float64(false)
}

func (d *OrderDecoder) decodeDontUseAutoPriceForHedge(msgDec *utils.MessageDecoder) {
	d.order.DontUseAutoPriceForHedge = msgDec.Bool()
}

func (d *OrderDecoder) decodeIsOmsContainer(msgDec *utils.MessageDecoder) {
	d.order.IsOmsContainer = msgDec.Bool()
}

func (d *OrderDecoder) decodeDiscretionaryUpToLimitPrice(msgDec *utils.MessageDecoder) {
	d.order.DiscretionaryUpToLimitPrice = msgDec.Bool()
}

func (d *OrderDecoder) decodeAutoCancelDate(msgDec *utils.MessageDecoder) {
	d.order.AutoCancelDate = msgDec.String(false)
}

func (d *OrderDecoder) decodeFilledQuantity(msgDec *utils.MessageDecoder) {
	d.order.FilledQuantity = NewDecimalFromMessageDecoder(msgDec, false)
}

func (d *OrderDecoder) decodeRefFuturesConId(msgDec *utils.MessageDecoder) {
	d.order.RefFuturesConID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeAutoCancelParent(msgDec *utils.MessageDecoder) {
	d.order.AutoCancelParent = msgDec.Bool()
}

func (d *OrderDecoder) decodeShareholder(msgDec *utils.MessageDecoder) {
	d.order.Shareholder = msgDec.String(false)
}

func (d *OrderDecoder) decodeImbalanceOnly(msgDec *utils.MessageDecoder) {
	d.order.ImbalanceOnly = msgDec.Bool()
}

func (d *OrderDecoder) decodeRouteMarketableToBbo(msgDec *utils.MessageDecoder) {
	d.order.RouteMarketableToBbo = msgDec.Bool()
}

func (d *OrderDecoder) decodeParentPermId(msgDec *utils.MessageDecoder) {
	d.order.ParentPermID = msgDec.Int64(false)
}

func (d *OrderDecoder) decodeCompletedTime(msgDec *utils.MessageDecoder) {
	d.orderState.CompletedTime = msgDec.String(false)
}

func (d *OrderDecoder) decodeCompletedStatus(msgDec *utils.MessageDecoder) {
	d.orderState.CompletedStatus = msgDec.String(false)
}

func (d *OrderDecoder) decodeUsePriceMgmtAlgo(msgDec *utils.MessageDecoder) {
	d.order.UsePriceMgmtAlgo = msgDec.Bool()
}

func (d *OrderDecoder) decodeDuration(msgDec *utils.MessageDecoder) {
	d.order.Duration = msgDec.Int64(true)
}

func (d *OrderDecoder) decodePostToAts(msgDec *utils.MessageDecoder) {
	d.order.PostToAts = msgDec.Int64(true)
}

func (d *OrderDecoder) decodePegBestPegMidOrderAttributes(msgDec *utils.MessageDecoder) {
	d.order.MinTradeQty = msgDec.Int64(true)
	d.order.MinCompeteSize = msgDec.Int64(true)
	d.order.CompeteAgainstBestOffset = msgDec.Float64(true)
	d.order.MidOffsetAtWhole = msgDec.Float64(true)
	d.order.MidOffsetAtHalf = msgDec.Float64(true)
}

func (d *OrderDecoder) decodeCustomerAccount(msgDec *utils.MessageDecoder) {
	d.order.CustomerAccount = msgDec.String(false)
}

func (d *OrderDecoder) decodeProfessionalCustomer(msgDec *utils.MessageDecoder) {
	d.order.ProfessionalCustomer = msgDec.Bool()
}

func (d *OrderDecoder) decodeBondAccruedInterest(msgDec *utils.MessageDecoder) {
	d.order.BondAccruedInterest = msgDec.String(false)
}

func (d *OrderDecoder) decodeIncludeOvernight(msgDec *utils.MessageDecoder) {
	d.order.IncludeOvernight = msgDec.Bool()
}

func (d *OrderDecoder) decodeCMETaggingFields(msgDec *utils.MessageDecoder) {
	d.order.ExtOperator = msgDec.String(false)
	d.order.ManualOrderIndicator = msgDec.Int64(true)
}

func (d *OrderDecoder) decodeSubmitter(msgDec *utils.MessageDecoder) {
	d.order.Submitter = msgDec.String(false)
}
