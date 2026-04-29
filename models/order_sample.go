package models

// -----------------------------------------------------------------------------

// These sample builders mirror common IB order templates.
// Test exchange-specific combinations in paper trading before using them elsewhere.

// AtAuction creates an auction order that works the open and may rest as a limit order afterward.
// Products: FUT, STK.
func AtAuction(action string, quantity *Decimal, price *float64) *Order {
	order := NewOrder()
	order.Action = action
	order.TIF = TimeInForceAuction
	order.OrderType = OrderTypeMarketToLimit
	order.TotalQuantity = quantity
	order.LmtPrice = price

	return order
}

// Discretionary creates a limit order with a hidden discretionary range around the displayed limit price.
// Products: STK.
func Discretionary(action string, quantity *Decimal, price *float64, discretionaryAmount float64) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = price
	order.DiscretionaryAmt = discretionaryAmount

	return order
}

// MarketOrder creates a market order with no explicit price protection.
// Products: BOND, CFD, EFP, CASH, FUND, FUT, FOP, OPT, STK, WAR.
func MarketOrder(action string, quantity *Decimal) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarket
	order.TotalQuantity = quantity

	return order
}

// MarketIfTouched creates a market-if-touched order that turns into a market order once the trigger is reached.
// Products: BOND, CFD, CASH, FUT, FOP, OPT, STK, WAR.
func MarketIfTouched(action string, quantity *Decimal, price *float64) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarketIfTouched
	order.TotalQuantity = quantity
	order.AuxPrice = price

	return order
}

// MarketOnClose creates a market-on-close order.
// Products: CFD, FUT, STK, WAR.
func MarketOnClose(action string, quantity *Decimal) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarketOnClose
	order.TotalQuantity = quantity

	return order
}

// MarketOnOpen creates a market-on-open order using the on-open time in force.
// Products: CFD, STK, OPT, WAR.
func MarketOnOpen(action string, quantity *Decimal) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarket
	order.TotalQuantity = quantity
	order.TIF = TimeInForceOnOpen

	return order
}

// MidpointMatch creates an ISE midpoint match order.
// Products: STK.
func MidpointMatch(action string, quantity *Decimal) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarket
	order.TotalQuantity = quantity

	return order
}

// Midprice creates a smart-routed midpoint order with an optional price cap.
func Midprice(action string, quantity *Decimal, priceCap *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMidPrice
	order.TotalQuantity = quantity
	order.LmtPrice = priceCap // optional

	return order
}

// PeggedToMarket creates a pegged-to-market order using the given market offset.
// The active price tracks the inside market while still respecting the order's cap.
//
//	Sell order price = Bid price + offset amount
//	Buy order price = Ask price - offset amount
//
// Products: STK.
func PeggedToMarket(action string, quantity *Decimal, marketOffset *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToMarket
	order.TotalQuantity = quantity
	order.AuxPrice = marketOffset // Offset price

	return order
}

// PeggedToStock creates a pegged-to-stock order that reprices from the underlying move and the supplied delta.
// Products: OPT.
func PeggedToStock(action string, quantity *Decimal, delta *float64, stockReferencePrice *float64, startingPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToStock
	order.TotalQuantity = quantity
	order.Delta = delta
	order.StockRefPrice = stockReferencePrice
	order.StartingPrice = startingPrice

	return order
}

// RelativePeggedToPrimary creates a relative order pegged aggressively to the primary quote with an offset.
// Products: CFD, STK, OPT, FUT.
func RelativePeggedToPrimary(action string, quantity *Decimal, priceCap *float64, offsetAmount *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeRelative
	order.TotalQuantity = quantity
	order.LmtPrice = priceCap
	order.AuxPrice = offsetAmount

	return order
}

// SweepToFill creates a limit order marked to sweep available liquidity level by level.
// Products: CFD, STK, WAR.
func SweepToFill(action string, quantity *Decimal, price *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = price
	order.SweepToFill = true

	return order
}

// AuctionLimit creates a BOX auction limit order.
// Products: OPT. Supported Exchanges: BOX.
func AuctionLimit(action string, quantity *Decimal, price *float64, auctionStrategy AuctionStrategy) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = price
	order.AuctionStrategy = auctionStrategy

	return order
}

