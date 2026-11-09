package ibkr

import (
	"context"
	"encoding/binary"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func TestWaitUntilNextReqIDHandlesProtobuf(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClient(t)
	defer cleanup()

	resultCh := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		resultCh <- client.waitUntilNextReqID(ctx, client.conn)
	}()

	writeHandshakeMessage(t, serverConn, buildProtoMessage(common.NEXT_VALID_ID, &protobuf.NextValidId{
		OrderId: protofmt.Int32(42),
	}))
	if err := <-resultCh; err != nil {
		t.Fatalf("waitUntilNextReqID failed: %v", err)
	}
	if got := atomic.LoadInt32(&client.nextValidReqID); got != 42 {
		t.Fatalf("unexpected next valid request id: got %d want 42", got)
	}
}

func TestWaitUntilNextReqIDHandlesLegacy(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClient(t)
	defer cleanup()

	resultCh := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		resultCh <- client.waitUntilNextReqID(ctx, client.conn)
	}()

	writeHandshakeMessage(t, serverConn, message.NewRawEncoder().RawUInt32(common.NEXT_VALID_ID).Int32(17).Bytes())
	if err := <-resultCh; err != nil {
		t.Fatalf("waitUntilNextReqID failed: %v", err)
	}
	if got := atomic.LoadInt32(&client.nextValidReqID); got != 17 {
		t.Fatalf("unexpected next valid request id: got %d want 17", got)
	}
}

func TestStartAPIUsesLegacyMessageBeforeProtobufVersion(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClientWithVersion(t, common.ServerVersionProtobufRestMessages3-1)
	defer cleanup()
	client.clientID = 17

	if err := client.startApi(client.conn, Options{OptionalCapabilities: "foo"}); err != nil {
		t.Fatalf("startApi failed: %v", err)
	}

	payload := readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.START_API {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	msgDec := message.NewDecoder(payload[4:])
	if got := msgDec.Int32(); got != 2 {
		t.Fatalf("unexpected start api version: got %d", got)
	}
	if got := msgDec.Int32(); got != 17 {
		t.Fatalf("unexpected client id: got %d", got)
	}
	if got := msgDec.String(); got != "foo" {
		t.Fatalf("unexpected optional capabilities: got %q", got)
	}
	if msgDec.Err() != nil {
		t.Fatalf("failed decoding legacy start api request: %v", msgDec.Err())
	}
}

func TestStartAPIUsesProtobufMessageAtProtobufVersion(t *testing.T) {
	client, serverConn, cleanup := newMarketRuleRequestTestClientWithVersion(t, common.ServerVersionProtobufRestMessages3)
	defer cleanup()
	client.clientID = 17

	if err := client.startApi(client.conn, Options{OptionalCapabilities: "foo"}); err != nil {
		t.Fatalf("startApi failed: %v", err)
	}

	payload := readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.START_API+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected outbound msg id: got %d", got)
	}
	pb := &protobuf.StartApiRequest{}
	if err := proto.Unmarshal(payload[4:], pb); err != nil {
		t.Fatalf("failed decoding protobuf start api request: %v", err)
	}
	if pb.GetClientId() != 17 || pb.GetOptionalCapabilities() != "foo" {
		t.Fatalf("unexpected protobuf start api request: %#v", pb)
	}
}

// -----------------------------------------------------------------------------

func writeHandshakeMessage(t *testing.T, conn net.Conn, payload []byte) {
	t.Helper()
	frame := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(frame, uint32(len(payload)))
	copy(frame[4:], payload)
	if _, err := conn.Write(frame); err != nil {
		t.Fatalf("write handshake message: %v", err)
	}
}
