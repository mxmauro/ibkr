package ibkr

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type implDepthMarketData struct {
	position  int
	operation int
	bidSide   bool
	entry     models.DepthMarketBookEntry
}

// -----------------------------------------------------------------------------

func (c *Client) getIncomingMessageID(msg []byte) (msgID int64, usingProtobuf bool, err error) {
	if len(msg) >= 4 {
		msgID = int64(binary.BigEndian.Uint32(msg))
		if msgID >= common.PROTOBUF_MSG_ID {
			usingProtobuf = true
			msgID -= common.PROTOBUF_MSG_ID
		}
	} else {
		err = errors.New("received invalid message")
	}
	return
}

func (c *Client) processIncomingMessage(msg []byte) error {
	msgID, _, err := c.getIncomingMessageID(msg)
	if err != nil {
		return err
	}

	msgDec := utils.NewMessageDecoder(msg[4:])
	switch msgID {
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

		/*
			case ORDER_STATUS:
				return c.processOrderStatusMsg(msgDec)
		*/

	case common.ERR_MSG:
		return c.processErrorMessage(msgDec)
		/*
			case OPEN_ORDER:
				return c.processOpenOrderMsg(msgDec)
			case ACCT_VALUE:
				return c.processAcctValueMsg(msgDec)
			case PORTFOLIO_VALUE:
				return c.processPortfolioValueMsg(msgDec)
			case ACCT_UPDATE_TIME:
				return c.processAcctUpdateTimeMsg(msgDec)
		*/

	case common.NEXT_VALID_ID:
		return nil // Ignore this message
	case common.CONTRACT_DATA:
		return c.processContractDataMsg(msgDec)
	case common.BOND_CONTRACT_DATA:
		return c.processBondContractDataMsg(msgDec)

		/*
			case EXECUTION_DATA:
				if useProtoBuf {
					return c.processExecutionDetailsMsgProtoBuf(msgDec)
				} else {
					return c.processExecutionDetailsMsg(msgDec)
				}
		*/

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
		/*
			case REAL_TIME_BARS:
				return c.processRealTimeBarsMsg(msgDec)
			case FUNDAMENTAL_DATA:
				return c.processFundamentalDataMsg(msgDec)
		*/

	case common.CONTRACT_DATA_END:
		return c.processContractDataEndMsg(msgDec)
	case common.OPEN_ORDER_END:
		return c.processOpenOrdersEndMsg(msgDec)

		/*
			case ACCT_DOWNLOAD_END:
				return c.processAcctDownloadEndMsg(msgDec)
			case EXECUTION_DATA_END:
				if useProtoBuf {
					return c.processExecutionDetailsEndMsgProtoBuf(msgDec)
				} else {
					return c.processExecutionDetailsEndMsg(msgDec)
				}
			case DELTA_NEUTRAL_VALIDATION:
				return c.processDeltaNeutralValidationMsg(msgDec)
		*/
	case common.TICK_SNAPSHOT_END:
		return c.processTickSnapshotEndMsg(msgDec)
	case common.MARKET_DATA_TYPE:
		return nil // Ignore this message. We don't use the ticker callback.
		/*
			case COMMISSION_AND_FEES_REPORT:
				return c.processCommissionAndFeesReportMsg(msgDec)
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
			case SECURITY_DEFINITION_OPTION_PARAMETER:
				return c.processSecurityDefinitionOptionalParameterMsg(msgDec)
			case SECURITY_DEFINITION_OPTION_PARAMETER_END:
				return c.processSecurityDefinitionOptionalParameterEndMsg(msgDec)
			case SOFT_DOLLAR_TIERS:
				return c.processSoftDollarTiersMsg(msgDec)
			case FAMILY_CODES:
				return c.processFamilyCodesMsg(msgDec)
			case SMART_COMPONENTS:
				return c.processSmartComponentsMsg(msgDec)
			case TICK_REQ_PARAMS:
				return c.processTickReqParamsMsg1(msgDec)
		*/
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
		return nil // KeepUpToDate is always false, ignore the message
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
		return c.processHistoricalTicks(msgDec)
	case common.HISTORICAL_TICKS_BID_ASK:
		return c.processHistoricalTicksBidAsk(msgDec)
	case common.HISTORICAL_TICKS_LAST:
		return c.processHistoricalTicksLast(msgDec)
		/*
			case TICK_BY_TICK:
				return c.processTickByTickDataMsg(msgDec)
			case ORDER_BOUND:
				return c.processOrderBoundMsg(msgDec)
			case COMPLETED_ORDER:
				return c.processCompletedOrderMsg(msgDec)
			case COMPLETED_ORDERS_END:
				return c.processCompletedOrdersEndMsg(msgDec)
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

	// Raise the event if an event handler is present
	if c.eventsHandler != nil {
		c.eventsHandler.ReceivedUnknownMessage(int(msgID))
	}

	// Done
	return nil
}

func (c *Client) processTickPriceMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		price := msgDec.Float64(false)
		size := models.NewDecimalFromMessageDecoder(msgDec, false) // ver 2 field
		attrMask := msgDec.Int64(false)                            // ver 3 field
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := models.NewTopMarketDataPrice(tickType)
		data.Price = price
		data.CanAutoExecute = attrMask&0x1 != 0
		data.PastLimit = attrMask&0x2 != 0
		data.PreOpen = attrMask&0x4 != 0
		resp.Channel <- data

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
			resp.Channel <- sizeData
		}

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processTickSizeMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		size := models.NewDecimalFromMessageDecoder(msgDec, false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := models.NewTopMarketDataSize(tickType)
		data.Size = size
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processTickOptionComputationMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		tickAttrib := msgDec.Int64(false)
		impliedVol := msgDec.Float64(false)
		if impliedVol < 0 {
			impliedVol = common.UNSET_FLOAT
		}
		delta := msgDec.Float64(false)
		if delta == -2 { // -2 is the "not computed" indicator
			delta = common.UNSET_FLOAT
		}
		price := common.UNSET_FLOAT
		pvDividend := common.UNSET_FLOAT
		if tickType == models.TickTypeModelOptionComputation || tickType == models.TickTypeDelayedModelOptionComputation {
			price = msgDec.Float64(false)
			if price == -1 { // -1 is the "not computed" indicator
				price = common.UNSET_FLOAT
			}
			pvDividend = msgDec.Float64(false)
			if pvDividend == -1 { // -1 is the "not computed" indicator
				pvDividend = common.UNSET_FLOAT
			}
		}
		gamma := msgDec.Float64(false)
		if gamma == -2 { // -2 is the "not yet computed" indicator
			gamma = common.UNSET_FLOAT
		}
		vega := msgDec.Float64(false)
		if vega == -2 { // -2 is the "not yet computed" indicator
			vega = common.UNSET_FLOAT
		}
		theta := msgDec.Float64(false)
		if theta == -2 { // -2 is the "not yet computed" indicator
			theta = common.UNSET_FLOAT
		}
		undPrice := msgDec.Float64(false)
		if undPrice == -1 { // -1 is the "not computed" indicator
			undPrice = common.UNSET_FLOAT
		}
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := models.NewTopMarketDataOptionComputation(tickType)
		data.IsPriceBased = tickAttrib != 0
		if impliedVol != common.UNSET_FLOAT {
			data.ImpliedVolatility = &impliedVol
		}
		if delta != common.UNSET_FLOAT {
			data.Delta = &delta
		}
		if price != common.UNSET_FLOAT {
			data.Price = &price
		}
		if pvDividend != common.UNSET_FLOAT {
			data.PvDividend = &pvDividend
		}
		if gamma != common.UNSET_FLOAT {
			data.Gamma = &gamma
		}
		if vega != common.UNSET_FLOAT {
			data.Vega = &vega
		}
		if theta != common.UNSET_FLOAT {
			data.Theta = &theta
		}
		if undPrice != common.UNSET_FLOAT {
			data.UnderlyingPrice = &undPrice
		}
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processTickGenericMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		value := msgDec.Float64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := models.NewTopMarketDataGeneric(tickType)
		data.Value = value
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processTickStringMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		value := msgDec.String(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		switch tickType {
		case models.TickTypeLastTimestamp:
			fallthrough
		case models.TickTypeDelayedLastTimestamp:
			fallthrough
		case models.TickTypeLastRegulatoryTime:
			secs, err2 := strconv.Atoi(value)
			if err2 != nil {
				return false, err2
			}

			// Notify
			data := models.NewTopMarketDataTimestamp(tickType)
			data.Timestamp = time.Unix(int64(secs), 0).UTC()
			resp.Channel <- data

		default:
			// Notify
			data := models.NewTopMarketDataString(tickType)
			data.Value = value
			resp.Channel <- data
		}

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processTickEfpMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.TopMarketDataResponse)

		tickType := models.TickType(msgDec.Int64(false))
		basisPoints := msgDec.Float64(false)
		formattedBasisPoints := msgDec.String(false)
		totalDividends := msgDec.Float64(false)
		holdDays := msgDec.Int64(false)
		futureLastTradeDate := msgDec.String(false)
		dividendImpact := msgDec.Float64(false)
		dividendsToLastTradeDate := msgDec.Float64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := models.NewTopMarketDataEFP(tickType)
		data.BasisPoints = basisPoints
		data.FormattedBasisPoints = formattedBasisPoints
		data.TotalDividends = totalDividends
		data.HoldDays = int(holdDays)
		data.FutureLastTradeDate = futureLastTradeDate
		data.DividendImpact = dividendImpact
		data.DividendsToLastTradeDate = dividendsToLastTradeDate
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processOrderStatusMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processErrorMessage(msgDec *utils.MessageDecoder) error {
	// Get the optional originating request ID
	reqID := c.decodeRequestID(msgDec, true)
	if reqID == 0 {
		return msgDec.Err()
	}
	code := msgDec.Int64(false)
	msg := msgDec.String(false)
	advancedOrderRejectJson := msgDec.String(false)
	ts := msgDec.EpochTimestamp(true)
	if msgDec.Err() != nil {
		return msgDec.Err()
	}

	// Ignore the following error codesL
	if code == 300 || code == 10167 {
		// 300 - "Can't find EId with tickerId: ###" messages because they are sent when we try to cancel a request,
		//       and it does no longer exist.
		// 10167 - "Requested market data is not subscribed. Delayed market data is not enabled" warning when delayed
		//         data is requested.
		return nil
	}

	if reqID > 0 {
		c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
			// Done

			return false, newRequestError(ts, code, msg, advancedOrderRejectJson)
		})
	} else {
		// Notify
		if c.eventsHandler != nil {
			c.eventsHandler.Error(ts, int(code), msg, advancedOrderRejectJson)
		}
	}

	// Done
	return nil
}

/*
func (c *Client) processOpenOrderMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processAcctValueMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		tag := msgDec.decodeString()
		val := msgDec.decodeString()
		currency := msgDec.decodeString()
		accountName := msgDec.decodeString()

		d.wrapper.UpdateAccountValue(tag, val, currency, accountName)
	}

func (c *Client) processPortfolioValueMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processAcctUpdateTimeMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		timeStamp := msgDec.decodeString()

		d.wrapper.UpdateAccountTime(timeStamp)
	}
*/
func (c *Client) processContractDataMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.ContractDetailsResponse)

		cd := models.NewContractDetails()
		cd.Contract.Symbol = msgDec.String(false)
		cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String(false))
		c.readLastTradeDate(msgDec, &cd, false)
		cd.Contract.LastTradeDate = msgDec.String(false)
		cd.Contract.Strike = msgDec.Float64(false)
		cd.Contract.Right = msgDec.String(false)
		cd.Contract.Exchange = msgDec.String(false)
		cd.Contract.Currency = msgDec.String(false)
		cd.Contract.LocalSymbol = msgDec.String(false)
		cd.MarketName = msgDec.String(false)
		cd.Contract.TradingClass = msgDec.String(false)
		cd.Contract.ConID = msgDec.Int64(false)
		cd.MinTick = msgDec.Float64(false)
		cd.Contract.Multiplier = msgDec.String(false)
		cd.OrderTypes = msgDec.String(false)
		cd.ValidExchanges = msgDec.String(false)
		cd.PriceMagnifier = msgDec.Int64(false)
		cd.UnderConID = msgDec.Int64(false)
		cd.LongName = msgDec.String(true)
		cd.Contract.PrimaryExchange = msgDec.String(false)
		cd.ContractMonth = msgDec.String(false)
		cd.Industry = msgDec.String(false)
		cd.Category = msgDec.String(false)
		cd.Subcategory = msgDec.String(false)
		cd.TimeZoneID = msgDec.String(false)
		cd.TradingHours = msgDec.String(false)
		cd.LiquidHours = msgDec.String(false)
		cd.EVRule = msgDec.String(false)
		cd.EVMultiplier = msgDec.Int64(false)
		secIDListCount := msgDec.Int64(false)
		if secIDListCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative security id count: %d", secIDListCount))
			return false, msgDec.Err()
		}
		cd.SecIDList = make(models.TagValueList, 0, secIDListCount)
		for i := int64(0); i < secIDListCount; i++ {
			tagValue := models.NewTagValue()
			tagValue.Tag = msgDec.String(false)
			tagValue.Value = msgDec.String(false)
			cd.SecIDList = append(cd.SecIDList, tagValue)
		}
		cd.AggGroup = msgDec.Int64(false)
		cd.UnderSymbol = msgDec.String(false)
		cd.UnderSecType = models.NewSecurityTypeFromString(msgDec.String(false))
		cd.MarketRuleIDs = msgDec.String(false)
		cd.RealExpirationDate = msgDec.String(false)
		cd.StockType = msgDec.String(false)
		cd.MinSize = models.NewDecimalFromMessageDecoder(msgDec, false)
		cd.SizeIncrement = models.NewDecimalFromMessageDecoder(msgDec, false)
		cd.SuggestedSizeIncrement = models.NewDecimalFromMessageDecoder(msgDec, false)
		if cd.Contract.SecType == models.SecurityTypeMutualFund {
			cd.FundName = msgDec.String(false)
			cd.FundFamily = msgDec.String(false)
			cd.FundType = msgDec.String(false)
			cd.FundFrontLoad = msgDec.String(false)
			cd.FundBackLoad = msgDec.String(false)
			cd.FundBackLoadTimeInterval = msgDec.String(false)
			cd.FundManagementFee = msgDec.String(false)
			cd.FundClosed = msgDec.Bool()
			cd.FundClosedForNewInvestors = msgDec.Bool()
			cd.FundClosedForNewMoney = msgDec.Bool()
			cd.FundNotifyAmount = msgDec.String(false)
			cd.FundMinimumInitialPurchase = msgDec.String(false)
			cd.FundSubsequentMinimumPurchase = msgDec.String(false)
			cd.FundBlueSkyStates = msgDec.String(false)
			cd.FundBlueSkyTerritories = msgDec.String(false)
			cd.FundDistributionPolicyIndicator = models.NewFundDistributionPolicyIndicatorFromString(msgDec.String(false))
			cd.FundAssetType = models.NewFundAssetFromString(msgDec.String(false))
		}
		ineligibilityReasonListCount := msgDec.Int64(false)
		if ineligibilityReasonListCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative ineligibility reason count: %d", ineligibilityReasonListCount))
			return false, msgDec.Err()
		}
		cd.IneligibilityReasonList = make([]models.IneligibilityReason, 0, ineligibilityReasonListCount)
		for i := int64(0); i < ineligibilityReasonListCount; i++ {
			ineligibilityReason := models.IneligibilityReason{}
			ineligibilityReason.ID = msgDec.String(false)
			ineligibilityReason.Description = msgDec.String(false)

			cd.IneligibilityReasonList = append(cd.IneligibilityReasonList, ineligibilityReason)
		}
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		resp.ContractDetails = append(resp.ContractDetails, cd)

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processBondContractDataMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.ContractDetailsResponse)

		cd := models.NewContractDetails()
		cd.Contract = *models.NewContract()
		cd.Contract.Symbol = msgDec.String(false)
		cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String(false))
		cd.Cusip = msgDec.String(false)
		cd.Coupon = msgDec.Float64(false)
		c.readLastTradeDate(msgDec, &cd, true)
		cd.IssueDate = msgDec.String(false)
		cd.Ratings = msgDec.String(false)
		cd.BondType = msgDec.String(false)
		cd.CouponType = msgDec.String(false)
		cd.Convertible = msgDec.Bool()
		cd.Callable = msgDec.Bool()
		cd.Putable = msgDec.Bool()
		cd.DescAppend = msgDec.String(false)
		cd.Contract.Exchange = msgDec.String(false)
		cd.Contract.Currency = msgDec.String(false)
		cd.MarketName = msgDec.String(false)
		cd.Contract.TradingClass = msgDec.String(false)
		cd.Contract.ConID = msgDec.Int64(false)
		cd.MinTick = msgDec.Float64(false)
		cd.OrderTypes = msgDec.String(false)
		cd.ValidExchanges = msgDec.String(false)
		cd.NextOptionDate = msgDec.String(false)
		cd.NextOptionType = msgDec.String(false)
		cd.NextOptionPartial = msgDec.Bool()
		cd.Notes = msgDec.String(false)
		cd.LongName = msgDec.String(false)
		cd.TimeZoneID = msgDec.String(false)
		cd.TradingHours = msgDec.String(false)
		cd.LiquidHours = msgDec.String(false)
		cd.EVRule = msgDec.String(false)
		cd.EVMultiplier = msgDec.Int64(false)
		secIDListCount := msgDec.Int64(false)
		if secIDListCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative security id count: %d", secIDListCount))
			return false, msgDec.Err()
		}
		cd.SecIDList = make(models.TagValueList, 0, secIDListCount)
		var i int64
		for i = 0; i < secIDListCount; i++ {
			tagValue := models.NewTagValue()
			tagValue.Tag = msgDec.String(false)
			tagValue.Value = msgDec.String(false)
			cd.SecIDList = append(cd.SecIDList, tagValue)
		}
		cd.AggGroup = msgDec.Int64(false)
		cd.MarketRuleIDs = msgDec.String(false)
		cd.MinSize = models.NewDecimalFromMessageDecoder(msgDec, false)
		cd.SizeIncrement = models.NewDecimalFromMessageDecoder(msgDec, false)
		cd.SuggestedSizeIncrement = models.NewDecimalFromMessageDecoder(msgDec, false)
		cd.IneligibilityReasonList = make([]models.IneligibilityReason, 0)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		resp.ContractDetails = append(resp.ContractDetails, cd)

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processExecutionDetailsMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processExecutionDetailsMsgProtoBuf(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processMarketDepthMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.MarketDepthDataResponse)

		position := msgDec.Int64(false)
		if position < 0 || position >= int64(resp.Book.Size) {
			return false, msgDec.Err()
		}
		operation := msgDec.Int64(false)
		side := msgDec.Int64(false)
		price := msgDec.Float64(false)
		size := models.NewDecimalFromMessageDecoder(msgDec, false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := &implDepthMarketData{
			position:  int(position),
			operation: int(operation),
			bidSide:   side == 0,
			entry: models.DepthMarketBookEntry{
				Price: price,
				Size:  size,
			},
		}
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processMarketDepthL2Msg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.MarketDepthDataResponse)

		position := msgDec.Int64(false)
		if position < 0 || position >= int64(resp.Book.Size) {
			return false, msgDec.Err()
		}
		_ = msgDec.String(false) // Market maker
		operation := msgDec.Int64(false)
		side := msgDec.Int64(false)
		price := msgDec.Float64(false)
		size := models.NewDecimalFromMessageDecoder(msgDec, false)
		_ = msgDec.Bool() // Is Smart Depth
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Notify
		data := &implDepthMarketData{
			position:  int(position),
			operation: int(operation),
			bidSide:   side == 0,
			entry: models.DepthMarketBookEntry{
				Price: price,
				Size:  size,
			},
		}
		resp.Channel <- data

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (idmd *implDepthMarketData) UpdateBook(book *models.MarketDepthBook) {
	var bookSide *[]models.DepthMarketBookEntry

	if idmd.bidSide {
		bookSide = &book.Bids
	} else {
		bookSide = &book.Asks
	}

	// Execute operation
	c := cap(*bookSide)
	l := len(*bookSide)
	switch idmd.operation {
	case 0: // Insert
		if idmd.position < c {
			book.Asks = (*bookSide)[:l+1]
			copy((*bookSide)[idmd.position+1:], (*bookSide)[idmd.position:l])
		} else {
			copy((*bookSide)[idmd.position+1:], (*bookSide)[idmd.position:l-1])
		}
		(*bookSide)[idmd.position] = idmd.entry

	case 1: // Update
		if idmd.position > l {
			*bookSide = (*bookSide)[:idmd.position+1]
		}
		(*bookSide)[idmd.position] = idmd.entry

	case 2: // Delete
		if idmd.position < l {
			copy((*bookSide)[idmd.position:], (*bookSide)[idmd.position+1:])
			(*bookSide)[l-1] = models.DepthMarketBookEntry{}
			*bookSide = (*bookSide)[:l-1]
		}
	}
}