// AuctionPeggedToStock creates a BOX auction order whose auction price is pegged from the
// underlying stock move and the supplied delta.
// Products: OPT. Supported Exchanges: BOX.
func AuctionPeggedToStock(action string, quantity *Decimal, startingPrice *float64, delta *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToStock
	order.TotalQuantity = quantity
	order.Delta = delta
	order.StartingPrice = startingPrice

	return order
}

// AuctionRelative creates a BOX auction-relative order priced from the supplied offset.
// Products: OPT. Supported Exchanges: BOX.
func AuctionRelative(action string, quantity *Decimal, offset *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeRelative
	order.TotalQuantity = quantity
	order.AuxPrice = offset

	return order
}

// Block creates a large-size option order with the block flag set.
// Products: OPT.
func Block(action string, quantity *Decimal, price *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity // Large volumes!
	order.LmtPrice = price
	order.BlockOrder = true

	return order
}

// BoxTop creates a BOX top order that may rest as a limit order after a partial fill.
// Products: OPT. Supported Exchanges: BOX.
func BoxTop(action string, quantity *Decimal) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeBoxTop
	order.TotalQuantity = quantity

	return order
}

// LimitOrder creates a standard limit order.
// Products: BOND, CFD, CASH, FUT, FOP, OPT, STK, WAR.
func LimitOrder(action string, quantity *Decimal, limitPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice

	return order
}

// LimitOrderWithCashQty .
// Forex orders can be placed in denomination of second currency in pair using cashQty field.
// Requires TWS or IBG 963+.
// https://www.interactivebrokers.com/en/index.php?f=23876#963-02
func LimitOrderWithCashQty(action string, limitPrice *float64, cashQty *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.LmtPrice = limitPrice
	order.CashQty = cashQty

	return order
}

// LimitIfTouched creates a limit-if-touched order with separate trigger and limit prices.
// Products: BOND, CFD, CASH, FUT, FOP, OPT, STK, WAR.
func LimitIfTouched(action string, quantity *Decimal, limitPrice *float64, triggerPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimitIfTouched
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice
	order.AuxPrice = triggerPrice

	return order
}

// LimitOnClose creates a limit-on-close order.
// Products: CFD, FUT, STK, WAR.
func LimitOnClose(action string, quantity *Decimal, limitPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimitOnClose
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice

	return order
}

// LimitOnOpen creates a limit order that is only eligible at the market open.
// Products: CFD, STK, OPT, WAR.
func LimitOnOpen(action string, quantity *Decimal, limitPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.TIF = TimeInForceOnOpen
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice

	return order
}

// PassiveRelative creates a less aggressive relative order that pegs away from the inside quote by the offset.
// Products: STK, WAR.
func PassiveRelative(action string, quantity *Decimal, offset *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePassiveRelative
	order.TotalQuantity = quantity
	order.AuxPrice = offset

	return order
}

// PeggedToMidpoint creates a midpoint-pegged order with an offset and optional cap.
// Products: STK.
func PeggedToMidpoint(action string, quantity *Decimal, offset *float64, limitPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToMidpoint
	order.TotalQuantity = quantity
	order.AuxPrice = offset
	order.LmtPrice = limitPrice

	return order
}

