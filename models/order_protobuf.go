package models

import (
	"strconv"
	"time"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// NewOrderFromProtobufDecoder decodes an order from protobuf data.
func NewOrderFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.Order) *Order {
	order := NewOrder()
	if pb == nil {
		return order
	}

	order.ClientID = msgDec.Int32(pb.ClientId)
	order.OrderID = msgDec.Int32(pb.OrderId)
	order.PermID = msgDec.Int64(pb.PermId)
	order.ParentID = msgDec.Int32(pb.ParentId)
	order.Action = msgDec.String(pb.Action)
	order.TotalQuantity = NewDecimalMaxFromProtobufDecoder(msgDec, pb.TotalQuantity)
	order.DisplaySize = msgDec.Int32(pb.DisplaySize)
	order.OrderType = msgDec.String(pb.OrderType)
	order.LmtPrice = msgDec.FloatMax(pb.LmtPrice)
	order.AuxPrice = msgDec.FloatMax(pb.AuxPrice)
	order.TIF = TimeInForce(msgDec.String(pb.Tif))
	order.Account = msgDec.String(pb.Account)
	order.SettlingFirm = msgDec.String(pb.SettlingFirm)
	order.ClearingAccount = msgDec.String(pb.ClearingAccount)
	order.ClearingIntent = msgDec.String(pb.ClearingIntent)
	order.AllOrNone = msgDec.Bool(pb.AllOrNone)
	order.BlockOrder = msgDec.Bool(pb.BlockOrder)
	order.Hidden = msgDec.Bool(pb.Hidden)
	order.OutsideRTH = msgDec.Bool(pb.OutsideRth)
	order.SweepToFill = msgDec.Bool(pb.SweepToFill)
	order.PercentOffset = msgDec.FloatMax(pb.PercentOffset)
	order.TrailingPercent = msgDec.FloatMax(pb.TrailingPercent)
	order.TrailStopPrice = msgDec.FloatMax(pb.TrailStopPrice)
	order.MinQty = msgDec.Int32Max(pb.MinQty)
	order.GoodAfterTime = msgDec.String(pb.GoodAfterTime)
	order.GoodTillDate = msgDec.String(pb.GoodTillDate)
	order.OCAGroup = msgDec.String(pb.OcaGroup)
	order.OrderRef = msgDec.String(pb.OrderRef)
	order.Rule80A = Rule80A(msgDec.String(pb.Rule80A))
	order.OCAType = Oca(msgDec.Int32(pb.OcaType))
	order.TriggerMethod = TriggerMethod(msgDec.Int32(pb.TriggerMethod))
	order.ActiveStartTime = msgDec.String(pb.ActiveStartTime)
	order.ActiveStopTime = msgDec.String(pb.ActiveStopTime)
	order.FAGroup = msgDec.String(pb.FaGroup)
	order.FAMethod = msgDec.String(pb.FaMethod)
	order.FAPercentage = msgDec.String(pb.FaPercentage)
	order.Volatility = msgDec.FloatMax(pb.Volatility)
	order.VolatilityType = Volatility(msgDec.Int32(pb.VolatilityType))
	order.ContinuousUpdate = msgDec.Bool(pb.ContinuousUpdate)
	order.ReferencePriceType = msgDec.Int32(pb.ReferencePriceType)
	order.DeltaNeutralOrderType = msgDec.String(pb.DeltaNeutralOrderType)
	order.DeltaNeutralAuxPrice = msgDec.FloatMax(pb.DeltaNeutralAuxPrice)
	order.DeltaNeutralConID = msgDec.Int32(pb.DeltaNeutralConId)
	order.DeltaNeutralOpenClose = msgDec.String(pb.DeltaNeutralOpenClose)
	order.DeltaNeutralShortSale = msgDec.Bool(pb.DeltaNeutralShortSale)
	order.DeltaNeutralShortSaleSlot = msgDec.Int32(pb.DeltaNeutralShortSaleSlot)
	order.DeltaNeutralDesignatedLocation = msgDec.String(pb.DeltaNeutralDesignatedLocation)
	order.ScaleInitLevelSize = msgDec.Int32Max(pb.ScaleInitLevelSize)
	order.ScaleSubsLevelSize = msgDec.Int32Max(pb.ScaleSubsLevelSize)
	order.ScalePriceIncrement = msgDec.FloatMax(pb.ScalePriceIncrement)
	order.ScalePriceAdjustValue = msgDec.FloatMax(pb.ScalePriceAdjustValue)
	order.ScalePriceAdjustInterval = msgDec.Int32Max(pb.ScalePriceAdjustInterval)
	order.ScaleProfitOffset = msgDec.FloatMax(pb.ScaleProfitOffset)
	order.ScaleAutoReset = msgDec.Bool(pb.ScaleAutoReset)
	order.ScaleInitPosition = msgDec.Int32Max(pb.ScaleInitPosition)
	order.ScaleInitFillQty = msgDec.Int32Max(pb.ScaleInitFillQty)
	order.ScaleRandomPercent = msgDec.Bool(pb.ScaleRandomPercent)
	order.ScaleTable = msgDec.String(pb.ScaleTable)
	order.HedgeType = msgDec.String(pb.HedgeType)
	order.HedgeParam = msgDec.String(pb.HedgeParam)
	order.AlgoStrategy = msgDec.String(pb.AlgoStrategy)
	order.AlgoParams = tagValueMapFromProtobuf(pb.AlgoParams)
	order.AlgoID = msgDec.String(pb.AlgoId)
	order.SmartComboRoutingParams = tagValueMapFromProtobuf(pb.SmartComboRoutingParams)
	order.WhatIf = msgDec.Bool(pb.WhatIf)
	order.Transmit = msgDec.Bool(pb.Transmit)
	order.OverridePercentageConstraints = msgDec.Bool(pb.OverridePercentageConstraints)
	order.OpenClose = msgDec.String(pb.OpenClose)
	order.Origin = msgDec.Int32(pb.Origin)
	order.ShortSaleSlot = msgDec.Int32(pb.ShortSaleSlot)
	order.DesignatedLocation = msgDec.String(pb.DesignatedLocation)
	order.ExemptCode = msgDec.Int32(pb.ExemptCode)
	order.DeltaNeutralSettlingFirm = msgDec.String(pb.DeltaNeutralSettlingFirm)
	order.DeltaNeutralClearingAccount = msgDec.String(pb.DeltaNeutralClearingAccount)
	order.DeltaNeutralClearingIntent = msgDec.String(pb.DeltaNeutralClearingIntent)
	order.DiscretionaryAmt = msgDec.Float(pb.DiscretionaryAmt)
	order.OptOutSmartRouting = msgDec.Bool(pb.OptOutSmartRouting)
	order.StartingPrice = msgDec.FloatMax(pb.StartingPrice)
	order.StockRefPrice = msgDec.FloatMax(pb.StockRefPrice)
	order.Delta = msgDec.FloatMax(pb.Delta)
	order.StockRangeLower = msgDec.FloatMax(pb.StockRangeLower)
	order.StockRangeUpper = msgDec.FloatMax(pb.StockRangeUpper)
	order.NotHeld = msgDec.Bool(pb.NotHeld)
	order.OrderMiscOptions = tagValueMapFromProtobuf(pb.OrderMiscOptions)
	order.Solicited = msgDec.Bool(pb.Solicited)
	order.RandomizeSize = msgDec.Bool(pb.RandomizeSize)
	order.RandomizePrice = msgDec.Bool(pb.RandomizePrice)
	order.ReferenceContractID = msgDec.Int32(pb.ReferenceContractId)
	order.PeggedChangeAmount = msgDec.Float(pb.PeggedChangeAmount)
	order.IsPeggedChangeAmountDecrease = msgDec.Bool(pb.IsPeggedChangeAmountDecrease)
	order.ReferenceChangeAmount = msgDec.Float(pb.ReferenceChangeAmount)
	order.ReferenceExchangeID = msgDec.String(pb.ReferenceExchangeId)
	order.AdjustedOrderType = msgDec.String(pb.AdjustedOrderType)
	order.TriggerPrice = msgDec.FloatMax(pb.TriggerPrice)
	order.AdjustedStopPrice = msgDec.FloatMax(pb.AdjustedStopPrice)
	order.AdjustedStopLimitPrice = msgDec.FloatMax(pb.AdjustedStopLimitPrice)
	order.AdjustedTrailingAmount = msgDec.FloatMax(pb.AdjustedTrailingAmount)
	order.AdjustableTrailingUnit = msgDec.Int32(pb.AdjustableTrailingUnit)
	order.LmtPriceOffset = msgDec.FloatMax(pb.LmtPriceOffset)
	order.Conditions = orderConditionsFromProtobuf(msgDec, pb.Conditions)
	order.ConditionsCancelOrder = msgDec.Bool(pb.ConditionsCancelOrder)
	order.ConditionsIgnoreRth = msgDec.Bool(pb.ConditionsIgnoreRth)
	order.ModelCode = msgDec.String(pb.ModelCode)
	order.ExtOperator = msgDec.String(pb.ExtOperator)
	order.SoftDollarTier = NewSoftDollarTierFromProtobufDecoder(msgDec, pb.SoftDollarTier)
	order.CashQty = msgDec.FloatMax(pb.CashQty)
	order.Mifid2DecisionMaker = msgDec.String(pb.Mifid2DecisionMaker)
	order.Mifid2DecisionAlgo = msgDec.String(pb.Mifid2DecisionAlgo)
	order.Mifid2ExecutionTrader = msgDec.String(pb.Mifid2ExecutionTrader)
	order.Mifid2ExecutionAlgo = msgDec.String(pb.Mifid2ExecutionAlgo)
	order.DontUseAutoPriceForHedge = msgDec.Bool(pb.DontUseAutoPriceForHedge)
	order.IsOmsContainer = msgDec.Bool(pb.IsOmsContainer)
	order.DiscretionaryUpToLimitPrice = msgDec.Bool(pb.DiscretionaryUpToLimitPrice)
	order.AutoCancelDate = msgDec.String(pb.AutoCancelDate)
	order.FilledQuantity = NewDecimalMaxFromProtobufDecoder(msgDec, pb.FilledQuantity)
	order.RefFuturesConID = msgDec.Int32(pb.RefFuturesConId)
	order.AutoCancelParent = msgDec.Bool(pb.AutoCancelParent)
	order.Shareholder = msgDec.String(pb.Shareholder)
	order.ImbalanceOnly = msgDec.Bool(pb.ImbalanceOnly)
	order.RouteMarketableToBbo = msgDec.Int32(pb.RouteMarketableToBbo)
	order.ParentPermID = msgDec.Int64(pb.ParentPermId)
	order.UsePriceMgmtAlgo = msgDec.Int32(pb.UsePriceMgmtAlgo) != 0
	order.Duration = msgDec.Int32Max(pb.Duration)
	order.PostToAts = msgDec.Int32Max(pb.PostToAts)
	order.AdvancedErrorOverride = msgDec.String(pb.AdvancedErrorOverride)
	order.ManualOrderTime = msgDec.String(pb.ManualOrderTime)
	order.MinTradeQty = msgDec.Int32Max(pb.MinTradeQty)
	order.MinCompeteSize = msgDec.Int32Max(pb.MinCompeteSize)
	order.CompeteAgainstBestOffset = msgDec.FloatMax(pb.CompeteAgainstBestOffset)
	order.MidOffsetAtWhole = msgDec.FloatMax(pb.MidOffsetAtWhole)
	order.MidOffsetAtHalf = msgDec.FloatMax(pb.MidOffsetAtHalf)
	order.CustomerAccount = msgDec.String(pb.CustomerAccount)
	order.ProfessionalCustomer = msgDec.Bool(pb.ProfessionalCustomer)
	order.BondAccruedInterest = msgDec.String(pb.BondAccruedInterest)
	order.IncludeOvernight = msgDec.Bool(pb.IncludeOvernight)
	order.ManualOrderIndicator = msgDec.Int32Max(pb.ManualOrderIndicator)
	order.Submitter = msgDec.String(pb.Submitter)
	order.Deactivate = msgDec.Bool(pb.Deactivate)
	order.PostOnly = msgDec.Bool(pb.PostOnly)
	order.AllowPreOpen = msgDec.Bool(pb.AllowPreOpen)
	order.IgnoreOpenAuction = msgDec.Bool(pb.IgnoreOpenAuction)
	order.SeekPriceImprovement = msgDec.Int32(pb.SeekPriceImprovement)
	order.WhatIfType = msgDec.Int32(pb.WhatIfType)
	order.HedgeMaxSize = msgDec.Int32Max(pb.HedgeMaxSize)

	return order
}