/*
func (c *Client) processNewsBulletinsMsg(msgDec *utils.MessageDecoder) error {

	msgDec.decode() // version

	msgID := msgDec.decodeInt64()
	msgType := msgDec.decodeInt64()
	newsMessage := msgDec.decodeString()
	originExch := msgDec.decodeString()

	d.wrapper.UpdateNewsBulletin(msgID, msgType, newsMessage, originExch)
}
*/

func (c *Client) processManagedAccountsMsg(msgDec *utils.MessageDecoder) error {
	c.reqMgr.withRequestWithoutID(common.REQ_MANAGED_ACCTS, func(_resp interface{}) error {
		resp := _resp.(*models.ManagedAccountsResponse)

		msgDec.Skip() // version
		accountsNames := msgDec.String(false)
		if msgDec.Err() != nil {
			return msgDec.Err()
		}

		resp.Accounts = strings.Split(accountsNames, ",")

		// Done
		return nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processReceiveFaMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		faDataType := FaDataType(msgDec.decodeInt64())
		cxml := msgDec.decodeString()

		d.wrapper.ReceiveFA(faDataType, cxml)
	}
*/

func (c *Client) decodeRequestID(msgDec *utils.MessageDecoder, canBeOptional bool) int32 {
	reqID := msgDec.Int64(false)
	if msgDec.Err() != nil {
		return 0
	}
	if canBeOptional && reqID <= 0 {
		return -1
	}
	if reqID < 1 || reqID > math.MaxInt32 {
		msgDec.SetErr(fmt.Errorf("received invalid request id: %d", reqID))
		return 0
	}
	return int32(reqID)
}

func (c *Client) processHistoricalDataMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.HistoricalDataResponse)

		barsCount := msgDec.Int64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}
		if barsCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative bars count: %d", barsCount))
			return false, msgDec.Err()
		}

		for i := int64(0); i < barsCount; i++ {
			bar := models.NewBar()
			// Epoch because the request of historical data has the date format equal to 2.
			bar.Date = msgDec.EpochTimestamp(false)
			bar.Open = msgDec.Float64(false)
			bar.High = msgDec.Float64(false)
			bar.Low = msgDec.Float64(false)
			bar.Close = msgDec.Float64(false)
			bar.Volume = models.NewDecimalFromMessageDecoder(msgDec, false)
			bar.Wap = models.NewDecimalFromMessageDecoder(msgDec, false)
			bar.BarCount = msgDec.Int64(false)

			resp.Bars = append(resp.Bars, bar)
		}

		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Done
		return false, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processHistoricalDataEndMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
		msgDec.Skip()
		msgDec.Skip()
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Done
		return true, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processScannerDataMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processScannerParametersMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		xml := msgDec.decodeString()

		d.wrapper.ScannerParameters(xml)
	}

