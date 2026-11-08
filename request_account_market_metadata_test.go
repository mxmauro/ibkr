package ibkr

import (
	"context"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func TestRequestPositionsCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.PositionsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestPositions(context.Background())
		resultCh <- struct {
			resp *models.PositionsResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	assertTextPositionsRequestPayload(t, payload)

	contract := testRequestContract()
	if err := client.processIncomingMessage(buildTextPositionMessage("DU123", contract, "5", 123.45)); err != nil {
		t.Fatalf("processIncomingMessage(position) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildTextPositionEndMessage()); err != nil {
		t.Fatalf("processIncomingMessage(position end) failed: %v", err)
	}

	result := waitPositionsResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestPositions failed: %v", result.err)
	}
	if got := len(result.resp.Positions); got != 1 {
		t.Fatalf("unexpected positions count: got %d", got)
	}
	assertPositionRow(t, result.resp.Positions[0], "DU123", contract, "5", 123.45)
}

func TestRequestPositionsCompletesFromProtobufReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.PositionsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestPositions(context.Background())
		resultCh <- struct {
			resp *models.PositionsResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	assertEmptyProtobufPositionsRequestPayload(t, payload)

	contract := testRequestContract()
	if err := client.processIncomingMessage(buildProtoMessage(common.POSITION_DATA, &protobuf.Position{
		Account:  protofmt.String("DU999"),
		Contract: contract.Proto(nil),
		Position: protofmt.String("7"),
		AvgCost:  protofmt.Float(77.7),
	})); err != nil {
		t.Fatalf("processIncomingMessage(position) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.POSITION_END, &protobuf.PositionEnd{})); err != nil {
		t.Fatalf("processIncomingMessage(position end) failed: %v", err)
	}

	result := waitPositionsResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestPositions failed: %v", result.err)
	}
	if got := len(result.resp.Positions); got != 1 {
		t.Fatalf("unexpected positions count: got %d", got)
	}
	assertPositionRow(t, result.resp.Positions[0], "DU999", contract, "7", 77.7)
}

func TestRequestAccountSummaryCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.AccountSummaryResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestAccountSummary(context.Background(), models.AccountSummaryRequestOptions{Group: "All", Tags: "NetLiquidation"})
		resultCh <- struct {
			resp *models.AccountSummaryResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	reqID := assertTextAccountSummaryRequestPayload(t, payload, "All", "NetLiquidation")

	if err := client.processIncomingMessage(buildTextAccountSummaryMessage(reqID, "DU123", "NetLiquidation", "1000.50", "USD")); err != nil {
		t.Fatalf("processIncomingMessage(account summary) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildTextAccountSummaryEndMessage(reqID)); err != nil {
		t.Fatalf("processIncomingMessage(account summary end) failed: %v", err)
	}

	result := waitAccountSummaryResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestAccountSummary failed: %v", result.err)
	}
	if got := len(result.resp.Values); got != 1 {
		t.Fatalf("unexpected account summary row count: got %d", got)
	}
	if !reflect.DeepEqual(result.resp.Values[0], models.AccountSummaryValue{Account: "DU123", Tag: "NetLiquidation", Value: "1000.50", Currency: "USD"}) {
		t.Fatalf("unexpected account summary row: %#v", result.resp.Values[0])
	}
}

func TestRequestAccountSummaryCompletesFromProtobufReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.AccountSummaryResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestAccountSummary(context.Background(), models.AccountSummaryRequestOptions{Group: "All", Tags: "BuyingPower"})
		resultCh <- struct {
			resp *models.AccountSummaryResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	req := assertProtobufAccountSummaryRequestPayload(t, payload, "All", "BuyingPower")

	if err := client.processIncomingMessage(buildProtoMessage(common.ACCOUNT_SUMMARY, &protobuf.AccountSummary{
		ReqId:    req.ReqId,
		Account:  protofmt.String("DU222"),
		Tag:      protofmt.String("BuyingPower"),
		Value:    protofmt.String("25000"),
		Currency: protofmt.String("USD"),
	})); err != nil {
		t.Fatalf("processIncomingMessage(account summary) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.ACCOUNT_SUMMARY_END, &protobuf.AccountSummaryEnd{ReqId: req.ReqId})); err != nil {
		t.Fatalf("processIncomingMessage(account summary end) failed: %v", err)
	}

	result := waitAccountSummaryResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestAccountSummary failed: %v", result.err)
	}
	if got := len(result.resp.Values); got != 1 {
		t.Fatalf("unexpected account summary row count: got %d", got)
	}
	if !reflect.DeepEqual(result.resp.Values[0], models.AccountSummaryValue{Account: "DU222", Tag: "BuyingPower", Value: "25000", Currency: "USD"}) {
		t.Fatalf("unexpected account summary row: %#v", result.resp.Values[0])
	}
}

func TestRequestHeadTimestampCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.HeadTimestampResponse
		err  error
	}, 1)
	contract := testRequestContract()
	go func() {
		resp, err := client.RequestHeadTimestamp(context.Background(), models.HeadTimestampRequestOptions{Contract: contract, WhatToShow: models.WhatToShowTrades, OnlyRegularTradingHours: true, FormatDate: 2})
		resultCh <- struct {
			resp *models.HeadTimestampResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	reqID := assertTextHeadTimestampRequestPayload(t, payload, contract, true, models.WhatToShowTrades, 2)
	wantTS := time.Unix(1712345678, 0).UTC()
	if err := client.processIncomingMessage(buildTextHeadTimestampMessage(reqID, wantTS)); err != nil {
		t.Fatalf("processIncomingMessage(head timestamp) failed: %v", err)
	}

	result := waitHeadTimestampResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestHeadTimestamp failed: %v", result.err)
	}
	if !result.resp.Timestamp.Equal(wantTS) {
		t.Fatalf("unexpected head timestamp: got %v want %v", result.resp.Timestamp, wantTS)
	}
}

