package ibkr

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protowire"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func (c *Client) getIncomingMessageID(msg []byte) (msgID uint32, usingProtobuf bool, err error) {
	if len(msg) >= 4 {
		msgID = binary.BigEndian.Uint32(msg)
		if msgID >= common.PROTOBUF_MSG_ID {
			usingProtobuf = true
			msgID -= common.PROTOBUF_MSG_ID
		}
	} else {
		err = errInvalidReceivedMessage
	}
	return
}

func wrapDecoderErr(decoderName string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s decode failed: %w", decoderName, err)
}

func isWarningMessage(errMsg string) bool {
	const warningPrefix = "Warning:"

	for start := 0; start < len(errMsg); {
		i := start
		for i < len(errMsg) && (errMsg[i] == ' ' || errMsg[i] == '\t') {
			i++
		}
		if len(errMsg)-i >= len(warningPrefix) && errMsg[i:i+len(warningPrefix)] == warningPrefix {
			return true
		}

		for start < len(errMsg) && errMsg[start] != '\n' && errMsg[start] != '\r' {
			start++
		}
		if start < len(errMsg) && errMsg[start] == '\r' {
			start++
			if start < len(errMsg) && errMsg[start] == '\n' {
				start++
			}
		} else if start < len(errMsg) && errMsg[start] == '\n' {
			start++
		}
	}
	return false
}

func (c *Client) processIncomingMessage(msg []byte) error {
	msgID, usingProtobuf, err := c.getIncomingMessageID(msg)
	if err != nil {
		return err
	}

	if usingProtobuf {
		msgDec := protofmt.NewDecoder(msg[4:])
		switch msgID {
		case common.TICK_PRICE:
			return c.processTickPriceProtobuf(msgDec)
		case common.TICK_SIZE:
			return c.processTickSizeProtobuf(msgDec)
		case common.TICK_OPTION_COMPUTATION:
			return c.processTickOptionComputationProtobuf(msgDec)
		case common.TICK_GENERIC:
			return c.processTickGenericProtobuf(msgDec)
		case common.TICK_STRING:
			return c.processTickStringProtobuf(msgDec)
		case common.ORDER_STATUS:
			return c.processOrderStatusProtobuf(msgDec)
		case common.ERR_MSG:
			return c.processErrorMessageProtobuf(msgDec)
		case common.OPEN_ORDER:
			return c.processOpenOrderProtobuf(msgDec)
		case common.ACCT_VALUE:
			return c.processAccountValueProtobuf(msgDec)
		case common.PORTFOLIO_VALUE:
			return c.processPortfolioValueProtobuf(msgDec)
		case common.ACCT_UPDATE_TIME:
			return c.processAccountUpdateTimeProtobuf(msgDec)
		case common.CONTRACT_DATA:
			return c.processContractDataProtobuf(msgDec)
		case common.BOND_CONTRACT_DATA:
			return c.processBondContractDataProtobuf(msgDec)
		case common.MARKET_DEPTH:
			return c.processMarketDepthProtobuf(msgDec)
		case common.MARKET_DEPTH_L2:
			return c.processMarketDepthL2Protobuf(msgDec)
		case common.MANAGED_ACCTS:
			return c.processManagedAccountsProtobuf(msgDec)
		case common.HISTORICAL_DATA:
			return c.processHistoricalDataProtobuf(msgDec)
		case common.REAL_TIME_BARS:
			return c.processRealTimeBarsProtobuf(msgDec)
		case common.CONTRACT_DATA_END:
			return c.processContractDataEndProtobuf(msgDec)
		case common.OPEN_ORDER_END:
			return c.processOpenOrdersEndProtobuf(msgDec)
		case common.TICK_SNAPSHOT_END:
			return c.processTickSnapshotEndProtobuf(msgDec)
		case common.MARKET_DATA_TYPE:
			return c.processMarketDataTypeProtobuf(msgDec)
		case common.TICK_REQ_PARAMS:
			return nil // Ignore this message. We don't use the tick request parameters.
		case common.SYMBOL_SAMPLES:
			return c.processSymbolSamplesProtobuf(msgDec)
		case common.HEAD_TIMESTAMP:
			return c.processHeadTimestampProtobuf(msgDec)
		case common.HISTORICAL_DATA_UPDATE:
			return c.processHistoricalDataUpdateProtobuf(msgDec)
		case common.EXECUTION_DATA:
			return c.processExecutionDetailsMsgProtoBuf(msgDec)
		case common.HISTORICAL_TICKS:
			return c.processHistoricalTicksProtobuf(msgDec)
		case common.HISTORICAL_TICKS_BID_ASK:
			return c.processHistoricalTicksBidAskProtobuf(msgDec)
		case common.HISTORICAL_TICKS_LAST:
			return c.processHistoricalTicksLastProtobuf(msgDec)
		case common.SECURITY_DEFINITION_OPTION_PARAMETER:
			return c.processSecurityDefinitionOptionalParameterProtobuf(msgDec)
		case common.SECURITY_DEFINITION_OPTION_PARAMETER_END:
			return c.processSecurityDefinitionOptionalParameterEndProtobuf(msgDec)
		case common.TICK_BY_TICK:
			return c.processTickByTickDataProtobuf(msgDec)
		case common.HISTORICAL_DATA_END:
			return c.processHistoricalDataEndProtobuf(msgDec)
		case common.EXECUTION_DATA_END:
			return c.processExecutionDetailsEndMsgProtoBuf(msgDec)
		case common.COMMISSION_AND_FEES_REPORT:
			return c.processCommissionAndFeesReportProtobuf(msgDec)
		case common.ORDER_BOUND:
			return c.processOrderBoundProtobuf(msgDec)
		case common.COMPLETED_ORDER:
			return c.processCompletedOrderProtobuf(msgDec)
		case common.COMPLETED_ORDERS_END:
			return c.processCompletedOrdersEndProtobuf(msgDec)
		case common.POSITION_DATA:
			return c.processPositionDataProtobuf(msgDec)
		case common.POSITION_END:
			return c.processPositionEndProtobuf(msgDec)
		case common.ACCOUNT_SUMMARY:
			return c.processAccountSummaryProtobuf(msgDec)
		case common.ACCOUNT_SUMMARY_END:
			return c.processAccountSummaryEndProtobuf(msgDec)
		case common.HISTOGRAM_DATA:
			return c.processHistogramDataProtobuf(msgDec)
		case common.MARKET_RULE:
			return c.processMarketRuleProtobuf(msgDec)
		}
	} else {
		msgDec := message.NewDecoder(msg[4:])
		switch msgID {
		case math.MaxUint32:
			return net.ErrClosed
		case common.TICK_PRICE:
			return c.processTickPriceMsg(msgDec)
		case common.TICK_SIZE:
			return c.processTickSizeMsg(msgDec)
		case common.TICK_OPTION_COMPUTATION:
			return c.processTickOptionComputationMsg(msgDec)
		case common.TICK_GENERIC:
			return c.processTickGenericMsg(msgDec)
		case common.TICK_STRING:
			return c.processTickStringMsg(msgDec)
		case common.TICK_EFP:
			return c.processTickEfpMsg(msgDec)
		case common.ORDER_STATUS:
			return c.processOrderStatusMsg(msgDec)
		case common.ERR_MSG:
			return c.processErrorMessageMsg(msgDec)
		case common.OPEN_ORDER:
			return c.processOpenOrderMsg(msgDec)
		case common.ACCT_VALUE:
			return c.processAccountValueMsg(msgDec)
		case common.PORTFOLIO_VALUE:
			return c.processPortfolioValueMsg(msgDec)
		case common.ACCT_UPDATE_TIME:
			return c.processAccountUpdateTimeMsg(msgDec)
		case common.NEXT_VALID_ID:
			return nil // Ignore this message
		case common.CONTRACT_DATA:
			return c.processContractDataMsg(msgDec)
		case common.BOND_CONTRACT_DATA:
			return c.processBondContractDataMsg(msgDec)
		case common.EXECUTION_DATA:
			return c.processExecutionDetailsMsg(msgDec)
		case common.MARKET_DEPTH:
			return c.processMarketDepthMsg(msgDec)
		case common.MARKET_DEPTH_L2:
			return c.processMarketDepthL2Msg(msgDec)
		case common.MANAGED_ACCTS:
			return c.processManagedAccountsMsg(msgDec)
		case common.HISTORICAL_DATA:
			return c.processHistoricalDataMsg(msgDec)
		case common.CURRENT_TIME:
			return nil // Ignore this message. We use the one having milliseconds
		case common.REAL_TIME_BARS:
			return c.processRealTimeBarsMsg(msgDec)
		case common.CONTRACT_DATA_END:
			return c.processContractDataEndMsg(msgDec)
		case common.OPEN_ORDER_END:
			return c.processOpenOrdersEndMsg(msgDec)
		case common.ACCT_DOWNLOAD_END:
			return c.processAccountDownloadEndMsg(msgDec)
		case common.TICK_SNAPSHOT_END:
			return c.processTickSnapshotEndMsg(msgDec)
		case common.MARKET_DATA_TYPE:
			return c.processMarketDataTypeMsg(msgDec)
		case common.COMMISSION_AND_FEES_REPORT:
			return c.processCommissionAndFeesReportMsg(msgDec)
		case common.POSITION_DATA:
			return c.processPositionDataMsg(msgDec)
		case common.POSITION_END:
			return c.processPositionEndMsg(msgDec)
		case common.ACCOUNT_SUMMARY:
			return c.processAccountSummaryMsg(msgDec)
		case common.ACCOUNT_SUMMARY_END:
			return c.processAccountSummaryEndMsg(msgDec)
		case common.TICK_REQ_PARAMS:
			return nil // Ignore this message. We don't use the tick request parameters.
		case common.SYMBOL_SAMPLES:
			return c.processSymbolSamplesMsg(msgDec)
		case common.MKT_DEPTH_EXCHANGES:
			return c.processMktDepthExchangesMsg(msgDec)
		case common.HEAD_TIMESTAMP:
			return c.processHeadTimestampMsg(msgDec)
		case common.HISTOGRAM_DATA:
			return c.processHistogramDataMsg(msgDec)
		case common.HISTORICAL_DATA_UPDATE:
			return c.processHistoricalDataUpdateMsg(msgDec)
		case common.SMART_COMPONENTS:
			return c.processSmartComponentsMsg(msgDec)
		case common.MARKET_RULE:
			return c.processMarketRuleMsg(msgDec)
		case common.HISTORICAL_TICKS:
			return c.processHistoricalTicksMsg(msgDec)
		case common.HISTORICAL_TICKS_BID_ASK:
			return c.processHistoricalTicksBidAskMsg(msgDec)
		case common.HISTORICAL_TICKS_LAST:
			return c.processHistoricalTicksLastMsg(msgDec)
		case common.TICK_BY_TICK:
			return c.processTickByTickDataMsg(msgDec)
		case common.ORDER_BOUND:
			return c.processOrderBoundMsg(msgDec)
		case common.COMPLETED_ORDER:
			return c.processCompletedOrderMsg(msgDec)
		case common.COMPLETED_ORDERS_END:
			return c.processCompletedOrdersEndMsg(msgDec)
		case common.HISTORICAL_DATA_END:
			return c.processHistoricalDataEndMsg(msgDec)
		case common.CURRENT_TIME_IN_MILLIS:
			return c.processCurrentTimeInMillisMsg(msgDec)
		}
	}

	// Raise the event if an event handler is present
	if c.eventsHandler != nil {
		c.eventsHandler.ReceivedUnknownMessage(msgID)
	}

	// Done
	return nil
}