func (c *Client) processRealTimeBarsMsg(msgDec *utils.MessageDecoder) error {

	msgDec.decode() // version

	reqID := msgDec.decodeInt64()

	time := msgDec.decodeInt64()
	open := msgDec.decodeFloat64()
	high := msgDec.decodeFloat64()
	low := msgDec.decodeFloat64()
	close := msgDec.decodeFloat64()
	volume := msgDec.decodeDecimal()
	wap := msgDec.decodeDecimal()
	count := msgDec.decodeInt64()

	d.wrapper.RealtimeBar(reqID, time, open, high, low, close, volume, wap, count)

}

func (c *Client) processFundamentalDataMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		data := msgDec.decodeString()

		d.wrapper.FundamentalData(reqID, data)
	}
*/
func (c *Client) processContractDataEndMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version

	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
		// Done
		return true, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processOpenOrdersEndMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	if msgDec.Err() != nil {
		return msgDec.Err()
	}

	// Notify
	c.eventsHandler.OpenOrdersEnd()

	// Done
	return nil
}

/*
func (c *Client) processAcctDownloadEndMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		accountName := msgDec.decodeString()

		d.wrapper.AccountDownloadEnd(accountName)
	}

func (c *Client) processExecutionDetailsEndMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.ExecDetailsEnd(reqID)
	}

func (c *Client) processExecutionDetailsEndMsgProtoBuf(msgDec *utils.MessageDecoder) error {

		var executionDetailsEndProto protobuf.ExecutionDetailsEnd
		err := proto.Unmarshal(msgDec.Bytes(), &executionDetailsEndProto)
		if err != nil {
			log.Panic().Err(err).Msg("processExecutionDetailsEndMsgProtoBuf unmarshal error")
		}

		var reqID int64 = int64(executionDetailsEndProto.GetReqId())

		d.wrapper.ExecDetailsEnd(reqID)
	}

func (c *Client) processDeltaNeutralValidationMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		deltaNeutralContract := NewDeltaNeutralContract()

		deltaNeutralContract.ConID = msgDec.decodeInt64()
		deltaNeutralContract.Delta = msgDec.decodeFloat64()
		deltaNeutralContract.Price = msgDec.decodeFloat64()

		d.wrapper.DeltaNeutralValidation(reqID, deltaNeutralContract)
	}
*/

