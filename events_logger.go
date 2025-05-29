package ibkr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type EventsLogger struct {
	cb EventsLoggerCallback
}

type EventsLoggerCallback func(msg string)

type eventsLoggerBuilder struct {
	cb        EventsLoggerCallback
	sb        strings.Builder
	atNewLine int
}

// -----------------------------------------------------------------------------

func NewIncomingMessageLogger(cb EventsLoggerCallback) Events {
	return &EventsLogger{
		cb: cb,
	}
}

func (el *EventsLogger) ConnectionClosed(err error) {
	el.build().
		str("err", err.Error()).
		msg("<ConnectionClosed>")
}

func (el *EventsLogger) ReceivedUnknownMessage(id int) {
	el.build().
		int64("MsgID", int64(id)).
		msg("<ReceivedUnknownMessage>")
}

func (el *EventsLogger) Error(ts time.Time, code int, message string, advancedOrderRejectJson string) {
	logger := el.build().
		str("Timestamp", ts.Format("2006-01-02 15:04:05")).
		int64("Code", int64(code)).
		str("Message", message)
	if len(advancedOrderRejectJson) > 0 {
		logger = logger.str("AdvancedOrderRejectJson", advancedOrderRejectJson)
	}
	logger.msg("<Error>")
}

func (el *EventsLogger) TickPrice(reqID models.TickerID, tickType models.TickType, price float64, attrib models.TickAttrib) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		str("Price", utils.FloatMaxString(price)).
		bool("CanAutoExecute", attrib.CanAutoExecute).
		bool("PastLimit", attrib.PastLimit).
		bool("PreOpen", attrib.PreOpen).
		msg("<TickPrice>")
}

func (el *EventsLogger) TickSize(reqID models.TickerID, tickType models.TickType, size models.Decimal) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		str("Size", size.StringMax()).
		msg("<TickSize>")
}

func (el *EventsLogger) TickOptionComputation(reqID models.TickerID, tickType models.TickType, tickAttrib int64, impliedVol float64, delta float64, optPrice float64, pvDividend float64, gamma float64, vega float64, theta float64, undPrice float64) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		str("TickAttrib", utils.IntMaxString(tickAttrib)).
		str("ImpliedVol", utils.FloatMaxString(impliedVol)).
		str("Delta", utils.FloatMaxString(delta)).
		str("OptPrice", utils.FloatMaxString(optPrice)).
		str("PvDividend", utils.FloatMaxString(pvDividend)).
		str("Gamma", utils.FloatMaxString(gamma)).
		str("Vega", utils.FloatMaxString(vega)).
		str("Theta", utils.FloatMaxString(theta)).
		str("UndPrice", utils.FloatMaxString(undPrice)).
		msg("<TickOptionComputation>")
}

func (el *EventsLogger) TickGeneric(reqID models.TickerID, tickType models.TickType, value float64) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		str("Value", utils.FloatMaxString(value)).
		msg("<TickGeneric>")
}

func (el *EventsLogger) TickString(reqID models.TickerID, tickType models.TickType, value string) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		str("Value", value).
		msg("<TickString>")
}

func (el *EventsLogger) TickEFP(reqID models.TickerID, tickType models.TickType, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", int64(tickType)).
		float64("BasisPoints", basisPoints).
		str("FormattedBasisPoints", formattedBasisPoints).
		float64("TotalDividends", totalDividends).
		int64("holdDays", holdDays).
		str("FutureLastTradeDate", futureLastTradeDate).
		float64("DividendImpact", dividendImpact).
		float64("DividendsToLastTradeDate", dividendsToLastTradeDate).
		msg("<TickEFP>")
}

func (el *EventsLogger) OrderStatus(orderID models.OrderID, status string, filled models.Decimal, remaining models.Decimal, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64) {
	logger := el.build().
		int64("OrderID", orderID).
		str("Status", status).
		str("Filled", filled.StringMax()).
		stringer("Remaining", &remaining).
		float64("AvgFillPrice", avgFillPrice).
		int64("PermID", permID).
		int64("ParentID", parentID).
		float64("LastFillPrice", lastFillPrice).
		int64("ClientID", clientID).
		str("WhyHeld", whyHeld).
		float64("MktCapPrice", mktCapPrice)
	logger.msg("<OrderStatus>")
}

