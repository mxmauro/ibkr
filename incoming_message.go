package ibkr

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

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
		err = errors.New("received invalid message")
	}
	return
}

func wrapDecoderErr(decoderName string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s decode failed: %w", decoderName, err)
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

			/*
				case RECEIVE_FA:
					return c.processReceiveFaMsg(msgDec)
			*/

		case common.HISTORICAL_DATA:
			return c.processHistoricalDataMsg(msgDec)

			/*
				case SCANNER_DATA:
					return c.processScannerDataMsg(msgDec)
				case SCANNER_PARAMETERS:
					return c.processScannerParametersMsg(msgDec)
			*/

		case common.CURRENT_TIME:
			return nil // Ignore this message. We use the one having milliseconds
		case common.REAL_TIME_BARS:
			return c.processRealTimeBarsMsg(msgDec)
			/*
				case FUNDAMENTAL_DATA:
					return c.processFundamentalDataMsg(msgDec)
			*/

		case common.CONTRACT_DATA_END:
			return c.processContractDataEndMsg(msgDec)
		case common.OPEN_ORDER_END:
			return c.processOpenOrdersEndMsg(msgDec)
		case common.TICK_SNAPSHOT_END:
			return c.processTickSnapshotEndMsg(msgDec)
		case common.MARKET_DATA_TYPE:
			return c.processMarketDataTypeMsg(msgDec)
		case common.COMMISSION_AND_FEES_REPORT:
			return c.processCommissionAndFeesReportMsg(msgDec)
			/*
				case POSITION_DATA:
					return c.processPositionDataMsg(msgDec)
				case POSITION_END:
					return c.processPositionEndMsg(msgDec)
				case ACCOUNT_SUMMARY:
					return c.processAccountSummaryMsg(msgDec)
				case ACCOUNT_SUMMARY_END:
					return c.processAccountSummaryEndMsg(msgDec)
				case VERIFY_MESSAGE_API:
					return c.processVerifyMessageApiMsg(msgDec)
				case VERIFY_COMPLETED:
					return c.processVerifyCompletedMsg(msgDec)
				case DISPLAY_GROUP_LIST:
					return c.processDisplayGroupListMsg(msgDec)
				case DISPLAY_GROUP_UPDATED:
					return c.processDisplayGroupUpdatedMsg(msgDec)
				case VERIFY_AND_AUTH_MESSAGE_API:
					return c.processVerifyAndAuthMessageApiMsg(msgDec)
				case VERIFY_AND_AUTH_COMPLETED:
					return c.processVerifyAndAuthCompletedMsg(msgDec)
				case POSITION_MULTI:
					return c.processPositionMultiMsg(msgDec)
				case POSITION_MULTI_END:
					return c.processPositionMultiEndMsg(msgDec)
				case ACCOUNT_UPDATE_MULTI:
					return c.processAccountUpdateMultiMsg(msgDec)
				case ACCOUNT_UPDATE_MULTI_END:
					return c.processAccountUpdateMultiEndMsg(msgDec)
				case SOFT_DOLLAR_TIERS:
					return c.processSoftDollarTiersMsg(msgDec)
				case FAMILY_CODES:
					return c.processFamilyCodesMsg(msgDec)
				case SMART_COMPONENTS:
					return c.processSmartComponentsMsg(msgDec)
			*/

		case common.TICK_REQ_PARAMS:
			return nil // Ignore this message. We don't use the tick request parameters.
		case common.SYMBOL_SAMPLES:
			return c.processSymbolSamplesMsg(msgDec)
			/*
				case MKT_DEPTH_EXCHANGES:
					return c.processMktDepthExchangesMsg(msgDec)
				case TICK_NEWS:
					return c.processTickNewsMsg(msgDec)
				case NEWS_PROVIDERS:
					return c.processNewsProvidersMsg(msgDec)
				case NEWS_ARTICLE:
					return c.processNewsArticleMsg(msgDec)
				case HISTORICAL_NEWS:
					return c.processHistoricalNewsMsg(msgDec)
				case HISTORICAL_NEWS_END:
					return c.processHistoricalNewsEndMsg(msgDec)
			*/
		case common.HEAD_TIMESTAMP:
			return c.processHeadTimestampMsg(msgDec)
			/*
				case HISTOGRAM_DATA:
					return c.processHistogramDataMsg(msgDec)
			*/
		case common.HISTORICAL_DATA_UPDATE:
			return c.processHistoricalDataUpdateMsg(msgDec)
			/*
				case REROUTE_MKT_DATA_REQ:
					return c.processRerouteMktDataReqMsg(msgDec)
				case REROUTE_MKT_DEPTH_REQ:
					return c.processRerouteMktDepthReqMsg(msgDec)
				case MARKET_RULE:
					return c.processMarketRuleMsg(msgDec)
				case PNL:
					return c.processPnLMsg(msgDec)
				case PNL_SINGLE:
					return c.processPnLSingleMsg(msgDec)
			*/
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
			/*
				case REPLACE_FA_END:
					return c.processReplaceFAEndMsg(msgDec)
				case WSH_META_DATA:
					return c.processWshMetaData(msgDec)
				case WSH_EVENT_DATA:
					return c.processWshEventData(msgDec)
				case HISTORICAL_SCHEDULE:
					return c.processHistoricalSchedule(msgDec)
				case USER_INFO:
					return c.processUserInfo(msgDec)
			*/
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

/*
func (c *Client) processOrderStatusMsg(msgDec *utils.Decoder) error {

	if d.serverVersion < MIN_SERVER_VER_MARKET_CAP_PRICE {
		msgDec.decode() // version
	}

	orderID := msgDec.decodeInt64()
	status := msgDec.decodeString()
	filled := msgDec.decodeDecimal()
	remaining := msgDec.decodeDecimal()
	avgFilledPrice := msgDec.decodeFloat64()

	permID := msgDec.decodeInt64()          // ver 2 field
	parentID := msgDec.decodeInt64()        // ver 3 field
	lastFillPrice := msgDec.decodeFloat64() // ver 4 field
	clientID := msgDec.decodeInt64()        // ver 5 field
	whyHeld := msgDec.decodeString()        // ver 6 field

	mktCapPrice := 0.0
	if d.serverVersion >= MIN_SERVER_VER_MARKET_CAP_PRICE {
		mktCapPrice = msgDec.decodeFloat64()
	}

	d.wrapper.OrderStatus(orderID, status, filled, remaining, avgFilledPrice, permID, parentID, lastFillPrice, clientID, whyHeld, mktCapPrice)
}
*/

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
		return errors.New("invalid request ID")
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
		return errors.New("invalid request ID")
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

	// Process the response
	if reqID > 0 {
		if c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
			// Done
			return false, newRequestError(ts, code, errMsg, advancedOrderRejectJson)
		}) {
			return nil
		}

		c.orderStore.RecordOrderError(reqID, ts, code, errMsg, advancedOrderRejectJson)
	}

	// Notify
	if c.eventsHandler != nil {
		c.eventsHandler.Error(ts, code, errMsg, advancedOrderRejectJson)
	}

	// Done
	return nil
}