func (c *Client) processTickSnapshotEndMsg(msgDec *utils.MessageDecoder) error {
	msgDec.Skip() // version
	// Gets the originating ticker ID
	tickerID := c.decodeRequestID(msgDec, false)
	if tickerID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(tickerID, func(_ interface{}) (bool, error) {
		// Signal end of snapshot
		return true, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processCommissionAndFeesReportMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processPositionDataMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processAccountSummaryMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		account := msgDec.decodeString()
		tag := msgDec.decodeString()
		value := msgDec.decodeString()
		currency := msgDec.decodeString()

		d.wrapper.AccountSummary(reqID, account, tag, value, currency)
	}

func (c *Client) processAccountSummaryEndMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.AccountSummaryEnd(reqID)
	}

func (c *Client) processVerifyMessageApiMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		apiData := msgDec.decodeString()

		d.wrapper.VerifyMessageAPI(apiData)
	}

func (c *Client) processVerifyCompletedMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		isSuccessful := msgDec.decodeBool()
		errorText := msgDec.decodeString()

		d.wrapper.VerifyCompleted(isSuccessful, errorText)
	}

func (c *Client) processDisplayGroupListMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		groups := msgDec.decodeString()

		d.wrapper.DisplayGroupList(reqID, groups)
	}

func (c *Client) processDisplayGroupUpdatedMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		contractInfo := msgDec.decodeString()

		d.wrapper.DisplayGroupUpdated(reqID, contractInfo)
	}