func (el *EventsLogger) OpenOrder(orderID models.OrderID, contract *models.Contract, order *models.Order, orderState *models.OrderState) {
	logger := el.build().
		str("PermID", utils.IntMaxString(order.PermID)).
		str("ClientID", utils.IntMaxString(order.ClientID)).
		str("OrderID", utils.IntMaxString(order.OrderID)).
		str("Account", order.Account).
		str("Symbol", contract.Symbol).
		str("SecType", string(contract.SecType)).
		str("Exchange", contract.Exchange).
		float64("Strike", contract.Strike).
		str("Action", order.Action).
		str("OrderType", order.OrderType).
		str("TotalQuantity", order.TotalQuantity.StringMax()).
		str("CashQty", utils.FloatMaxString(order.CashQty)).
		str("LmtPrice", utils.FloatMaxString(order.LmtPrice)).
		str("AuxPrice", utils.FloatMaxString(order.AuxPrice)).
		str("Status", orderState.Status).
		str("MinTradeQty", utils.IntMaxString(order.MinTradeQty)).
		str("MinCompeteSize", utils.IntMaxString(order.MinCompeteSize))
	if order.CompeteAgainstBestOffset == models.COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID {
		logger.str("CompeteAgainstBestOffset", "UpToMid")
	} else {
		logger.str("CompeteAgainstBestOffset", utils.FloatMaxString(order.CompeteAgainstBestOffset))
	}
	logger.str("MidOffsetAtWhole", utils.FloatMaxString(order.MidOffsetAtWhole)).
		str("MidOffsetAtHalf", utils.FloatMaxString(order.MidOffsetAtHalf)).
		str("FAGroup", order.FAGroup).
		str("CustomerAccount", order.CustomerAccount).
		bool("ProfessionalCustomer", order.ProfessionalCustomer).
		str("ManualOrderIndicator", utils.IntMaxString(order.ManualOrderIndicator)).
		str("Submitter", order.Submitter).
		bool("ImbalanceOnly", order.ImbalanceOnly).
		msg("<OpenOrder>")
}

func (el *EventsLogger) OpenOrdersEnd() {
	el.build().
		msg("<OpenOrdersEnd>")
}

func (el *EventsLogger) WinError(text string, lastError int64) {
	el.build().
		str("Text", text).
		int64("LastError", lastError).
		msg("<WinError>")
}

func (el *EventsLogger) UpdateAccountValue(tag string, value string, currency string, accountName string) {
	el.build().
		str("Tag", tag).
		str("Value", value).
		str("Currency", currency).
		str("AccountName", accountName).
		msg("<UpdateAccountValue>")
}

func (el *EventsLogger) UpdatePortfolio(contract *models.Contract, position models.Decimal, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accountName string) {
	el.build().
		str("Symbol", contract.Symbol).
		str("SecType", string(contract.SecType)).
		str("Exchange", contract.Exchange).
		str("Position", position.StringMax()).
		str("MarketPrice", utils.FloatMaxString(marketPrice)).
		str("MarketValue", utils.FloatMaxString(marketValue)).
		str("AverageCost", utils.FloatMaxString(averageCost)).
		str("UnrealizedPNL", utils.FloatMaxString(unrealizedPNL)).
		str("RealizedPNL", utils.FloatMaxString(realizedPNL)).
		str("AccountName", accountName).
		msg("<UpdatePortfolio>")
}

func (el *EventsLogger) UpdateAccountTime(timeStamp string) {
	el.build().
		str("TimeStamp", timeStamp).
		msg("<UpdateAccountTime>")
}

func (el *EventsLogger) AccountDownloadEnd(accountName string) {
	el.build().
		str("AccountName", accountName).
		msg("<AccountDownloadEnd>")
}

func (el *EventsLogger) ContractDetails(reqID int64, contractDetails *models.ContractDetails) {
	el.build().
		int64("ReqID", reqID).
		stringer("ContractDetails", contractDetails).
		msg("<ContractDetails>")
}

func (el *EventsLogger) BondContractDetails(reqID int64, contractDetails *models.ContractDetails) {
	el.build().
		int64("ReqID", reqID).
		stringer("ContractDetails", contractDetails).
		msg("<BondContractDetails>")
}

func (el *EventsLogger) ContractDetailsEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<ContractDetailsEnd>")
}

func (el *EventsLogger) ExecDetails(reqID int64, contract *models.Contract, execution *models.Execution) {
	el.build().
		int64("ReqID", reqID).
		stringer("Contract", contract).
		stringer("Execution", execution).
		msg("<ExecDetails>")
}