func (c *Client) processTickPriceMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	price := msgDec.Float()
	size := models.NewDecimalFromMessageDecoder(msgDec) // ver 2 field
	attrMask := msgDec.Int32()                          // ver 3 field
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickPriceMsg", msgDec.Err())
	}

	// Done
	return c.processTickPriceCommon(tickerID, tickType, price, size, attrMask)
}

func (c *Client) processTickPriceProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickPrice{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickType(msgDec.Int32(pb.TickType))
	price := msgDec.Float(pb.Price)
	size := models.NewDecimalFromProtobufDecoder(msgDec, pb.Size)
	attrMask := msgDec.Int32(pb.AttrMask)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickPriceProtobuf", msgDec.Err())
	}

	// Done
	return c.processTickPriceCommon(tickerID, tickType, price, size, attrMask)
}

func (c *Client) processTickPriceCommon(
	tickerID int32, tickType models.TickType, price float64, size models.Decimal, attrMask int32,
) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		// Notify
		data := models.NewTopMarketDataPrice(tickType)
		data.Price = price
		data.CanAutoExecute = attrMask&0x1 != 0
		data.PastLimit = attrMask&0x2 != 0
		data.PreOpen = attrMask&0x4 != 0
		resp.enqueue(data)

		var sizeData *models.TopMarketDataSize
		switch tickType {
		case models.TickTypeBid:
			sizeData = models.NewTopMarketDataSize(models.TickTypeBidSize)
		case models.TickTypeAsk:
			sizeData = models.NewTopMarketDataSize(models.TickTypeAskSize)
		case models.TickTypeLast:
			sizeData = models.NewTopMarketDataSize(models.TickTypeLastSize)
		case models.TickTypeDelayedBid:
			sizeData = models.NewTopMarketDataSize(models.TickTypeDelayedBidSize)
		case models.TickTypeDelayedAsk:
			sizeData = models.NewTopMarketDataSize(models.TickTypeDelayedAskSize)
		case models.TickTypeDelayedLast:
			sizeData = models.NewTopMarketDataSize(models.TickTypeDelayedLastSize)
		}
		if sizeData != nil {
			sizeData.Size = size
			resp.enqueue(sizeData)
		}

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processTickSizeMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	size := models.NewDecimalFromMessageDecoder(msgDec)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickSizeMsg", msgDec.Err())
	}

	// Done
	return c.processTickSizeCommon(tickerID, tickType, size)
}

func (c *Client) processTickSizeProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickSize{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickType(msgDec.Int32(pb.TickType))
	size := models.NewDecimalFromProtobufDecoder(msgDec, pb.Size)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickSizeProtobuf", msgDec.Err())
	}

	// Done
	return c.processTickSizeCommon(tickerID, tickType, size)
}

func (c *Client) processTickSizeCommon(tickerID int32, tickType models.TickType, size models.Decimal) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		// Notify
		data := models.NewTopMarketDataSize(tickType)
		data.Size = size
		resp.enqueue(data)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processTickOptionComputationMsg(msgDec *message.Decoder) error {
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	tickAttrib := msgDec.Int32()
	impliedVol := msgDec.FloatMax()
	if impliedVol != nil && utils.EqualFloat(*impliedVol, -1) { // -1 is the "not computed" indicator
		impliedVol = nil
	}
	delta := msgDec.FloatMax()
	if delta != nil && utils.EqualFloat(*delta, -2) { // -2 is the "not computed" indicator
		delta = nil
	}
	var price *float64
	var pvDividend *float64
	if tickType == models.TickTypeModelOptionComputation || tickType == models.TickTypeDelayedModelOptionComputation {
		price = msgDec.FloatMax()
		if price != nil && utils.EqualFloat(*price, -1) { // -1 is the "not computed" indicator
			price = nil
		}
		pvDividend = msgDec.FloatMax()
		if pvDividend != nil && utils.EqualFloat(*pvDividend, -1) { // -1 is the "not computed" indicator
			pvDividend = nil
		}
	}
	gamma := msgDec.FloatMax()
	if gamma != nil && utils.EqualFloat(*gamma, -2) { // -2 is the "not yet computed" indicator
		gamma = nil
	}
	vega := msgDec.FloatMax()
	if vega != nil && utils.EqualFloat(*vega, -2) { // -2 is the "not yet computed" indicator
		vega = nil
	}
	theta := msgDec.FloatMax()
	if theta != nil && utils.EqualFloat(*theta, -2) { // -2 is the "not yet computed" indicator
		theta = nil
	}
	undPrice := msgDec.FloatMax()
	if undPrice != nil && utils.EqualFloat(*undPrice, -1) { // -1 is the "not computed" indicator
		undPrice = nil
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickOptionComputationMsg", msgDec.Err())
	}

	// Done
	return c.processTickOptionComputationCommon(
		tickerID, tickType, tickAttrib, impliedVol, delta, price, pvDividend, gamma, vega, theta, undPrice,
	)
}