func (c *Client) processVerifyAndAuthMessageApiMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		apiData := msgDec.decodeString()
		xyzChallange := msgDec.decodeString()

		d.wrapper.VerifyAndAuthMessageAPI(apiData, xyzChallange)
	}

func (c *Client) processVerifyAndAuthCompletedMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		isSuccessful := msgDec.decodeBool()
		errorText := msgDec.decodeString()

		d.wrapper.VerifyAndAuthCompleted(isSuccessful, errorText)
	}

func (c *Client) processPositionMultiMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processPositionMultiEndMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.PositionMultiEnd(reqID)
	}

func (c *Client) processAccountUpdateMultiMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()
		account := msgDec.decodeString()
		modelCode := msgDec.decodeString()
		key := msgDec.decodeString()
		value := msgDec.decodeString()
		currency := msgDec.decodeString()

		d.wrapper.AccountUpdateMulti(reqID, account, modelCode, key, value, currency)
	}

func (c *Client) processAccountUpdateMultiEndMsg(msgDec *utils.MessageDecoder) error {

		msgDec.decode() // version

		reqID := msgDec.decodeInt64()

		d.wrapper.AccountUpdateMultiEnd(reqID)
	}

func (c *Client) processSecurityDefinitionOptionalParameterMsg(msgDec *utils.MessageDecoder) error {

	reqID := msgDec.decodeInt64()

	exchange := msgDec.decodeString()
	underlyingConID := msgDec.decodeInt64()
	tradingClass := msgDec.decodeString()
	multiplier := msgDec.decodeString()

	expCount := msgDec.decodeInt64()
	expirations := make([]string, 0, expCount)
	var i int64
	for i = 0; i < expCount; i++ {
		expiration := msgDec.decodeString()
		expirations = append(expirations, expiration)
	}

	strikeCount := msgDec.decodeInt64()
	strikes := make([]float64, 0, strikeCount)
	for i = 0; i < strikeCount; i++ {
		strike := msgDec.decodeFloat64()
		strikes = append(strikes, strike)
	}

	d.wrapper.SecurityDefinitionOptionParameter(reqID, exchange, underlyingConID, tradingClass, multiplier, expirations, strikes)

}