func TestRequestHeadTimestampCompletesFromProtobufReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.HeadTimestampResponse
		err  error
	}, 1)
	contract := testRequestContract()
	go func() {
		resp, err := client.RequestHeadTimestamp(context.Background(), models.HeadTimestampRequestOptions{Contract: contract, WhatToShow: models.WhatToShowBid, OnlyRegularTradingHours: false, FormatDate: 2})
		resultCh <- struct {
			resp *models.HeadTimestampResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	req := assertProtobufHeadTimestampRequestPayload(t, payload, contract, false, models.WhatToShowBid, 2)
	if err := client.processIncomingMessage(buildProtoMessage(common.HEAD_TIMESTAMP, &protobuf.HeadTimestamp{ReqId: req.ReqId, HeadTimestamp: protofmt.String("1712345678")})); err != nil {
		t.Fatalf("processIncomingMessage(head timestamp) failed: %v", err)
	}

	result := waitHeadTimestampResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestHeadTimestamp failed: %v", result.err)
	}
	wantTS := time.Unix(1712345678, 0).UTC()
	if !result.resp.Timestamp.Equal(wantTS) {
		t.Fatalf("unexpected head timestamp: got %v want %v", result.resp.Timestamp, wantTS)
	}
}
func TestRequestHistogramDataCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.HistogramDataResponse
		err  error
	}, 1)
	contract := testRequestContract()
	go func() {
		resp, err := client.RequestHistogramData(context.Background(), models.HistogramDataRequestOptions{Contract: contract, OnlyRegularTradingHours: true, TimePeriod: "3 days"})
		resultCh <- struct {
			resp *models.HistogramDataResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	reqID := assertTextHistogramRequestPayload(t, payload, contract, true, "3 days")
	if err := client.processIncomingMessage(buildTextHistogramDataMessage(reqID, []histogramPoint{{price: 10.5, size: "4"}, {price: 11.0, size: "9"}})); err != nil {
		t.Fatalf("processIncomingMessage(histogram) failed: %v", err)
	}

	result := waitHistogramResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestHistogramData failed: %v", result.err)
	}
	assertHistogramData(t, result.resp.Data, []histogramPoint{{price: 10.5, size: "4"}, {price: 11.0, size: "9"}})
}