// Proto converts the order into its protobuf representation.
func (o *Order) Proto() *protobuf.Order {
	if o == nil {
		return nil
	}

	pb := &protobuf.Order{
		ClientId:                       protofmt.Int32(o.ClientID),
		OrderId:                        protofmt.Int32(o.OrderID),
		PermId:                         protofmt.Int64(o.PermID),
		ParentId:                       protofmt.Int32(o.ParentID),
		Action:                         protofmt.String(o.Action),
		TotalQuantity:                  decimalToProtoString(o.TotalQuantity),
		DisplaySize:                    protofmt.Int32(o.DisplaySize),
		OrderType:                      protofmt.String(o.OrderType),
		LmtPrice:                       protofmt.FloatMax(o.LmtPrice),
		AuxPrice:                       protofmt.FloatMax(o.AuxPrice),
		Tif:                            protofmt.String(string(o.TIF)),
		Account:                        protofmt.String(o.Account),
		SettlingFirm:                   protofmt.String(o.SettlingFirm),
		ClearingAccount:                protofmt.String(o.ClearingAccount),
		ClearingIntent:                 protofmt.String(o.ClearingIntent),
		AllOrNone:                      protofmt.Bool(o.AllOrNone),
		BlockOrder:                     protofmt.Bool(o.BlockOrder),
		Hidden:                         protofmt.Bool(o.Hidden),
		OutsideRth:                     protofmt.Bool(o.OutsideRTH),
		SweepToFill:                    protofmt.Bool(o.SweepToFill),
		PercentOffset:                  protofmt.FloatMax(o.PercentOffset),
		TrailingPercent:                protofmt.FloatMax(o.TrailingPercent),
		TrailStopPrice:                 protofmt.FloatMax(o.TrailStopPrice),
		MinQty:                         protofmt.Int32Max(o.MinQty),
		GoodAfterTime:                  protofmt.String(o.GoodAfterTime),
		GoodTillDate:                   protofmt.String(o.GoodTillDate),
		OcaGroup:                       protofmt.String(o.OCAGroup),
		OrderRef:                       protofmt.String(o.OrderRef),
		Rule80A:                        protofmt.String(string(o.Rule80A)),
		OcaType:                        protofmt.Int32(int32(o.OCAType)),
		TriggerMethod:                  protofmt.Int32(int32(o.TriggerMethod)),
		ActiveStartTime:                protofmt.String(o.ActiveStartTime),
		ActiveStopTime:                 protofmt.String(o.ActiveStopTime),
		FaGroup:                        protofmt.String(o.FAGroup),
		FaMethod:                       protofmt.String(o.FAMethod),
		FaPercentage:                   protofmt.String(o.FAPercentage),
		Volatility:                     protofmt.FloatMax(o.Volatility),
		VolatilityType:                 protofmt.Int32(int32(o.VolatilityType)),
		ContinuousUpdate:               protofmt.Bool(o.ContinuousUpdate),
		ReferencePriceType:             protofmt.Int32(o.ReferencePriceType),
		DeltaNeutralOrderType:          protofmt.String(o.DeltaNeutralOrderType),
		DeltaNeutralAuxPrice:           protofmt.FloatMax(o.DeltaNeutralAuxPrice),
		DeltaNeutralConId:              protofmt.Int32(o.DeltaNeutralConID),
		DeltaNeutralOpenClose:          protofmt.String(o.DeltaNeutralOpenClose),
		DeltaNeutralShortSale:          protofmt.Bool(o.DeltaNeutralShortSale),
		DeltaNeutralShortSaleSlot:      protofmt.Int32(o.DeltaNeutralShortSaleSlot),
		DeltaNeutralDesignatedLocation: protofmt.String(o.DeltaNeutralDesignatedLocation),
		ScaleInitLevelSize:             protofmt.Int32Max(o.ScaleInitLevelSize),
		ScaleSubsLevelSize:             protofmt.Int32Max(o.ScaleSubsLevelSize),
		ScalePriceIncrement:            protofmt.FloatMax(o.ScalePriceIncrement),
		ScalePriceAdjustValue:          protofmt.FloatMax(o.ScalePriceAdjustValue),
		ScalePriceAdjustInterval:       protofmt.Int32Max(o.ScalePriceAdjustInterval),
		ScaleProfitOffset:              protofmt.FloatMax(o.ScaleProfitOffset),
		ScaleAutoReset:                 protofmt.Bool(o.ScaleAutoReset),
		ScaleInitPosition:              protofmt.Int32Max(o.ScaleInitPosition),
		ScaleInitFillQty:               protofmt.Int32Max(o.ScaleInitFillQty),
		ScaleRandomPercent:             protofmt.Bool(o.ScaleRandomPercent),
		ScaleTable:                     protofmt.String(o.ScaleTable),
		HedgeType:                      protofmt.String(o.HedgeType),
		HedgeParam:                     protofmt.String(o.HedgeParam),
		AlgoStrategy:                   protofmt.String(o.AlgoStrategy),
		AlgoParams:                     tagValueListToProtoMap(o.AlgoParams),
		AlgoId:                         protofmt.String(o.AlgoID),
		SmartComboRoutingParams:        tagValueListToProtoMap(o.SmartComboRoutingParams),
		WhatIf:                         protofmt.Bool(o.WhatIf),
		Transmit:                       protofmt.Bool(o.Transmit),
		OverridePercentageConstraints:  protofmt.Bool(o.OverridePercentageConstraints),
		OpenClose:                      protofmt.String(o.OpenClose),
		Origin:                         protofmt.Int32(o.Origin),
		ShortSaleSlot:                  protofmt.Int32(o.ShortSaleSlot),
		DesignatedLocation:             protofmt.String(o.DesignatedLocation),
		ExemptCode:                     protofmt.Int32(o.ExemptCode),
		DeltaNeutralSettlingFirm:       protofmt.String(o.DeltaNeutralSettlingFirm),
		DeltaNeutralClearingAccount:    protofmt.String(o.DeltaNeutralClearingAccount),
		DeltaNeutralClearingIntent:     protofmt.String(o.DeltaNeutralClearingIntent),
		DiscretionaryAmt:               protofmt.Float(o.DiscretionaryAmt),
		OptOutSmartRouting:             protofmt.Bool(o.OptOutSmartRouting),
		StartingPrice:                  protofmt.FloatMax(o.StartingPrice),
		StockRefPrice:                  protofmt.FloatMax(o.StockRefPrice),
		Delta:                          protofmt.FloatMax(o.Delta),
		StockRangeLower:                protofmt.FloatMax(o.StockRangeLower),
		StockRangeUpper:                protofmt.FloatMax(o.StockRangeUpper),
		NotHeld:                        protofmt.Bool(o.NotHeld),
		OrderMiscOptions:               tagValueListToProtoMap(o.OrderMiscOptions),
		Solicited:                      protofmt.Bool(o.Solicited),
		RandomizeSize:                  protofmt.Bool(o.RandomizeSize),
		RandomizePrice:                 protofmt.Bool(o.RandomizePrice),
		ReferenceContractId:            protofmt.Int32(o.ReferenceContractID),
		PeggedChangeAmount:             protofmt.Float(o.PeggedChangeAmount),
		IsPeggedChangeAmountDecrease:   protofmt.Bool(o.IsPeggedChangeAmountDecrease),
		ReferenceChangeAmount:          protofmt.Float(o.ReferenceChangeAmount),
		ReferenceExchangeId:            protofmt.String(o.ReferenceExchangeID),
		AdjustedOrderType:              protofmt.String(o.AdjustedOrderType),
		TriggerPrice:                   protofmt.FloatMax(o.TriggerPrice),
		AdjustedStopPrice:              protofmt.FloatMax(o.AdjustedStopPrice),
		AdjustedStopLimitPrice:         protofmt.FloatMax(o.AdjustedStopLimitPrice),
		AdjustedTrailingAmount:         protofmt.FloatMax(o.AdjustedTrailingAmount),
		AdjustableTrailingUnit:         protofmt.Int32(o.AdjustableTrailingUnit),
		LmtPriceOffset:                 protofmt.FloatMax(o.LmtPriceOffset),
		Conditions:                     orderConditionsToProtobuf(o.Conditions),
		ConditionsCancelOrder:          protofmt.Bool(o.ConditionsCancelOrder),
		ConditionsIgnoreRth:            protofmt.Bool(o.ConditionsIgnoreRth),
		ModelCode:                      protofmt.String(o.ModelCode),
		ExtOperator:                    protofmt.String(o.ExtOperator),
		SoftDollarTier:                 o.SoftDollarTier.Proto(),
		CashQty:                        protofmt.FloatMax(o.CashQty),
		Mifid2DecisionMaker:            protofmt.String(o.Mifid2DecisionMaker),
		Mifid2DecisionAlgo:             protofmt.String(o.Mifid2DecisionAlgo),
		Mifid2ExecutionTrader:          protofmt.String(o.Mifid2ExecutionTrader),
		Mifid2ExecutionAlgo:            protofmt.String(o.Mifid2ExecutionAlgo),
		DontUseAutoPriceForHedge:       protofmt.Bool(o.DontUseAutoPriceForHedge),
		IsOmsContainer:                 protofmt.Bool(o.IsOmsContainer),
		DiscretionaryUpToLimitPrice:    protofmt.Bool(o.DiscretionaryUpToLimitPrice),
		AutoCancelDate:                 protofmt.String(o.AutoCancelDate),
		FilledQuantity:                 decimalToProtoString(o.FilledQuantity),
		RefFuturesConId:                protofmt.Int32(o.RefFuturesConID),
		AutoCancelParent:               protofmt.Bool(o.AutoCancelParent),
		Shareholder:                    protofmt.String(o.Shareholder),
		ImbalanceOnly:                  protofmt.Bool(o.ImbalanceOnly),
		RouteMarketableToBbo:           protofmt.Int32(o.RouteMarketableToBbo),
		ParentPermId:                   protofmt.Int64(o.ParentPermID),
		AdvancedErrorOverride:          protofmt.String(o.AdvancedErrorOverride),
		ManualOrderTime:                protofmt.String(o.ManualOrderTime),
		MinTradeQty:                    protofmt.Int32Max(o.MinTradeQty),
		MinCompeteSize:                 protofmt.Int32Max(o.MinCompeteSize),
		CompeteAgainstBestOffset:       protofmt.FloatMax(o.CompeteAgainstBestOffset),
		MidOffsetAtWhole:               protofmt.FloatMax(o.MidOffsetAtWhole),
		MidOffsetAtHalf:                protofmt.FloatMax(o.MidOffsetAtHalf),
		CustomerAccount:                protofmt.String(o.CustomerAccount),
		ProfessionalCustomer:           protofmt.Bool(o.ProfessionalCustomer),
		BondAccruedInterest:            protofmt.String(o.BondAccruedInterest),
		IncludeOvernight:               protofmt.Bool(o.IncludeOvernight),
		ManualOrderIndicator:           protofmt.Int32Max(o.ManualOrderIndicator),
		Submitter:                      protofmt.String(o.Submitter),
		Deactivate:                     protofmt.Bool(o.Deactivate),
		PostOnly:                       protofmt.Bool(o.PostOnly),
		AllowPreOpen:                   protofmt.Bool(o.AllowPreOpen),
		IgnoreOpenAuction:              protofmt.Bool(o.IgnoreOpenAuction),
		SeekPriceImprovement:           protofmt.Int32(o.SeekPriceImprovement),
		WhatIfType:                     protofmt.Int32(o.WhatIfType),
		HedgeMaxSize:                   protofmt.Int32Max(o.HedgeMaxSize),
	}
	if o.UsePriceMgmtAlgo {
		pb.UsePriceMgmtAlgo = protofmt.Int32(1)
	}
	if o.Duration != nil {
		pb.Duration = protofmt.Int32Max(o.Duration)
	}
	if o.PostToAts != nil {
		pb.PostToAts = protofmt.Int32Max(o.PostToAts)
	}
	return pb
}