func (el *EventsLogger) ExecDetailsEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<ExecDetailsEnd>")
}

func (el *EventsLogger) UpdateMktDepth(TickerID models.TickerID, position int64, operation int64, side int64, price float64, size models.Decimal) {
	el.build().
		int64("TickerID", TickerID).
		int64("Position", position).
		int64("Operation", operation).
		int64("Side", side).
		str("Price", utils.FloatMaxString(price)).
		str("Size", size.StringMax()).
		msg("<UpdateMktDepth>")
}

func (el *EventsLogger) UpdateMktDepthL2(TickerID models.TickerID, position int64, marketMaker string, operation int64, side int64, price float64, size models.Decimal, isSmartDepth bool) {
	el.build().
		int64("TickerID", TickerID).
		int64("Position", position).
		str("MarketMaker", marketMaker).
		int64("Operation", operation).
		int64("Side", side).
		str("Price", utils.FloatMaxString(price)).
		str("Size", size.StringMax()).
		bool("IsSmartDepth", isSmartDepth).
		msg("<UpdateMktDepthL2>")
}

func (el *EventsLogger) ReceiveFA(faDataType models.FaData, cxml string) {
	el.build().
		stringer("FaDataType", faDataType).
		str("Cxml", cxml).
		msg("<ReceiveFA>")
}

func (el *EventsLogger) ScannerParameters(xml string) {
	el.build().
		str("Xml", xml[:50]).
		msg("<ScannerParameters>")
}

func (el *EventsLogger) ScannerData(reqID int64, rank int64, contractDetails *models.ContractDetails, distance string, benchmark string, projection string, legsStr string) {
	el.build().
		int64("ReqID", reqID).
		int64("Rank", rank).
		stringer("ContractDetails", contractDetails).
		str("Distance", distance).
		str("Benchmark", benchmark).
		str("Projection", projection).
		str("LegsStr", legsStr).
		msg("<ScannerData>")
}

func (el *EventsLogger) ScannerDataEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<ScannerDataEnd>")
}

func (el *EventsLogger) RealtimeBar(reqID int64, time int64, open float64, high float64, low float64, close float64, volume models.Decimal, wap models.Decimal, count int64) {
	el.build().
		int64("ReqID", reqID).
		int64("Bar time", time).
		float64("Open", open).
		float64("High", high).
		float64("Low", low).
		float64("Close", close).
		stringer("Volume", &volume).
		stringer("Wap", &wap).
		int64("Count", count).
		msg("<RealtimeBar>")
}

func (el *EventsLogger) CurrentTime(t time.Time) {
	el.build().
		time("Server Time", t).
		msg("<CurrentTime>")
}

func (el *EventsLogger) FundamentalData(reqID int64, data string) {
	el.build().
		int64("ReqID", reqID).
		str("Data", data).
		msg("<FundamentalData>")
}

func (el *EventsLogger) DeltaNeutralValidation(reqID int64, deltaNeutralContract models.DeltaNeutralContract) {
	el.build().
		int64("ReqID", reqID).
		stringer("DeltaNeutralContract", deltaNeutralContract).
		msg("<DeltaNeutralValidation>")
}

func (el *EventsLogger) TickSnapshotEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<TickSnapshotEnd>")
}

func (el *EventsLogger) MarketDataType(reqID int64, marketDataType int64) {
	el.build().
		int64("ReqID", reqID).
		int64("MarketDataType", marketDataType).
		msg("<MarketDataType>")
}

func (el *EventsLogger) CommissionAndFeesReport(commissionAndFeesReport models.CommissionAndFeesReport) {
	el.build().
		stringer("CommissionAndFeesReport", commissionAndFeesReport).
		msg("<CommissionAndFeesReport>")
}

func (el *EventsLogger) Position(account string, contract *models.Contract, position models.Decimal, avgCost float64) {
	el.build().
		str("Account", account).
		stringer("Contract", contract).
		str("Position", position.StringMax()).
		str("AvgCost", utils.FloatMaxString(avgCost)).
		msg("<Position>")
}

func (el *EventsLogger) PositionEnd() {
	el.build().
		msg("<PositionEnd>")
}