func TestRequestHistogramDataCompletesFromProtobufReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.HistogramDataResponse
		err  error
	}, 1)
	contract := testRequestContract()
	go func() {
		resp, err := client.RequestHistogramData(context.Background(), models.HistogramDataRequestOptions{Contract: contract, OnlyRegularTradingHours: false, TimePeriod: "1 week"})
		resultCh <- struct {
			resp *models.HistogramDataResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	req := assertProtobufHistogramRequestPayload(t, payload, contract, false, "1 week")
	if err := client.processIncomingMessage(buildProtoMessage(common.HISTOGRAM_DATA, &protobuf.HistogramData{ReqId: req.ReqId, HistogramDataEntries: []*protobuf.HistogramDataEntry{{Price: protofmt.Float(20.25), Size: protofmt.String("3")}, {Price: protofmt.Float(21.75), Size: protofmt.String("8")}}})); err != nil {
		t.Fatalf("processIncomingMessage(histogram) failed: %v", err)
	}

	result := waitHistogramResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestHistogramData failed: %v", result.err)
	}
	assertHistogramData(t, result.resp.Data, []histogramPoint{{price: 20.25, size: "3"}, {price: 21.75, size: "8"}})
}

func TestRequestMktDepthExchangesCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.MktDepthExchangesResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestMktDepthExchanges(context.Background())
		resultCh <- struct {
			resp *models.MktDepthExchangesResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	assertTextMktDepthExchangesRequestPayload(t, payload)

	agg := int32(7)
	if err := client.processIncomingMessage(buildTextMktDepthExchangesMessage([]models.DepthMktDataDescription{{Exchange: "ARCA", SecType: models.SecurityTypeStock, ListingExchange: "NYSE", ServiceDataType: "Deep2", AggGroup: &agg}})); err != nil {
		t.Fatalf("processIncomingMessage(mkt depth exchanges) failed: %v", err)
	}

	result := waitMktDepthExchangesResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestMktDepthExchanges failed: %v", result.err)
	}
	if got := len(result.resp.Exchanges); got != 1 {
		t.Fatalf("unexpected exchanges count: got %d", got)
	}
	if !reflect.DeepEqual(result.resp.Exchanges[0], models.DepthMktDataDescription{Exchange: "ARCA", SecType: models.SecurityTypeStock, ListingExchange: "NYSE", ServiceDataType: "Deep2", AggGroup: &agg}) {
		t.Fatalf("unexpected exchange row: %#v", result.resp.Exchanges[0])
	}
}

func TestRequestSmartComponentsCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufHistoricalData)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.SmartComponentsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestSmartComponents(context.Background(), models.SmartComponentsRequestOptions{BBOExchange: "a6"})
		resultCh <- struct {
			resp *models.SmartComponentsResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	reqID := assertTextSmartComponentsRequestPayload(t, payload, "a6")
	if err := client.processIncomingMessage(buildTextSmartComponentsMessage(reqID, []models.SmartComponent{{BitNumber: 1, Exchange: "ISLAND", ExchangeLetter: "I"}})); err != nil {
		t.Fatalf("processIncomingMessage(smart components) failed: %v", err)
	}

	result := waitSmartComponentsResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestSmartComponents failed: %v", result.err)
	}
	if result.resp.ReqID != reqID {
		t.Fatalf("unexpected smart components req id: got %d want %d", result.resp.ReqID, reqID)
	}
	if !reflect.DeepEqual(result.resp.Components, []models.SmartComponent{{BitNumber: 1, Exchange: "ISLAND", ExchangeLetter: "I"}}) {
		t.Fatalf("unexpected smart components: %#v", result.resp.Components)
	}
}