func (c *Client) processTickOptionComputationProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickOptionComputation{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickType(msgDec.Int32(pb.TickType))
	tickAttrib := msgDec.Int32(pb.TickAttrib)
	impliedVol := msgDec.FloatMax(pb.ImpliedVol)
	if impliedVol != nil && utils.EqualFloat(*impliedVol, -1) { // -1 is the "not computed" indicator
		impliedVol = nil
	}
	delta := msgDec.FloatMax(pb.Delta)
	if delta != nil && utils.EqualFloat(*delta, -2) { // -2 is the "not computed" indicator
		delta = nil
	}
	var price *float64
	var pvDividend *float64
	if tickType == models.TickTypeModelOptionComputation || tickType == models.TickTypeDelayedModelOptionComputation {
		price = msgDec.FloatMax(pb.OptPrice)
		if price != nil && utils.EqualFloat(*price, -1) { // -1 is the "not computed" indicator
			price = nil
		}
		pvDividend = msgDec.FloatMax(pb.PvDividend)
		if pvDividend != nil && utils.EqualFloat(*pvDividend, -1) { // -1 is the "not computed" indicator
			pvDividend = nil
		}
	}
	gamma := msgDec.FloatMax(pb.Gamma)
	if gamma != nil && utils.EqualFloat(*gamma, -2) { // -2 is the "not yet computed" indicator
		gamma = nil
	}
	vega := msgDec.FloatMax(pb.Vega)
	if vega != nil && utils.EqualFloat(*vega, -2) { // -2 is the "not yet computed" indicator
		vega = nil
	}
	theta := msgDec.FloatMax(pb.Theta)
	if theta != nil && utils.EqualFloat(*theta, -2) { // -2 is the "not yet computed" indicator
		theta = nil
	}
	undPrice := msgDec.FloatMax(pb.UndPrice)
	if undPrice != nil && utils.EqualFloat(*undPrice, -1) { // -1 is the "not computed" indicator
		undPrice = nil
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickOptionComputationProtobuf", msgDec.Err())
	}

	// Done
	return c.processTickOptionComputationCommon(
		tickerID, tickType, tickAttrib, impliedVol, delta, price, pvDividend, gamma, vega, theta, undPrice,
	)
}

func (c *Client) processTickOptionComputationCommon(
	tickerID int32, tickType models.TickType, tickAttrib int32, impliedVol *float64, delta *float64, price *float64,
	pvDividend *float64, gamma *float64, vega *float64, theta *float64, undPrice *float64,
) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		// Notify
		data := models.NewTopMarketDataOptionComputation(tickType)
		data.IsPriceBased = tickAttrib != 0
		data.ImpliedVolatility = impliedVol
		data.Delta = delta
		data.Price = price
		data.PvDividend = pvDividend
		data.Gamma = gamma
		data.Vega = vega
		data.Theta = theta
		data.UnderlyingPrice = undPrice
		resp.enqueue(data)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processTickGenericMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	value := msgDec.Float()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickGenericMsg", msgDec.Err())
	}

	// Done
	return c.processTickGenericCommon(tickerID, tickType, value)
}

func (c *Client) processTickGenericProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickGeneric{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickType(msgDec.Int32(pb.TickType))
	value := msgDec.Float(pb.Value)

	// Done
	return c.processTickGenericCommon(tickerID, tickType, value)
}

func (c *Client) processTickGenericCommon(tickerID int32, tickType models.TickType, value float64) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		// Notify
		data := models.NewTopMarketDataGeneric(tickType)
		data.Value = value
		resp.enqueue(data)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processTickStringMsg(msgDec *message.Decoder) error {
	var ts time.Time

	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	value := msgDec.String()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickStringMsg", msgDec.Err())
	}

	switch tickType {
	case models.TickTypeLastTimestamp:
		fallthrough
	case models.TickTypeDelayedLastTimestamp:
		fallthrough
	case models.TickTypeLastRegulatoryTime:
		secs, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		ts = time.Unix(secs, 0).UTC()
	}

	// Done
	return c.processTickStringCommon(tickerID, tickType, value, ts)
}

func (c *Client) processTickStringProtobuf(msgDec *protofmt.Decoder) error {
	var ts time.Time

	pb := protobuf.TickString{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickType(msgDec.Int32(pb.TickType))
	value := msgDec.String(pb.Value)

	switch tickType {
	case models.TickTypeLastTimestamp:
		fallthrough
	case models.TickTypeDelayedLastTimestamp:
		fallthrough
	case models.TickTypeLastRegulatoryTime:
		secs, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		ts = time.Unix(secs, 0).UTC()
	}

	// Done
	return c.processTickStringCommon(tickerID, tickType, value, ts)
}

func (c *Client) processTickStringCommon(tickerID int32, tickType models.TickType, value string, ts time.Time) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		switch tickType {
		case models.TickTypeLastTimestamp:
			fallthrough
		case models.TickTypeDelayedLastTimestamp:
			fallthrough
		case models.TickTypeLastRegulatoryTime:
			// Notify
			data := models.NewTopMarketDataTimestamp(tickType)
			data.Timestamp = ts
			resp.enqueue(data)

		default:
			// Notify
			data := models.NewTopMarketDataString(tickType)
			data.Value = value
			resp.enqueue(data)
		}

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processTickEfpMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	tickType := models.TickType(msgDec.Int32())
	basisPoints := msgDec.Float()
	formattedBasisPoints := msgDec.String()
	impliedFuturesPrice := msgDec.Float()
	holdDays := msgDec.Int32()
	futureLastTradeDate := msgDec.String()
	dividendImpact := msgDec.Float()
	dividendsToLastTradeDate := msgDec.Float()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickEfpMsg", msgDec.Err())
	}

	// Done
	return c.processTickEfpCommon(
		tickerID, tickType, basisPoints, formattedBasisPoints, impliedFuturesPrice, holdDays, futureLastTradeDate,
		dividendImpact, dividendsToLastTradeDate,
	)
}

func (c *Client) processTickEfpCommon(
	tickerID int32, tickType models.TickType, basisPoints float64, formattedBasisPoints string,
	impliedFuturesPrice float64, holdDays int32, futureLastTradeDate string, dividendImpact float64,
	dividendsToLastTradeDate float64,
) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*topMarketDataStreamResponse)

		// Notify
		data := models.NewTopMarketDataEFP(tickType)
		data.BasisPoints = basisPoints
		data.FormattedBasisPoints = formattedBasisPoints
		data.ImpliedFuturesPrice = impliedFuturesPrice
		data.HoldDays = holdDays
		data.FutureLastTradeDate = futureLastTradeDate
		data.DividendImpact = dividendImpact
		data.DividendsToLastTradeDate = dividendsToLastTradeDate
		resp.enqueue(data)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processErrorMessageMsg(msgDec *message.Decoder) error {
	// Get the optional originating request ID
	reqID := msgDec.RequestID(true)
	code := int(msgDec.Int32())
	errMsg := msgDec.String()
	advancedOrderRejectJson := msgDec.String()
	ts := msgDec.EpochTimestamp(true)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processErrorMessageMsg", msgDec.Err())
	}
	if reqID == 0 {
		return errInvalidRequestID
	}

	// Done
	return c.processErrorMessageCommon(reqID, code, errMsg, advancedOrderRejectJson, ts)
}

func (c *Client) processErrorMessageProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ErrorMessage{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.Id, true)
	code := int(msgDec.Int32(pb.ErrorCode))
	errMsg := msgDec.String(pb.ErrorMsg)
	advancedOrderRejectJson := msgDec.String(pb.AdvancedOrderRejectJson)
	ts := msgDec.EpochTimestamp(pb.ErrorTime, true)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processErrorMessageProtobuf", msgDec.Err())
	}
	if reqID == 0 {
		return errInvalidRequestID
	}

	// Done
	return c.processErrorMessageCommon(reqID, code, errMsg, advancedOrderRejectJson, ts)
}

func (c *Client) processErrorMessageCommon(
	reqID int32, code int, errMsg string, advancedOrderRejectJson string, ts time.Time,
) error {
	// Ignore the following error codes
	if code == 300 || code == 10167 {
		// 300 - "Can't find EId with tickerId: ###" messages because they are sent when we try to cancel a request,
		//       and it does no longer exist.
		// 10167 - "Requested market data is not subscribed. Delayed market data is not enabled" warning when delayed
		//         data is requested.
		return nil
	}

	if code == 317 {
		// 317 - Market depth data has been RESET. Please empty deep book contents before applying any new entries.
		return c.processMarketDepthCommon(reqID, 0, models.MarketDepthDataOperationClear, false, 0, models.DecimalZero)
	}

	// Classify multiline order messages so warning-only updates do not fail requests.
	isWarning := isWarningMessage(errMsg)

	// Process the response
	if reqID > 0 {
		if (!isWarning) && c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
			// Request-scoped hard errors should complete the pending request immediately.
			return false, newRequestError(ts, code, errMsg, advancedOrderRejectJson)
		}) {
			return nil
		}

		if !isWarning {
			c.orderStore.RecordOrderError(reqID, ts, code, errMsg, advancedOrderRejectJson)
		} else {
			c.orderStore.RecordOrderWarning(reqID, ts, code, errMsg, advancedOrderRejectJson)
		}
	}

	// Notify
	if c.eventsHandler != nil {
		if !isWarning {
			c.eventsHandler.Error(ts, code, errMsg, advancedOrderRejectJson)
		} else {
			c.eventsHandler.Warning(ts, code, errMsg, advancedOrderRejectJson)
		}
	}

	// Done
	return nil
}