func (el *EventsLogger) AccountSummary(reqID int64, account string, tag string, value string, currency string) {
	el.build().
		int64("ReqID", reqID).
		str("Account", account).
		str("Tag", tag).
		str("Value", value).
		str("Currency", currency).
		msg("<AccountSummary>")
}

func (el *EventsLogger) AccountSummaryEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<AccountSummaryEnd>")
}

func (el *EventsLogger) VerifyMessageAPI(apiData string) {
	el.build().
		str("ApiData", apiData).
		msg("<VerifyMessageAPI>")
}

func (el *EventsLogger) VerifyCompleted(isSuccessful bool, errorText string) {
	el.build().
		bool("IsSuccessful", isSuccessful).
		str("ErrorText", errorText).
		msg("<VerifyCompleted>")
}

func (el *EventsLogger) DisplayGroupList(reqID int64, groups string) {
	el.build().
		int64("ReqID", reqID).
		str("Groups", groups).
		msg("<DisplayGroupList>")
}

func (el *EventsLogger) DisplayGroupUpdated(reqID int64, contractInfo string) {
	el.build().
		int64("ReqID", reqID).
		str("ContractInfo", contractInfo).
		msg("<DisplayGroupUpdated>")
}

func (el *EventsLogger) VerifyAndAuthMessageAPI(apiData string, xyzChallange string) {
	el.build().
		str("ApiData", apiData).
		str("XyzChallange", xyzChallange).
		msg("<VerifyAndAuthMessageAPI>")
}

func (el *EventsLogger) VerifyAndAuthCompleted(isSuccessful bool, errorText string) {
	el.build().
		bool("IsSuccessful", isSuccessful).
		str("ErrorText", errorText).
		msg("<VerifyAndAuthCompleted>")
}

func (el *EventsLogger) ConnectAck() {
	el.build().
		msg("<ConnectAck>...")
}

func (el *EventsLogger) PositionMulti(reqID int64, account string, modelCode string, contract *models.Contract, pos models.Decimal, avgCost float64) {
	el.build().
		int64("ReqID", reqID).
		str("Account", account).
		str("ModelCode", modelCode).
		stringer("Contract", contract).
		str("Position", pos.StringMax()).
		str("AvgCost", utils.FloatMaxString(avgCost)).
		msg("<PositionMulti>")
}

func (el *EventsLogger) PositionMultiEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<PositionMultiEnd>")
}

func (el *EventsLogger) AccountUpdateMulti(reqID int64, account string, modelCode string, key string, value string, currency string) {
	el.build().
		int64("ReqID", reqID).
		str("Account", account).
		str("ModelCode", modelCode).
		str("Key", key).
		str("Value", value).
		str("Currency", currency).
		msg("<AccountUpdateMulti>")
}

func (el *EventsLogger) AccountUpdateMultiEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<AccountUpdateMultiEnd>")
}

func (el *EventsLogger) SecurityDefinitionOptionParameter(reqID int64, exchange string, underlyingConID int64, tradingClass string, multiplier string, expirations []string, strikes []float64) {
	el.build().
		int64("ReqID", reqID).
		str("Exchange", exchange).
		str("UnderlyingConID", utils.IntMaxString(underlyingConID)).
		str("TradingClass", tradingClass).
		str("Multiplier", multiplier).
		strArray("Expirations", expirations).
		float64Array("Strikes", strikes).
		msg("<SecurityDefinitionOptionParameter>")
}

func (el *EventsLogger) SecurityDefinitionOptionParameterEnd(reqID int64) {
	el.build().
		int64("ReqID", reqID).
		msg("<SecurityDefinitionOptionParameterEnd>")
}

func (el *EventsLogger) SoftDollarTiers(reqID int64, tiers []models.SoftDollarTier) {
	tiersStr := make([]string, 0, len(tiers))
	for _, sdt := range tiers {
		tiersStr = append(tiersStr, sdt.String())
	}
	el.build().
		int64("ReqID", reqID).
		strArray("Tiers", tiersStr).
		msg("<SoftDollarTiers>")
}

func (el *EventsLogger) FamilyCodes(familyCodes []models.FamilyCode) {
	codesStr := make([]string, 0, len(familyCodes))
	for _, fc := range familyCodes {
		codesStr = append(codesStr, fc.String())
	}
	el.build().
		strArray("Codes", codesStr).
		msg("<FamilyCodes>")
}