/*
func (c *Client) processOpenOrderMsg(msgDec *utils.Decoder) error {

		order := NewOrder()
		contract := NewContract()
		orderState := NewOrderState()

		version := d.serverVersion
		if d.serverVersion < MIN_SERVER_VER_ORDER_CONTAINER {
			version = Version(msgDec.decodeInt64())
		}

		orderDecoder := &OrderDecoder{order, contract, orderState, version, d.serverVersion}

		// read orderID
		orderDecoder.decodeOrderId(msgDec)

		// read contract fields
		orderDecoder.decodeContractFields(msgDec)

		// read order fields
		orderDecoder.decodeAction(msgDec)
		orderDecoder.decodeTotalQuantity(msgDec)
		orderDecoder.decodeOrderType(msgDec)
		orderDecoder.decodeLmtPrice(msgDec)
		orderDecoder.decodeAuxPrice(msgDec)
		orderDecoder.decodeTIF(msgDec)
		orderDecoder.decodeOcaGroup(msgDec)
		orderDecoder.decodeAccount(msgDec)
		orderDecoder.decodeOpenClose(msgDec)
		orderDecoder.decodeOrigin(msgDec)
		orderDecoder.decodeOrderRef(msgDec)
		orderDecoder.decodeClientId(msgDec)
		orderDecoder.decodePermId(msgDec)
		orderDecoder.decodeOutsideRth(msgDec)
		orderDecoder.decodeHidden(msgDec)
		orderDecoder.decodeDiscretionaryAmount(msgDec)
		orderDecoder.decodeGoodAfterTime(msgDec)
		orderDecoder.skipSharesAllocation(msgDec)
		orderDecoder.decodeFAParams(msgDec)
		orderDecoder.decodeModelCode(msgDec)
		orderDecoder.decodeGoodTillDate(msgDec)
		orderDecoder.decodeRule80A(msgDec)
		orderDecoder.decodePercentOffset(msgDec)
		orderDecoder.decodeSettlingFirm(msgDec)
		orderDecoder.decodeShortSaleParams(msgDec)
		orderDecoder.decodeAuctionStrategy(msgDec)
		orderDecoder.decodeBoxOrderParams(msgDec)
		orderDecoder.decodePegToStkOrVolOrderParams(msgDec)
		orderDecoder.decodeDisplaySize(msgDec)
		orderDecoder.decodeBlockOrder(msgDec)
		orderDecoder.decodeSweepToFill(msgDec)
		orderDecoder.decodeAllOrNone(msgDec)
		orderDecoder.decodeMinQty(msgDec)
		orderDecoder.decodeOcaType(msgDec)
		orderDecoder.skipETradeOnly(msgDec)
		orderDecoder.skipFirmQuoteOnly(msgDec)
		orderDecoder.skipNbboPriceCap(msgDec)
		orderDecoder.decodeParentId(msgDec)
		orderDecoder.decodeTriggerMethod(msgDec)
		orderDecoder.decodeVolOrderParams(msgDec, true)
		orderDecoder.decodeTrailParams(msgDec)
		orderDecoder.decodeBasisPoints(msgDec)
		orderDecoder.decodeComboLegs(msgDec)
		orderDecoder.decodeSmartComboRoutingParams(msgDec)
		orderDecoder.decodeScaleOrderParams(msgDec)
		orderDecoder.decodeHedgeParams(msgDec)
		orderDecoder.decodeOptOutSmartRouting(msgDec)
		orderDecoder.decodeClearingParams(msgDec)
		orderDecoder.decodeNotHeld(msgDec)
		orderDecoder.decodeDeltaNeutral(msgDec)
		orderDecoder.decodeAlgoParams(msgDec)
		orderDecoder.decodeSolicited(msgDec)
		orderDecoder.decodeWhatIfInfoAndCommissionAndFees(msgDec)
		orderDecoder.decodeVolRandomizeFlags(msgDec)
		orderDecoder.decodePegBenchParams(msgDec)
		orderDecoder.decodeConditions(msgDec)
		orderDecoder.decodeAdjustedOrderParams(msgDec)
		orderDecoder.decodeSoftDollarTier(msgDec)
		orderDecoder.decodeCashQty(msgDec)
		orderDecoder.decodeDontUseAutoPriceForHedge(msgDec)
		orderDecoder.decodeIsOmsContainer(msgDec)
		orderDecoder.decodeDiscretionaryUpToLimitPrice(msgDec)
		orderDecoder.decodeUsePriceMgmtAlgo(msgDec)
		orderDecoder.decodeDuration(msgDec)
		orderDecoder.decodePostToAts(msgDec)
		orderDecoder.decodeAutoCancelParent(msgDec, MIN_SERVER_VER_AUTO_CANCEL_PARENT)
		orderDecoder.decodePegBestPegMidOrderAttributes(msgDec)
		orderDecoder.decodeCustomerAccount(msgDec)
		orderDecoder.decodeProfessionalCustomer(msgDec)
		orderDecoder.decodeBondAccruedInterest(msgDec)
		orderDecoder.decodeIncludeOvernight(msgDec)
		orderDecoder.decodeCMETaggingFields(msgDec)
		orderDecoder.decodeSubmitter(msgDec)
		orderDecoder.decodeImbalanceOnly(msgDec, MIN_SERVER_VER_IMBALANCE_ONLY)

		d.wrapper.OpenOrder(order.OrderID, contract, order, orderState)
	}

func (c *Client) processAcctValueMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		tag := msgDec.decodeString()
		val := msgDec.decodeString()
		currency := msgDec.decodeString()
		accountName := msgDec.decodeString()

		d.wrapper.UpdateAccountValue(tag, val, currency, accountName)
	}

func (c *Client) processPortfolioValueMsg(msgDec *utils.Decoder) error {

	version := msgDec.decodeInt64()

	// read contract fields
	contract := NewContract()
	contract.ConID = msgDec.decodeInt64() // ver 6 field
	contract.Symbol = msgDec.decodeString()
	contract.SecType = msgDec.decodeString()
	contract.LastTradeDateOrContractMonth = msgDec.decodeString()
	contract.Strike = msgDec.decodeFloat64()
	contract.Right = msgDec.decodeString()

	if version >= 7 {
		contract.Multiplier = msgDec.decodeString()
		contract.PrimaryExchange = msgDec.decodeString()
	}

	contract.Currency = msgDec.decodeString()
	contract.LocalSymbol = msgDec.decodeString() // ver 2 field
	if version >= 8 {
		contract.TradingClass = msgDec.decodeString()
	}
	position := msgDec.decodeDecimal()

	marketPrice := msgDec.decodeFloat64()
	marketValue := msgDec.decodeFloat64()
	averageCost := msgDec.decodeFloat64()   // ver 3 field
	unrealizedPNL := msgDec.decodeFloat64() // ver 3 field
	realizedPNL := msgDec.decodeFloat64()   // ver 3 field

	accountName := msgDec.decodeString() // ver 4 field

	if version == 6 && d.serverVersion == 39 {
		contract.PrimaryExchange = msgDec.decodeString()
	}

	d.wrapper.UpdatePortfolio(contract, position, marketPrice, marketValue, averageCost, unrealizedPNL, realizedPNL, accountName)

}

func (c *Client) processAcctUpdateTimeMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		timeStamp := msgDec.decodeString()

		d.wrapper.UpdateAccountTime(timeStamp)
	}
*/
func (c *Client) processContractDataMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	cd := models.NewContractDetails()
	cd.Contract.Symbol = msgDec.String()
	cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String())
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
	cd.UnderSecType = models.NewSecurityTypeFromString(msgDec.String())
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
		cd.FundDistributionPolicyIndicator = models.NewFundDistributionPolicyIndicatorFromString(msgDec.String())
		cd.FundAssetType = models.NewFundAssetFromString(msgDec.String())
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

