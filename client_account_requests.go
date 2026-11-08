package ibkr

import (
	"context"
	"net"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func (c *Client) RequestPositions(ctx context.Context) (*models.PositionsResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.PositionsResponse{Positions: make([]models.Position, 0)}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.POSITION_DATA),
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_POSITIONS) {
		pb := protobuf.PositionsRequest{}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_POSITIONS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_POSITIONS).
			Int(1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CancelPositions() error {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_POSITIONS) {
		pb := protobuf.CancelPositions{}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_POSITIONS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_POSITIONS).
			Int(1)
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

func (c *Client) RequestAccountSummary(
	ctx context.Context, opts models.AccountSummaryRequestOptions,
) (*models.AccountSummaryResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.AccountSummaryResponse{Values: make([]models.AccountSummaryValue, 0)}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_ACCOUNT_SUMMARY,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_ACCOUNT_SUMMARY) {
		pb := protobuf.AccountSummaryRequest{
			ReqId: protofmt.Int32(req.ID()),
			Group: protofmt.String(opts.Group),
			Tags:  protofmt.String(opts.Tags),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCOUNT_SUMMARY + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCOUNT_SUMMARY).
			Int(1).
			RequestID(req.ID()).
			String(opts.Group).
			String(opts.Tags)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CancelAccountSummary(reqID int32) error {
	var msgEnc *message.Encoder

	if reqID <= 0 {
		return errInvalidRequestID
	}
	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_ACCOUNT_SUMMARY) {
		pb := protobuf.CancelAccountSummary{ReqId: protofmt.Int32(reqID)}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_ACCOUNT_SUMMARY + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_ACCOUNT_SUMMARY).
			Int(1).
			RequestID(reqID)
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

func (c *Client) RequestAccountUpdates(
	ctx context.Context, opts models.AccountUpdatesRequestOptions,
) (*models.AccountUpdatesResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newAccountUpdatesStreamResponse()
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.ACCT_VALUE),
		Response: streamResp,
		CompleteCB: func(_ *Request, _ error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelAccountUpdates(req, opts.AccountCode)
	})
	resp.Err = func() error {
		return req.Err()
	}

	if c.isProtoBufAvailable(common.REQ_ACCT_DATA) {
		pb := protobuf.AccountDataRequest{
			Subscribe: protofmt.Bool(true),
			AcctCode:  protofmt.String(opts.AccountCode),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCT_DATA).
			Int(2).
			Bool(true).
			String(opts.AccountCode)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}

	bindStreamRequestContext(ctx, req, resp.Cancel)
	return resp, nil
}

func (c *Client) cancelAccountUpdates(req *Request, accountCode string) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.REQ_ACCT_DATA) {
		pb := protobuf.AccountDataRequest{
			Subscribe: protofmt.Bool(false),
			AcctCode:  protofmt.String(accountCode),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_ACCT_DATA).
			Int(2).
			Bool(false).
			String(accountCode)
	}
	if msgEnc.Err() == nil {
		_ = c.sendMessage(msgEnc.Bytes())
	}

	c.reqMgr.removeRequest(req, nil)
}