// NewOrderStateFromProtobufDecoder decodes an order state from protobuf data.
func NewOrderStateFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.OrderState) *OrderState {
	orderState := NewOrderState()
	if pb == nil {
		return orderState
	}

	orderState.Status = msgDec.String(pb.Status)
	orderState.InitMarginBefore = floatToStateString(msgDec.FloatMax(pb.InitMarginBefore))
	orderState.MaintMarginBefore = floatToStateString(msgDec.FloatMax(pb.MaintMarginBefore))
	orderState.EquityWithLoanBefore = floatToStateString(msgDec.FloatMax(pb.EquityWithLoanBefore))
	orderState.InitMarginChange = floatToStateString(msgDec.FloatMax(pb.InitMarginChange))
	orderState.MaintMarginChange = floatToStateString(msgDec.FloatMax(pb.MaintMarginChange))
	orderState.EquityWithLoanChange = floatToStateString(msgDec.FloatMax(pb.EquityWithLoanChange))
	orderState.InitMarginAfter = floatToStateString(msgDec.FloatMax(pb.InitMarginAfter))
	orderState.MaintMarginAfter = floatToStateString(msgDec.FloatMax(pb.MaintMarginAfter))
	orderState.EquityWithLoanAfter = floatToStateString(msgDec.FloatMax(pb.EquityWithLoanAfter))
	orderState.CommissionAndFees = msgDec.FloatMax(pb.CommissionAndFees)
	orderState.MinCommissionAndFees = msgDec.FloatMax(pb.MinCommissionAndFees)
	orderState.MaxCommissionAndFees = msgDec.FloatMax(pb.MaxCommissionAndFees)
	orderState.CommissionAndFeesCurrency = msgDec.String(pb.CommissionAndFeesCurrency)
	orderState.MarginCurrency = msgDec.String(pb.MarginCurrency)
	orderState.InitMarginBeforeOutsideRTH = msgDec.FloatMax(pb.InitMarginBeforeOutsideRTH)
	orderState.MaintMarginBeforeOutsideRTH = msgDec.FloatMax(pb.MaintMarginBeforeOutsideRTH)
	orderState.EquityWithLoanBeforeOutsideRTH = msgDec.FloatMax(pb.EquityWithLoanBeforeOutsideRTH)
	orderState.InitMarginChangeOutsideRTH = msgDec.FloatMax(pb.InitMarginChangeOutsideRTH)
	orderState.MaintMarginChangeOutsideRTH = msgDec.FloatMax(pb.MaintMarginChangeOutsideRTH)
	orderState.EquityWithLoanChangeOutsideRTH = msgDec.FloatMax(pb.EquityWithLoanChangeOutsideRTH)
	orderState.InitMarginAfterOutsideRTH = msgDec.FloatMax(pb.InitMarginAfterOutsideRTH)
	orderState.MaintMarginAfterOutsideRTH = msgDec.FloatMax(pb.MaintMarginAfterOutsideRTH)
	orderState.EquityWithLoanAfterOutsideRTH = msgDec.FloatMax(pb.EquityWithLoanAfterOutsideRTH)
	orderState.SuggestedSize = NewDecimalMaxFromProtobufDecoder(msgDec, pb.SuggestedSize)
	orderState.RejectReason = msgDec.String(pb.RejectReason)
	orderState.WarningText = msgDec.String(pb.WarningText)
	orderState.CompletedTime = msgDec.String(pb.CompletedTime)
	orderState.CompletedStatus = msgDec.String(pb.CompletedStatus)
	if len(pb.OrderAllocations) > 0 {
		orderState.OrderAllocations = make([]*OrderAllocation, 0, len(pb.OrderAllocations))
		for _, oa := range pb.OrderAllocations {
			orderState.OrderAllocations = append(orderState.OrderAllocations, NewOrderAllocationFromProtobufDecoder(msgDec, oa))
		}
	}

	return orderState
}