/*
		private void processContractDataMsgProtoBuf() throws IOException {
		    byte[] byteArray = readByteArray();

		    ContractDataProto.ContractData contractDataProto = ContractDataProto.ContractData.parseFrom(byteArray);
		    m_EWrapper.contractDataProtoBuf(contractDataProto);

		    int reqId = contractDataProto.hasReqId() ? contractDataProto.getReqId() : EClientErrors.NO_VALID_ID;

		    if (!contractDataProto.hasContract() || !contractDataProto.hasContractDetails()) {
		        return;
		    }
		    // set contract details fields
		    ContractDetails contractDetails = EDecoderUtils.decodeContractDetails(contractDataProto.getContract(), contractDataProto.getContractDetails(), false);

		    m_EWrapper.contractDetails(reqId, contractDetails);
		}

	    private void processBondContractDataMsgProtoBuf() throws IOException {
	        byte[] byteArray = readByteArray();

	        ContractDataProto.ContractData contractDataProto = ContractDataProto.ContractData.parseFrom(byteArray);
	        m_EWrapper.bondContractDataProtoBuf(contractDataProto);

	        int reqId = contractDataProto.hasReqId() ? contractDataProto.getReqId() : EClientErrors.NO_VALID_ID;

	        if (!contractDataProto.hasContract() || !contractDataProto.hasContractDetails()) {
	            return;
	        }
	        // set contract details fields
	        ContractDetails contractDetails = EDecoderUtils.decodeContractDetails(contractDataProto.getContract(), contractDataProto.getContractDetails(), true);

	        m_EWrapper.bondContractDetails(reqId, contractDetails);
	    }
*/
func (c *Client) processBondContractDataMsg(msgDec *message.Decoder) error {
	// Gets the originating request ID
	reqID := msgDec.RequestID(false)
	cd := models.NewContractDetails()
	cd.Contract.Symbol = msgDec.String()
	cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String())
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

