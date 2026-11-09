package ibkr

import (
	"context"
	"math"
	"net"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// RequestMarketRule retrieves the price increment table for a market rule id.
func (c *Client) RequestMarketRule(ctx context.Context, opts models.MarketRuleRequestOptions) (*models.MarketRuleResponse, error) {
	if opts.MarketRuleID <= 0 || opts.MarketRuleID > math.MaxInt32 {
		return nil, errInvalidMarketRuleID
	}

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.MarketRuleResponse{}
	req := c.createRequest(RequestOptions{
		Type:         RequestTypeRequestWithoutID,
		MarketRuleID: opts.MarketRuleID,
		MsgCode:      int(common.MARKET_RULE),
		Response:     resp,
	})

	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_MARKET_RULE) {
		pb := protobuf.MarketRuleRequest{MarketRuleId: protofmt.Int32(int32(opts.MarketRuleID))}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MARKET_RULE + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MARKET_RULE).
			Int64(opts.MarketRuleID)
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

func (c *Client) RequestHeadTimestamp(
	ctx context.Context, opts models.HeadTimestampRequestOptions,
) (*models.HeadTimestampResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if len(opts.WhatToShow.String()) == 0 {
		return nil, errInvalidWhatToShow
	}
	if opts.FormatDate == 0 {
		opts.FormatDate = 2
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.HeadTimestampResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HEAD_TIMESTAMP,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_HEAD_TIMESTAMP) {
		pb := protobuf.HeadTimestampRequest{
			ReqId:      protofmt.Int32(req.ID()),
			Contract:   opts.Contract.Proto(nil),
			UseRTH:     protofmt.Bool(opts.OnlyRegularTradingHours),
			WhatToShow: protofmt.String(opts.WhatToShow.String()),
			FormatDate: protofmt.Int32(int32(opts.FormatDate)),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HEAD_TIMESTAMP + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HEAD_TIMESTAMP).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			Bool(opts.OnlyRegularTradingHours).
			String(opts.WhatToShow.String()).
			Int(opts.FormatDate)
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

func (c *Client) CancelHeadTimestamp(reqID int32) error {
	var msgEnc *message.Encoder

	if reqID <= 0 {
		return errInvalidRequestID
	}
	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_HEAD_TIMESTAMP) {
		pb := protobuf.CancelHeadTimestamp{ReqId: protofmt.Int32(reqID)}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_HEAD_TIMESTAMP + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_HEAD_TIMESTAMP).
			Int(1).
			RequestID(reqID)
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

func (c *Client) RequestHistogramData(ctx context.Context, opts models.HistogramDataRequestOptions) (*models.HistogramDataResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.TimePeriod == "" {
		return nil, errInvalidDuration
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.HistogramDataResponse{Data: make([]models.HistogramData, 0)}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTOGRAM_DATA,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_HISTOGRAM_DATA) {
		pb := protobuf.HistogramDataRequest{
			ReqId:      protofmt.Int32(req.ID()),
			Contract:   opts.Contract.Proto(nil),
			UseRTH:     protofmt.Bool(opts.OnlyRegularTradingHours),
			TimePeriod: protofmt.String(opts.TimePeriod),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HISTOGRAM_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HISTOGRAM_DATA).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			Bool(opts.OnlyRegularTradingHours).
			String(opts.TimePeriod)
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

func (c *Client) CancelHistogramData(reqID int32) error {
	var msgEnc *message.Encoder

	if reqID <= 0 {
		return errInvalidRequestID
	}
	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_HISTOGRAM_DATA) {
		pb := protobuf.CancelHistogramData{ReqId: protofmt.Int32(reqID)}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_HISTOGRAM_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_HISTOGRAM_DATA).
			Int(1).
			RequestID(reqID)
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	return c.sendMessage(msgEnc.Bytes())
}

func (c *Client) RequestMktDepthExchanges(ctx context.Context) (*models.MktDepthExchangesResponse, error) {
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.MktDepthExchangesResponse{Exchanges: make([]models.DepthMktDataDescription, 0)}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.MKT_DEPTH_EXCHANGES),
		Response: resp,
	})

	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_MKT_DEPTH_EXCHANGES) {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MKT_DEPTH_EXCHANGES + common.PROTOBUF_MSG_ID).
			Proto(&protobuf.MarketDepthExchangesRequest{})
	} else {
		msgEnc = message.NewEncoder().RawUInt32(common.REQ_MKT_DEPTH_EXCHANGES)
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

func (c *Client) RequestSmartComponents(
	ctx context.Context, opts models.SmartComponentsRequestOptions,
) (*models.SmartComponentsResponse, error) {
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.SmartComponentsResponse{Components: make([]models.SmartComponent, 0)}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_SMART_COMPONENTS,
		Response: resp,
	})

	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_SMART_COMPONENTS) {
		pb := protobuf.SmartComponentsRequest{
			ReqId:       protofmt.Int32(req.ID()),
			BboExchange: protofmt.String(opts.BBOExchange),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_SMART_COMPONENTS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_SMART_COMPONENTS).
			RequestID(req.ID()).
			String(opts.BBOExchange)
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