// NewExecutionFromProtobufDecoder decodes an execution from protobuf data.
func NewExecutionFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.Execution) *Execution {
	execution := NewExecution()
	if pb == nil {
		return execution
	}

	execution.OrderID = msgDec.Int32(pb.OrderId)
	execution.ExecID = msgDec.String(pb.ExecId)
	execution.Time = msgDec.String(pb.Time)
	execution.AcctNumber = msgDec.String(pb.AcctNumber)
	execution.Exchange = msgDec.String(pb.Exchange)
	execution.Side = msgDec.String(pb.Side)
	execution.Shares = NewDecimalMaxFromProtobufDecoder(msgDec, pb.Shares)
	execution.Price = msgDec.Float(pb.Price)
	execution.PermID = msgDec.Int64(pb.PermId)
	execution.ClientID = msgDec.Int32(pb.ClientId)
	if msgDec.Bool(pb.IsLiquidation) {
		execution.Liquidation = 1
	}
	execution.CumQty = NewDecimalMaxFromProtobufDecoder(msgDec, pb.CumQty)
	execution.AvgPrice = msgDec.Float(pb.AvgPrice)
	execution.OrderRef = msgDec.String(pb.OrderRef)
	execution.EVRule = msgDec.String(pb.EvRule)
	execution.EVMultiplier = msgDec.Float(pb.EvMultiplier)
	execution.ModelCode = msgDec.String(pb.ModelCode)
	execution.LastLiquidity = Liquidities(msgDec.Int32(pb.LastLiquidity))
	execution.PendingPriceRevision = msgDec.Bool(pb.IsPriceRevisionPending)
	execution.Submitter = msgDec.String(pb.Submitter)
	execution.OptExerciseOrLapse = OptionExercise(msgDec.Int32(pb.OptExerciseOrLapseType))

	return execution
}