func (c *Client) processSecurityDefinitionOptionalParameterEndMsg(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()

		d.wrapper.SecurityDefinitionOptionParameterEnd(reqID)
	}

func (c *Client) processSoftDollarTiersMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processFamilyCodesMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processSymbolSamplesMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.MatchingSymbolsResponse)

		contractDescriptionsCount := msgDec.Int64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}
		if contractDescriptionsCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative contract descriptions count: %d", contractDescriptionsCount))
			return false, msgDec.Err()
		}

		for i := int64(0); i < contractDescriptionsCount; i++ {
			cd := models.NewContractDescription()
			cd.Contract.ConID = msgDec.Int64(false)
			cd.Contract.Symbol = msgDec.String(false)
			cd.Contract.SecType = models.NewSecurityTypeFromString(msgDec.String(false))
			cd.Contract.PrimaryExchange = msgDec.String(false)
			cd.Contract.Currency = msgDec.String(false)
			derivativeSecTypesCount := msgDec.Int64(false)
			if derivativeSecTypesCount < 0 {
				msgDec.SetErr(fmt.Errorf("negative derivative security types count: %d", derivativeSecTypesCount))
				return false, msgDec.Err()
			}
			cd.DerivativeSecTypes = make([]models.SecurityType, 0, derivativeSecTypesCount)
			for j := int64(0); j < derivativeSecTypesCount; j++ {
				derivativeSecType := models.NewSecurityTypeFromString(msgDec.String(false))

				cd.DerivativeSecTypes = append(cd.DerivativeSecTypes, derivativeSecType)
			}
			cd.Contract.Description = msgDec.String(false)
			cd.Contract.IssuerID = msgDec.String(false)

			resp.ContractDescriptions = append(resp.ContractDescriptions, cd)
		}

		// Done
		return true, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processMktDepthExchangesMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processTickNewsMsg(msgDec *utils.MessageDecoder) error {

	tickerID := msgDec.decodeInt64()

	timeStamp := msgDec.decodeInt64()
	providerCode := msgDec.decodeString()
	articleID := msgDec.decodeString()
	headline := msgDec.decodeString()
	extraData := msgDec.decodeString()

	d.wrapper.TickNews(tickerID, timeStamp, providerCode, articleID, headline, extraData)
}