func (c *Client) processContractDataMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	cd := models.NewContractDetails()
	cd.Contract.Symbol = msgDec.String()
	cd.Contract.SecType = models.SecurityType(msgDec.String())
	cd.Contract.LastTradeDateOrContractMonth = msgDec.String()
	cd.Contract.LastTradeDate = msgDec.String()
	cd.Contract.Strike = msgDec.FloatMax()
	cd.Contract.Right = msgDec.String()
	cd.Contract.Exchange = msgDec.String()
	cd.Contract.Currency = msgDec.String()
	cd.Contract.LocalSymbol = msgDec.String()
	cd.MarketName = msgDec.String()
	cd.Contract.TradingClass = msgDec.String()
	cd.Contract.ConID = msgDec.Int32()
	cd.MinTick = msgDec.FloatMax()
	cd.Contract.Multiplier = msgDec.FloatMax()
	cd.OrderTypes = msgDec.String()
	cd.ValidExchanges = msgDec.String()
	cd.PriceMagnifier = msgDec.Int32Max()
	cd.UnderConID = msgDec.Int32Max()
	cd.LongName = msgDec.EscapedString()
	cd.Contract.PrimaryExchange = msgDec.String()
	cd.ContractMonth = msgDec.String()
	cd.Industry = msgDec.String()
	cd.Category = msgDec.String()
	cd.Subcategory = msgDec.String()
	cd.TimeZoneID = msgDec.String()
	cd.TradingHours = msgDec.String()
	cd.LiquidHours = msgDec.String()
	cd.EVRule = msgDec.String()
	cd.EVMultiplier = msgDec.Float()
	secIDListCount := int(msgDec.Int32())
	if secIDListCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative security id count: %d", secIDListCount))
		return wrapDecoderErr("processContractDataMsg", msgDec.Err())
	}
	cd.SecIDList = make(models.TagValueList, 0, secIDListCount)
	for i := 0; i < secIDListCount; i++ {
		tagValue := models.NewTagValue()
		tagValue.Tag = msgDec.String()
		tagValue.Value = msgDec.String()
		cd.SecIDList = append(cd.SecIDList, tagValue)
	}
	cd.AggGroup = msgDec.Int32()
	cd.UnderSymbol = msgDec.String()
	cd.UnderSecType = models.SecurityType(msgDec.String())
	cd.MarketRuleIDs = msgDec.String()
	cd.RealExpirationDate = msgDec.String()
	cd.StockType = msgDec.String()
	cd.MinSize = models.NewDecimalMaxFromMessageDecoder(msgDec)
	cd.SizeIncrement = models.NewDecimalMaxFromMessageDecoder(msgDec)
	cd.SuggestedSizeIncrement = models.NewDecimalMaxFromMessageDecoder(msgDec)
	if cd.Contract.SecType == models.SecurityTypeMutualFund {
		cd.FundName = msgDec.String()
		cd.FundFamily = msgDec.String()
		cd.FundType = msgDec.String()
		cd.FundFrontLoad = msgDec.String()
		cd.FundBackLoad = msgDec.String()
		cd.FundBackLoadTimeInterval = msgDec.String()
		cd.FundManagementFee = msgDec.String()
		cd.FundClosed = msgDec.Bool()
		cd.FundClosedForNewInvestors = msgDec.Bool()
		cd.FundClosedForNewMoney = msgDec.Bool()
		cd.FundNotifyAmount = msgDec.String()
		cd.FundMinimumInitialPurchase = msgDec.String()
		cd.FundSubsequentMinimumPurchase = msgDec.String()
		cd.FundBlueSkyStates = msgDec.String()
		cd.FundBlueSkyTerritories = msgDec.String()
		cd.FundDistributionPolicyIndicator = models.FundDistributionPolicyIndicator(msgDec.String())
		cd.FundAssetType = models.FundAsset(msgDec.String())
	}
	ineligibilityReasonListCount := int(msgDec.Int32())
	if ineligibilityReasonListCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative ineligibility reason count: %d", ineligibilityReasonListCount))
		return wrapDecoderErr("processContractDataMsg", msgDec.Err())
	}
	cd.IneligibilityReasonList = make([]models.IneligibilityReason, 0, ineligibilityReasonListCount)
	for i := 0; i < ineligibilityReasonListCount; i++ {
		ineligibilityReason := models.IneligibilityReason{}
		ineligibilityReason.ID = msgDec.String()
		ineligibilityReason.Description = msgDec.String()
		cd.IneligibilityReasonList = append(cd.IneligibilityReasonList, ineligibilityReason)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processContractDataMsg", msgDec.Err())
	}

	cd.OnContractLastTradeDateOrContractMonthUpdated(false)

	// Done
	return c.processContractDataCommon(reqID, cd)
}

func (c *Client) processContractDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ContractData{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processContractDataProtobuf", msgDec.Err())
	}
	if pb.ContractDetails == nil || pb.Contract == nil {
		return nil
	}
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	cd := models.NewContractDetailsFromProtobufDecoder(msgDec, pb.ContractDetails, pb.Contract)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processContractDataProtobuf", msgDec.Err())
	}

	cd.OnContractLastTradeDateOrContractMonthUpdated(false)

	// Done
	return c.processContractDataCommon(reqID, cd)
}

func (c *Client) processBondContractDataMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	cd := models.NewContractDetails()
	cd.Contract.Symbol = msgDec.String()
	cd.Contract.SecType = models.SecurityType(msgDec.String())
	cd.Cusip = msgDec.String()
	cd.Coupon = msgDec.Float()
	cd.Contract.LastTradeDateOrContractMonth = msgDec.String()
	cd.IssueDate = msgDec.String()
	cd.Ratings = msgDec.String()
	cd.BondType = msgDec.String()
	cd.CouponType = msgDec.String()
	cd.Convertible = msgDec.Bool()
	cd.Callable = msgDec.Bool()
	cd.Puttable = msgDec.Bool()
	cd.DescAppend = msgDec.String()
	cd.Contract.Exchange = msgDec.String()
	cd.Contract.Currency = msgDec.String()
	cd.MarketName = msgDec.String()
	cd.Contract.TradingClass = msgDec.String()
	cd.Contract.ConID = msgDec.Int32()
	cd.MinTick = msgDec.FloatMax()
	cd.OrderTypes = msgDec.String()
	cd.ValidExchanges = msgDec.String()
	cd.NextOptionDate = msgDec.String()
	cd.NextOptionType = msgDec.String()
	cd.NextOptionPartial = msgDec.Bool()
	cd.BondNotes = msgDec.String()
	cd.LongName = msgDec.String()
	cd.TimeZoneID = msgDec.String()
	cd.TradingHours = msgDec.String()
	cd.LiquidHours = msgDec.String()
	cd.EVRule = msgDec.String()
	cd.EVMultiplier = msgDec.Float()
	secIDListCount := int(msgDec.Int32())
	if secIDListCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative security id count: %d", secIDListCount))
		return wrapDecoderErr("processBondContractDataMsg", msgDec.Err())
	}
	cd.SecIDList = make(models.TagValueList, 0, secIDListCount)
	for i := 0; i < secIDListCount; i++ {
		tagValue := models.NewTagValue()
		tagValue.Tag = msgDec.String()
		tagValue.Value = msgDec.String()
		cd.SecIDList = append(cd.SecIDList, tagValue)
	}
	cd.AggGroup = msgDec.Int32()
	cd.MarketRuleIDs = msgDec.String()
	cd.MinSize = models.NewDecimalMaxFromMessageDecoder(msgDec)
	cd.SizeIncrement = models.NewDecimalMaxFromMessageDecoder(msgDec)
	cd.SuggestedSizeIncrement = models.NewDecimalMaxFromMessageDecoder(msgDec)
	cd.IneligibilityReasonList = make([]models.IneligibilityReason, 0)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processBondContractDataMsg", msgDec.Err())
	}

	cd.OnContractLastTradeDateOrContractMonthUpdated(true)

	// Done
	return c.processContractDataCommon(reqID, cd)
}

func (c *Client) processBondContractDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ContractData{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processBondContractDataProtobuf", msgDec.Err())
	}
	if pb.ContractDetails == nil || pb.Contract == nil {
		return nil
	}
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	cd := models.NewContractDetailsFromProtobufDecoder(msgDec, pb.ContractDetails, pb.Contract)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processBondContractDataProtobuf", msgDec.Err())
	}

	cd.OnContractLastTradeDateOrContractMonthUpdated(true)

	// Done
	return c.processContractDataCommon(reqID, cd)
}