func (el *EventsLogger) SymbolSamples(reqID int64, contractDescriptions []models.ContractDescription) {
	descsStr := make([]string, 0, len(contractDescriptions))
	for _, cd := range contractDescriptions {
		descsStr = append(descsStr, cd.Contract.String())
	}
	el.build().
		int64("Nb samples", int64(len(contractDescriptions))).
		int64("ReqID", reqID).
		strArray("Contracts", descsStr).
		msg("<SymbolSamples>")
}

func (el *EventsLogger) MktDepthExchanges(depthMktDataDescriptions []models.DepthMktDataDescription) {
	el.build().
		any("DepthMktDataDescriptions", depthMktDataDescriptions).
		msg("<MktDepthExchanges>")
}

func (el *EventsLogger) SmartComponents(reqID int64, smartComponents []models.SmartComponent) {
	compsStr := make([]string, 0, len(smartComponents))
	for _, sc := range smartComponents {
		compsStr = append(compsStr, sc.String())
	}
	el.build().
		int64("ReqID", reqID).
		strArray("Components", compsStr).
		msg("<SmartComponents>")
}

func (el *EventsLogger) TickReqParams(TickerID models.TickerID, minTick float64, bboExchange string, snapshotPermissions int64) {
	el.build().
		int64("TickerID", TickerID).
		str("MinTick", utils.FloatMaxString(minTick)).
		str("BboExchange", bboExchange).
		str("SnapshotPermissions", utils.IntMaxString(snapshotPermissions)).
		msg("<TickReqParams>")
}

func (el *EventsLogger) HeadTimestamp(reqID int64, headTimestamp string) {
	el.build().
		int64("ReqID", reqID).
		str("HeadTimestamp", headTimestamp).
		msg("<HeadTimestamp>")
}

func (el *EventsLogger) HistogramData(reqID int64, data []models.HistogramData) {
	el.build().
		int64("ReqID", reqID).any("Data", data).
		msg("<HistogramData>")
}

func (el *EventsLogger) RerouteMktDataReq(reqID int64, conID int64, exchange string) {
	el.build().
		int64("ReqID", reqID).
		int64("ConID", conID).
		str("Exchange", exchange).
		msg("<RerouteMktDataReq>")
}

func (el *EventsLogger) RerouteMktDepthReq(reqID int64, conID int64, exchange string) {
	el.build().
		int64("ReqID", reqID).
		int64("ConID", conID).
		str("Exchange", exchange).
		msg("<RerouteMktDepthReq>")
}

func (el *EventsLogger) MarketRule(marketRuleID int64, priceIncrements []models.PriceIncrement) {
	el.build().
		int64("MarketRuleID", marketRuleID).any("PriceIncrements", priceIncrements).
		msg("<MarketRule>")
}

func (el *EventsLogger) Pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64) {
	el.build().
		int64("ReqID", reqID).
		str("DailyPnL", utils.FloatMaxString(dailyPnL)).
		str("UnrealizedPnL", utils.FloatMaxString(unrealizedPnL)).
		str("RealizedPnL", utils.FloatMaxString(realizedPnL)).
		msg("<Pnl>")
}

func (el *EventsLogger) PnlSingle(reqID int64, pos models.Decimal, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64) {
	el.build().
		int64("ReqID", reqID).
		str("Position", pos.StringMax()).
		str("DailyPnL", utils.FloatMaxString(dailyPnL)).
		str("UnrealizedPnL", utils.FloatMaxString(unrealizedPnL)).
		str("RealizedPnL", utils.FloatMaxString(realizedPnL)).
		str("Value", utils.FloatMaxString(value)).
		msg("<PnlSingle>")
}

func (el *EventsLogger) TickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size models.Decimal, tickAttribLast models.TickAttribLast, exchange string, specialConditions string) {
	el.build().
		int64("ReqID", reqID).
		int64("TickType", tickType).
		int64("Tick time", time).
		str("Price", utils.FloatMaxString(price)).
		str("Size", size.StringMax()).
		bool("PastLimit", tickAttribLast.PastLimit).
		bool("Unreported", tickAttribLast.Unreported).
		str("Exchange", exchange).
		str("SpecialConditions", specialConditions).
		msg("<TickByTickAllLast>")
}