func (c *Client) processTickReqParamsMsg(msgDec *utils.MessageDecoder) error {

	tickerID := msgDec.decodeInt64()

	minTick := msgDec.decodeFloat64()
	bboExchange := msgDec.decodeString()
	snapshotPermissions := msgDec.decodeInt64()

	d.wrapper.TickReqParams(tickerID, minTick, bboExchange, snapshotPermissions)
}

func (c *Client) processSmartComponentsMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processNewsProvidersMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processNewsArticleMsg(msgDec *utils.MessageDecoder) error {

	reqID := msgDec.decodeInt64()

	articleType := msgDec.decodeInt64()
	articleText := msgDec.decodeString()

	d.wrapper.NewsArticle(reqID, articleType, articleText)
}

func (c *Client) processHistoricalNewsMsg(msgDec *utils.MessageDecoder) error {

	requestID := msgDec.decodeInt64()

	time := msgDec.decodeString()
	providerCode := msgDec.decodeString()
	articleID := msgDec.decodeString()
	headline := msgDec.decodeString()

	d.wrapper.HistoricalNews(requestID, time, providerCode, articleID, headline)
}

func (c *Client) processHistoricalNewsEndMsg(msgDec *utils.MessageDecoder) error {

	requestID := msgDec.decodeInt64()

	hasMore := msgDec.decodeBool()

	d.wrapper.HistoricalNewsEnd(requestID, hasMore)
}
*/

func (c *Client) processHeadTimestampMsg(msgDec *utils.MessageDecoder) error {
	// Gets the originating ticker ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	_ = msgDec.String(false)
	// d.wrapper.HeadTimestamp(reqID, headTimestamp)

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processHistogramDataMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processRerouteMktDataReqMsg(msgDec *utils.MessageDecoder) error {

	reqID := msgDec.decodeInt64()

	conID := msgDec.decodeInt64()
	exchange := msgDec.decodeString()

	d.wrapper.RerouteMktDataReq(reqID, conID, exchange)
}

func (c *Client) processRerouteMktDepthReqMsg(msgDec *utils.MessageDecoder) error {

	reqID := msgDec.decodeInt64()

	conID := msgDec.decodeInt64()
	exchange := msgDec.decodeString()

	d.wrapper.RerouteMktDepthReq(reqID, conID, exchange)
}

func (c *Client) processMarketRuleMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processPnLMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processPnLSingleMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processHistoricalTicks(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.HistoricalTicksResponse)

		ticksCount := msgDec.Int64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}
		if ticksCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
			return false, msgDec.Err()
		}

		for i := int64(0); i < ticksCount; i++ {
			historicalTick := models.NewHistoricalTick()
			historicalTick.Time = msgDec.EpochTimestamp(false)
			msgDec.Skip()
			historicalTick.Price = msgDec.Float64(false)
			historicalTick.Size = models.NewDecimalFromMessageDecoder(msgDec, false)

			resp.Ticks = append(resp.Ticks, historicalTick)
		}

		done := msgDec.Bool()
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Done
		return done, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processHistoricalTicksBidAsk(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.HistoricalTicksResponse)

		ticksCount := msgDec.Int64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}
		if ticksCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
			return false, msgDec.Err()
		}

		for i := int64(0); i < ticksCount; i++ {
			historicalTickBidAsk := models.NewHistoricalTickBidAsk()
			historicalTickBidAsk.Time = msgDec.EpochTimestamp(false)
			mask := msgDec.Int64(false)
			historicalTickBidAsk.TickAttribBidAsk = models.NewTickAttribBidAsk()
			historicalTickBidAsk.TickAttribBidAsk.AskPastHigh = mask&1 != 0
			historicalTickBidAsk.TickAttribBidAsk.BidPastLow = mask&2 != 0
			historicalTickBidAsk.PriceBid = msgDec.Float64(false)
			historicalTickBidAsk.PriceAsk = msgDec.Float64(false)
			historicalTickBidAsk.SizeBid = models.NewDecimalFromMessageDecoder(msgDec, false)
			historicalTickBidAsk.SizeAsk = models.NewDecimalFromMessageDecoder(msgDec, false)

			resp.TicksBidAsk = append(resp.TicksBidAsk, historicalTickBidAsk)
		}

		done := msgDec.Bool()
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Done
		return done, nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) processHistoricalTicksLast(msgDec *utils.MessageDecoder) error {
	// Gets the originating request ID
	reqID := c.decodeRequestID(msgDec, false)
	if reqID == 0 {
		return msgDec.Err()
	}

	c.reqMgr.withRequestWithID(reqID, func(_resp interface{}) (bool, error) {
		resp := _resp.(*models.HistoricalTicksResponse)

		ticksCount := msgDec.Int64(false)
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}
		if ticksCount < 0 {
			msgDec.SetErr(fmt.Errorf("negative ticks count: %d", ticksCount))
			return false, msgDec.Err()
		}

		for i := int64(0); i < ticksCount; i++ {
			historicalTickLast := models.NewHistoricalTickLast()
			historicalTickLast.Time = msgDec.EpochTimestamp(false)
			mask := msgDec.Int64(false)
			historicalTickLast.TickAttribLast = models.NewTickAttribLast()
			historicalTickLast.TickAttribLast.PastLimit = mask&1 != 0
			historicalTickLast.TickAttribLast.Unreported = mask&2 != 0
			historicalTickLast.Price = msgDec.Float64(false)
			historicalTickLast.Size = models.NewDecimalFromMessageDecoder(msgDec, false)
			historicalTickLast.Exchange = msgDec.String(false)
			historicalTickLast.SpecialConditions = msgDec.String(false)

			resp.TicksLast = append(resp.TicksLast, historicalTickLast)
		}

		done := msgDec.Bool()
		if msgDec.Err() != nil {
			return false, msgDec.Err()
		}

		// Done
		return done, nil
	})

	// Done
	return msgDec.Err()
}