func (c *Client) processContractDataCommon(reqID int32, cd *models.ContractDetails) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.ContractDetailsResponse)

		resp.ContractDetails = append(resp.ContractDetails, cd)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processMarketDepthMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	position := int(msgDec.Int32())
	if position < 0 {
		msgDec.SetErr(fmt.Errorf("invalid position: %d", position))
		return wrapDecoderErr("processMarketDepthMsg", msgDec.Err())
	}
	operation := int(msgDec.Int32())
	bidSide := msgDec.Int32() == 0
	price := msgDec.Float()
	size := models.NewDecimalFromMessageDecoder(msgDec)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthMsg", msgDec.Err())
	}

	// Done
	return c.processMarketDepthCommon(
		tickerID, position, models.MarketDepthDataOperation(operation), bidSide, price, size,
	)
}

func (c *Client) processMarketDepthProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.MarketDepth{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthProtobuf", msgDec.Err())
	}
	if pb.MarketDepthData == nil {
		return nil
	}
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	position := int(msgDec.Int32(pb.MarketDepthData.Position))
	if position < 0 {
		msgDec.SetErr(fmt.Errorf("invalid position: %d", position))
		return wrapDecoderErr("processMarketDepthProtobuf", msgDec.Err())
	}
	operation := int(msgDec.Int32(pb.MarketDepthData.Operation))
	side := msgDec.Int32(pb.MarketDepthData.Side)
	price := msgDec.Float(pb.MarketDepthData.Price)
	size := models.NewDecimalFromProtobufDecoder(msgDec, pb.MarketDepthData.Size)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthProtobuf", msgDec.Err())
	}

	// Done
	return c.processMarketDepthCommon(
		tickerID, position, models.MarketDepthDataOperation(operation), side == 0, price, size,
	)
}

func (c *Client) processMarketDepthCommon(
	tickerID int32, position int, operation models.MarketDepthDataOperation, bidSide bool, price float64, size models.Decimal,
) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*marketDepthDataStreamResponse)

		if position >= resp.resp.Book.Size {
			return false, nil // Ignore
		}

		// Notify
		data := models.MarketDepthData{
			Position:  position,
			Operation: operation,
			BidSide:   bidSide,
			Entry: models.DepthMarketBookEntry{
				Price: price,
				Size:  size,
			},
		}
		resp.enqueue(data)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processMarketDepthL2Msg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	position := int(msgDec.Int32())
	if position < 0 {
		msgDec.SetErr(fmt.Errorf("invalid position: %d", position))
		return wrapDecoderErr("processMarketDepthL2Msg", msgDec.Err())
	}
	_ = msgDec.String() // Market maker
	operation := int(msgDec.Int32())
	bidSide := msgDec.Int32() == 0
	price := msgDec.Float()
	size := models.NewDecimalFromMessageDecoder(msgDec)
	_ = msgDec.Bool() // Is Smart Depth
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthL2Msg", msgDec.Err())
	}

	// Done
	return c.processMarketDepthCommon(
		tickerID, position, models.MarketDepthDataOperation(operation), bidSide, price, size,
	)
}

func (c *Client) processMarketDepthL2Protobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.MarketDepthL2{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthL2Protobuf", msgDec.Err())
	}
	if pb.MarketDepthData == nil {
		return nil
	}
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	position := int(msgDec.Int32(pb.MarketDepthData.Position))
	if position < 0 {
		msgDec.SetErr(fmt.Errorf("invalid position: %d", position))
		return wrapDecoderErr("processMarketDepthL2Protobuf", msgDec.Err())
	}
	operation := int(msgDec.Int32(pb.MarketDepthData.Operation))
	side := msgDec.Int32(pb.MarketDepthData.Side)
	price := msgDec.Float(pb.MarketDepthData.Price)
	size := models.NewDecimalFromProtobufDecoder(msgDec, pb.MarketDepthData.Size)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDepthL2Protobuf", msgDec.Err())
	}

	// Done
	return c.processMarketDepthCommon(
		tickerID, position, models.MarketDepthDataOperation(operation), side == 0, price, size,
	)
}

func (c *Client) processManagedAccountsMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	accountsNames := msgDec.String()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processManagedAccountsMsg", msgDec.Err())
	}

	accounts := strings.Split(accountsNames, ",")

	// Done
	return c.processManagedAccountsCommon(accounts)
}

func (c *Client) processManagedAccountsProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ManagedAccounts{}
	msgDec.Unmarshal(&pb)
	accountsNames := msgDec.String(pb.AccountsList)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processManagedAccountsProtobuf", msgDec.Err())
	}

	accounts := strings.Split(accountsNames, ",")

	// Done
	return c.processManagedAccountsCommon(accounts)
}

func (c *Client) processManagedAccountsCommon(accounts []string) error {
	c.reqMgr.withAllRequestsWithoutID(common.REQ_MANAGED_ACCTS, func(_resp interface{}) error {
		resp := _resp.(*models.ManagedAccountsResponse)

		resp.Accounts = accounts

		// Done
		return nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalDataMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	barsCount := int(msgDec.Int32())
	if barsCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative bars count: %d", barsCount))
		return wrapDecoderErr("processHistoricalDataMsg", msgDec.Err())
	}
	bars := make([]models.HistoricalDataBar, 0, barsCount)
	for i := 0; i < barsCount; i++ {
		bar := models.NewHistoricalDataBar()
		// Epoch because the request of historical data has the date format equal to 2.
		bar.Date = msgDec.EpochTimestamp(false)
		bar.Open = msgDec.Float()
		bar.High = msgDec.Float()
		bar.Low = msgDec.Float()
		bar.Close = msgDec.Float()
		bar.Volume = models.NewDecimalMaxFromMessageDecoder(msgDec)
		bar.Wap = models.NewDecimalMaxFromMessageDecoder(msgDec)
		bar.Count = msgDec.Int32()

		bars = append(bars, bar)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataMsg", msgDec.Err())
	}

	// Done
	return c.processHistoricalDataCommon(reqID, bars)
}

func (c *Client) processHistoricalDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalData{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataProtobuf", msgDec.Err())
	}
	if pb.HistoricalDataBars == nil {
		return nil
	}
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	barsCount := len(pb.HistoricalDataBars)
	bars := make([]models.HistoricalDataBar, 0, barsCount)
	for i := 0; i < barsCount; i++ {
		bar := models.NewHistoricalDataBarFromProtobufDecoder(msgDec, pb.HistoricalDataBars[i])

		bars = append(bars, bar)
	}

	// Done
	return c.processHistoricalDataCommon(reqID, bars)
}

func (c *Client) processHistoricalDataCommon(reqID int32, bars []models.HistoricalDataBar) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*historicalDataStreamResponse)

		resp.appendBars(bars)

		// Done
		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalDataEndMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	msgDec.Skip()
	msgDec.Skip()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataEndMsg", msgDec.Err())
	}

	// Done
	return c.processHistoricalDataEndCommon(reqID)
}

func (c *Client) processHistoricalDataEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalDataEnd{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataEndProtobuf", msgDec.Err())
	}

	// Done
	return c.processHistoricalDataEndCommon(tickerID)
}

func (c *Client) processHistoricalDataEndCommon(tickerID int32) error {
	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*historicalDataStreamResponse)

		// Signal end of snapshot
		return !resp.keepUpToDate, nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalDataUpdateMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	bar := models.NewHistoricalDataBar()
	bar.Date = msgDec.EpochTimestamp(false)
	bar.Open = msgDec.Float()
	bar.Close = msgDec.Float()
	bar.High = msgDec.Float()
	bar.Low = msgDec.Float()
	bar.Wap = models.NewDecimalMaxFromMessageDecoder(msgDec)
	bar.Volume = models.NewDecimalMaxFromMessageDecoder(msgDec)
	bar.Count = msgDec.Int32()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataUpdateMsg", msgDec.Err())
	}

	return c.processHistoricalDataUpdateCommon(reqID, bar)
}

func (c *Client) processHistoricalDataUpdateProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalDataUpdate{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalDataUpdateProtobuf", msgDec.Err())
	}
	if pb.HistoricalDataBar == nil {
		return nil
	}

	bar := models.NewHistoricalDataBarFromProtobufDecoder(msgDec, pb.HistoricalDataBar)

	return c.processHistoricalDataUpdateCommon(reqID, bar)
}

func (c *Client) processHistoricalDataUpdateCommon(reqID int32, bar models.HistoricalDataBar) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*historicalDataStreamResponse)
		resp.applyUpdate(bar)

		return false, nil
	})

	return nil
}

func (c *Client) processRealTimeBarsMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	reqID := msgDec.RequestID(false)
	data := models.NewRealTimeBar()
	data.Time = msgDec.EpochTimestamp(false)
	data.Open = msgDec.Float()
	data.High = msgDec.Float()
	data.Low = msgDec.Float()
	data.Close = msgDec.Float()
	data.Volume = models.NewDecimalMaxFromMessageDecoder(msgDec)
	data.Wap = models.NewDecimalMaxFromMessageDecoder(msgDec)
	data.Count = msgDec.Int32()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processRealTimeBarsMsg", msgDec.Err())
	}

	return c.processRealTimeBarsCommon(reqID, data)
}