/*
func (c *Client) processExecutionDetailsMsg(msgDec *utils.Decoder) error {

	version := d.serverVersion
	if d.serverVersion < MIN_SERVER_VER_LAST_LIQUIDITY {
		version = Version(msgDec.decodeInt64())
	}

	var reqID int64 = -1
	if version >= 7 {
		reqID = msgDec.decodeInt64()
	}

	orderID := msgDec.decodeInt64()

	// decode contact fields
	contract := NewContract()
	contract.ConID = msgDec.decodeInt64()
	contract.Symbol = msgDec.decodeString()
	contract.SecType = msgDec.decodeString()
	contract.LastTradeDateOrContractMonth = msgDec.decodeString()
	contract.Strike = msgDec.decodeFloat64()
	contract.Right = msgDec.decodeString()

	if version >= 9 {
		contract.Multiplier = msgDec.decodeString()
	}

	contract.Exchange = msgDec.decodeString()
	contract.Currency = msgDec.decodeString()
	contract.LocalSymbol = msgDec.decodeString()

	if version >= 10 {
		contract.TradingClass = msgDec.decodeString()
	}

	// read execution fields
	execution := NewExecution()
	execution.OrderID = orderID
	execution.ExecID = msgDec.decodeString()
	execution.Time = msgDec.decodeString()
	execution.AcctNumber = msgDec.decodeString()
	execution.Exchange = msgDec.decodeString()
	execution.Side = msgDec.decodeString()
	execution.Shares = msgDec.decodeDecimal()
	execution.Price = msgDec.decodeFloat64()
	execution.PermID = msgDec.decodeInt64()
	execution.ClientID = msgDec.decodeInt64()
	execution.Liquidation = msgDec.decodeInt64()

	if version >= 6 {
		execution.CumQty = msgDec.decodeDecimal()
		execution.AvgPrice = msgDec.decodeFloat64()
	}

	if version >= 8 {
		execution.OrderRef = msgDec.decodeString()
	}

	if version >= 9 {
		execution.EVRule = msgDec.decodeString()
		execution.EVMultiplier = msgDec.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		execution.ModelCode = msgDec.decodeString()
	}

	if d.serverVersion >= MIN_SERVER_VER_LAST_LIQUIDITY {
		execution.LastLiquidity = msgDec.decodeInt64()
	}

	if d.serverVersion >= MIN_SERVER_VER_PENDING_PRICE_REVISION {
		execution.PendingPriceRevision = msgDec.decodeBool()
	}

	if d.serverVersion >= MIN_SERVER_VER_SUBMITTER {
		execution.Submitter = msgDec.decodeString()
	}

	d.wrapper.ExecDetails(reqID, contract, execution)
}

func (c *Client) processExecutionDetailsMsgProtoBuf(msgDec *utils.Decoder) error {

	var executionDetailsProto protobuf.ExecutionDetails
	err := proto.Unmarshal(msgDec.Bytes(), &executionDetailsProto)
	if err != nil {
		log.Panic().Err(err).Msg("processExecutionDetailsMsgProtoBuf unmarshal error")
	}

	var reqID int64 = int64(executionDetailsProto.GetReqId())

	var contract *Contract
	if executionDetailsProto.Contract != nil {
		contract = decodeContract(executionDetailsProto.GetContract())
	}

	var execution *Execution
	if executionDetailsProto.Execution != nil {
		execution = decodeExecution(executionDetailsProto.GetExecution())
	}

	d.wrapper.ExecDetails(reqID, contract, execution)
}
*/

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

/*
func (c *Client) processNewsBulletinsMsg(msgDec *utils.Decoder) error {

	msgDec.decode() // version

	msgID := msgDec.decodeInt64()
	msgType := msgDec.decodeInt64()
	newsMessage := msgDec.decodeString()
	originExch := msgDec.decodeString()

	d.wrapper.UpdateNewsBulletin(msgID, msgType, newsMessage, originExch)
}
*/

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
	c.reqMgr.withRequestWithoutID(common.REQ_MANAGED_ACCTS, func(_resp interface{}) error {
		resp := _resp.(*models.ManagedAccountsResponse)

		resp.Accounts = accounts

		// Done
		return nil
	})

	// Done
	return nil
}

/*
func (c *Client) processReceiveFaMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		faDataType := FaDataType(msgDec.decodeInt64())
		cxml := msgDec.decodeString()

		d.wrapper.ReceiveFA(faDataType, cxml)
	}
*/

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

/*
func (c *Client) processScannerDataMsg(msgDec *utils.Decoder) error {

		msgDec.Skip() // version

		reqID := msgDec.decodeInt64()

		numberOfElements := msgDec.decodeInt64()

		var i int64
		for i = 0; i < numberOfElements; i++ {

			contractDetails := NewContractDetails()

			rank := msgDec.decodeInt64()
			contractDetails.Contract.ConID = msgDec.decodeInt64()
			contractDetails.Contract.Symbol = msgDec.decodeString()
			contractDetails.Contract.SecType = msgDec.decodeString()
			contractDetails.Contract.LastTradeDateOrContractMonth = msgDec.decodeString()
			contractDetails.Contract.Strike = msgDec.decodeFloat64()
			contractDetails.Contract.Right = msgDec.decodeString()
			contractDetails.Contract.Exchange = msgDec.decodeString()
			contractDetails.Contract.Currency = msgDec.decodeString()
			contractDetails.Contract.LocalSymbol = msgDec.decodeString()
			contractDetails.MarketName = msgDec.decodeString()
			contractDetails.Contract.TradingClass = msgDec.decodeString()
			distance := msgDec.decodeString()
			benchmark := msgDec.decodeString()
			projection := msgDec.decodeString()
			legsStr := msgDec.decodeString()

			d.wrapper.ScannerData(reqID, rank, contractDetails, distance, benchmark, projection, legsStr)

		}

		d.wrapper.ScannerDataEnd(reqID)
	}

func (c *Client) processScannerParametersMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		xml := msgDec.decodeString()

		d.wrapper.ScannerParameters(xml)
	}

func (c *Client) processFundamentalDataMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		data := msgDec.decodeString()

		d.wrapper.FundamentalData(reqID, data)
}
*/

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

