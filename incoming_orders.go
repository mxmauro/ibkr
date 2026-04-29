package ibkr

import (
	"strconv"
	"time"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func (c *Client) processOrderStatusMsg(msgDec *message.Decoder) error {
	orderID := msgDec.Int32()
	status := msgDec.String()
	filled := models.NewDecimalFromMessageDecoder(msgDec)
	remaining := models.NewDecimalFromMessageDecoder(msgDec)
	avgFilledPrice := msgDec.Float()
	permID := msgDec.Int64()
	parentID := msgDec.Int32()
	lastFillPrice := msgDec.Float()
	clientID := msgDec.Int32()
	whyHeld := msgDec.String()
	marketCapPrice := msgDec.Float()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOrderStatusMsg", msgDec.Err())
	}

	return c.processOrderStatusCommon(models.OrderStatusUpdate{
		OrderID:        orderID,
		Status:         status,
		Filled:         filled,
		Remaining:      remaining,
		AvgFilledPrice: avgFilledPrice,
		PermID:         permID,
		ParentID:       parentID,
		LastFillPrice:  lastFillPrice,
		ClientID:       clientID,
		WhyHeld:        whyHeld,
		MarketCapPrice: marketCapPrice,
	})
}

func (c *Client) processOrderStatusProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.OrderStatus{}
	msgDec.Unmarshal(&pb)
	update := models.OrderStatusUpdate{
		OrderID:        msgDec.Int32(pb.OrderId),
		Status:         msgDec.String(pb.Status),
		Filled:         models.NewDecimalFromProtobufDecoder(msgDec, pb.Filled),
		Remaining:      models.NewDecimalFromProtobufDecoder(msgDec, pb.Remaining),
		AvgFilledPrice: msgDec.Float(pb.AvgFillPrice),
		PermID:         msgDec.Int64(pb.PermId),
		ParentID:       msgDec.Int32(pb.ParentId),
		LastFillPrice:  msgDec.Float(pb.LastFillPrice),
		ClientID:       msgDec.Int32(pb.ClientId),
		WhyHeld:        msgDec.String(pb.WhyHeld),
		MarketCapPrice: msgDec.Float(pb.MktCapPrice),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOrderStatusProtobuf", msgDec.Err())
	}

	return c.processOrderStatusCommon(update)
}

func (c *Client) processOrderStatusCommon(update models.OrderStatusUpdate) error {
	if update.PermID != 0 {
		c.orderStore.SetOrderPermanentID(update.OrderID, update.ClientID, update.PermID)
	}
	c.orderStore.UpdateOrderStatus(update)
	return nil
}

func (c *Client) processOpenOrderMsg(msgDec *message.Decoder) error {
	contract, order, orderState := models.DecodeOpenOrderFromMessage(msgDec, int(c.serverVersion))
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOpenOrderMsg", msgDec.Err())
	}

	return c.processOpenOrderCommon(contract, order, orderState)
}

func (c *Client) processOpenOrderProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.OpenOrder{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOpenOrderProtobuf", msgDec.Err())
	}

	contract := models.NewContractFromProtobufDecoder(msgDec, pb.Contract)
	order := models.NewOrderFromProtobufDecoder(msgDec, pb.Order)
	orderState := models.NewOrderStateFromProtobufDecoder(msgDec, pb.OrderState)
	if order.OrderID == 0 {
		order.OrderID = msgDec.Int32(pb.OrderId)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOpenOrderProtobuf", msgDec.Err())
	}

	return c.processOpenOrderCommon(contract, order, orderState)
}

func (c *Client) processOpenOrderCommon(
	contract *models.Contract, order *models.Order, orderState *models.OrderState,
) error {
	if order.PermID != 0 {
		c.orderStore.SetOrderPermanentID(order.OrderID, order.ClientID, order.PermID)
	}
	c.orderStore.RecordOpenOrder(contract, order, orderState)
	c.reqMgr.withFirstRequestWithoutID(int(common.OPEN_ORDER_END), func(resp interface{}) (bool, error) {
		r := resp.(*models.OpenOrdersResponse)
		r.Orders = append(r.Orders, models.OpenOrderData{
			Contract:   contract,
			Order:      order,
			OrderState: orderState,
		})
		return false, nil
	})
	return nil
}