func TestRequestAccountUpdatesRoutesProtobufEventsAndCompletesWithoutHandler(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions)
	defer cleanup()

	resp, err := client.RequestAccountUpdates(context.Background(), models.AccountUpdatesRequestOptions{AccountCode: "DU111"})
	if err != nil {
		t.Fatalf("RequestAccountUpdates failed: %v", err)
	}

	payload := readFramedPayload(t, serverConn)
	assertProtobufAccountUpdatesRequestPayload(t, payload, true, "DU111")

	contract := testRequestContract()
	if err := client.processIncomingMessage(buildProtoMessage(common.ACCT_VALUE, &protobuf.AccountValue{Key: protofmt.String("NetLiquidation"), Value: protofmt.String("1000"), Currency: protofmt.String("USD"), AccountName: protofmt.String("DU111")})); err != nil {
		t.Fatalf("processIncomingMessage(account value) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.PORTFOLIO_VALUE, &protobuf.PortfolioValue{Contract: contract.Proto(nil), Position: protofmt.String("2"), MarketPrice: protofmt.Float(101.25), MarketValue: protofmt.Float(202.5), AverageCost: protofmt.Float(100.0), UnrealizedPNL: protofmt.Float(2.5), RealizedPNL: protofmt.Float(1.0), AccountName: protofmt.String("DU111")})); err != nil {
		t.Fatalf("processIncomingMessage(portfolio) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.ACCT_UPDATE_TIME, &protobuf.AccountUpdateTime{TimeStamp: protofmt.String("12:34")})); err != nil {
		t.Fatalf("processIncomingMessage(account time) failed: %v", err)
	}
	if err := client.processIncomingMessage(buildTextAccountDownloadEndMessage("DU111")); err != nil {
		t.Fatalf("processIncomingMessage(account download end) failed: %v", err)
	}

	events := collectAccountUpdateEvents(t, resp.Channel, 4)
	assertAccountUpdateEvents(t, events, contract)
}

func TestRequestAccountUpdatesContextCancellationAdvancesQueuedNoIDRequests(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufAccountsPositions-1)
	defer cleanup()

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	resp1, err := client.RequestAccountUpdates(ctx1, models.AccountUpdatesRequestOptions{AccountCode: "DU111"})
	if err != nil {
		t.Fatalf("RequestAccountUpdates(first) failed: %v", err)
	}
	assertTextAccountUpdatesRequestPayload(t, readFramedPayload(t, serverConn), true, "DU111")

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	resp2, err := client.RequestAccountUpdates(ctx2, models.AccountUpdatesRequestOptions{AccountCode: "DU222"})
	if err != nil {
		t.Fatalf("RequestAccountUpdates(second) failed: %v", err)
	}
	assertTextAccountUpdatesRequestPayload(t, readFramedPayload(t, serverConn), true, "DU222")

	if err := client.processIncomingMessage(buildTextAccountValueMessage("NetLiquidation", "1", "USD", "DU111")); err != nil {
		t.Fatalf("processIncomingMessage(first account value) failed: %v", err)
	}
	first := waitAccountUpdateEvent(t, resp1.Channel)
	if value, ok := first.(models.AccountValueUpdate); !ok || value.AccountName != "DU111" {
		t.Fatalf("unexpected first stream event: %#v", first)
	}
	assertNoAccountUpdateEvent(t, resp2.Channel)

	cancel1()
	assertTextAccountUpdatesRequestPayload(t, readFramedPayload(t, serverConn), false, "DU111")
	waitAccountUpdatesClosed(t, resp1.Channel)

	if err := client.processIncomingMessage(buildTextAccountValueMessage("NetLiquidation", "2", "USD", "DU222")); err != nil {
		t.Fatalf("processIncomingMessage(second account value) failed: %v", err)
	}
	second := waitAccountUpdateEvent(t, resp2.Channel)
	if value, ok := second.(models.AccountValueUpdate); !ok || value.AccountName != "DU222" || value.Value != "2" {
		t.Fatalf("unexpected second stream event: %#v", second)
	}

	resp2.Cancel()
	assertTextAccountUpdatesRequestPayload(t, readFramedPayload(t, serverConn), false, "DU222")
	waitAccountUpdatesClosed(t, resp2.Channel)
}

// -----------------------------------------------------------------------------

type histogramPoint struct {
	price float64
	size  string
}

func newRequestTestClientWithVersion(t *testing.T, serverVersion int32) (*Client, net.Conn, func()) {
	t.Helper()
	client, serverConn, cleanup := newMarketRuleRequestTestClient(t)
	client.serverVersion = serverVersion
	client.nextValidReqID = 1
	client.nextValidReqWithoutID = 1
	return client, serverConn, cleanup
}

func testRequestContract() *models.Contract {
	return &models.Contract{ConID: 1001, Symbol: "AAPL", SecType: models.SecurityTypeStock, Exchange: "SMART", Currency: "USD", LocalSymbol: "AAPL", TradingClass: "NMS"}
}