/*
func (c *Client) processOpenOrdersEndMsg(msgDec *utils.Decoder) error {
	msgDec.Skip() // version
	if msgDec.Err() != nil {
		return msgDec.Err()
	}

	// Notify
	c.eventsHandler.OpenOrdersEnd()

	// Done
	return nil
}

func (c *Client) processAcctDownloadEndMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		accountName := msgDec.decodeString()

		d.wrapper.AccountDownloadEnd(accountName)
	}

func (c *Client) processExecutionDetailsEndMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.ExecDetailsEnd(reqID)
	}

func (c *Client) processExecutionDetailsEndMsgProtoBuf(msgDec *utils.Decoder) error {

		var executionDetailsEndProto protobuf.ExecutionDetailsEnd
		err := proto.Unmarshal(msgDec.Bytes(), &executionDetailsEndProto)
		if err != nil {
			log.Panic().Err(err).Msg("processExecutionDetailsEndMsgProtoBuf unmarshal error")
		}

		var reqID int64 = int64(executionDetailsEndProto.GetReqId())

		d.wrapper.ExecDetailsEnd(reqID)
	}

func (c *Client) processDeltaNeutralValidationMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		deltaNeutralContract := NewDeltaNeutralContract()

		deltaNeutralContract.ConID = msgDec.decodeInt64()
		deltaNeutralContract.Delta = msgDec.decodeFloat64()
		deltaNeutralContract.Price = msgDec.decodeFloat64()

		d.wrapper.DeltaNeutralValidation(reqID, deltaNeutralContract)
	}
*/

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