func (el *EventsLogger) TickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize models.Decimal, askSize models.Decimal, tickAttribBidAsk models.TickAttribBidAsk) {
	el.build().
		int64("ReqID", reqID).
		int64("Tick time", time).
		str("BidPrice", utils.FloatMaxString(bidPrice)).
		str("AskPrice", utils.FloatMaxString(askPrice)).
		str("BidSize", bidSize.StringMax()).
		str("AskSize", askSize.StringMax()).
		bool("AskPastHigh", tickAttribBidAsk.AskPastHigh).
		bool("BidPastLow", tickAttribBidAsk.BidPastLow).
		msg("<TickByTickBidAsk>")
}

func (el *EventsLogger) TickByTickMidPoint(reqID int64, time int64, midPoint float64) {
	el.build().
		int64("ReqID", reqID).
		int64("Tick time", time).
		str("MidPoint", utils.FloatMaxString(midPoint)).
		msg("<TickByTickMidPoint>")
}

func (el *EventsLogger) OrderBound(permID int64, clientID int64, orderID int64) {
	el.build().
		str("PermID", utils.IntMaxString(permID)).
		str("ClientID", utils.IntMaxString(clientID)).
		str("OrderID", utils.IntMaxString(orderID)).
		msg("<OrderBound>")
}

func (el *EventsLogger) CompletedOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState) {
	logger := el.build().
		str("Account", order.Account).
		str("PermID", utils.IntMaxString(order.PermID)).
		str("ParentPermID", utils.IntMaxString(order.ParentPermID)).
		str("Symbol", contract.Symbol).
		str("SecType", string(contract.SecType)).
		str("Exchange", contract.Exchange).
		str("Action", order.Action).
		str("OrderType", order.OrderType).
		str("TotalQuantity", order.TotalQuantity.StringMax()).
		str("CashQty", utils.FloatMaxString(order.CashQty)).
		str("FilledQuantity", order.FilledQuantity.StringMax()).
		str("LmtPrice", utils.FloatMaxString(order.LmtPrice)).
		str("AuxPrice", utils.FloatMaxString(order.AuxPrice)).
		str("Status", orderState.Status).
		str("CompletedTime", orderState.CompletedTime).
		str("CompletedStatus", orderState.CompletedStatus).
		str("MinTradeQty", utils.IntMaxString(order.MinTradeQty)).
		str("MinCompeteSize", utils.IntMaxString(order.MinCompeteSize))
	if order.CompeteAgainstBestOffset == models.COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID {
		logger.str("CompeteAgainstBestOffset", "UpToMid")
	} else {
		logger.str("CompeteAgainstBestOffset", utils.FloatMaxString(order.CompeteAgainstBestOffset))
	}
	logger.str("MidOffsetAtWhole", utils.FloatMaxString(order.MidOffsetAtWhole)).
		str("MidOffsetAtHalf", utils.FloatMaxString(order.MidOffsetAtHalf)).
		str("CustomerAccount", order.CustomerAccount).
		bool("ProfessionalCustomer", order.ProfessionalCustomer).
		str("Submitter", order.Submitter).
		bool("ImbalanceOnly", order.ImbalanceOnly).
		msg("<CompletedOrder>")
}

func (el *EventsLogger) CompletedOrdersEnd() {
	el.build().
		msg("<CompletedOrdersEnd>")
}

func (el *EventsLogger) ReplaceFAEnd(reqID int64, text string) {
	el.build().
		int64("ReqID", reqID).
		str("Text", text).
		msg("<ReplaceFAEnd>")
}

func (el *EventsLogger) WshMetaData(reqID int64, dataJson string) {
	el.build().
		int64("ReqID", reqID).
		str("DataJson", dataJson).
		msg("<WshMetaData>")
}

func (el *EventsLogger) WshEventData(reqID int64, dataJson string) {
	el.build().
		int64("ReqID", reqID).
		str("DataJson", dataJson).
		msg("<WshEventData>")
}

func (el *EventsLogger) HistoricalSchedule(reqID int64, startDarteTime, endDateTime, timeZone string, sessions []models.HistoricalSession) {
	el.build().
		int64("ReqID", reqID).
		str("StartDarteTime", startDarteTime).
		str("EndDateTime", endDateTime).
		str("TimeZone", timeZone).
		msg("<HistoricalSchedule>")
}

func (el *EventsLogger) UserInfo(reqID int64, whiteBrandingId string) {
	el.build().
		int64("ReqID", reqID).
		str("WhiteBrandingId", whiteBrandingId).
		msg("<UserInfo>")
}