// BracketOrder builds a parent order with attached take-profit and stop-loss children.
// Products: CFD, BAG, FOP, CASH, FUT, OPT, STK, WAR.
func BracketOrder(parentOrderId int32, action string, quantity *Decimal, limitPrice *float64, takeProfitLimitPrice *float64, stopLossPrice *float64) (parent, takeProfit, stopLoss *Order) {

	// This will be our main or "parent" order
	parent = NewOrder()
	parent.OrderID = parentOrderId
	parent.Action = action
	parent.OrderType = OrderTypeLimit
	parent.TotalQuantity = quantity
	parent.LmtPrice = limitPrice
	// The parent and children orders will need this attribute set to False to prevent accidental executions.
	// The LAST CHILD will have it set to True,
	parent.Transmit = false

	takeProfit = NewOrder()
	takeProfit.OrderID = parent.OrderID + 1
	if action == "BUY" {
		takeProfit.Action = OrderActionSell
	} else {
		takeProfit.Action = OrderActionBuy
	}
	takeProfit.OrderType = OrderTypeLimit
	takeProfit.TotalQuantity = quantity
	takeProfit.LmtPrice = takeProfitLimitPrice
	takeProfit.ParentID = parentOrderId
	takeProfit.Transmit = false

	stopLoss = NewOrder()
	stopLoss.OrderID = parent.OrderID + 2
	if action == "BUY" {
		stopLoss.Action = OrderActionSell
	} else {
		stopLoss.Action = OrderActionBuy
	}
	stopLoss.OrderType = OrderTypeStop
	// Stop trigger price
	stopLoss.AuxPrice = stopLossPrice
	stopLoss.TotalQuantity = quantity
	stopLoss.ParentID = parentOrderId
	// In this case, the low side order will be the last child being sent. Therefore, it needs to set this attribute to True
	// to activate all its predecessors
	stopLoss.Transmit = true
	return
}

// MarketToLimit creates a market-to-limit order that may rest as a limit order after a partial fill.
// Products: CFD, FUT, FOP, OPT, STK, WAR.
func MarketToLimit(action string, quantity *Decimal) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarketToLimit
	order.TotalQuantity = quantity

	return order
}

// MarketWithProtection creates a market-with-protection order.
// Products: FUT, FOP.
func MarketWithProtection(action string, quantity *Decimal) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarketWithProtection
	order.TotalQuantity = quantity

	return order
}

// Stop creates a stop order that becomes a market order when the stop price is triggered.
// Products: CFD, BAG, CASH, FUT, FOP, OPT, STK, WAR.
func Stop(action string, quantity *Decimal, stopPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeStop
	order.AuxPrice = stopPrice
	order.TotalQuantity = quantity

	return order
}

// StopLimit creates a stop-limit order with explicit stop and limit prices.
// Products: CFD, CASH, FUT, FOP, OPT, STK, WAR.
func StopLimit(action string, quantity *Decimal, limitPrice *float64, stopPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeStopLimit
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice
	order.AuxPrice = stopPrice

	return order
}

// StopWithProtection creates a stop order that converts into a protected market order.
// Products: FUT.
func StopWithProtection(action string, quantity *Decimal, stopPrice *float64) *Order {

	order := NewOrder()
	order.TotalQuantity = quantity
	order.Action = action
	order.OrderType = OrderTypeStopWithProtection
	order.AuxPrice = stopPrice

	return order
}

// TrailingStop creates a trailing stop order.
// Products: CFD, CASH, FOP, FUT, OPT, STK, WAR.
func TrailingStop(action string, quantity *Decimal, trailingPercent *float64, trailStopPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeTrailingStop
	order.TotalQuantity = quantity
	order.TrailingPercent = trailingPercent
	order.TrailStopPrice = trailStopPrice

	return order
}

// TrailingStopLimit creates a trailing stop limit order.
// Products: BOND, CFD, CASH, FUT, FOP, OPT, STK, WAR.
func TrailingStopLimit(action string, quantity *Decimal, lmtPriceOffset *float64, trailingAmount *float64, trailStopPrice *float64) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeTrailingStopLimit
	order.TotalQuantity = quantity
	order.TrailStopPrice = trailStopPrice
	order.LmtPriceOffset = lmtPriceOffset
	order.AuxPrice = trailingAmount

	return order
}

// ComboLimitOrder creates a combination limit order and optionally marks it as non-guaranteed.
// Products: OPT, STK, FUT.
func ComboLimitOrder(action string, quantity *Decimal, limitPrice *float64, nonGuaranteed bool) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = quantity
	order.LmtPrice = limitPrice
	if nonGuaranteed {
		order.SmartComboRoutingParams = []TagValue{}
		order.SmartComboRoutingParams = append(order.SmartComboRoutingParams, TagValue{Tag: "NonGuaranteed", Value: "1"})
	}

	return order
}

