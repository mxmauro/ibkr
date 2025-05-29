package ibkr

import (
	"time"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type Events interface {
	ConnectionClosed(err error)
	ReceivedUnknownMessage(id int)

	Error(ts time.Time, code int, message string, advancedOrderRejectJson string)

	// TickPrice handles all price related ticks. Every tickPrice callback is followed by a tickSize.
	// A tickPrice value of -1 or 0 followed by a tickSize of 0 indicates there is no data for this field currently available, whereas a tickPrice with a positive tickSize indicates an active quote of 0 (typically for a combo contract).
	TickPrice(reqID models.TickerID, tickType models.TickType, price float64, attrib models.TickAttrib)
	// TickSize handles all size related ticks.
	TickSize(reqID models.TickerID, tickType models.TickType, size models.Decimal)
	// TickOptionComputation is called when the market in an option or its underlier moves.
	// TWS's option model volatilities, prices, and deltas, along with the present value of dividends expected on that options underlier are received.
	TickOptionComputation(reqID models.TickerID, tickType models.TickType, tickAttrib int64, impliedVol float64, delta float64, optPrice float64, pvDividend float64, gamma float64, vega float64, theta float64, undPrice float64)
	// TickGeneric .
	TickGeneric(reqID models.TickerID, tickType models.TickType, value float64)
	// TickString .
	TickString(reqID models.TickerID, tickType models.TickType, value string)
	// TickEFP handles market for Exchange for Physical.
	// tickerId is the request's identifier.
	// tickType is the type of tick being received.
	// basisPoints is the annualized basis points, which is representative of the financing rate that can be directly compared to broker rates.
	// formattedBasisPoints is the annualized basis points as a formatted string that depicts them in percentage form.
	// impliedFuture is the implied Futures price.
	// holdDays is the number of hold days until the lastTradeDate of the EFP.
	// futureLastTradeDate is the expiration date of the single stock future.
	// dividendImpact is the dividend impact upon the annualized basis points interest rate.
	// dividendsToLastTradeDate is the dividends expected until the expiration of the single stock future.
	TickEFP(reqID models.TickerID, tickType models.TickType, basisPoints float64, formattedBasisPoints string, totalDividends float64, holdDays int64, futureLastTradeDate string, dividendImpact float64, dividendsToLastTradeDate float64)
	// OrderStatus is called whenever the status of an order changes.
	// It is also fired after reconnecting to TWS if the client has any open orders.
	// OrderID is the order ID that was specified previously in the	call to placeOrder().
	// status is the order status. Possible values include:
	//		PendingSubmit - indicates that you have transmitted the order, but have not  yet received confirmation that it has been accepted by the order destination. NOTE: This order status is not sent by TWS and should be explicitly set by the API developer when an order is submitted.
	//		PendingCancel - indicates that you have sent a request to cancel the order but have not yet received cancel confirmation from the order destination. At this point, your order is not confirmed canceled. You may still receive an execution while your cancellation request is pending. NOTE: This order status is not sent by TWS and should be explicitly set by the API developer when an order is canceled.
	//		PreSubmitted - indicates that a simulated order type has been accepted by the IB system and that this order has yet to be elected. The order is held in the IB system until the election criteria are met. At that time the order is transmitted to the order destination as specified.
	//		Submitted - indicates that your order has been accepted at the order destination and is working.
	//		Cancelled - indicates that the balance of your order has been confirmed canceled by the IB system. This could occur unexpectedly when IB or the destination has rejected your order.
	//		Filled - indicates that the order has been completely filled.
	//		Inactive - indicates that the order has been accepted by the system (simulated orders) or an exchange (native orders) but that currently the order is inactive due to system, exchange or other issues.
	// filled specifies the number of shares that have been executed. For more information about partial fills, see Order Status for Partial Fills.
	// remaining specifies the number of shares still outstanding.
	// avgFillPrice is the average price of the shares that have been executed. This parameter is valid only if the filled parameter value is greater than zero. Otherwise, the price parameter will be zero.
	// permId is the TWS id used to identify orders. Remains the same over TWS sessions.
	// parentId is the order ID of the parent order, used for bracket and auto trailing stop orders.
	// lastFilledPrice is the last price of the shares that have been executed. This parameter is valid only if the filled parameter value is greater than zero. Otherwise, the price parameter will be zero.
	// clientId is the ID of the client (or TWS) that placed the order. Note that TWS orders have a fixed clientId and OrderID of 0 that distinguishes them from API orders.
	// whyHeld is the field used to identify an order held when TWS is trying to locate shares for a short sell. The value used to indicate this is 'locate'.
	OrderStatus(orderID models.OrderID, status string, filled models.Decimal, remaining models.Decimal, avgFillPrice float64, permID int64, parentID int64, lastFillPrice float64, clientID int64, whyHeld string, mktCapPrice float64)
	// OpenOrder is called to feed in open orders.
	// orderID: OrderId - The order ID assigned by TWS. Use to cancel or update TWS order.
	// contract: Contract - The Contract class attributes describe the contract.
	// order: Order - The Order class gives the details of the open order.
	// orderState: OrderState - The orderState class includes attributes Used for both pre and post trade margin and commission and fees data.
	OpenOrder(orderID models.OrderID, contract *models.Contract, order *models.Order, orderState *models.OrderState)
	// OpenOrdersEnd is called at the end of a given request for open orders.
	OpenOrdersEnd()
	// WinError .
	WinError(text string, lastError int64)

	// UpdateAccountValue is called only when reqAccountUpdates() has been called.
	UpdateAccountValue(tag string, val string, currency string, accountName string)
	// UpdatePortfolio is called only when reqAccountUpdates() has been called.
	UpdatePortfolio(contract *models.Contract, position models.Decimal, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accountName string)
	// UpdateAccountTime .
	UpdateAccountTime(timeStamp string)
	// AccountDownloadEnd is called after a batch updateAccountValue() and updatePortfolio() is sent.
	AccountDownloadEnd(accountName string)

	// ContractDetails Receives the full contract's definitions. This method will return all contracts matching the requested via reqContractDetails().
	// For example, one can obtain the whole option chain with it.
	ContractDetails(reqID int64, contractDetails *models.ContractDetails)
	// BondContractDetails is called when reqContractDetails function has been called for bonds.
	BondContractDetails(reqID int64, contractDetails *models.ContractDetails)
	// ContractDetailsEnd is called once all contract details for a given request are received.
	// This helps to define the end of an option chain.
	ContractDetailsEnd(reqID int64)
	// ExecDetails is called when the reqExecutions() functions is invoked, or when an order is filled.
	ExecDetails(reqID int64, contract *models.Contract, execution *models.Execution)
	// ExecDetailsEnd is called once all executions have been sent to a client in response to reqExecutions().
	ExecDetailsEnd(reqID int64)
	// Error is called when there is an error with the communication or when TWS wants to send a message to the client.
	// UpdateMktDepth returns the order book.
	// 	TickerID -  the request's identifier.
	// 	position -  the order book's row being updated.
	// 	operation - how to refresh the row:
	// 		0 = insert (insert this new order into the row identified by 'position').
	// 		1 = update (update the existing order in the row identified by 'position').
	// 		2 = delete (delete the existing order at the row identified by 'position').
	// 	side -  0 for ask, 1 for bid.
	// 	price - the order's price.
	// 	size -  the order's size.
	UpdateMktDepth(TickerID models.TickerID, position int64, operation int64, side int64, price float64, size models.Decimal)
	// UpdateMktDepthL2 returns the order book.
	// 	TickerID -  the request's identifier.
	//  position -  the order book's row being updated.
	//  marketMaker - the exchange holding the order.
	//  operation - how to refresh the row:
	//  	0 = insert (insert this new order into the row identified by 'position').
	//      1 = update (update the existing order in the row identified by 'position').
	//      2 = delete (delete the existing order at the row identified by 'position').
	//  side -  0 for ask, 1 for bid.
	//  price - the order's price.
	//  size -  the order's size.
	//  isSmartDepth - is SMART Depth request.
	UpdateMktDepthL2(TickerID models.TickerID, position int64, marketMaker string, operation int64, side int64, price float64, size models.Decimal, isSmartDepth bool)

	// ReceiveFA receives the Financial Advisor's configuration available in the TWS
	//  faData - one of:
	// 		GROUPS: offer traders a way to create a group of accounts and apply a single allocation method to all accounts in the group.
	// 		ALIASES: let you easily identify the accounts by meaningful names rather than account numbers.
	// faXmlData -  the xml-formatted configuration
	ReceiveFA(faData models.FaData, cxml string)

	// ScannerParameters Provides the xml-formatted parameters available to create a market scanner.
	ScannerParameters(xml string)
	// ScannerData Provides the data resulting from the market scanner request.
	// reqID - the request's identifier.
	// rank -  the ranking within the response of this bar.
	// contractDetails - the data's ContractDetails
	// distance -      according to query.
	// benchmark -     according to query.
	// projection -    according to query.
	// legStr - describes the combo legs when the scanner is returning EFP
	ScannerData(reqID int64, rank int64, contractDetails *models.ContractDetails, distance string, benchmark string, projection string, legsStr string)
	// ScannerDataEnd indicates that the scanner data reception has terminated.
	ScannerDataEnd(reqID int64)
	// RealtimeBar updates the real time 5 seconds bars
	// reqID - the request's identifier
	// time  - start of bar in unix (or 'epoch') time
	// open_  - the bar's open value
	// high  - the bar's high value
	// low   - the bar's low value
	// close - the bar's closing value
	// volume - the bar's traded volume if available
	// wap   - the bar's Weighted Average Price
	// count - the number of trades during the bar's timespan (only available for TRADES).
	RealtimeBar(reqID models.TickerID, time int64, open float64, high float64, low float64, close float64, volume models.Decimal, wap models.Decimal, count int64)

	// FundamentalData
	FundamentalData(reqID models.TickerID, data string)
	// DeltaNeutralValidation
	DeltaNeutralValidation(reqID int64, deltaNeutralContract models.DeltaNeutralContract)
	// TickSnapshotEnd indicates the snapshot reception is finished.
	TickSnapshotEnd(reqID int64)
	// MarketDataType is called when market data switches between real-time and frozen.
	// The marketDataType( ) callback accepts a reqId parameter and is sent per every subscription because different contracts can generally trade on a different schedule
	MarketDataType(reqID models.TickerID, marketDataType int64)
	// CommissionAndFeesReport is called immediately after a trade execution or by calling reqExecutions().
	CommissionAndFeesReport(commissionAndFeesReport models.CommissionAndFeesReport)
	// Position returns real-time positions for all accounts in response to the reqPositions() method.
	Position(account string, contract *models.Contract, position models.Decimal, avgCost float64)
	// PositionEnd is called once all position data for a given request are received and functions as an end marker for the position() data.
	PositionEnd()
	// AccountSummary returns the data from the TWS Account Window Summary tab in response to reqAccountSummary().
	AccountSummary(reqID int64, account string, tag string, value string, currency string)
	// AccountSummaryEnd is called once all account summary data for a given request are received.
	AccountSummaryEnd(reqID int64)
	// VerifyMessageAPI .
	VerifyMessageAPI(apiData string)
	// VerifyCompleted .
	VerifyCompleted(isSuccessful bool, errorText string)
	// DisplayGroupList is a one-time response to queryDisplayGroups().
	// reqID - The reqID specified in queryDisplayGroups().
	// groups - A list of integers representing visible group ID separated by the | character, and sorted by most used group first. This list will
	//      not change during TWS session (in other words, user cannot add a new group; sorting can change though).
	DisplayGroupList(reqID int64, groups string)
	// DisplayGroupUpdated .
	DisplayGroupUpdated(reqID int64, contractInfo string)
	// VerifyAndAuthMessageAPI .
	VerifyAndAuthMessageAPI(apiData string, xyzChallange string)
	// VerifyAndAuthCompleted .
	VerifyAndAuthCompleted(isSuccessful bool, errorText string)

	// PositionMulti .
	PositionMulti(reqID int64, account string, modelCode string, contract *models.Contract, pos models.Decimal, avgCost float64)
	// PositionMultiEnd .
	PositionMultiEnd(reqID int64)
	// AccountUpdateMulti .
	AccountUpdateMulti(reqID int64, account string, modleCode string, key string, value string, currency string)
	// AccountUpdateMultiEnd .
	AccountUpdateMultiEnd(reqID int64)
	// SecurityDefinitionOptionParameter returns the option chain for an underlying on an exchange specified in reqSecDefOptParams.
	// There will be multiple callbacks to securityDefinitionOptionParameter if multiple exchanges are specified in reqSecDefOptParams.
	// reqId - ID of the request initiating the callback.
	// underlyingConId - The conID of the underlying security.
	// tradingClass -  the option trading class.
	// multiplier -    the option multiplier.
	// expirations - a list of the expiries for the options of this underlying on this exchange.
	// strikes - a list of the possible strikes for options of this underlying on this exchange.
	SecurityDefinitionOptionParameter(reqID int64, exchange string, underlyingConID int64, tradingClass string, multiplier string, expirations []string, strikes []float64)
	// SecurityDefinitionOptionParameterEnd is called when all callbacks to securityDefinitionOptionParameter are completed.
	SecurityDefinitionOptionParameterEnd(reqID int64)
	// SoftDollarTiers is called when receives Soft Dollar Tier configuration information.
	// reqID - The request ID used in the call to reqSoftDollarTiers()
	// tiers - Stores a list of SoftDollarTier that contains all Soft Dollar Tiers information
	SoftDollarTiers(reqID int64, tiers []models.SoftDollarTier)
	// FamilyCodes .
	FamilyCodes(familyCodes []models.FamilyCode)
	// SymbolSamples .
	SymbolSamples(reqID int64, contractDescriptions []models.ContractDescription)
	// MktDepthExchanges .
	MktDepthExchanges(depthMktDataDescriptions []models.DepthMktDataDescription)
	// SmartComponents .
	SmartComponents(reqID int64, smartComponents []models.SmartComponent)
	// TickReqParams .
	TickReqParams(TickerID models.TickerID, minTick float64, bboExchange string, snapshotPermissions int64)
	// HeadTimestamp returns earliest available data of a type of data for a particular contract.
	HeadTimestamp(reqID int64, headTimestamp string)
	// HistogramData returns histogram data for a contract.
	HistogramData(reqID int64, data []models.HistogramData)

	// RerouteMktDataReq .
	RerouteMktDataReq(reqID int64, conID int64, exchange string)
	// RerouteMktDepthReq .
	RerouteMktDepthReq(reqID int64, conID int64, exchange string)
	// MarketRule .
	MarketRule(marketRuleID int64, priceIncrements []models.PriceIncrement)
	// Pnl returns the daily PnL for the account.
	Pnl(reqID int64, dailyPnL float64, unrealizedPnL float64, realizedPnL float64)
	// PnlSingle returns the daily PnL for a single position in the account.
	PnlSingle(reqID int64, pos models.Decimal, dailyPnL float64, unrealizedPnL float64, realizedPnL float64, value float64)

	// TickByTickAllLast returns tick-by-tick data for tickType = "Last" or "AllLast".
	TickByTickAllLast(reqID int64, tickType int64, time int64, price float64, size models.Decimal, tickAttribLast models.TickAttribLast, exchange string, specialConditions string)
	// TickByTickBidAsk .
	TickByTickBidAsk(reqID int64, time int64, bidPrice float64, askPrice float64, bidSize models.Decimal, askSize models.Decimal, tickAttribBidAsk models.TickAttribBidAsk)
	// TickByTickMidPoint .
	TickByTickMidPoint(reqID int64, time int64, midPoint float64)
	// OrderBound returns orderBound notification.
	OrderBound(permID int64, clientID int64, orderID int64)
	// CompletedOrder is called to feed in completed orders.
	CompletedOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState)
	// CompletedOrdersEnd is called at the end of a given request for completed orders.
	CompletedOrdersEnd()
	// ReplaceFAEnd is called at the end of a replace FA.
	ReplaceFAEnd(reqID int64, text string)
	// WshMetaData .
	WshMetaData(reqID int64, dataJson string)
	// WshEventData .
	WshEventData(reqID int64, dataJson string)

	// UserInfo returns user info.
	UserInfo(reqID int64, whiteBrandingId string)
}