func (c *Client) processRealTimeBarsProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.RealTimeBarTick{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	data := models.NewRealTimeBar()
	data.Time = msgDec.EpochTimestamp(pb.Time, false)
	data.Open = msgDec.Float(pb.Open)
	data.High = msgDec.Float(pb.High)
	data.Low = msgDec.Float(pb.Low)
	data.Close = msgDec.Float(pb.Close)
	data.Volume = models.NewDecimalMaxFromProtobufDecoder(msgDec, pb.Volume)
	data.Wap = models.NewDecimalMaxFromProtobufDecoder(msgDec, pb.WAP)
	data.Count = msgDec.Int32(pb.Count)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processRealTimeBarsProtobuf", msgDec.Err())
	}

	return c.processRealTimeBarsCommon(reqID, data)
}

func (c *Client) processRealTimeBarsCommon(reqID int32, data models.RealTimeBar) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*realTimeBarsStreamResponse)
		if resp.barDurationSecs > 0 {
			data.EndTime = data.Time.Add(time.Duration(resp.barDurationSecs) * time.Second)
		}
		resp.enqueue(data)

		return false, nil
	})

	return nil
}

func (c *Client) processContractDataEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processContractDataEndMsg", msgDec.Err())
	}

	// Done
	return c.processContractDataEndCommon(reqID)
}

func (c *Client) processContractDataEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.ContractDataEnd{}
	msgDec.Unmarshal(&pb)
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processContractDataEndProtobuf", msgDec.Err())
	}

	// Done
	return c.processContractDataEndCommon(reqID)
}

func (c *Client) processContractDataEndCommon(reqID int32) error {
	c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
		// Done
		return true, nil
	})

	// Done
	return nil
}

func (c *Client) processTickSnapshotEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickSnapshotEndMsg", msgDec.Err())
	}

	// Done
	return c.processTickSnapshotEndCommon(tickerID)
}

func (c *Client) processTickSnapshotEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickSnapshotEnd{}
	msgDec.Unmarshal(&pb)
	// Gets the originating ticker ID
	tickerID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickSnapshotEndProtobuf", msgDec.Err())
	}

	// Done
	return c.processTickSnapshotEndCommon(tickerID)
}

func (c *Client) processTickSnapshotEndCommon(tickerID int32) error {
	c.reqMgr.withRequestWithID(tickerID, func(_ interface{}) (bool, error) {
		// Signal end of snapshot
		return true, nil
	})

	// Done
	return nil
}

func (c *Client) processMarketDataTypeMsg(msgDec *message.Decoder) error {
	msgDec.Skip() // version
	reqID := msgDec.RequestID(false)
	dataType := models.MarketDataType(msgDec.Int32())
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDataTypeMsg", msgDec.Err())
	}

	return c.processMarketDataTypeCommon(reqID, dataType)
}

func (c *Client) processMarketDataTypeProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.MarketDataType{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	dataType := models.MarketDataType(msgDec.Int32(pb.MarketDataType))
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketDataTypeProtobuf", msgDec.Err())
	}

	return c.processMarketDataTypeCommon(reqID, dataType)
}

func (c *Client) processMarketDataTypeCommon(reqID int32, dataType models.MarketDataType) error {
	if len(dataType.String()) == 0 {
		return fmt.Errorf("invalid market data type: %d", dataType)
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		if resp, ok := _resp.(*topMarketDataStreamResponse); ok {
			resp.setMarketDataType(dataType)
		}
		return false, nil
	})

	return nil
}

func (c *Client) processSymbolSamplesMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	contractDescriptionsCount := int(msgDec.Int32())
	if contractDescriptionsCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative contract descriptions count: %d", contractDescriptionsCount))
		return wrapDecoderErr("processSymbolSamplesMsg", msgDec.Err())
	}
	cds := make([]*models.ContractDescription, 0, contractDescriptionsCount)
	for i := 0; i < contractDescriptionsCount; i++ {
		cd := models.NewContractDescription()
		cd.Contract.ConID = msgDec.Int32()
		cd.Contract.Symbol = msgDec.String()
		cd.Contract.SecType = models.SecurityType(msgDec.String())
		cd.Contract.PrimaryExchange = msgDec.String()
		cd.Contract.Currency = msgDec.String()
		derivativeSecTypesCount := int(msgDec.Int32())
		if derivativeSecTypesCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative derivative security types count: %d", derivativeSecTypesCount))
			return wrapDecoderErr("processSymbolSamplesMsg", msgDec.Err())
		}
		cd.DerivativeSecTypes = make([]models.SecurityType, 0, derivativeSecTypesCount)
		for j := 0; j < derivativeSecTypesCount; j++ {
			derivativeSecType := models.SecurityType(msgDec.String())
			cd.DerivativeSecTypes = append(cd.DerivativeSecTypes, derivativeSecType)
		}
		cd.Contract.Description = msgDec.String()
		cd.Contract.IssuerID = msgDec.String()

		cds = append(cds, cd)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSymbolSamplesMsg", msgDec.Err())
	}

	// Done
	return c.processSymbolSamplesCommon(reqID, cds)
}

func (c *Client) processSymbolSamplesProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.SymbolSamples{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	cds := make([]*models.ContractDescription, len(pb.ContractDescriptions))
	for idx := range pb.ContractDescriptions {
		cds[idx] = models.NewContractDescriptionFromProtobufDecoder(msgDec, pb.ContractDescriptions[idx])
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSymbolSamplesProtobuf", msgDec.Err())
	}

	return c.processSymbolSamplesCommon(reqID, cds)
}

func (c *Client) processSymbolSamplesCommon(reqID int32, cds []*models.ContractDescription) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.MatchingSymbolsResponse)

		resp.ContractDescriptions = append(resp.ContractDescriptions, cds...)

		// Done
		return true, nil
	})

	// Done
	return nil
}

func (c *Client) processHeadTimestampMsg(msgDec *message.Decoder) error {
	// Gets the originating ticker ID
	reqID := msgDec.RequestID(false)
	ts := msgDec.EpochTimestamp(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHeadTimestampMsg", msgDec.Err())
	}

	// Done
	return c.processHeadTimestampCommon(reqID, ts)
}

func (c *Client) processHeadTimestampProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HeadTimestamp{}
	msgDec.Unmarshal(&pb)
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	ts := msgDec.EpochTimestampFromString(pb.HeadTimestamp, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHeadTimestampProtobuf", msgDec.Err())
	}

	// Done
	return c.processHeadTimestampCommon(reqID, ts)
}