func (el *EventsLogger) build() *eventsLoggerBuilder {
	b := &eventsLoggerBuilder{
		cb:        el.cb,
		sb:        strings.Builder{},
		atNewLine: -1,
	}
	return b
}

func (b *eventsLoggerBuilder) int64(name string, value int64) *eventsLoggerBuilder {
	b.addSeparator()
	_, _ = b.sb.WriteString(name)
	_, _ = b.sb.WriteRune('=')
	_, _ = b.sb.Write(strconv.AppendInt(nil, value, 10))
	return b
}

func (b *eventsLoggerBuilder) float64(name string, value float64) *eventsLoggerBuilder {
	b.addSeparator()
	_, _ = b.sb.WriteString(name)
	_, _ = b.sb.WriteRune('=')
	_, _ = b.sb.Write(b.float2str(value))
	return b
}

func (b *eventsLoggerBuilder) float64Array(name string, value []float64) *eventsLoggerBuilder {
	b.addSeparator()
	_, _ = b.sb.WriteString(name)
	_, _ = b.sb.Write([]byte("=["))
	for i, v := range value {
		if i > 0 {
			_, _ = b.sb.WriteRune(',')
		}
		_, _ = b.sb.Write(b.float2str(v))

	}
	_, _ = b.sb.WriteRune(']')
	return b
}

func (b *eventsLoggerBuilder) str(name string, value string) *eventsLoggerBuilder {
	b.addSeparator()
	_, _ = b.sb.WriteString(name)
	_, _ = b.sb.WriteRune('=')
	_, _ = b.sb.WriteString(value)
	return b
}

func (b *eventsLoggerBuilder) strArray(name string, value []string) *eventsLoggerBuilder {
	return b.str(name, "["+strings.Join(value, ",")+"]")
}

func (b *eventsLoggerBuilder) stringer(name string, value fmt.Stringer) *eventsLoggerBuilder {
	return b.str(name, value.String())
}

func (b *eventsLoggerBuilder) bool(name string, value bool) *eventsLoggerBuilder {
	if value {
		return b.str(name, "true")
	}
	return b.str(name, "false")
}

func (b *eventsLoggerBuilder) time(name string, value time.Time) *eventsLoggerBuilder {
	return b.str(name, value.Format(time.RFC3339))
}

func (b *eventsLoggerBuilder) any(name string, value interface{}) *eventsLoggerBuilder {
	if value == nil {
		return b.str(name, "<nil>")
	}
	switch v := value.(type) {
	case fmt.Stringer:
		return b.str(name, v.String())
	case *fmt.Stringer:
		return b.str(name, (*v).String())
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("    ", "  ")
	err := enc.Encode(value)
	if err == nil {
		b.newLine()
		_, _ = b.sb.WriteString(name)
		_, _ = b.sb.WriteRune('=')
		_, _ = b.sb.Write(bytes.TrimRight(buf.Bytes(), "\r\n"))
	}
	return b
}

func (b *eventsLoggerBuilder) addSeparator() {
	if b.atNewLine <= 0 {
		_, _ = b.sb.Write([]byte(" | "))
		b.atNewLine = 0
	} else {
		_, _ = b.sb.Write([]byte("    "))
		b.atNewLine = 0
	}
}

func (b *eventsLoggerBuilder) newLine() *eventsLoggerBuilder {
	if b.atNewLine <= 0 {
		_, _ = b.sb.WriteRune('\n')
		b.atNewLine = 1
	}
	return b
}

func (b *eventsLoggerBuilder) msg(m string) {
	s := "IB" + m + b.sb.String()
	b.cb(s)
}

func (b *eventsLoggerBuilder) msgf(format string, args ...interface{}) {
	b.msg(fmt.Sprintf(format, args...))
}

func (b *eventsLoggerBuilder) float2str(val float64) []byte {
	switch {
	case math.IsNaN(val):
		return []byte("NaN")
	case math.IsInf(val, 1):
		return []byte("+Inf")
	case math.IsInf(val, -1):
		return []byte("-Inf")
	}

	strFmt := byte('f')
	if abs := math.Abs(val); abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			strFmt = 'e'
		}
	}
	dst := strconv.AppendFloat(nil, val, strFmt, -1, 64)
	if strFmt == 'e' {
		// Clean up e-09 to e-9
		n := len(dst)
		if n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
			dst[n-2] = dst[n-1]
			dst = dst[:n-1]
		}
	}
	return dst
}
