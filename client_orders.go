package ibkr

import (
	"context"
	"net"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// PlaceOrder submits a new order or updates an existing order when opts.Order.OrderID is set.
func (c *Client) PlaceOrder(opts models.PlaceOrderRequestOptions) (int32, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return 0, errInvalidContract
	}
	if opts.Order == nil {
		return 0, errInvalidOrder
	}
	if len(opts.Order.Action) == 0 {
		return 0, errInvalidOrderAction
	}
	if len(opts.Order.OrderType) == 0 {
		return 0, errInvalidOrderType
	}
	if opts.Order.TotalQuantity == nil {
		return 0, errInvalidTotalQuantity
	}
	if opts.Order.ConditionsIncludeOvernight && c.serverVersion < common.ServerVersionConditionsIncludeOvernight {
		return 0, errConditionsIncludeOvernightUnsupported
	}

	if !c.rp.Acquire() {
		return 0, net.ErrClosed
	}
	defer c.rp.Release()

	orderID := opts.Order.OrderID
	if orderID <= 0 {
		orderID = c.getNextRequestID(RequestTypeRequestWithID)
		if orderID <= 0 {
			return 0, net.ErrClosed
		}
		opts.Order.OrderID = orderID
	}
	if opts.Order.ClientID == 0 {
		opts.Order.ClientID = c.clientID
	}

	// Persist the outbound order intent before sending so reconnect recovery can
	// reconcile later status, open-order, and execution callbacks.
	c.orderStore.CreateOrder(opts.Contract, opts.Order)

	if c.isProtoBufAvailable(common.PLACE_ORDER) {
		pb := protobuf.PlaceOrderRequest{
			OrderId:  protofmt.Int32(orderID),
			Contract: opts.Contract.Proto(opts.Order),
			Order:    opts.Order.Proto(),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.PLACE_ORDER + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = encodeLegacyPlaceOrderRequest(orderID, opts.Contract, opts.Order)
	}
	if msgEnc.Err() != nil {
		return 0, msgEnc.Err()
	}

	if err := c.sendMessage(msgEnc.Bytes()); err != nil {
		return 0, err
	}

	return orderID, nil
}

// -----------------------------------------------------------------------------

// CancelOrder requests cancellation of an active order.
func (c *Client) CancelOrder(opts models.CancelOrderRequestOptions) error {
	var msgEnc *message.Encoder

	if opts.OrderID <= 0 {
		return errInvalidOrderID
	}

	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_ORDER) {
		pb := protobuf.CancelOrderRequest{
			OrderId: protofmt.Int32(opts.OrderID),
		}
		if opts.OrderCancel != nil {
			pb.OrderCancel = opts.OrderCancel.Proto()
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_ORDER + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(6).
			RawUInt32(common.CANCEL_ORDER).
			Int(1).
			RequestID(opts.OrderID)
		if opts.OrderCancel != nil {
			msgEnc.String(opts.OrderCancel.ManualOrderCancelTime)
			msgEnc.String(opts.OrderCancel.ExtOperator)
			msgEnc.Int32Max(opts.OrderCancel.ManualOrderIndicator)
		}
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

// -----------------------------------------------------------------------------

// RequestOpenOrders retrieves the current open orders owned by this client ID.
func (c *Client) RequestOpenOrders(ctx context.Context) (*models.OpenOrdersResponse, error) {
	return c.requestOpenOrdersCommon(ctx, common.REQ_OPEN_ORDERS, false)
}

// -----------------------------------------------------------------------------

// RequestAllOpenOrders retrieves all open orders visible to the current connection.
func (c *Client) RequestAllOpenOrders(ctx context.Context) (*models.OpenOrdersResponse, error) {
	return c.requestOpenOrdersCommon(ctx, common.REQ_ALL_OPEN_ORDERS, false)
}

// -----------------------------------------------------------------------------

// RequestAutoOpenOrders configures whether new TWS orders are auto-bound to this client.
func (c *Client) RequestAutoOpenOrders(opts models.AutoOpenOrdersRequestOptions) error {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.REQ_AUTO_OPEN_ORDERS) {
		pb := protobuf.AutoOpenOrdersRequest{
			AutoBind: protofmt.Bool(opts.AutoBind),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_AUTO_OPEN_ORDERS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_AUTO_OPEN_ORDERS).
			Int(1).
			Bool(opts.AutoBind)
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

// -----------------------------------------------------------------------------

// RequestExecutions retrieves executions matching the provided filter.
func (c *Client) RequestExecutions(
	ctx context.Context, opts models.ExecutionsRequestOptions,
) (*models.ExecutionsResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.ExecutionsResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_EXECUTIONS,
		Response: resp,
	})

	filter := opts.Filter
	if filter == nil {
		filter = models.NewExecutionFilter()
	}

	if c.isProtoBufAvailable(common.REQ_EXECUTIONS) {
		pb := protobuf.ExecutionRequest{
			ReqId:           protofmt.Int32(req.ID()),
			ExecutionFilter: filter.Proto(),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_EXECUTIONS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(16).
			RawUInt32(common.REQ_EXECUTIONS).
			Int(3).
			RequestID(req.ID()).
			Int32(filter.ClientID).
			String(filter.AcctCode).
			String(filter.Time).
			String(filter.Symbol).
			String(string(filter.SecType)).
			String(filter.Exchange).
			String(filter.Side).
			Int32Max(filter.LastNDays).
			Int(len(filter.SpecificDates))
		for _, date := range filter.SpecificDates {
			msgEnc.Int64(date.Unix())
		}
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------

// RequestCompletedOrders retrieves completed orders visible to the current connection.
func (c *Client) RequestCompletedOrders(
	ctx context.Context, opts models.CompletedOrdersRequestOptions,
) (*models.CompletedOrdersResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.CompletedOrdersResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.COMPLETED_ORDERS_END),
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_COMPLETED_ORDERS) {
		pb := protobuf.CompletedOrdersRequest{
			ApiOnly: protofmt.Bool(opts.APIOnly),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_COMPLETED_ORDERS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(2).
			RawUInt32(common.REQ_COMPLETED_ORDERS).
			Bool(opts.APIOnly)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------

// RequestGlobalCancel requests cancellation of all active visible orders.
func (c *Client) RequestGlobalCancel(orderCancel *models.OrderCancel) error {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.REQ_GLOBAL_CANCEL) {
		pb := protobuf.GlobalCancelRequest{}
		if orderCancel != nil {
			pb.OrderCancel = orderCancel.Proto()
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_GLOBAL_CANCEL + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(5).
			RawUInt32(common.REQ_GLOBAL_CANCEL)
		if orderCancel != nil {
			msgEnc.String(orderCancel.ManualOrderCancelTime)
			msgEnc.String(orderCancel.ExtOperator)
			msgEnc.Int32Max(orderCancel.ManualOrderIndicator)
		}
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

// -----------------------------------------------------------------------------

func (c *Client) requestOpenOrdersCommon(
	ctx context.Context, msgCode uint32, autoBind bool,
) (*models.OpenOrdersResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.OpenOrdersResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.OPEN_ORDER_END),
		Response: resp,
	})

	if c.isProtoBufAvailable(msgCode) {
		switch msgCode {
		case common.REQ_OPEN_ORDERS:
			pb := protobuf.OpenOrdersRequest{}
			msgEnc = message.NewEncoder().
				RawUInt32(common.REQ_OPEN_ORDERS + common.PROTOBUF_MSG_ID).
				Proto(&pb)
		case common.REQ_ALL_OPEN_ORDERS:
			pb := protobuf.AllOpenOrdersRequest{}
			msgEnc = message.NewEncoder().
				RawUInt32(common.REQ_ALL_OPEN_ORDERS + common.PROTOBUF_MSG_ID).
				Proto(&pb)
		}
	} else {
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(msgCode).
			Int(1)
		if msgCode == common.REQ_AUTO_OPEN_ORDERS {
			msgEnc.Bool(autoBind)
		}
	}
	if msgEnc == nil {
		return nil, errInvalidOpenOrdersRequest
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------

func encodeLegacyPlaceOrderRequest(orderID int32, contract *models.Contract, order *models.Order) *message.Encoder {
	msgEnc := message.NewEncoder().Reserve(170).
		RawUInt32(common.PLACE_ORDER).
		RequestID(orderID).
		Int32(contract.ConID).
		String(contract.Symbol).
		String(string(contract.SecType)).
		String(contract.LastTradeDateOrContractMonth).
		FloatMax(contract.Strike).
		String(contract.Right).
		FloatMax(contract.Multiplier).
		String(contract.Exchange).
		String(contract.PrimaryExchange).
		String(contract.Currency).
		String(contract.LocalSymbol).
		String(contract.TradingClass).
		String(contract.SecIDType).
		String(contract.SecID).
		String(order.Action).
		Marshal(order.TotalQuantity, 1).
		String(order.OrderType).
		FloatMax(order.LmtPrice).
		FloatMax(order.AuxPrice).
		String(string(order.TIF)).
		String(order.OCAGroup).
		String(order.Account).
		String(order.OpenClose).
		Int32(order.Origin).
		String(order.OrderRef).
		Bool(order.Transmit).
		Int32(order.ParentID).
		Bool(order.BlockOrder).
		Bool(order.SweepToFill).
		Int32(order.DisplaySize).
		Int32(int32(order.TriggerMethod)).
		Bool(order.OutsideRTH).
		Bool(order.Hidden)

	if contract.SecType == models.SecurityTypePair {
		msgEnc.Int(len(contract.ComboLegs))
		for _, comboLeg := range contract.ComboLegs {
			msgEnc.Int32(comboLeg.ConID).
				Int32(comboLeg.Ratio).
				String(comboLeg.Action).
				String(comboLeg.Exchange).
				Int32(comboLeg.OpenClose).
				Int32(comboLeg.ShortSalesSlot).
				String(comboLeg.DesignatedLocation).
				Int32(comboLeg.ExemptCode)
		}
		msgEnc.Int(len(order.OrderComboLegs))
		for _, orderComboLeg := range order.OrderComboLegs {
			msgEnc.FloatMax(orderComboLeg.Price)
		}
		msgEnc.Int(len(order.SmartComboRoutingParams))
		for _, tv := range order.SmartComboRoutingParams {
			msgEnc.String(tv.Tag).String(tv.Value)
		}
	} else {
		msgEnc.Int(0).Int(0).Int(0)
	}

	msgEnc.String("").
		Float(order.DiscretionaryAmt).
		String(order.GoodAfterTime).
		String(order.GoodTillDate).
		String(order.FAGroup).
		String(order.FAMethod).
		String(order.FAPercentage).
		String(order.ModelCode).
		Int32(order.ShortSaleSlot).
		String(order.DesignatedLocation).
		Int32(order.ExemptCode).
		Int32(int32(order.OCAType)).
		String(string(order.Rule80A)).
		String(order.SettlingFirm).
		Bool(order.AllOrNone).
		Int32Max(order.MinQty).
		FloatMax(order.PercentOffset).
		Bool(false).
		Bool(false).
		FloatMax(nil).
		Int32(int32(order.AuctionStrategy)).
		FloatMax(order.StartingPrice).
		FloatMax(order.StockRefPrice).
		FloatMax(order.Delta).
		FloatMax(order.StockRangeLower).
		FloatMax(order.StockRangeUpper).
		Bool(order.OverridePercentageConstraints).
		FloatMax(order.Volatility).
		Int32(int32(order.VolatilityType)).
		String(order.DeltaNeutralOrderType).
		FloatMax(order.DeltaNeutralAuxPrice)

	if len(order.DeltaNeutralOrderType) > 0 {
		msgEnc.Int32(order.DeltaNeutralConID).
			String(order.DeltaNeutralSettlingFirm).
			String(order.DeltaNeutralClearingAccount).
			String(order.DeltaNeutralClearingIntent).
			String(order.DeltaNeutralOpenClose).
			Bool(order.DeltaNeutralShortSale).
			Int32(order.DeltaNeutralShortSaleSlot).
			String(order.DeltaNeutralDesignatedLocation)
	}

	msgEnc.Bool(order.ContinuousUpdate).
		Int32(order.ReferencePriceType).
		FloatMax(order.TrailStopPrice).
		FloatMax(order.TrailingPercent).
		Int32Max(order.ScaleInitLevelSize).
		Int32Max(order.ScaleSubsLevelSize).
		FloatMax(order.ScalePriceIncrement)

	if order.ScalePriceIncrement != nil && *order.ScalePriceIncrement > 0 {
		msgEnc.FloatMax(order.ScalePriceAdjustValue).
			Int32Max(order.ScalePriceAdjustInterval).
			FloatMax(order.ScaleProfitOffset).
			Bool(order.ScaleAutoReset).
			Int32Max(order.ScaleInitPosition).
			Int32Max(order.ScaleInitFillQty).
			Bool(order.ScaleRandomPercent)
	}

	msgEnc.String(order.ScaleTable).
		String(order.ActiveStartTime).
		String(order.ActiveStopTime).
		String(order.HedgeType)
	if len(order.HedgeType) > 0 {
		msgEnc.String(order.HedgeParam)
	}

	msgEnc.Bool(order.OptOutSmartRouting).
		String(order.ClearingAccount).
		String(order.ClearingIntent).
		Bool(order.NotHeld)

	if contract.DeltaNeutralContract != nil {
		msgEnc.Bool(true).
			Int32(contract.DeltaNeutralContract.ConID).
			Float(contract.DeltaNeutralContract.Delta).
			Float(contract.DeltaNeutralContract.Price)
	} else {
		msgEnc.Bool(false)
	}

	msgEnc.String(order.AlgoStrategy)
	if len(order.AlgoStrategy) > 0 {
		msgEnc.Int(len(order.AlgoParams))
		for _, tv := range order.AlgoParams {
			msgEnc.String(tv.Tag).String(tv.Value)
		}
	}

	msgEnc.String(order.AlgoID).
		Bool(order.WhatIf).
		Marshal((*models.TagValueList)(&order.OrderMiscOptions), 1).
		Bool(order.Solicited).
		Bool(order.RandomizeSize).
		Bool(order.RandomizePrice)

	msgEnc.Int(len(order.Conditions))
	for _, cond := range order.Conditions {
		msgEnc.Int32(int32(cond.Type()))
		encodeOrderConditionFields(msgEnc, models.OrderConditionFields(cond))
	}
	if len(order.Conditions) > 0 {
		msgEnc.Bool(order.ConditionsIgnoreRth).
			Bool(order.ConditionsCancelOrder)
	}

	msgEnc.String(order.AdjustedOrderType).
		FloatMax(order.TriggerPrice).
		FloatMax(order.LmtPriceOffset).
		FloatMax(order.AdjustedStopPrice).
		FloatMax(order.AdjustedStopLimitPrice).
		FloatMax(order.AdjustedTrailingAmount).
		Int32(order.AdjustableTrailingUnit).
		String(order.ExtOperator).
		String(order.SoftDollarTier.Name).
		String(order.SoftDollarTier.Value).
		FloatMax(order.CashQty).
		String(order.Mifid2DecisionMaker).
		String(order.Mifid2DecisionAlgo).
		String(order.Mifid2ExecutionTrader).
		String(order.Mifid2ExecutionAlgo).
		Bool(order.DontUseAutoPriceForHedge).
		Bool(order.IsOmsContainer).
		Bool(order.DiscretionaryUpToLimitPrice).
		Bool(order.UsePriceMgmtAlgo)

	if order.Duration != nil {
		msgEnc.Int32Max(order.Duration)
	} else {
		msgEnc.AddDelim()
	}
	if order.PostToAts != nil {
		msgEnc.Int32Max(order.PostToAts)
	} else {
		msgEnc.AddDelim()
	}

	msgEnc.Bool(order.AutoCancelParent).
		String(order.AdvancedErrorOverride).
		String(order.ManualOrderTime)

	if contract.Exchange == "IBKRATS" {
		msgEnc.Int32Max(order.MinTradeQty)
	}
	switch order.OrderType {
	case "PEG BEST":
		msgEnc.Int32Max(order.MinCompeteSize).
			FloatMax(order.CompeteAgainstBestOffset)
		if order.CompeteAgainstBestOffset != nil && *order.CompeteAgainstBestOffset == models.CompeteAgainstBestOffsetUpToMid {
			msgEnc.FloatMax(order.MidOffsetAtWhole).
				FloatMax(order.MidOffsetAtHalf)
		}
	case "PEG MID":
		msgEnc.FloatMax(order.MidOffsetAtWhole).
			FloatMax(order.MidOffsetAtHalf)
	}

	msgEnc.String(order.CustomerAccount).
		Bool(order.ProfessionalCustomer).
		Bool(order.IncludeOvernight).
		Int32Max(order.ManualOrderIndicator).
		Bool(order.ImbalanceOnly)

	return msgEnc
}

func encodeOrderConditionFields(msgEnc *message.Encoder, fields []any) {
	for _, field := range fields {
		switch v := field.(type) {
		case string:
			msgEnc.String(v)
		case bool:
			msgEnc.Bool(v)
		case int:
			msgEnc.Int(v)
		case int32:
			msgEnc.Int32(v)
		case int64:
			msgEnc.Int64(v)
		case float64:
			msgEnc.Float(v)
		case models.SecurityType:
			msgEnc.String(string(v))
		case models.TriggerMethod:
			msgEnc.Int32(int32(v))
		default:
			msgEnc.AddDelim()
		}
	}
}