// ComboMarketOrder creates a combination market order and optionally marks it as non-guaranteed.
// Products: OPT, STK, FUT.
func ComboMarketOrder(action string, quantity *Decimal, nonGuaranteed bool) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeMarket
	order.TotalQuantity = quantity
	if nonGuaranteed {
		order.SmartComboRoutingParams = []TagValue{}
		order.SmartComboRoutingParams = append(order.SmartComboRoutingParams, TagValue{Tag: "NonGuaranteed", Value: "1"})
	}

	return order
}

// LimitOrderForComboWithLegPrices creates a combo limit order with explicit leg prices.
// Products: OPT, STK, FUT.
func LimitOrderForComboWithLegPrices(action string, quantity Decimal, legPrices []float64, nonGuaranteed bool) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.TotalQuantity = &quantity
	order.OrderComboLegs = []OrderComboLeg{}
	for _, price := range legPrices {
		comboLeg := NewOrderComboLeg()
		comboLeg.Price = &price
		order.OrderComboLegs = append(order.OrderComboLegs, comboLeg)
	}
	if nonGuaranteed {
		order.SmartComboRoutingParams = []TagValue{}
		order.SmartComboRoutingParams = append(order.SmartComboRoutingParams, TagValue{Tag: "NonGuaranteed", Value: "1"})
	}

	return order
}

// RelativeLimitCombo creates a relative-plus-limit combo order.
// Products: OPT, STK, FUT.
func RelativeLimitCombo(action string, quantity Decimal, limitPrice float64, nonGuaranteed bool) *Order {

	order := NewOrder()
	order.Action = action
	order.TotalQuantity = &quantity
	order.OrderType = OrderTypeRelativePlusLimit
	order.LmtPrice = &limitPrice
	if nonGuaranteed {
		order.SmartComboRoutingParams = []TagValue{}
		order.SmartComboRoutingParams = append(order.SmartComboRoutingParams, TagValue{Tag: "NonGuaranteed", Value: "1"})
	}

	return order
}

// NewRelativeMarketComboOrder creates a relative-plus-market combo order.
// Products: OPT, STK, FUT.
func NewRelativeMarketComboOrder(action string, quantity Decimal, nonGuaranteed bool) *Order {
	order := NewOrder()
	order.Action = action
	order.TotalQuantity = &quantity
	order.OrderType = OrderTypeRelativePlusMarket
	if nonGuaranteed {
		order.SmartComboRoutingParams = []TagValue{}
		order.SmartComboRoutingParams = append(order.SmartComboRoutingParams, TagValue{Tag: "NonGuaranteed", Value: "1"})
	}
	return order
}

// OneCancelsAll assigns OCA fields so filling one order can cancel the others in the group.
// Products: BOND, CASH, FUT, FOP, STK, OPT, WAR.
func OneCancelsAll(ocaGroup string, ocaOrder *Order, ocaType Oca) {
	ocaOrder.OCAGroup = ocaGroup
	ocaOrder.OCAType = ocaType
}

// NewVolatilityOrder creates a volatility-priced order for options or futures options.
// Products: FOP, OPT.
func NewVolatilityOrder(action string, quantity Decimal, volatilityPercent float64, volatilityType Volatility) *Order {

	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeVolatility
	order.TotalQuantity = &quantity
	order.Volatility = &volatilityPercent // Expressed in percentage (40%)
	order.VolatilityType = volatilityType

	return order
}

// NewMarketFHedgeOrder ...
func NewMarketFHedgeOrder(parentOrderId int32, action string) *Order {
	// FX Hedge orders can only have a quantity of 0
	order := MarketOrder(action, &DecimalZero)
	order.ParentID = parentOrderId
	order.HedgeType = "F"
	return order
}