// Proto converts the soft dollar tier into its protobuf representation.
func (sdt SoftDollarTier) Proto() *protobuf.SoftDollarTier {
	return &protobuf.SoftDollarTier{
		Name:        protofmt.String(sdt.Name),
		Value:       protofmt.String(sdt.Value),
		DisplayName: protofmt.String(sdt.DisplayName),
	}
}

// NewSoftDollarTierFromProtobufDecoder decodes a soft dollar tier from protobuf data.
func NewSoftDollarTierFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.SoftDollarTier) SoftDollarTier {
	if pb == nil {
		return SoftDollarTier{}
	}

	return SoftDollarTier{
		Name:        msgDec.String(pb.Name),
		Value:       msgDec.String(pb.Value),
		DisplayName: msgDec.String(pb.DisplayName),
	}
}

// Proto converts the cancel request into its protobuf representation.
func (oc OrderCancel) Proto() *protobuf.OrderCancel {
	return &protobuf.OrderCancel{
		ManualOrderCancelTime: protofmt.String(oc.ManualOrderCancelTime),
		ExtOperator:           protofmt.String(oc.ExtOperator),
		ManualOrderIndicator:  protofmt.Int32Max(oc.ManualOrderIndicator),
	}
}

// Proto converts the execution filter into its protobuf representation.
func (ef *ExecutionFilter) Proto() *protobuf.ExecutionFilter {
	if ef == nil {
		return &protobuf.ExecutionFilter{}
	}

	pb := &protobuf.ExecutionFilter{
		ClientId:  protofmt.Int32(ef.ClientID),
		AcctCode:  protofmt.String(ef.AcctCode),
		Time:      protofmt.String(ef.Time),
		Symbol:    protofmt.String(ef.Symbol),
		SecType:   protofmt.String(string(ef.SecType)),
		Exchange:  protofmt.String(ef.Exchange),
		Side:      protofmt.String(ef.Side),
		LastNDays: protofmt.Int32Max(ef.LastNDays),
	}
	for _, date := range ef.SpecificDates {
		pb.SpecificDates = append(pb.SpecificDates, int32(date.Unix()))
	}
	return pb
}