func (c *Client) processHeadTimestampCommon(reqID int32, ts time.Time) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.HeadTimestampResponse)

		resp.Timestamp = ts

		// Done
		return true, nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalTicksMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	ticksCount := int(msgDec.Int32())
	if ticksCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
		return wrapDecoderErr("processHistoricalTicksMsg", msgDec.Err())
	}
	ticks := make([]models.HistoricalTick, ticksCount)
	for idx := 0; idx < ticksCount; idx++ {
		ht := models.NewHistoricalTick()
		ht.Time = msgDec.EpochTimestamp(false)
		msgDec.Skip()
		ht.Price = msgDec.Float()
		ht.Size = models.NewDecimalMaxFromMessageDecoder(msgDec)
		ticks[idx] = ht
	}
	done := msgDec.Bool()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksMsg", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalTicks{}
	msgDec.Unmarshal(&pb)
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	ticksCount := len(pb.HistoricalTicks)
	ticks := make([]models.HistoricalTick, ticksCount)
	for idx := range pb.HistoricalTicks {
		ticks[idx] = models.NewHistoricalTickFromProtobufDecoder(msgDec, pb.HistoricalTicks[idx])
	}
	done := msgDec.Bool(pb.IsDone)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksProtobuf", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksCommon(reqID int32, ticks []models.HistoricalTick, done bool) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		switch resp := _resp.(type) {
		case *models.HistoricalTicksResponse:
			resp.Ticks = append(resp.Ticks, ticks...)
			return done, nil

		case *tickByTickStreamResponse:
			for idx := range ticks {
				data := models.NewTickByTickMidPoint()
				data.Time = ticks[idx].Time
				data.MidPoint = ticks[idx].Price
				resp.enqueue(data)
			}
			return false, nil
		}

		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalTicksBidAskMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	ticksCount := int(msgDec.Int32())
	if ticksCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
		return wrapDecoderErr("processHistoricalTicksBidAskMsg", msgDec.Err())
	}
	ticks := make([]models.HistoricalTickBidAsk, ticksCount)
	for idx := 0; idx < ticksCount; idx++ {
		htba := models.NewHistoricalTickBidAsk()
		htba.Time = msgDec.EpochTimestamp(false)
		mask := msgDec.Int32()
		htba.TickAttribBidAsk = models.NewTickAttribBidAsk()
		htba.TickAttribBidAsk.AskPastHigh = mask&1 != 0
		htba.TickAttribBidAsk.BidPastLow = mask&2 != 0
		htba.PriceBid = msgDec.Float()
		htba.PriceAsk = msgDec.Float()
		htba.SizeBid = models.NewDecimalMaxFromMessageDecoder(msgDec)
		htba.SizeAsk = models.NewDecimalMaxFromMessageDecoder(msgDec)
		ticks[idx] = htba
	}
	done := msgDec.Bool()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksBidAskMsg", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksBidAskCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksBidAskProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalTicksBidAsk{}
	msgDec.Unmarshal(&pb)
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	ticksCount := len(pb.HistoricalTicksBidAsk)
	ticks := make([]models.HistoricalTickBidAsk, ticksCount)
	for idx := range pb.HistoricalTicksBidAsk {
		ticks[idx] = models.NewHistoricalTickBidAskFromProtobufDecoder(msgDec, pb.HistoricalTicksBidAsk[idx])
	}
	done := msgDec.Bool(pb.IsDone)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksBidAskProtobuf", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksBidAskCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksBidAskCommon(reqID int32, ticks []models.HistoricalTickBidAsk, done bool) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		switch resp := _resp.(type) {
		case *models.HistoricalTicksResponse:
			resp.TicksBidAsk = append(resp.TicksBidAsk, ticks...)
			return done, nil

		case *tickByTickStreamResponse:
			for idx := range ticks {
				data := models.NewTickByTickBidAsk()
				data.Time = ticks[idx].Time
				data.PriceBid = ticks[idx].PriceBid
				data.PriceAsk = ticks[idx].PriceAsk
				data.SizeBid = ticks[idx].SizeBid
				data.SizeAsk = ticks[idx].SizeAsk
				data.TickAttribBidAsk = ticks[idx].TickAttribBidAsk
				resp.enqueue(data)
			}
			return false, nil
		}

		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processHistoricalTicksLastMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	ticksCount := int(msgDec.Int32())
	if ticksCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
		return wrapDecoderErr("processHistoricalTicksLastMsg", msgDec.Err())
	}
	ticks := make([]models.HistoricalTickLast, ticksCount)
	for idx := 0; idx < ticksCount; idx++ {
		htl := models.NewHistoricalTickLast()
		htl.Time = msgDec.EpochTimestamp(false)
		mask := msgDec.Int32()
		htl.TickAttribLast = models.NewTickAttribLast()
		htl.TickAttribLast.PastLimit = mask&1 != 0
		htl.TickAttribLast.Unreported = mask&2 != 0
		htl.Price = msgDec.Float()
		htl.Size = models.NewDecimalMaxFromMessageDecoder(msgDec)
		htl.Exchange = msgDec.String()
		htl.SpecialConditions = msgDec.String()
		ticks[idx] = htl
	}
	done := msgDec.Bool()
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksLastMsg", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksLastCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksLastProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistoricalTicksLast{}
	msgDec.Unmarshal(&pb)
	// Gets the originating request ID
	reqID := msgDec.RequestID(pb.ReqId, false)
	ticksCount := len(pb.HistoricalTicksLast)
	ticks := make([]models.HistoricalTickLast, ticksCount)
	for idx := range pb.HistoricalTicksLast {
		ticks[idx] = models.NewHistoricalTickLastFromProtobufDecoder(msgDec, pb.HistoricalTicksLast[idx])
	}
	done := msgDec.Bool(pb.IsDone)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistoricalTicksLastProtobuf", msgDec.Err())
	}

	// Done
	return c.processHistoricalTicksLastCommon(reqID, ticks, done)
}

func (c *Client) processHistoricalTicksLastCommon(reqID int32, ticks []models.HistoricalTickLast, done bool) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		switch resp := _resp.(type) {
		case *models.HistoricalTicksResponse:
			resp.TicksLast = append(resp.TicksLast, ticks...)
			return done, nil

		case *tickByTickStreamResponse:
			for idx := range ticks {
				data := models.NewTickByTickAllLast(models.TickByTickTypeLast)
				data.Time = ticks[idx].Time
				data.Price = ticks[idx].Price
				data.Size = ticks[idx].Size
				data.TickAttribLast = ticks[idx].TickAttribLast
				data.Exchange = ticks[idx].Exchange
				data.SpecialConditions = ticks[idx].SpecialConditions
				resp.enqueue(data)
			}
			return false, nil
		}

		return false, nil
	})

	// Done
	return nil
}

func (c *Client) processSecurityDefinitionOptionalParameterMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	exchange := msgDec.String()
	underlyingConID := msgDec.Int32()
	tradingClass := msgDec.String()
	multiplier := msgDec.String()

	expCount := int(msgDec.Int32())
	if expCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative expirations count: %d", expCount))
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterMsg", msgDec.Err())
	}
	expirations := make([]string, expCount)
	for idx := 0; idx < expCount; idx++ {
		expirations[idx] = msgDec.String()
	}

	strikeCount := int(msgDec.Int32())
	if strikeCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative strikes count: %d", strikeCount))
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterMsg", msgDec.Err())
	}
	strikes := make([]float64, strikeCount)
	for idx := 0; idx < strikeCount; idx++ {
		strikes[idx] = msgDec.Float()
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterMsg", msgDec.Err())
	}

	data := models.SecurityDefinitionOptionParameter{
		Exchange:             exchange,
		UnderlyingContractID: underlyingConID,
		TradingClass:         tradingClass,
		Multiplier:           multiplier,
		Expirations:          expirations,
		Strikes:              strikes,
	}

	return c.processSecurityDefinitionOptionalParameterCommon(reqID, data)
}

func (c *Client) processSecurityDefinitionOptionalParameterProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.SecDefOptParameter{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	data := models.SecurityDefinitionOptionParameter{
		Exchange:             msgDec.String(pb.Exchange),
		UnderlyingContractID: msgDec.Int32(pb.UnderlyingConId),
		TradingClass:         msgDec.String(pb.TradingClass),
		Multiplier:           msgDec.String(pb.Multiplier),
		Expirations:          append(make([]string, 0, len(pb.Expirations)), pb.Expirations...),
		Strikes:              append(make([]float64, 0, len(pb.Strikes)), pb.Strikes...),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterProtobuf", msgDec.Err())
	}

	return c.processSecurityDefinitionOptionalParameterCommon(reqID, data)
}

func (c *Client) processSecurityDefinitionOptionalParameterCommon(
	reqID int32, data models.SecurityDefinitionOptionParameter,
) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.SecurityDefinitionOptionParametersResponse)
		resp.Parameters = append(resp.Parameters, data)

		return false, nil
	})

	return nil
}

func (c *Client) processSecurityDefinitionOptionalParameterEndMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterEndMsg", msgDec.Err())
	}

	return c.processSecurityDefinitionOptionalParameterEndCommon(reqID)
}

func (c *Client) processSecurityDefinitionOptionalParameterEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.SecDefOptParameterEnd{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSecurityDefinitionOptionalParameterEndProtobuf", msgDec.Err())
	}

	return c.processSecurityDefinitionOptionalParameterEndCommon(reqID)
}

func (c *Client) processSecurityDefinitionOptionalParameterEndCommon(reqID int32) error {
	c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
		return true, nil
	})

	return nil
}

func (c *Client) processTickByTickDataMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	tickType := models.TickByTickType(msgDec.Int32())
	ts := msgDec.EpochTimestamp(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickByTickDataMsg", msgDec.Err())
	}

	switch tickType {
	case models.TickByTickTypeLast, models.TickByTickTypeAllLast:
		data := models.NewTickByTickAllLast(tickType)
		data.Time = ts
		data.Price = msgDec.Float()
		data.Size = models.NewDecimalMaxFromMessageDecoder(msgDec)
		mask := msgDec.Int32()
		data.TickAttribLast = models.NewTickAttribLast()
		data.TickAttribLast.PastLimit = mask&0x1 != 0
		data.TickAttribLast.Unreported = mask&0x2 != 0
		data.Exchange = msgDec.String()
		data.SpecialConditions = msgDec.String()
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataMsg", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)

	case models.TickByTickTypeBidAsk:
		data := models.NewTickByTickBidAsk()
		data.Time = ts
		data.PriceBid = msgDec.Float()
		data.PriceAsk = msgDec.Float()
		data.SizeBid = models.NewDecimalMaxFromMessageDecoder(msgDec)
		data.SizeAsk = models.NewDecimalMaxFromMessageDecoder(msgDec)
		mask := msgDec.Int32()
		data.TickAttribBidAsk = models.NewTickAttribBidAsk()
		data.TickAttribBidAsk.BidPastLow = mask&0x1 != 0
		data.TickAttribBidAsk.AskPastHigh = mask&0x2 != 0
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataMsg", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)

	case models.TickByTickTypeMidPoint:
		data := models.NewTickByTickMidPoint()
		data.Time = ts
		data.MidPoint = msgDec.Float()
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataMsg", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)
	}

	return nil
}