func buildProtoMessage(msgID uint32, pb proto.Message) []byte {
	payload, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	msg := make([]byte, 4, 4+len(payload))
	binary.BigEndian.PutUint32(msg, msgID+common.PROTOBUF_MSG_ID)
	return append(msg, payload...)
}

func buildTextPositionMessage(account string, contract *models.Contract, position string, avgCost float64) []byte {
	return message.NewRawEncoder().RawUInt32(common.POSITION_DATA).Int(1).String(account).Marshal(contract, 1).String(position).Float(avgCost).Bytes()
}

func buildTextPositionEndMessage() []byte {
	return message.NewRawEncoder().RawUInt32(common.POSITION_END).Bytes()
}

func buildTextAccountSummaryMessage(reqID int32, account, tag, value, currency string) []byte {
	return message.NewRawEncoder().RawUInt32(common.ACCOUNT_SUMMARY).Int(1).RequestID(reqID).String(account).String(tag).String(value).String(currency).Bytes()
}

func buildTextAccountSummaryEndMessage(reqID int32) []byte {
	return message.NewRawEncoder().RawUInt32(common.ACCOUNT_SUMMARY_END).Int(1).RequestID(reqID).Bytes()
}

func buildTextHeadTimestampMessage(reqID int32, ts time.Time) []byte {
	return message.NewRawEncoder().RawUInt32(common.HEAD_TIMESTAMP).RequestID(reqID).Int64(ts.Unix()).Bytes()
}
func buildTextHistogramDataMessage(reqID int32, points []histogramPoint) []byte {
	msgEnc := message.NewRawEncoder().RawUInt32(common.HISTOGRAM_DATA).RequestID(reqID).Int64(int64(len(points)))
	for _, point := range points {
		msgEnc.Float(point.price).String(point.size)
	}
	return msgEnc.Bytes()
}

func buildTextMktDepthExchangesMessage(values []models.DepthMktDataDescription) []byte {
	msgEnc := message.NewRawEncoder().RawUInt32(common.MKT_DEPTH_EXCHANGES).Int64(int64(len(values)))
	for _, value := range values {
		msgEnc.String(value.Exchange).String(string(value.SecType)).String(value.ListingExchange).String(value.ServiceDataType).Int32Max(value.AggGroup)
	}
	return msgEnc.Bytes()
}

func buildTextSmartComponentsMessage(reqID int32, values []models.SmartComponent) []byte {
	msgEnc := message.NewRawEncoder().RawUInt32(common.SMART_COMPONENTS).RequestID(reqID).Int64(int64(len(values)))
	for _, value := range values {
		msgEnc.Int64(value.BitNumber).String(value.Exchange).String(value.ExchangeLetter)
	}
	return msgEnc.Bytes()
}

func buildTextAccountValueMessage(tag, value, currency, accountName string) []byte {
	return message.NewRawEncoder().RawUInt32(common.ACCT_VALUE).Int(1).String(tag).String(value).String(currency).String(accountName).Bytes()
}

func buildTextAccountDownloadEndMessage(accountName string) []byte {
	return message.NewRawEncoder().RawUInt32(common.ACCT_DOWNLOAD_END).Int(1).String(accountName).Bytes()
}

func assertTextPositionsRequestPayload(t *testing.T, payload []byte) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_POSITIONS {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_POSITIONS)
	}
	msgDec := message.NewDecoder(payload[4:])
	if got := msgDec.Int32(); got != 1 {
		t.Fatalf("unexpected positions version: got %d", got)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding positions request: %v", msgDec.Err())
	}
}

func assertEmptyProtobufPositionsRequestPayload(t *testing.T, payload []byte) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_POSITIONS+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.PositionsRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf positions request: %v", err)
	}
}

