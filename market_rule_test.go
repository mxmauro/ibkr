package ibkr

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/connection"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/mxmauro/resetevent"
)

// -----------------------------------------------------------------------------

func TestRequestMarketRuleCompletesFromTextReply(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClientWithVersion(t, common.ServerVersionProtobufRestMessages2-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.MarketRuleResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestMarketRule(context.Background(), models.MarketRuleRequestOptions{MarketRuleID: 26})
		resultCh <- struct {
			resp *models.MarketRuleResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	assertMarketRuleRequestPayload(t, payload, 26)

	if err := client.processIncomingMessage(buildTextMarketRuleMessage(26, []models.PriceIncrement{{LowEdge: 0, Increment: 0.01}, {LowEdge: 1, Increment: 0.05}})); err != nil {
		t.Fatalf("processIncomingMessage failed: %v", err)
	}

	result := waitMarketRuleResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestMarketRule failed: %v", result.err)
	}
	if result.resp.MarketRuleID != 26 {
		t.Fatalf("unexpected market rule id: got %d", result.resp.MarketRuleID)
	}
	if !reflect.DeepEqual(result.resp.PriceIncrements, []models.PriceIncrement{{LowEdge: 0, Increment: 0.01}, {LowEdge: 1, Increment: 0.05}}) {
		t.Fatalf("unexpected price increments: %#v", result.resp.PriceIncrements)
	}
}

func TestRequestMarketRuleCompletesFromProtobufReply(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClient(t)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.MarketRuleResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestMarketRule(context.Background(), models.MarketRuleRequestOptions{MarketRuleID: 42})
		resultCh <- struct {
			resp *models.MarketRuleResponse
			err  error
		}{resp: resp, err: err}
	}()

	payload := readFramedPayload(t, serverConn)
	assertProtobufMarketRuleRequestPayload(t, payload, 42)

	if err := client.processIncomingMessage(buildProtoMarketRuleMessage(42, []models.PriceIncrement{{LowEdge: 0.5, Increment: 0.01}})); err != nil {
		t.Fatalf("processIncomingMessage failed: %v", err)
	}

	result := waitMarketRuleResult(t, resultCh)
	if result.err != nil {
		t.Fatalf("RequestMarketRule failed: %v", result.err)
	}
	if result.resp.MarketRuleID != 42 {
		t.Fatalf("unexpected market rule id: got %d", result.resp.MarketRuleID)
	}
	if !reflect.DeepEqual(result.resp.PriceIncrements, []models.PriceIncrement{{LowEdge: 0.5, Increment: 0.01}}) {
		t.Fatalf("unexpected price increments: %#v", result.resp.PriceIncrements)
	}
}