/*
func (c *Client) processTickByTickDataMsg(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()

		tickType := msgDec.decodeInt64()
		time := msgDec.decodeInt64()

		switch tickType {
		case 0: // None
		case 1, 2: // Last or AllLast
			price := msgDec.decodeFloat64()
			size := msgDec.decodeDecimal()
			mask := msgDec.decodeInt64()

			tickAttribLast := NewTickAttribLast()
			tickAttribLast.PastLimit = mask&1 != 0
			tickAttribLast.Unreported = mask&2 != 0

			exchange := msgDec.decodeString()
			specialConditions := msgDec.decodeString()

			d.wrapper.TickByTickAllLast(reqID, tickType, time, price, size, tickAttribLast, exchange, specialConditions)

		case 3: // BidAsk
			bidPrice := msgDec.decodeFloat64()
			askPrice := msgDec.decodeFloat64()
			bidSize := msgDec.decodeDecimal()
			askSize := msgDec.decodeDecimal()
			mask := msgDec.decodeInt64()

			tickAttribBidAsk := NewTickAttribBidAsk()
			tickAttribBidAsk.BidPastLow = mask&1 != 0
			tickAttribBidAsk.AskPastHigh = mask&2 != 0

			d.wrapper.TickByTickBidAsk(reqID, time, bidPrice, askPrice, bidSize, askSize, tickAttribBidAsk)

		case 4: // MidPoint
			midPoint := msgDec.decodeFloat64()

			d.wrapper.TickByTickMidPoint(reqID, time, midPoint)
		}
	}

func (c *Client) processOrderBoundMsg(msgDec *utils.MessageDecoder) error {

		permID := msgDec.decodeInt64()
		clientId := msgDec.decodeInt64()
		orderId := msgDec.decodeInt64()

		d.wrapper.OrderBound(permID, clientId, orderId)
	}

func (c *Client) processCompletedOrderMsg(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processReplaceFAEndMsg(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()
		text := msgDec.decodeString()

		d.wrapper.ReplaceFAEnd(reqID, text)
	}

func (c *Client) processWshMetaData(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()
		dataJSON := msgDec.decodeString()

		d.wrapper.WshMetaData(reqID, dataJSON)
	}

func (c *Client) processWshEventData(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()
		dataJSON := msgDec.decodeString()

		d.wrapper.WshEventData(reqID, dataJSON)
	}

func (c *Client) processHistoricalSchedule(msgDec *utils.MessageDecoder) error {

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

func (c *Client) processUserInfo(msgDec *utils.MessageDecoder) error {

		reqID := msgDec.decodeInt64()
		whiteBrandingId := msgDec.decodeString()

		d.wrapper.UserInfo(reqID, whiteBrandingId)
	}
*/

func (c *Client) processCurrentTimeInMillisMsg(msgDec *utils.MessageDecoder) error {
	c.reqMgr.withRequestWithoutID(common.REQ_CURRENT_TIME_IN_MILLIS, func(_resp interface{}) error {
		resp := _resp.(*models.CurrentTimeResponse)

		resp.CurrentTime = msgDec.EpochTimestamp(true)
		if msgDec.Err() != nil {
			return msgDec.Err()
		}

		// Done
		return nil
	})

	// Done
	return msgDec.Err()
}

func (c *Client) readLastTradeDate(msgDec *utils.MessageDecoder, contract *models.ContractDetails, isBond bool) {
	lastTradeDateOrContractMonth := msgDec.String(false) // YYYYMM or YYYYMMDD
	if len(lastTradeDateOrContractMonth) > 0 {
		var split []string

		if strings.Contains(lastTradeDateOrContractMonth, "-") {
			split = strings.Split(lastTradeDateOrContractMonth, "-")
		} else {
			split = strings.Split(lastTradeDateOrContractMonth, " ")
		}
		if len(split) > 0 {
			if isBond {
				contract.Maturity = split[0]
			} else {
				contract.Contract.LastTradeDateOrContractMonth = split[0]
			}

			if len(split) > 1 {
				contract.LastTradeTime = split[1]
			}
			if isBond && len(split) > 2 {
				contract.TimeZoneID = split[2]
			}
		}
	}
}