func assertTextAccountSummaryRequestPayload(t *testing.T, payload []byte, group, tags string) int32 {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_ACCOUNT_SUMMARY {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_ACCOUNT_SUMMARY)
	}
	msgDec := message.NewDecoder(payload[4:])
	if got := msgDec.Int32(); got != 1 {
		t.Fatalf("unexpected account summary version: got %d", got)
	}
	reqID := msgDec.RequestID(false)
	if got := msgDec.String(); got != group {
		t.Fatalf("unexpected account summary group: got %q want %q", got, group)
	}
	if got := msgDec.String(); got != tags {
		t.Fatalf("unexpected account summary tags: got %q want %q", got, tags)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding account summary request: %v", msgDec.Err())
	}
	return reqID
}

func assertProtobufAccountSummaryRequestPayload(t *testing.T, payload []byte, group, tags string) *protobuf.AccountSummaryRequest {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_ACCOUNT_SUMMARY+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.AccountSummaryRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf account summary request: %v", err)
	}
	if pb.GetGroup() != group || pb.GetTags() != tags {
		t.Fatalf("unexpected protobuf account summary request: %#v", pb)
	}
	return pb
}

func assertTextHeadTimestampRequestPayload(t *testing.T, payload []byte, contract *models.Contract, useRTH bool, what models.WhatToShow, formatDate int) int32 {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_HEAD_TIMESTAMP {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_HEAD_TIMESTAMP)
	}
	msgDec := message.NewDecoder(payload[4:])
	reqID := msgDec.RequestID(false)
	assertDecodedContract(t, msgDec, contract)
	if got := msgDec.Bool(); got != useRTH {
		t.Fatalf("unexpected head timestamp useRTH: got %t want %t", got, useRTH)
	}
	if got := msgDec.String(); got != what.String() {
		t.Fatalf("unexpected head timestamp whatToShow: got %q want %q", got, what.String())
	}
	if got := msgDec.Int32(); got != int32(formatDate) {
		t.Fatalf("unexpected head timestamp formatDate: got %d want %d", got, formatDate)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding head timestamp request: %v", msgDec.Err())
	}
	return reqID
}

func assertProtobufHeadTimestampRequestPayload(t *testing.T, payload []byte, contract *models.Contract, useRTH bool, what models.WhatToShow, formatDate int) *protobuf.HeadTimestampRequest {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_HEAD_TIMESTAMP+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.HeadTimestampRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf head timestamp request: %v", err)
	}
	if pb.GetUseRTH() != useRTH || pb.GetWhatToShow() != what.String() || pb.GetFormatDate() != int32(formatDate) || pb.GetContract().GetConId() != contract.ConID {
		t.Fatalf("unexpected protobuf head timestamp request: %#v", pb)
	}
	return pb
}

func assertTextHistogramRequestPayload(t *testing.T, payload []byte, contract *models.Contract, useRTH bool, period string) int32 {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_HISTOGRAM_DATA {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_HISTOGRAM_DATA)
	}
	msgDec := message.NewDecoder(payload[4:])
	reqID := msgDec.RequestID(false)
	assertDecodedContract(t, msgDec, contract)
	if got := msgDec.Bool(); got != useRTH {
		t.Fatalf("unexpected histogram useRTH: got %t want %t", got, useRTH)
	}
	if got := msgDec.String(); got != period {
		t.Fatalf("unexpected histogram period: got %q want %q", got, period)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding histogram request: %v", msgDec.Err())
	}
	return reqID
}

func assertProtobufHistogramRequestPayload(t *testing.T, payload []byte, contract *models.Contract, useRTH bool, period string) *protobuf.HistogramDataRequest {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_HISTOGRAM_DATA+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.HistogramDataRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf histogram request: %v", err)
	}
	if pb.GetUseRTH() != useRTH || pb.GetTimePeriod() != period || pb.GetContract().GetConId() != contract.ConID {
		t.Fatalf("unexpected protobuf histogram request: %#v", pb)
	}
	return pb
}

func assertTextMktDepthExchangesRequestPayload(t *testing.T, payload []byte) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_MKT_DEPTH_EXCHANGES {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_MKT_DEPTH_EXCHANGES)
	}
}

func assertTextSmartComponentsRequestPayload(t *testing.T, payload []byte, exchange string) int32 {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_SMART_COMPONENTS {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_SMART_COMPONENTS)
	}
	msgDec := message.NewDecoder(payload[4:])
	reqID := msgDec.RequestID(false)
	if got := msgDec.String(); got != exchange {
		t.Fatalf("unexpected smart components exchange: got %q want %q", got, exchange)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding smart components request: %v", msgDec.Err())
	}
	return reqID
}