// NewPeggedToBenchmarkOrder
func NewPeggedToBenchmarkOrder(
	action string, quantity Decimal, startingPrice float64, peggedChangeAmountDecrease bool,
	peggedChangeAmount float64, referenceChangeAmount float64, referenceConId int32, referenceExchange string,
	stockReferencePrice float64, referenceContractLowerRange float64, referenceContractUpperRange float64,
) *Order {
	order := NewOrder()
	order.OrderType = OrderTypePeggedToBenchmark
	// BUY or SELL
	order.Action = action
	order.TotalQuantity = &quantity
	// Beginning with price...
	order.StartingPrice = &startingPrice
	// increase/decrease price..
	order.IsPeggedChangeAmountDecrease = peggedChangeAmountDecrease
	// by... (and likewise for price moving in opposite direction)
	order.PeggedChangeAmount = peggedChangeAmount
	// whenever there is a price change of...
	order.ReferenceChangeAmount = referenceChangeAmount
	// in the reference contract...
	order.ReferenceContractID = referenceConId
	// being traded at...
	order.ReferenceExchangeID = referenceExchange
	// starting reference price is...
	order.StockRefPrice = &stockReferencePrice
	// Keep order active as long as reference contract trades between...
	order.StockRangeLower = &referenceContractLowerRange
	// and...
	order.StockRangeUpper = &referenceContractUpperRange

	return order
}

// NewAttachAdjustableToStopOrder
func NewAttachAdjustableToStopOrder(
	parent *Order, attachedOrderStopPrice float64, triggerPrice float64, adjustStopPrice float64,
) *Order {
	// Attached order is a conventional STP order in opposite direction
	action := OrderActionSell
	if parent.Action == OrderActionSell {
		action = OrderActionBuy
	}
	order := Stop(action, parent.TotalQuantity, &attachedOrderStopPrice)
	order.ParentID = parent.OrderID
	// When trigger price is penetrated
	order.TriggerPrice = &triggerPrice
	// The parent order will be turned into a STP order
	order.AdjustedOrderType = OrderTypeStop
	// With the given STP price
	order.AdjustedStopPrice = &adjustStopPrice

	return order
}

// NewAttachAdjustableToStopLimitOrder
func NewAttachAdjustableToStopLimitOrder(
	parent *Order, attachedOrderStopPrice float64, triggerPrice float64, adjustedStopPrice float64,
	adjustedStopLimitPrice float64,
) *Order {
	// Attached order is a conventional STP order
	action := OrderActionSell
	if parent.Action == OrderActionSell {
		action = OrderActionBuy
	}
	order := Stop(action, parent.TotalQuantity, &attachedOrderStopPrice)
	order.ParentID = parent.OrderID
	// When trigger price is penetrated
	order.TriggerPrice = &triggerPrice
	// The parent order will be turned into a STP LMT order
	order.AdjustedOrderType = OrderTypeStopLimit
	// With the given stop price
	order.AdjustedStopPrice = &adjustedStopPrice
	// And the given limit price
	order.AdjustedStopLimitPrice = &adjustedStopLimitPrice

	return order
}

// NewAttachAdjustableToTrailOrder
func NewAttachAdjustableToTrailOrder(
	parent *Order, attachedOrderStopPrice float64, triggerPrice float64, adjustedStopPrice float64,
	adjustedTrailAmount float64, trailUnit int32,
) *Order {

	// Attached order is a conventional STP order
	action := OrderActionSell
	if parent.Action == OrderActionSell {
		action = OrderActionBuy
	}
	order := Stop(action, parent.TotalQuantity, &attachedOrderStopPrice)
	order.ParentID = parent.OrderID
	// When trigger price is penetrated
	order.TriggerPrice = &triggerPrice
	// The parent order will be turned into a TRAIL order
	order.AdjustedOrderType = OrderTypeTrailingStop
	// With a stop price of...
	order.AdjustedStopPrice = &adjustedStopPrice
	// traling by and amount (0) or a percent (100)...
	order.AdjustableTrailingUnit = trailUnit
	// of...
	order.AdjustedTrailingAmount = &adjustedTrailAmount

	return order
}

// WhatIfLimitOrder ...
func WhatIfLimitOrder(action string, quantity Decimal, limitPrice float64) *Order {

	order := LimitOrder(action, &quantity, &limitPrice)
	order.WhatIf = true

	return order
}