// NewOrderAllocationFromProtobufDecoder decodes an order allocation from protobuf data.
func NewOrderAllocationFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.OrderAllocation) *OrderAllocation {
	if pb == nil {
		return NewOrderAllocation()
	}

	return &OrderAllocation{
		Account:         msgDec.String(pb.Account),
		Position:        NewDecimalMaxFromProtobufDecoder(msgDec, pb.Position),
		PositionDesired: NewDecimalMaxFromProtobufDecoder(msgDec, pb.PositionDesired),
		PositionAfter:   NewDecimalMaxFromProtobufDecoder(msgDec, pb.PositionAfter),
		DesiredAllocQty: NewDecimalMaxFromProtobufDecoder(msgDec, pb.DesiredAllocQty),
		AllowedAllocQty: NewDecimalMaxFromProtobufDecoder(msgDec, pb.AllowedAllocQty),
		IsMonetary:      msgDec.Bool(pb.IsMonetary),
	}
}

func orderConditionsFromProtobuf(msgDec *protofmt.Decoder, pbs []*protobuf.OrderCondition) []OrderCondition {
	if len(pbs) == 0 {
		return nil
	}

	conditions := make([]OrderCondition, 0, len(pbs))
	for _, pb := range pbs {
		condType := OrderConditionType(msgDec.Int32(pb.Type))
		cond, err := NewOrderCondition(condType)
		if err != nil {
			msgDec.SetErr(err)
			return nil
		}

		switch c := cond.(type) {
		case *OrderPriceCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.IsMore = msgDec.Bool(pb.IsMore)
			c.ConID = msgDec.Int32(pb.ConId)
			c.Exchange = msgDec.String(pb.Exchange)
			c.Price = msgDec.Float(pb.Price)
			c.TriggerMethod = TriggerMethod(msgDec.Int32(pb.TriggerMethod))
		case *OrderTimeCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.IsMore = msgDec.Bool(pb.IsMore)
			c.Time = msgDec.String(pb.Time)
		case *OrderMarginCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.IsMore = msgDec.Bool(pb.IsMore)
			c.Percent = msgDec.Int32(pb.Percent)
		case *OrderExecutionCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.SecType = SecurityType(msgDec.String(pb.SecType))
			c.Exchange = msgDec.String(pb.Exchange)
			c.Symbol = msgDec.String(pb.Symbol)
		case *OrderVolumeCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.IsMore = msgDec.Bool(pb.IsMore)
			c.ConID = msgDec.Int32(pb.ConId)
			c.Exchange = msgDec.String(pb.Exchange)
			c.Volume = msgDec.Int32(pb.Volume)
		case *OrderPercentChangeCondition:
			c.IsConjunctionConnection = msgDec.Bool(pb.IsConjunctionConnection)
			c.IsMore = msgDec.Bool(pb.IsMore)
			c.ConID = msgDec.Int32(pb.ConId)
			c.Exchange = msgDec.String(pb.Exchange)
			c.ChangePercent = msgDec.Float(pb.ChangePercent)
		}

		conditions = append(conditions, cond)
	}

	return conditions
}