func assertProtobufAccountUpdatesRequestPayload(t *testing.T, payload []byte, subscribe bool, accountCode string) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_ACCT_DATA+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.AccountDataRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf account updates request: %v", err)
	}
	if pb.GetSubscribe() != subscribe || pb.GetAcctCode() != accountCode {
		t.Fatalf("unexpected protobuf account updates request: %#v", pb)
	}
}

func assertTextAccountUpdatesRequestPayload(t *testing.T, payload []byte, subscribe bool, accountCode string) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_ACCT_DATA {
		t.Fatalf("unexpected outbound msg id: got %d want %d", got, common.REQ_ACCT_DATA)
	}
	msgDec := message.NewDecoder(payload[4:])
	if got := msgDec.Int32(); got != 2 {
		t.Fatalf("unexpected account updates version: got %d", got)
	}
	if got := msgDec.Bool(); got != subscribe {
		t.Fatalf("unexpected account updates subscribe flag: got %t want %t", got, subscribe)
	}
	if got := msgDec.String(); got != accountCode {
		t.Fatalf("unexpected account updates account code: got %q want %q", got, accountCode)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding account updates request: %v", msgDec.Err())
	}
}

func assertDecodedContract(t *testing.T, msgDec *message.Decoder, contract *models.Contract) {
	t.Helper()
	if got := msgDec.Int32(); got != contract.ConID {
		t.Fatalf("unexpected contract conid: got %d want %d", got, contract.ConID)
	}
	if got := msgDec.String(); got != contract.Symbol {
		t.Fatalf("unexpected contract symbol: got %q want %q", got, contract.Symbol)
	}
	if got := msgDec.String(); got != string(contract.SecType) {
		t.Fatalf("unexpected contract sec type: got %q want %q", got, contract.SecType)
	}
	_ = msgDec.String()
	_ = msgDec.FloatMax()
	_ = msgDec.String()
	_ = msgDec.FloatMax()
	if got := msgDec.String(); got != contract.Exchange {
		t.Fatalf("unexpected contract exchange: got %q want %q", got, contract.Exchange)
	}
	if got := msgDec.String(); got != contract.PrimaryExchange {
		t.Fatalf("unexpected contract primary exchange: got %q want %q", got, contract.PrimaryExchange)
	}
	if got := msgDec.String(); got != contract.Currency {
		t.Fatalf("unexpected contract currency: got %q want %q", got, contract.Currency)
	}
	if got := msgDec.String(); got != contract.LocalSymbol {
		t.Fatalf("unexpected contract local symbol: got %q want %q", got, contract.LocalSymbol)
	}
	if got := msgDec.String(); got != contract.TradingClass {
		t.Fatalf("unexpected contract trading class: got %q want %q", got, contract.TradingClass)
	}
}

func assertPositionRow(t *testing.T, got models.Position, account string, contract *models.Contract, position string, avgCost float64) {
	t.Helper()
	if got.Account != account {
		t.Fatalf("unexpected account: got %q want %q", got.Account, account)
	}
	if !got.Contract.Equal(contract) {
		t.Fatalf("unexpected contract: %#v", got.Contract)
	}
	if got.Position.String() != position {
		t.Fatalf("unexpected position: got %s want %s", got.Position.String(), position)
	}
	if got.AvgCost != avgCost {
		t.Fatalf("unexpected avg cost: got %f want %f", got.AvgCost, avgCost)
	}
}

func assertHistogramData(t *testing.T, got []models.HistogramData, want []histogramPoint) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("unexpected histogram length: got %d want %d", len(got), len(want))
	}
	for idx := range want {
		if got[idx].Price != want[idx].price {
			t.Fatalf("unexpected histogram price at %d: got %f want %f", idx, got[idx].Price, want[idx].price)
		}
		if got[idx].Size == nil || got[idx].Size.String() != want[idx].size {
			t.Fatalf("unexpected histogram size at %d: got %v want %s", idx, got[idx].Size, want[idx].size)
		}
	}
}
func collectAccountUpdateEvents(t *testing.T, ch <-chan models.AccountUpdateEvent, count int) []models.AccountUpdateEvent {
	t.Helper()
	events := make([]models.AccountUpdateEvent, 0, count)
	deadline := time.After(5 * time.Second)
	for len(events) < count {
		select {
		case ev, ok := <-ch:
			if !ok {
				return events
			}
			events = append(events, ev)
		case <-deadline:
			t.Fatalf("timed out waiting for %d account update events", count)
		}
	}
	return events
}