func TestRequestMarketRuleCompletesAllDuplicatePendingRequestsOnFirstReply(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClient(t)
	defer cleanup()

	firstResultCh := make(chan struct {
		resp *models.MarketRuleResponse
		err  error
	}, 1)
	secondResultCh := make(chan struct {
		resp *models.MarketRuleResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.RequestMarketRule(context.Background(), models.MarketRuleRequestOptions{MarketRuleID: 15})
		firstResultCh <- struct {
			resp *models.MarketRuleResponse
			err  error
		}{resp: resp, err: err}
	}()
	go func() {
		resp, err := client.RequestMarketRule(context.Background(), models.MarketRuleRequestOptions{MarketRuleID: 15})
		secondResultCh <- struct {
			resp *models.MarketRuleResponse
			err  error
		}{resp: resp, err: err}
	}()

	assertProtobufMarketRuleRequestPayload(t, readFramedPayload(t, serverConn), 15)
	assertProtobufMarketRuleRequestPayload(t, readFramedPayload(t, serverConn), 15)

	increments := []models.PriceIncrement{{LowEdge: 0, Increment: 0.25}}
	if err := client.processIncomingMessage(buildProtoMarketRuleMessage(15, increments)); err != nil {
		t.Fatalf("processIncomingMessage failed: %v", err)
	}

	first := waitMarketRuleResult(t, firstResultCh)
	second := waitMarketRuleResult(t, secondResultCh)
	for _, result := range []struct {
		name string
		data struct {
			resp *models.MarketRuleResponse
			err  error
		}
	}{
		{name: "first", data: first},
		{name: "second", data: second},
	} {
		if result.data.err != nil {
			t.Fatalf("%s RequestMarketRule failed: %v", result.name, result.data.err)
		}
		if result.data.resp.MarketRuleID != 15 {
			t.Fatalf("%s response returned unexpected market rule id: %d", result.name, result.data.resp.MarketRuleID)
		}
		if !reflect.DeepEqual(result.data.resp.PriceIncrements, increments) {
			t.Fatalf("%s response returned unexpected increments: %#v", result.name, result.data.resp.PriceIncrements)
		}
	}
}

func TestOrderAlgoStrategyRoundTripPreservesArbitraryValues(t *testing.T) {
	order := models.NewOrder()
	order.AlgoStrategy = "ArrivalPxV2"
	order.AlgoParams = []models.TagValue{{Tag: "priority", Value: "Urgent"}, {Tag: "maxPctVol", Value: "0.2"}}

	msgDec := protofmt.NewDecoder(nil)
	decoded := models.NewOrderFromProtobufDecoder(msgDec, order.Proto())
	if msgDec.Err() != nil {
		t.Fatalf("protobuf decode failed: %v", msgDec.Err())
	}
	if decoded.AlgoStrategy != order.AlgoStrategy {
		t.Fatalf("unexpected algo strategy: got %q want %q", decoded.AlgoStrategy, order.AlgoStrategy)
	}
	if !reflect.DeepEqual(tagValuesToMap(decoded.AlgoParams), tagValuesToMap(order.AlgoParams)) {
		t.Fatalf("unexpected algo params: got %#v want %#v", decoded.AlgoParams, order.AlgoParams)
	}
}

func TestFundamentalsRemainUnsupportedAndDocumented(t *testing.T) {
	clientContent := readRepoFiles(t, "client_*.go")
	if strings.Contains(clientContent, "RequestFundamentalData(") {
		t.Fatal("unexpected fundamentals request API found")
	}

	eventsContent := readRepoFile(t, "events.go")
	if strings.Contains(eventsContent, "FundamentalData(") {
		t.Fatal("unexpected fundamentals event callback found")
	}

	incomingContent := readRepoFile(t, "incoming_general.go")
	if strings.Contains(incomingContent, "processFundamentalDataMsg") || strings.Contains(incomingContent, "case FUNDAMENTAL_DATA") {
		t.Fatal("unexpected fundamentals message handling found in incoming_general.go")
	}

	constsContent := readRepoFile(t, "common/consts.go")
	if !strings.Contains(constsContent, "FUNDAMENTAL_DATA                         uint32 = 51 // Deprecated in v1047.") {
		t.Fatal("missing deprecated fundamentals incoming message constant comment")
	}
	if !strings.Contains(constsContent, "REQ_FUNDAMENTAL_DATA          = 52 // Deprecated in v1047.") {
		t.Fatal("missing deprecated fundamentals request constant comment")
	}
	if !strings.Contains(constsContent, "CANCEL_FUNDAMENTAL_DATA       = 53 // Deprecated in v1047.") {
		t.Fatal("missing deprecated fundamentals cancel constant comment")
	}
}

func newMarketRuleRequestTestClient(t *testing.T) (*Client, net.Conn, func()) {
	return newMarketRuleRequestTestClientWithVersion(t, common.ServerVersionProtobufRestMessages2)
}

func newMarketRuleRequestTestClientWithVersion(t *testing.T, serverVersion int32) (*Client, net.Conn, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	acceptCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errCh <- err
			return
		}
		acceptCh <- conn
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := connection.New(ctx, listener.Addr().String())
	if err != nil {
		listener.Close()
		t.Fatalf("connection.New failed: %v", err)
	}

	var serverConn net.Conn
	select {
	case serverConn = <-acceptCh:
	case err = <-errCh:
		conn.Close()
		listener.Close()
		t.Fatalf("accept failed: %v", err)
	case <-time.After(5 * time.Second):
		conn.Close()
		listener.Close()
		t.Fatal("timed out waiting for accepted connection")
	}

	client := &Client{
		conn:             conn,
		isDisconnectedEv: resetevent.NewManualResetEvent(),
		serverVersion:    serverVersion,
	}
	client.rp.Initialize()
	client.initRequestManager()

	cleanup := func() {
		_ = serverConn.Close()
		conn.Close()
		_ = listener.Close()
	}

	return client, serverConn, cleanup
}