func (c *Client) processExecutionDetailsMsg(msgDec *message.Decoder) error {
	reqID := msgDec.Int32()
	orderID := msgDec.Int32()

	contract := models.NewContract()
	contract.ConID = msgDec.Int32()
	contract.Symbol = msgDec.String()
	contract.SecType = models.SecurityType(msgDec.String())
	contract.LastTradeDateOrContractMonth = msgDec.String()
	contract.Strike = msgDec.FloatMax()
	contract.Right = msgDec.String()
	contract.Multiplier = msgDec.FloatMax()
	contract.Exchange = msgDec.String()
	contract.Currency = msgDec.String()
	contract.LocalSymbol = msgDec.String()
	contract.TradingClass = msgDec.String()

	execution := models.NewExecution()
	execution.OrderID = orderID
	execution.ExecID = msgDec.String()
	execution.Time = msgDec.String()
	execution.AcctNumber = msgDec.String()
	execution.Exchange = msgDec.String()
	execution.Side = msgDec.String()
	execution.Shares = models.NewDecimalMaxFromMessageDecoder(msgDec)
	execution.Price = msgDec.Float()
	execution.PermID = msgDec.Int64()
	execution.ClientID = msgDec.Int32()
	if msgDec.Bool() {
		execution.Liquidation = 1
	}
	execution.CumQty = models.NewDecimalMaxFromMessageDecoder(msgDec)
	execution.AvgPrice = msgDec.Float()
	execution.OrderRef = msgDec.String()
	execution.EVRule = msgDec.String()
	execution.EVMultiplier = msgDec.Float()
	execution.ModelCode = msgDec.String()
	execution.LastLiquidity = models.Liquidities(msgDec.Int32())
	execution.PendingPriceRevision = msgDec.Bool()
	execution.Submitter = msgDec.String()
	execution.OptExerciseOrLapse = models.OptionExercise(msgDec.Int32())
	if msgDec.Err() != nil {
		return wrapDecoderErr("processExecutionDetailsMsg", msgDec.Err())
	}

	return c.processExecutionDetailsCommon(reqID, contract, execution)
}

func (c *Client) processExecutionDetailsMsgProtoBuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ExecutionDetails{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.Int32(pb.ReqId)
	contract := models.NewContractFromProtobufDecoder(msgDec, pb.Contract)
	execution := models.NewExecutionFromProtobufDecoder(msgDec, pb.Execution)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processExecutionDetailsMsgProtoBuf", msgDec.Err())
	}

	return c.processExecutionDetailsCommon(reqID, contract, execution)
}

func (c *Client) processExecutionDetailsCommon(
	reqID int32, contract *models.Contract, execution *models.Execution,
) error {
	if execution == nil {
		return nil
	}
	c.orderStore.RecordExecution(contract, execution)
	if reqID > 0 {
		c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
			r := resp.(*models.ExecutionsResponse)
			r.Executions = append(r.Executions, models.ExecutionData{
				RequestID: reqID,
				Contract:  contract,
				Execution: execution,
			})
			return false, nil
		})
	}
	return nil
}

func (c *Client) processOpenOrdersEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOpenOrdersEndMsg", msgDec.Err())
	}
	return c.processOpenOrdersEndCommon()
}

func (c *Client) processOpenOrdersEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.OpenOrdersEnd{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOpenOrdersEndProtobuf", msgDec.Err())
	}
	return c.processOpenOrdersEndCommon()
}

func (c *Client) processOpenOrdersEndCommon() error {
	c.reqMgr.withFirstRequestWithoutID(int(common.OPEN_ORDER_END), func(_ interface{}) (bool, error) {
		return true, nil
	})
	return nil
}

func (c *Client) processExecutionDetailsEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	reqID := msgDec.Int32()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processExecutionDetailsEndMsg", msgDec.Err())
	}
	return c.processExecutionDetailsEndCommon(reqID)
}

func (c *Client) processExecutionDetailsEndMsgProtoBuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ExecutionDetailsEnd{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.Int32(pb.ReqId)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processExecutionDetailsEndMsgProtoBuf", msgDec.Err())
	}
	return c.processExecutionDetailsEndCommon(reqID)
}

func (c *Client) processExecutionDetailsEndCommon(reqID int32) error {
	if reqID > 0 {
		c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
			return true, nil
		})
	}
	return nil
}

func (c *Client) processCommissionAndFeesReportMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	report := models.NewCommissionAndFeesReport()
	report.ExecID = msgDec.String()
	report.CommissionAndFees = msgDec.Float()
	report.Currency = msgDec.String()
	report.RealizedPNL = msgDec.Float()
	report.Yield = msgDec.Float()
	redemptionDate := msgDec.String()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCommissionAndFeesReportMsg", msgDec.Err())
	}
	report.YieldRedemptionDate = parseCommissionAndFeesReportDate(redemptionDate)
	c.orderStore.RecordCommissionAndFeesReport(report)
	return nil
}

func (c *Client) processCommissionAndFeesReportProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.CommissionAndFeesReport{}
	msgDec.Unmarshal(&pb)
	report := models.NewCommissionAndFeesReport()
	report.ExecID = msgDec.String(pb.ExecId)
	report.CommissionAndFees = msgDec.Float(pb.CommissionAndFees)
	report.Currency = msgDec.String(pb.Currency)
	report.RealizedPNL = msgDec.Float(pb.RealizedPNL)
	report.Yield = msgDec.Float(pb.BondYield)
	report.YieldRedemptionDate = parseCommissionAndFeesReportDate(msgDec.String(pb.YieldRedemptionDate))
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCommissionAndFeesReportProtobuf", msgDec.Err())
	}
	c.orderStore.RecordCommissionAndFeesReport(report)
	return nil
}

func (c *Client) processOrderBoundMsg(msgDec *message.Decoder) error {
	permID := msgDec.Int64()
	clientID := msgDec.Int32()
	orderID := msgDec.Int32()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOrderBoundMsg", msgDec.Err())
	}
	c.orderStore.BindOrder(permID, clientID, orderID)
	return nil
}

func (c *Client) processOrderBoundProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.OrderBound{}
	msgDec.Unmarshal(&pb)
	permID := msgDec.Int64(pb.PermId)
	clientID := msgDec.Int32(pb.ClientId)
	orderID := msgDec.Int32(pb.OrderId)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processOrderBoundProtobuf", msgDec.Err())
	}
	c.orderStore.BindOrder(permID, clientID, orderID)
	return nil
}

func (c *Client) processCompletedOrderMsg(msgDec *message.Decoder) error {
	contract, order, orderState := models.DecodeCompletedOrderFromMessage(msgDec, int(c.serverVersion))
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCompletedOrderMsg", msgDec.Err())
	}

	return c.processCompletedOrderCommon(contract, order, orderState)
}

func (c *Client) processCompletedOrderProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.CompletedOrder{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCompletedOrderProtobuf", msgDec.Err())
	}

	contract := models.NewContractFromProtobufDecoder(msgDec, pb.Contract)
	order := models.NewOrderFromProtobufDecoder(msgDec, pb.Order)
	orderState := models.NewOrderStateFromProtobufDecoder(msgDec, pb.OrderState)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCompletedOrderProtobuf", msgDec.Err())
	}

	return c.processCompletedOrderCommon(contract, order, orderState)
}

func (c *Client) processCompletedOrderCommon(
	contract *models.Contract, order *models.Order, orderState *models.OrderState,
) error {
	if order.PermID != 0 {
		c.orderStore.SetOrderPermanentID(order.OrderID, order.ClientID, order.PermID)
	}
	c.orderStore.RecordCompletedOrder(contract, order, orderState)
	c.reqMgr.withFirstRequestWithoutID(int(common.COMPLETED_ORDERS_END), func(resp interface{}) (bool, error) {
		r := resp.(*models.CompletedOrdersResponse)
		r.Orders = append(r.Orders, models.CompletedOrderData{
			Contract:   contract,
			Order:      order,
			OrderState: orderState,
		})
		return false, nil
	})
	return nil
}

func (c *Client) processCompletedOrdersEndMsg(msgDec *message.Decoder) error {
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCompletedOrdersEndMsg", msgDec.Err())
	}
	return c.processCompletedOrdersEndCommon()
}

func (c *Client) processCompletedOrdersEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.CompletedOrdersEnd{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCompletedOrdersEndProtobuf", msgDec.Err())
	}
	return c.processCompletedOrdersEndCommon()
}

func (c *Client) processCompletedOrdersEndCommon() error {
	c.reqMgr.withFirstRequestWithoutID(int(common.COMPLETED_ORDERS_END), func(_ interface{}) (bool, error) {
		return true, nil
	})
	return nil
}

func parseCommissionAndFeesReportDate(value string) time.Time {
	if len(value) == 0 {
		return time.Time{}
	}

	if len(value) == 8 {
		if t, err := time.Parse("20060102", value); err == nil {
			return t
		}
	}

	if unixDate, err := strconv.ParseInt(value, 10, 64); err == nil && unixDate > 0 {
		return time.Unix(unixDate, 0).UTC()
	}

	return time.Time{}
}
