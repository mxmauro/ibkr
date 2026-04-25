package ibkr

import (
	"testing"
	"time"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

func TestHistoricalDataStreamResponseApplyUpdateReplacesLastBar(t *testing.T) {
	streamResp := newHistoricalDataStreamResponse(true, int64(models.BarSizeOneMinute.Seconds()))

	firstBar := models.HistoricalDataBar{
		Date:  time.Date(2026, time.January, 2, 15, 4, 0, 0, time.UTC),
		Open:  10,
		Close: 11,
	}
	streamResp.appendBars([]models.HistoricalDataBar{firstBar})

	updateBar := firstBar
	updateBar.Close = 12
	streamResp.applyUpdate(updateBar)

	if got := len(streamResp.resp.Bars); got != 1 {
		t.Fatalf("unexpected bars count: got %d, want 1", got)
	}
	if got := streamResp.resp.Bars[0].Close; got != 12 {
		t.Fatalf("unexpected updated close: got %f, want 12", got)
	}

	streamResp.close()
}

func TestHistoricalDataEndCommonKeepsStreamingRequestOpenWhenKeepUpToDate(t *testing.T) {
	c := &Client{}
	c.initRequestManager()

	streamResp := newHistoricalDataStreamResponse(true, int64(models.BarSizeOneMinute.Seconds()))
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_DATA,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	c.reqMgr.addRequest(req)

	if err := c.processHistoricalDataEndCommon(req.ID()); err != nil {
		t.Fatal(err)
	}

	select {
	case <-req.CompleteCh():
		t.Fatal("request completed unexpectedly")
	default:
	}

	c.reqMgr.removeRequest(req, nil)
}