func readFramedPayload(t *testing.T, conn net.Conn) []byte {
	t.Helper()

	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline failed: %v", err)
	}

	var header [4]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		t.Fatalf("failed reading frame header: %v", err)
	}
	payload := make([]byte, binary.BigEndian.Uint32(header[:]))
	if _, err := io.ReadFull(conn, payload); err != nil {
		t.Fatalf("failed reading frame payload: %v", err)
	}
	return payload
}

func assertMarketRuleRequestPayload(t *testing.T, payload []byte, marketRuleID int64) {
	t.Helper()

	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_MARKET_RULE {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	msgDec := message.NewDecoder(payload[4:])
	if got := msgDec.Int64(); got != marketRuleID {
		t.Fatalf("unexpected outbound market rule id: got %d want %d", got, marketRuleID)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding outbound request: %v", msgDec.Err())
	}
}

func assertProtobufMarketRuleRequestPayload(t *testing.T, payload []byte, marketRuleID int64) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_MARKET_RULE+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.MarketRuleRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf market rule request: %v", err)
	}
	if got := pb.GetMarketRuleId(); got != int32(marketRuleID) {
		t.Fatalf("unexpected outbound market rule id: got %d want %d", got, marketRuleID)
	}
}

func buildTextMarketRuleMessage(marketRuleID int64, increments []models.PriceIncrement) []byte {
	msgEnc := message.NewRawEncoder().
		RawUInt32(common.MARKET_RULE).
		Int64(marketRuleID).
		Int64(int64(len(increments)))
	for _, increment := range increments {
		msgEnc.Float(increment.LowEdge).Float(increment.Increment)
	}
	return msgEnc.Bytes()
}

func buildProtoMarketRuleMessage(marketRuleID int64, increments []models.PriceIncrement) []byte {
	pb := &protobuf.MarketRule{MarketRuleId: protofmt.Int32(int32(marketRuleID))}
	for _, increment := range increments {
		pb.PriceIncrements = append(pb.PriceIncrements, &protobuf.PriceIncrement{
			LowEdge:   protofmt.Float(increment.LowEdge),
			Increment: protofmt.Float(increment.Increment),
		})
	}
	return buildProtoMessage(common.MARKET_RULE, pb)
}

func waitMarketRuleResult(t *testing.T, ch <-chan struct {
	resp *models.MarketRuleResponse
	err  error
}) struct {
	resp *models.MarketRuleResponse
	err  error
} {
	t.Helper()

	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for market rule result")
		return struct {
			resp *models.MarketRuleResponse
			err  error
		}{}
	}
}

func tagValuesToMap(values []models.TagValue) map[string]string {
	out := make(map[string]string, len(values))
	for _, value := range values {
		out[value.Tag] = value.Value
	}
	return out
}

func readRepoFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s) failed: %v", path, err)
	}
	return string(data)
}

func readRepoFiles(t *testing.T, pattern string) string {
	t.Helper()

	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Glob(%s) failed: %v", pattern, err)
	}
	if len(matches) == 0 {
		t.Fatalf("Glob(%s) returned no matches", pattern)
	}

	var b strings.Builder
	for _, match := range matches {
		b.WriteString(readRepoFile(t, match))
		b.WriteByte('\n')
	}
	return b.String()
}