// NewPriceCondition ...
func NewPriceCondition(
	triggerMethod TriggerMethod, conId int32, exchange string, price float64, isMore bool, isConjunction bool,
) *OrderPriceCondition {
	cond, _ := NewOrderCondition(OrderConditionTypePrice)
	priceCondition := cond.(*OrderPriceCondition)
	// When this contract...
	priceCondition.ConID = conId
	// traded on this exchange
	priceCondition.Exchange = exchange
	// has a price above/below
	priceCondition.IsMore = isMore
	priceCondition.TriggerMethod = triggerMethod
	// this quantity
	priceCondition.Price = price
	// AND | OR next condition (will be ignored if no more conditions are added)
	priceCondition.IsConjunctionConnection = isConjunction

	return priceCondition

}

// NewExecutionCondition ...
func NewExecutionCondition(
	symbol string, secType SecurityType, exchange string, isConjunction bool,
) *OrderExecutionCondition {
	cond, _ := NewOrderCondition(OrderConditionTypeExecution)
	execCondition := cond.(*OrderExecutionCondition)
	// When an execution on symbol
	execCondition.Symbol = symbol
	// at exchange
	execCondition.Exchange = exchange
	// for this secType
	execCondition.SecType = secType
	// AND | OR next condition (will be ignored if no more conditions are added)
	execCondition.IsConjunctionConnection = isConjunction

	return execCondition
}

// NewMarginCondition ...
func NewMarginCondition(percent int32, isMore bool, isConjunction bool) *OrderMarginCondition {
	cond, _ := NewOrderCondition(OrderConditionTypeMargin)
	marginCondition := cond.(*OrderMarginCondition)
	// If margin is above/below
	marginCondition.IsMore = isMore
	// given percent
	marginCondition.Percent = percent
	// AND | OR next condition (will be ignored if no more conditions are added)
	marginCondition.IsConjunctionConnection = isConjunction

	return marginCondition
}

// NewPercentageChangeCondition ...
func NewPercentageChangeCondition(
	pctChange float64, conId int32, exchange string, isMore bool, isConjunction bool,
) *OrderPercentChangeCondition {
	cond, _ := NewOrderCondition(OrderConditionTypePercentChange)
	pctChangeCondition := cond.(*OrderPercentChangeCondition)
	// If there is a price percent change measured against last close price above or below...
	pctChangeCondition.IsMore = isMore
	// this amount...
	pctChangeCondition.ChangePercent = pctChange
	// on this contract
	pctChangeCondition.ConID = conId
	// when traded on this exchange...
	pctChangeCondition.Exchange = exchange
	// AND | OR next condition (will be ignored if no more conditions are added)
	pctChangeCondition.IsConjunctionConnection = isConjunction

	return pctChangeCondition
}

// NewTimeCondition ...
func NewTimeCondition(time string, isMore bool, isConjunction bool) *OrderTimeCondition {
	cond, _ := NewOrderCondition(OrderConditionTypeTime)
	timeCondition := cond.(*OrderTimeCondition)
	// Before or after...
	timeCondition.IsMore = isMore
	// this time..
	timeCondition.Time = time
	// AND | OR next condition (will be ignored if no more conditions are added)
	timeCondition.IsConjunctionConnection = isConjunction

	return timeCondition
}

// NewVolumeCondition ...
func NewVolumeCondition(conId int32, exchange string, isMore bool, volume int32, isConjunction bool) *OrderVolumeCondition {
	cond, _ := NewOrderCondition(OrderConditionTypeVolume)
	volCond := cond.(*OrderVolumeCondition)
	// Whenever contract...
	volCond.ConID = conId
	// When traded at
	volCond.Exchange = exchange
	// reaches a volume higher/lower
	volCond.IsMore = isMore
	// than this...
	volCond.Volume = volume
	// AND | OR next condition (will be ignored if no more conditions are added)
	volCond.IsConjunctionConnection = isConjunction

	return volCond
}

// LimitIBKRATS .
func LimitIBKRATS(action string, quantity Decimal, limitPrice float64) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypeLimit
	order.LmtPrice = &limitPrice
	order.TotalQuantity = &quantity
	order.NotHeld = true

	return order
}

// LimitOrderWithManualOrderTime .
func LimitOrderWithManualOrderTime(action string, quantity Decimal, limitPrice float64, manualOrderTime string) *Order {
	order := LimitOrder(action, &quantity, &limitPrice)
	order.ManualOrderTime = manualOrderTime

	return order
}