func assertAccountUpdateEvents(t *testing.T, events []models.AccountUpdateEvent, contract *models.Contract) {
	t.Helper()
	if got := len(events); got != 4 {
		t.Fatalf("unexpected account update event count: got %d", got)
	}
	value, ok := events[0].(models.AccountValueUpdate)
	if !ok || value.Tag != "NetLiquidation" || value.Value != "1000" || value.AccountName != "DU111" {
		t.Fatalf("unexpected account value update: %#v", events[0])
	}
	portfolio, ok := events[1].(models.PortfolioUpdate)
	if !ok || !portfolio.Contract.Equal(contract) || portfolio.Position.String() != "2" || portfolio.AccountName != "DU111" {
		t.Fatalf("unexpected portfolio update: %#v", events[1])
	}
	acctTime, ok := events[2].(models.AccountTimeUpdate)
	if !ok || acctTime.TimeStamp != "12:34" {
		t.Fatalf("unexpected account time update: %#v", events[2])
	}
	done, ok := events[3].(models.AccountDownloadEndUpdate)
	if !ok || done.AccountName != "DU111" {
		t.Fatalf("unexpected account download end update: %#v", events[3])
	}
}

func waitAccountUpdateEvent(t *testing.T, ch <-chan models.AccountUpdateEvent) models.AccountUpdateEvent {
	t.Helper()
	select {
	case ev := <-ch:
		return ev
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for account update event")
		return nil
	}
}

func assertNoAccountUpdateEvent(t *testing.T, ch <-chan models.AccountUpdateEvent) {
	t.Helper()
	select {
	case ev := <-ch:
		t.Fatalf("unexpected account update event: %#v", ev)
	case <-time.After(150 * time.Millisecond):
	}
}

func waitAccountUpdatesClosed(t *testing.T, ch <-chan models.AccountUpdateEvent) {
	t.Helper()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected account updates channel to be closed")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for account updates channel close")
	}
}

func waitPositionsResult(t *testing.T, ch <-chan struct {
	resp *models.PositionsResponse
	err  error
}) struct {
	resp *models.PositionsResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for positions result")
		return struct {
			resp *models.PositionsResponse
			err  error
		}{}
	}
}

func waitAccountSummaryResult(t *testing.T, ch <-chan struct {
	resp *models.AccountSummaryResponse
	err  error
}) struct {
	resp *models.AccountSummaryResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for account summary result")
		return struct {
			resp *models.AccountSummaryResponse
			err  error
		}{}
	}
}

func waitHeadTimestampResult(t *testing.T, ch <-chan struct {
	resp *models.HeadTimestampResponse
	err  error
}) struct {
	resp *models.HeadTimestampResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for head timestamp result")
		return struct {
			resp *models.HeadTimestampResponse
			err  error
		}{}
	}
}

func waitHistogramResult(t *testing.T, ch <-chan struct {
	resp *models.HistogramDataResponse
	err  error
}) struct {
	resp *models.HistogramDataResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for histogram result")
		return struct {
			resp *models.HistogramDataResponse
			err  error
		}{}
	}
}

func waitMktDepthExchangesResult(t *testing.T, ch <-chan struct {
	resp *models.MktDepthExchangesResponse
	err  error
}) struct {
	resp *models.MktDepthExchangesResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for mkt depth exchanges result")
		return struct {
			resp *models.MktDepthExchangesResponse
			err  error
		}{}
	}
}

func waitSmartComponentsResult(t *testing.T, ch <-chan struct {
	resp *models.SmartComponentsResponse
	err  error
}) struct {
	resp *models.SmartComponentsResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for smart components result")
		return struct {
			resp *models.SmartComponentsResponse
			err  error
		}{}
	}
}
