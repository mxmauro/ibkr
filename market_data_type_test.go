package ibkr

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

func TestTopMarketDataStreamResponseUsesConfiguredMarketDataType(t *testing.T) {
	streamResp := newTopMarketDataStreamResponse(models.MarketDataTypeDelayed)
	defer streamResp.close()

	data := models.NewTopMarketDataPrice(models.TickTypeBid)
	streamResp.enqueue(data)

	got := waitTopMarketData(t, streamResp.resp.Channel)
	if got.MarketDataType() != models.MarketDataTypeDelayed {
		t.Fatalf("unexpected market data type: got %v", got.MarketDataType())
	}
	if !got.Delayed() {
		t.Fatal("expected delayed flag to be set")
	}
	if got.Frozen() {
		t.Fatal("did not expect frozen flag to be set")
	}
}

func TestTopMarketDataStreamResponseDelayedTickOverridesConfiguredType(t *testing.T) {
	streamResp := newTopMarketDataStreamResponse(models.MarketDataTypeRealtime)
	defer streamResp.close()

	data := models.NewTopMarketDataPrice(models.TickTypeDelayedBid)
	streamResp.enqueue(data)

	got := waitTopMarketData(t, streamResp.resp.Channel)
	if got.MarketDataType() != models.MarketDataTypeDelayed {
		t.Fatalf("unexpected market data type: got %v", got.MarketDataType())
	}
	if !got.Delayed() {
		t.Fatal("expected delayed flag to be set")
	}
	if got.Frozen() {
		t.Fatal("did not expect frozen flag to be set")
	}
}

func TestProcessMarketDataTypeCommonUpdatesTopMarketStream(t *testing.T) {
	c := &Client{}
	c.initRequestManager()
	atomic.StoreInt32(&c.nextValidReqID, 1)

	streamResp := newTopMarketDataStreamResponse(models.MarketDataTypeDelayed)
	defer streamResp.close()

	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_MKT_DATA,
		Response: streamResp,
	})
	c.reqMgr.addRequest(req)
	defer c.reqMgr.removeRequest(req, nil)

	if err := c.processMarketDataTypeCommon(req.ID(), models.MarketDataTypeFrozen); err != nil {
		t.Fatalf("processMarketDataTypeCommon failed: %v", err)
	}

	data := models.NewTopMarketDataPrice(models.TickTypeBid)
	streamResp.enqueue(data)

	got := waitTopMarketData(t, streamResp.resp.Channel)
	if got.MarketDataType() != models.MarketDataTypeFrozen {
		t.Fatalf("unexpected market data type: got %v", got.MarketDataType())
	}
	if got.Delayed() {
		t.Fatal("did not expect delayed flag to be set")
	}
	if !got.Frozen() {
		t.Fatal("expected frozen flag to be set")
	}
}

func waitTopMarketData(t *testing.T, ch <-chan models.TopMarketData) models.TopMarketData {
	t.Helper()

	select {
	case data := <-ch:
		return data
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for top market data")
		return nil
	}
}