// NewPegBestUpToMidOrder ...
func NewPegBestUpToMidOrder(
	action string, quantity Decimal, limitPrice float64, minTradeQty int32, minCompeteSize int32,
	midOffsetAtWhole float64, midOffsetAtHalf float64,
) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToBest
	order.LmtPrice = &limitPrice
	order.TotalQuantity = &quantity
	order.NotHeld = true
	order.MinTradeQty = &minTradeQty
	order.MinCompeteSize = &minCompeteSize
	order.CompeteAgainstBestOffset = &CompeteAgainstBestOffsetUpToMid
	order.MidOffsetAtWhole = &midOffsetAtWhole
	order.MidOffsetAtHalf = &midOffsetAtHalf

	return order
}

// NewPegBestOrder ...
func NewPegBestOrder(
	action string, quantity Decimal, limitPrice float64, minTradeQty int32, minCompeteSize int32,
	competeAgainstBestOffset float64,
) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToBest
	order.LmtPrice = &limitPrice
	order.TotalQuantity = &quantity
	order.NotHeld = true
	order.MinTradeQty = &minTradeQty
	order.MinCompeteSize = &minCompeteSize
	order.CompeteAgainstBestOffset = &competeAgainstBestOffset

	return order
}

// NewPegMidOrder ...
func NewPegMidOrder(
	action string, quantity Decimal, limitPrice float64, minTradeQty int32, midOffsetAtWhole float64,
	midOffsetAtHalf float64,
) *Order {
	order := NewOrder()
	order.Action = action
	order.OrderType = OrderTypePeggedToMidpoint
	order.LmtPrice = &limitPrice
	order.TotalQuantity = &quantity
	order.NotHeld = true
	order.MinTradeQty = &minTradeQty
	order.MidOffsetAtWhole = &midOffsetAtWhole
	order.MidOffsetAtHalf = &midOffsetAtHalf

	return order
}

// NewLimitOrderWithCustomerAccountOrder ...
func NewLimitOrderWithCustomerAccountOrder(
	action string, quantity Decimal, limitPrice float64, customerAccount string,
) *Order {
	order := LimitOrder(action, &quantity, &limitPrice)
	order.CustomerAccount = customerAccount
	return order
}

func NewLimitOrderWithIncludeOvernight(action string, quantity Decimal, limitPrice float64) *Order {
	order := LimitOrder(action, &quantity, &limitPrice)
	order.IncludeOvernight = true
	return order
}

// NewEmptyCancelOrder ...
func NewEmptyCancelOrder() OrderCancel {
	orderCancel := NewOrderCancel()
	return orderCancel
}

// NewCancelOrderWithManualTime ...
func NewCancelOrderWithManualTime(manualOrderCancelTime string) OrderCancel {
	orderCancel := NewOrderCancel()
	orderCancel.ManualOrderCancelTime = manualOrderCancelTime
	return orderCancel
}

// NewLimitOrderWithCmeTaggingFieldsOrder ...
func NewLimitOrderWithCmeTaggingFieldsOrder(
	action string, quantity Decimal, limitPrice float64, extOperator string, manualOrderIndicator int32,
) *Order {
	order := LimitOrder(action, &quantity, &limitPrice)
	order.ExtOperator = extOperator
	order.ManualOrderIndicator = &manualOrderIndicator
	return order
}

// NewOrderCancelWithCmeTaggingFields ...
func NewOrderCancelWithCmeTaggingFields(extOperator string, manualOrderIndicator int32) OrderCancel {
	orderCancel := NewOrderCancel()
	orderCancel.ExtOperator = extOperator
	orderCancel.ManualOrderIndicator = &manualOrderIndicator
	return orderCancel
}

// NewLimitOnCloseOrderWithImbalanceOnlyOrder ...
func NewLimitOnCloseOrderWithImbalanceOnlyOrder(action string, quantity Decimal, limitPrice float64) *Order {
	order := LimitOnClose(action, &quantity, &limitPrice)
	order.ImbalanceOnly = true
	return order
}