func (c *Client) processTickByTickDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickByTickData{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	tickType := models.TickByTickType(msgDec.Int32(pb.TickType))
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickByTickDataProtobuf", msgDec.Err())
	}

	switch tickType {
	case models.TickByTickTypeLast, models.TickByTickTypeAllLast:
		data := models.NewTickByTickAllLast(tickType)
		htl := models.NewHistoricalTickLastFromProtobufDecoder(msgDec, pb.GetHistoricalTickLast())
		data.Time = htl.Time
		data.Price = htl.Price
		data.Size = htl.Size
		data.TickAttribLast = htl.TickAttribLast
		data.Exchange = htl.Exchange
		data.SpecialConditions = htl.SpecialConditions
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataProtobuf", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)

	case models.TickByTickTypeBidAsk:
		data := models.NewTickByTickBidAsk()
		htba := models.NewHistoricalTickBidAskFromProtobufDecoder(msgDec, pb.GetHistoricalTickBidAsk())
		data.Time = htba.Time
		data.PriceBid = htba.PriceBid
		data.PriceAsk = htba.PriceAsk
		data.SizeBid = htba.SizeBid
		data.SizeAsk = htba.SizeAsk
		data.TickAttribBidAsk = htba.TickAttribBidAsk
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataProtobuf", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)

	case models.TickByTickTypeMidPoint:
		data := models.NewTickByTickMidPoint()
		ht := models.NewHistoricalTickFromProtobufDecoder(msgDec, pb.GetHistoricalTickMidPoint())
		data.Time = ht.Time
		data.MidPoint = ht.Price
		if msgDec.Err() != nil {
			return wrapDecoderErr("processTickByTickDataProtobuf", msgDec.Err())
		}
		return c.processTickByTickDataCommon(reqID, data)
	}

	return nil
}

func (c *Client) processTickByTickDataCommon(reqID int32, data models.TickByTickData) error {
	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*tickByTickStreamResponse)
		resp.enqueue(data)

		return false, nil
	})

	return nil
}

// -----------------------------------------------------------------------------

type marketRuleEventHandler interface {
	MarketRule(marketRuleID int64, priceIncrements []models.PriceIncrement)
}

func (c *Client) processMarketRuleMsg(msgDec *message.Decoder) error {
	marketRuleID := msgDec.Int64()
	priceIncrementsCount := int(msgDec.Int64())
	if priceIncrementsCount < 0 {
		msgDec.SetErr(fmt.Errorf("negative price increments count: %d", priceIncrementsCount))
		return wrapDecoderErr("processMarketRuleMsg", msgDec.Err())
	}

	priceIncrements := make([]models.PriceIncrement, priceIncrementsCount)
	for idx := 0; idx < priceIncrementsCount; idx++ {
		priceIncrements[idx] = models.PriceIncrement{
			LowEdge:   msgDec.Float(),
			Increment: msgDec.Float(),
		}
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketRuleMsg", msgDec.Err())
	}

	return c.processMarketRuleCommon(marketRuleID, priceIncrements)
}

func (c *Client) processMarketRuleProtobuf(msgDec *protofmt.Decoder) error {
	marketRuleID, priceIncrements, err := decodeMarketRuleProtoMessage(msgDec)
	if err != nil {
		return wrapDecoderErr("processMarketRuleProtobuf", err)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMarketRuleProtobuf", msgDec.Err())
	}

	return c.processMarketRuleCommon(marketRuleID, priceIncrements)
}

func decodeMarketRuleProtoMessage(msgDec *protofmt.Decoder) (int64, []models.PriceIncrement, error) {
	var marketRuleID int64
	var priceIncrements []models.PriceIncrement
	msg := msgDecRawBytes(msgDec)

	for len(msg) > 0 {
		num, typ, n := protowire.ConsumeTag(msg)
		if n < 0 {
			return 0, nil, protowire.ParseError(n)
		}
		msg = msg[n:]

		switch num {
		case 1:
			if typ != protowire.VarintType {
				return 0, nil, fmt.Errorf("invalid market rule id wire type: %d", typ)
			}
			value, m := protowire.ConsumeVarint(msg)
			if m < 0 {
				return 0, nil, protowire.ParseError(m)
			}
			msg = msg[m:]
			marketRuleID = int64(value)
		case 2:
			if typ != protowire.BytesType {
				return 0, nil, fmt.Errorf("invalid price increments wire type: %d", typ)
			}
			value, m := protowire.ConsumeBytes(msg)
			if m < 0 {
				return 0, nil, protowire.ParseError(m)
			}
			msg = msg[m:]
			priceIncrement, err := decodeMarketRulePriceIncrementProto(value)
			if err != nil {
				return 0, nil, err
			}
			priceIncrements = append(priceIncrements, priceIncrement)
		default:
			m := protowire.ConsumeFieldValue(num, typ, msg)
			if m < 0 {
				return 0, nil, protowire.ParseError(m)
			}
			msg = msg[m:]
		}
	}

	if marketRuleID <= 0 {
		return 0, nil, fmt.Errorf("received invalid market rule id: %d", marketRuleID)
	}

	return marketRuleID, priceIncrements, nil
}

func decodeMarketRulePriceIncrementProto(msg []byte) (models.PriceIncrement, error) {
	priceIncrement := models.PriceIncrement{}

	for len(msg) > 0 {
		num, typ, n := protowire.ConsumeTag(msg)
		if n < 0 {
			return models.PriceIncrement{}, protowire.ParseError(n)
		}
		msg = msg[n:]

		switch num {
		case 1:
			if typ != protowire.Fixed64Type {
				return models.PriceIncrement{}, fmt.Errorf("invalid low edge wire type: %d", typ)
			}
			value, m := protowire.ConsumeFixed64(msg)
			if m < 0 {
				return models.PriceIncrement{}, protowire.ParseError(m)
			}
			msg = msg[m:]
			priceIncrement.LowEdge = math.Float64frombits(value)
		case 2:
			if typ != protowire.Fixed64Type {
				return models.PriceIncrement{}, fmt.Errorf("invalid increment wire type: %d", typ)
			}
			value, m := protowire.ConsumeFixed64(msg)
			if m < 0 {
				return models.PriceIncrement{}, protowire.ParseError(m)
			}
			msg = msg[m:]
			priceIncrement.Increment = math.Float64frombits(value)
		default:
			m := protowire.ConsumeFieldValue(num, typ, msg)
			if m < 0 {
				return models.PriceIncrement{}, protowire.ParseError(m)
			}
			msg = msg[m:]
		}
	}

	return priceIncrement, nil
}

func msgDecRawBytes(msgDec *protofmt.Decoder) []byte {
	return msgDec.RawBytes()
}

func (c *Client) processMarketRuleCommon(marketRuleID int64, priceIncrements []models.PriceIncrement) error {
	c.reqMgr.withAllMatchingRequestsWithoutID(int(common.MARKET_RULE), func(req *Request) bool {
		return req.marketRuleID == marketRuleID
	}, func(_resp interface{}) error {
		resp := _resp.(*models.MarketRuleResponse)
		resp.MarketRuleID = marketRuleID
		resp.PriceIncrements = append(resp.PriceIncrements[:0], priceIncrements...)
		return nil
	})

	if handler, ok := c.eventsHandler.(marketRuleEventHandler); ok {
		handler.MarketRule(marketRuleID, priceIncrements)
	}

	return nil
}

// -----------------------------------------------------------------------------

func (c *Client) processCurrentTimeInMillisMsg(msgDec *message.Decoder) error {
	ts := msgDec.EpochTimestamp(true)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCurrentTimeInMillisMsg", msgDec.Err())
	}

	return c.processCurrentTimeInMillisCommon(ts)
}

func (c *Client) processCurrentTimeInMillisCommon(ts time.Time) error {
	c.reqMgr.withAllRequestsWithoutID(common.REQ_CURRENT_TIME_IN_MILLIS, func(_resp interface{}) error {
		resp := _resp.(*models.CurrentTimeResponse)

		resp.CurrentTime = ts

		// Done
		return nil
	})

	// Done
	return nil
}