func orderConditionsToProtobuf(conditions []OrderCondition) []*protobuf.OrderCondition {
	if len(conditions) == 0 {
		return nil
	}

	out := make([]*protobuf.OrderCondition, 0, len(conditions))
	for _, cond := range conditions {
		pb := &protobuf.OrderCondition{
			Type:                    protofmt.Int32(int32(cond.Type())),
			IsConjunctionConnection: protofmt.Bool(cond.makeFields()[0] == "a"),
		}

		switch c := cond.(type) {
		case *OrderPriceCondition:
			pb.IsMore = protofmt.Bool(c.IsMore)
			pb.ConId = protofmt.Int32(c.ConID)
			pb.Exchange = protofmt.String(c.Exchange)
			pb.Price = protofmt.Float(c.Price)
			pb.TriggerMethod = protofmt.Int32(int32(c.TriggerMethod))
		case *OrderTimeCondition:
			pb.IsMore = protofmt.Bool(c.IsMore)
			pb.Time = protofmt.String(c.Time)
		case *OrderMarginCondition:
			pb.IsMore = protofmt.Bool(c.IsMore)
			pb.Percent = protofmt.Int32(c.Percent)
		case *OrderExecutionCondition:
			pb.SecType = protofmt.String(string(c.SecType))
			pb.Exchange = protofmt.String(c.Exchange)
			pb.Symbol = protofmt.String(c.Symbol)
		case *OrderVolumeCondition:
			pb.IsMore = protofmt.Bool(c.IsMore)
			pb.ConId = protofmt.Int32(c.ConID)
			pb.Exchange = protofmt.String(c.Exchange)
			pb.Volume = protofmt.Int32(c.Volume)
		case *OrderPercentChangeCondition:
			pb.IsMore = protofmt.Bool(c.IsMore)
			pb.ConId = protofmt.Int32(c.ConID)
			pb.Exchange = protofmt.String(c.Exchange)
			pb.ChangePercent = protofmt.Float(c.ChangePercent)
		}

		out = append(out, pb)
	}

	return out
}

func tagValueMapFromProtobuf(m map[string]string) []TagValue {
	if len(m) == 0 {
		return nil
	}

	values := make([]TagValue, 0, len(m))
	for k, v := range m {
		values = append(values, TagValue{Tag: k, Value: v})
	}
	return values
}

func tagValueListToProtoMap(values []TagValue) map[string]string {
	if len(values) == 0 {
		return nil
	}

	m := make(map[string]string, len(values))
	for _, value := range values {
		m[value.Tag] = value.Value
	}
	return m
}

func decimalToProtoString(val *Decimal) *string {
	if val == nil {
		return nil
	}
	s := val.StringMax()
	return &s
}

func floatToStateString(val *float64) string {
	if val == nil {
		return ""
	}
	return strconv.FormatFloat(*val, 'f', -1, 64)
}

func parseCommissionRedemptionDate(value string) time.Time {
	if len(value) == 0 {
		return time.Time{}
	}
	t, err := time.Parse("20060102", value)
	if err != nil {
		return time.Time{}
	}
	return t
}