/*
func (c *Client) processCommissionAndFeesReportMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		commissionAndFeesReport := NewCommissionAndFeesReport()
		commissionAndFeesReport.ExecID = msgDec.decodeString()
		commissionAndFeesReport.CommissionAndFees = msgDec.decodeFloat64()
		commissionAndFeesReport.Currency = msgDec.decodeString()
		commissionAndFeesReport.RealizedPNL = msgDec.decodeFloat64()
		commissionAndFeesReport.Yield = msgDec.decodeFloat64()
		commissionAndFeesReport.YieldRedemptionDate = msgDec.decodeInt64()

		d.wrapper.CommissionAndFeesReport(commissionAndFeesReport)
	}

func (c *Client) processPositionDataMsg(msgDec *utils.Decoder) error {

		version := msgDec.decodeInt64()

		account := msgDec.decodeString()

		// decode contract fields
		contract := NewContract()
		contract.ConID = msgDec.decodeInt64()
		contract.Symbol = msgDec.decodeString()
		contract.SecType = msgDec.decodeString()
		contract.LastTradeDateOrContractMonth = msgDec.decodeString()
		contract.Strike = msgDec.decodeFloat64()
		contract.Right = msgDec.decodeString()
		contract.Multiplier = msgDec.decodeString()
		contract.Exchange = msgDec.decodeString()
		contract.Currency = msgDec.decodeString()
		contract.LocalSymbol = msgDec.decodeString()
		if version >= 2 {
			contract.TradingClass = msgDec.decodeString()
		}

		position := msgDec.decodeDecimal()

		var avgCost float64
		if version >= 3 {
			avgCost = msgDec.decodeFloat64()
		}

		d.wrapper.Position(account, contract, position, avgCost)
	}

func (c *Client) processPositionEndMsg(*msgDecfer) {

		d.wrapper.PositionEnd()
	}

func (c *Client) processAccountSummaryMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		account := msgDec.decodeString()
		tag := msgDec.decodeString()
		value := msgDec.decodeString()
		currency := msgDec.decodeString()

		d.wrapper.AccountSummary(reqID, account, tag, value, currency)
	}

func (c *Client) processAccountSummaryEndMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.AccountSummaryEnd(reqID)
	}

func (c *Client) processVerifyMessageApiMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		apiData := msgDec.decodeString()

		d.wrapper.VerifyMessageAPI(apiData)
	}

func (c *Client) processVerifyCompletedMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		isSuccessful := msgDec.decodeBool()
		errorText := msgDec.decodeString()

		d.wrapper.VerifyCompleted(isSuccessful, errorText)
	}

func (c *Client) processDisplayGroupListMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		groups := msgDec.decodeString()

		d.wrapper.DisplayGroupList(reqID, groups)
	}

func (c *Client) processDisplayGroupUpdatedMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		contractInfo := msgDec.decodeString()

		d.wrapper.DisplayGroupUpdated(reqID, contractInfo)
	}

func (c *Client) processVerifyAndAuthMessageApiMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		apiData := msgDec.decodeString()
		xyzChallange := msgDec.decodeString()

		d.wrapper.VerifyAndAuthMessageAPI(apiData, xyzChallange)
	}

func (c *Client) processVerifyAndAuthCompletedMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		isSuccessful := msgDec.decodeBool()
		errorText := msgDec.decodeString()

		d.wrapper.VerifyAndAuthCompleted(isSuccessful, errorText)
	}

func (c *Client) processPositionMultiMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		account := msgDec.decodeString()

		// decode contract fields
		contract := &Contract{}
		contract.ConID = msgDec.decodeInt64()
		contract.Symbol = msgDec.decodeString()
		contract.SecType = msgDec.decodeString()
		contract.LastTradeDateOrContractMonth = msgDec.decodeString()
		contract.Strike = msgDec.decodeFloat64()
		contract.Right = msgDec.decodeString()
		contract.Multiplier = msgDec.decodeString()
		contract.Exchange = msgDec.decodeString()
		contract.Currency = msgDec.decodeString()
		contract.LocalSymbol = msgDec.decodeString()
		contract.TradingClass = msgDec.decodeString()

		pos := msgDec.decodeDecimal()

		avgCost := msgDec.decodeFloat64()
		modelCode := msgDec.decodeString()

		d.wrapper.PositionMulti(reqID, account, modelCode, contract, pos, avgCost)
	}

func (c *Client) processPositionMultiEndMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.PositionMultiEnd(reqID)
	}

func (c *Client) processAccountUpdateMultiMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		account := msgDec.decodeString()
		modelCode := msgDec.decodeString()
		key := msgDec.decodeString()
		value := msgDec.decodeString()
		currency := msgDec.decodeString()

		d.wrapper.AccountUpdateMulti(reqID, account, modelCode, key, value, currency)
	}

func (c *Client) processAccountUpdateMultiEndMsg(msgDec *utils.Decoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.AccountUpdateMultiEnd(reqID)
	}

func (c *Client) processSoftDollarTiersMsg(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()

		tiersCount := msgDec.decodeInt64()
		tiers := make([]SoftDollarTier, 0, tiersCount)
		var i int64
		for i = 0; i < tiersCount; i++ {
			tier := NewSoftDollarTier()
			tier.Name = msgDec.decodeString()
			tier.Value = msgDec.decodeString()
			tier.DisplayName = msgDec.decodeString()
			tiers = append(tiers, tier)
		}

		d.wrapper.SoftDollarTiers(reqID, tiers)
	}

func (c *Client) processFamilyCodesMsg(msgDec *utils.Decoder) error {

		familyCodesCount := msgDec.decodeInt64()
		familyCodes := make([]FamilyCode, 0, familyCodesCount)
		var i int64
		for i = 0; i < familyCodesCount; i++ {
			familyCode := NewFamilyCode()
			familyCode.AccountID = msgDec.decodeString()
			familyCode.FamilyCodeStr = msgDec.decodeString()
			familyCodes = append(familyCodes, familyCode)
		}

		d.wrapper.FamilyCodes(familyCodes)
	}
*/

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
		cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String())
		cd.Contract.PrimaryExchange = msgDec.String()
		cd.Contract.Currency = msgDec.String()
		derivativeSecTypesCount := int(msgDec.Int32())
		if derivativeSecTypesCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative derivative security types count: %d", derivativeSecTypesCount))
			return wrapDecoderErr("processSymbolSamplesMsg", msgDec.Err())
		}
		cd.DerivativeSecTypes = make([]models.SecurityType, 0, derivativeSecTypesCount)
		for j := 0; j < derivativeSecTypesCount; j++ {
			derivativeSecType := models.NewSecurityTypeFromString(msgDec.String())
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

/*
func (c *Client) processMktDepthExchangesMsg(msgDec *utils.Decoder) error {

	depthMktDataDescriptionsCount := msgDec.decodeInt64()
	depthMktDataDescriptions := make([]DepthMktDataDescription, 0, depthMktDataDescriptionsCount)

	var i int64
	for i = 0; i < depthMktDataDescriptionsCount; i++ {
		desc := NewDepthMktDataDescription()
		desc.Exchange = msgDec.decodeString()
		desc.SecType = msgDec.decodeString()
		if d.serverVersion >= MIN_SERVER_VER_SERVICE_DATA_TYPE {
			desc.ListingExch = msgDec.decodeString()
			desc.SecType = msgDec.decodeString()
			desc.AggGroup = msgDec.decodeInt64()
		} else {
			_ = msgDec.decodeInt64() // boolean notSuppIsL2
		}

		depthMktDataDescriptions = append(depthMktDataDescriptions, desc)
	}

	d.wrapper.MktDepthExchanges(depthMktDataDescriptions)
}

func (c *Client) processTickNewsMsg(msgDec *utils.Decoder) error {

	tickerID := msgDec.decodeInt64()

	timeStamp := msgDec.decodeInt64()
	providerCode := msgDec.decodeString()
	articleID := msgDec.decodeString()
	headline := msgDec.decodeString()
	extraData := msgDec.decodeString()

	d.wrapper.TickNews(tickerID, timeStamp, providerCode, articleID, headline, extraData)
}

func (c *Client) processTickReqParamsMsg(msgDec *utils.Decoder) error {

	tickerID := msgDec.decodeInt64()

	minTick := msgDec.decodeFloat64()
	bboExchange := msgDec.decodeString()
	snapshotPermissions := msgDec.decodeInt64()

	d.wrapper.TickReqParams(tickerID, minTick, bboExchange, snapshotPermissions)
}

func (c *Client) processSmartComponentsMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	smartComponentsCount := msgDec.decodeInt64()
	smartComponents := make([]SmartComponent, 0, smartComponentsCount)
	var i int64
	for i = 0; i < smartComponentsCount; i++ {
		smartComponent := NewSmartComponent()
		smartComponent.BitNumber = msgDec.decodeInt64()
		smartComponent.Exchange = msgDec.decodeString()
		smartComponent.ExchangeLetter = msgDec.decodeString()
		smartComponents = append(smartComponents, smartComponent)
	}

	d.wrapper.SmartComponents(reqID, smartComponents)
}

func (c *Client) processNewsProvidersMsg(msgDec *utils.Decoder) error {

	newsProvidersCount := msgDec.decodeInt64()
	newsProviders := make([]NewsProvider, 0, newsProvidersCount)
	var i int64
	for i = 0; i < newsProvidersCount; i++ {
		provider := NewNewsProvider()
		provider.Name = msgDec.decodeString()
		provider.Code = msgDec.decodeString()
		newsProviders = append(newsProviders, provider)
	}

	d.wrapper.NewsProviders(newsProviders)
}

func (c *Client) processNewsArticleMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	articleType := msgDec.decodeInt64()
	articleText := msgDec.decodeString()

	d.wrapper.NewsArticle(reqID, articleType, articleText)
}

func (c *Client) processHistoricalNewsMsg(msgDec *utils.Decoder) error {

	requestID := msgDec.decodeInt64()

	time := msgDec.decodeString()
	providerCode := msgDec.decodeString()
	articleID := msgDec.decodeString()
	headline := msgDec.decodeString()

	d.wrapper.HistoricalNews(requestID, time, providerCode, articleID, headline)
}

func (c *Client) processHistoricalNewsEndMsg(msgDec *utils.Decoder) error {

	requestID := msgDec.decodeInt64()

	hasMore := msgDec.decodeBool()

	d.wrapper.HistoricalNewsEnd(requestID, hasMore)
}
*/

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

/*
func (c *Client) processHistogramDataMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	numPoints := msgDec.decodeInt64()
	data := make([]HistogramData, 0, numPoints)
	var i int64
	for i = 0; i < numPoints; i++ {
		p := HistogramData{}
		p.Price = msgDec.decodeFloat64()
		p.Size = msgDec.decodeDecimal()
		data = append(data, p)
	}

	d.wrapper.HistogramData(reqID, data)
}

func (c *Client) processRerouteMktDataReqMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	conID := msgDec.decodeInt64()
	exchange := msgDec.decodeString()

	d.wrapper.RerouteMktDataReq(reqID, conID, exchange)
}

func (c *Client) processRerouteMktDepthReqMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	conID := msgDec.decodeInt64()
	exchange := msgDec.decodeString()

	d.wrapper.RerouteMktDepthReq(reqID, conID, exchange)
}

func (c *Client) processMarketRuleMsg(msgDec *utils.Decoder) error {

	marketRuleID := msgDec.decodeInt64()

	priceIncrementsCount := msgDec.decodeInt64()
	priceIncrements := make([]PriceIncrement, 0, priceIncrementsCount)

	var i int64
	for i = 0; i < priceIncrementsCount; i++ {
		priceInc := NewPriceIncrement()
		priceInc.LowEdge = msgDec.decodeFloat64()
		priceInc.Increment = msgDec.decodeFloat64()
		priceIncrements = append(priceIncrements, priceInc)
	}

	d.wrapper.MarketRule(marketRuleID, priceIncrements)
}

func (c *Client) processPnLMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	dailyPnL := msgDec.decodeFloat64()
	var unrealizedPnL float64
	var realizedPnL float64

	if d.serverVersion >= MIN_SERVER_VER_UNREALIZED_PNL {
		unrealizedPnL = msgDec.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_REALIZED_PNL {
		realizedPnL = msgDec.decodeFloat64()
	}

	d.wrapper.Pnl(reqID, dailyPnL, unrealizedPnL, realizedPnL)
}

func (c *Client) processPnLSingleMsg(msgDec *utils.Decoder) error {

	reqID := msgDec.decodeInt64()

	pos := msgDec.decodeDecimal()
	dailyPnL := msgDec.decodeFloat64()
	var unrealizedPnL float64
	var realizedPnL float64

	if d.serverVersion >= MIN_SERVER_VER_UNREALIZED_PNL {
		unrealizedPnL = msgDec.decodeFloat64()
	}

	if d.serverVersion >= MIN_SERVER_VER_REALIZED_PNL {
		realizedPnL = msgDec.decodeFloat64()
	}

	value := msgDec.decodeFloat64()

	d.wrapper.PnlSingle(reqID, pos, dailyPnL, unrealizedPnL, realizedPnL, value)
}
*/

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

/*
func (c *Client) processOrderBoundMsg(msgDec *utils.Decoder) error {

		permID := msgDec.decodeInt64()
		clientId := msgDec.decodeInt64()
		orderId := msgDec.decodeInt64()

		d.wrapper.OrderBound(permID, clientId, orderId)
	}

func (c *Client) processCompletedOrderMsg(msgDec *utils.Decoder) error {

		order := NewOrder()
		contract := NewContract()
		orderState := NewOrderState()

		orderDecoder := &OrderDecoder{order, contract, orderState, Version(UNSET_INT), d.serverVersion}

		// read contract fields
		orderDecoder.decodeContractFields(msgDec)

		// read order fields
		orderDecoder.decodeAction(msgDec)
		orderDecoder.decodeTotalQuantity(msgDec)
		orderDecoder.decodeOrderType(msgDec)
		orderDecoder.decodeLmtPrice(msgDec)
		orderDecoder.decodeAuxPrice(msgDec)
		orderDecoder.decodeTIF(msgDec)
		orderDecoder.decodeOcaGroup(msgDec)
		orderDecoder.decodeAccount(msgDec)
		orderDecoder.decodeOpenClose(msgDec)
		orderDecoder.decodeOrigin(msgDec)
		orderDecoder.decodeOrderRef(msgDec)
		orderDecoder.decodePermId(msgDec)
		orderDecoder.decodeOutsideRth(msgDec)
		orderDecoder.decodeHidden(msgDec)
		orderDecoder.decodeDiscretionaryAmount(msgDec)
		orderDecoder.decodeGoodAfterTime(msgDec)
		orderDecoder.decodeFAParams(msgDec)
		orderDecoder.decodeModelCode(msgDec)
		orderDecoder.decodeGoodTillDate(msgDec)
		orderDecoder.decodeRule80A(msgDec)
		orderDecoder.decodePercentOffset(msgDec)
		orderDecoder.decodeSettlingFirm(msgDec)
		orderDecoder.decodeShortSaleParams(msgDec)
		orderDecoder.decodeBoxOrderParams(msgDec)
		orderDecoder.decodePegToStkOrVolOrderParams(msgDec)
		orderDecoder.decodeDisplaySize(msgDec)
		orderDecoder.decodeSweepToFill(msgDec)
		orderDecoder.decodeAllOrNone(msgDec)
		orderDecoder.decodeMinQty(msgDec)
		orderDecoder.decodeOcaType(msgDec)
		orderDecoder.decodeTriggerMethod(msgDec)
		orderDecoder.decodeVolOrderParams(msgDec, false)
		orderDecoder.decodeTrailParams(msgDec)
		orderDecoder.decodeComboLegs(msgDec)
		orderDecoder.decodeSmartComboRoutingParams(msgDec)
		orderDecoder.decodeScaleOrderParams(msgDec)
		orderDecoder.decodeHedgeParams(msgDec)
		orderDecoder.decodeClearingParams(msgDec)
		orderDecoder.decodeNotHeld(msgDec)
		orderDecoder.decodeDeltaNeutral(msgDec)
		orderDecoder.decodeAlgoParams(msgDec)
		orderDecoder.decodeSolicited(msgDec)
		orderDecoder.decodeOrderStatus(msgDec)
		orderDecoder.decodeVolRandomizeFlags(msgDec)
		orderDecoder.decodePegBenchParams(msgDec)
		orderDecoder.decodeConditions(msgDec)
		orderDecoder.decodeStopPriceAndLmtPriceOffset(msgDec)
		orderDecoder.decodeCashQty(msgDec)
		orderDecoder.decodeDontUseAutoPriceForHedge(msgDec)
		orderDecoder.decodeIsOmsContainer(msgDec)
		orderDecoder.decodeAutoCancelDate(msgDec)
		orderDecoder.decodeFilledQuantity(msgDec)
		orderDecoder.decodeRefFuturesConId(msgDec)
		orderDecoder.decodeAutoCancelParent(msgDec, MIN_CLIENT_VER)
		orderDecoder.decodeShareholder(msgDec)
		orderDecoder.decodeImbalanceOnly(msgDec, MIN_CLIENT_VER)
		orderDecoder.decodeRouteMarketableToBbo(msgDec)
		orderDecoder.decodeParentPermId(msgDec)
		orderDecoder.decodeCompletedTime(msgDec)
		orderDecoder.decodeCompletedStatus(msgDec)
		orderDecoder.decodePegBestPegMidOrderAttributes(msgDec)
		orderDecoder.decodeCustomerAccount(msgDec)
		orderDecoder.decodeProfessionalCustomer(msgDec)
		orderDecoder.decodeSubmitter(msgDec)

		d.wrapper.CompletedOrder(contract, order, orderState)
	}

	func (c *Client) processCompletedOrdersEndMsg(*msgDecfer) {
		d.wrapper.CompletedOrdersEnd()
	}

func (c *Client) processReplaceFAEndMsg(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()
		text := msgDec.decodeString()

		d.wrapper.ReplaceFAEnd(reqID, text)
	}

func (c *Client) processWshMetaData(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()
		dataJSON := msgDec.decodeString()

		d.wrapper.WshMetaData(reqID, dataJSON)
	}

func (c *Client) processWshEventData(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()
		dataJSON := msgDec.decodeString()

		d.wrapper.WshEventData(reqID, dataJSON)
	}

func (c *Client) processHistoricalSchedule(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()
		startDateTime := msgDec.decodeString()
		endDateTime := msgDec.decodeString()
		timeZone := msgDec.decodeString()
		sessionsCount := msgDec.decodeInt64()
		sessions := make([]HistoricalSession, 0, sessionsCount)
		var i int64
		for i = 0; i < sessionsCount; i++ {
			historicalSession := NewHistoricalSession()
			historicalSession.StartDateTime = msgDec.decodeString()
			historicalSession.EndDateTime = msgDec.decodeString()
			historicalSession.RefDate = msgDec.decodeString()
			sessions = append(sessions, historicalSession)
		}

		d.wrapper.HistoricalSchedule(reqID, startDateTime, endDateTime, timeZone, sessions)
	}

func (c *Client) processUserInfo(msgDec *utils.Decoder) error {

		reqID := msgDec.decodeInt64()
		whiteBrandingId := msgDec.decodeString()

		d.wrapper.UserInfo(reqID, whiteBrandingId)
	}
*/

func (c *Client) processCurrentTimeInMillisMsg(msgDec *message.Decoder) error {
	ts := msgDec.EpochTimestamp(true)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processCurrentTimeInMillisMsg", msgDec.Err())
	}

	return c.processCurrentTimeInMillisCommon(ts)
}

func (c *Client) processCurrentTimeInMillisCommon(ts time.Time) error {
	c.reqMgr.withRequestWithoutID(common.REQ_CURRENT_TIME_IN_MILLIS, func(_resp interface{}) error {
		resp := _resp.(*models.CurrentTimeResponse)

		resp.CurrentTime = ts

		// Done
		return nil
	})

	// Done
	return nil
}
