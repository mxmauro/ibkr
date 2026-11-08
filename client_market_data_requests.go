package ibkr

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// RequestMarketDataType sets the preferred market data mode for subsequent subscriptions.
func (c *Client) RequestMarketDataType(_ context.Context, opts models.MarketDataTypeRequestOptions) error {
	var msgEnc *message.Encoder

	if len(opts.Type.String()) == 0 {
		return errInvalidMarketDataType
	}
	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.REQ_MARKET_DATA_TYPE) {
		pb := protobuf.MarketDataTypeRequest{
			MarketDataType: protofmt.Int32(int32(opts.Type)),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MARKET_DATA_TYPE + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 1
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_MARKET_DATA_TYPE).
			Int(version).
			Int(int(opts.Type))
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	err := c.sendMessage(msgEnc.Bytes())
	if err != nil {
		return err
	}
	atomic.StoreInt32(&c.marketDataType, int32(opts.Type))

	return nil
}

// RequestTopMarketData subscribes to top-of-book market data for a contract.
func (c *Client) RequestTopMarketData(ctx context.Context, opts models.TopMarketDataRequestOptions) (*models.TopMarketDataResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.Snapshot && len(opts.AdditionalGenericTicks) > 0 {
		return nil, errGenericTicksWithSnapshot
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newTopMarketDataStreamResponse(c.currentMarketDataType())
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_MKT_DATA,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelTopMarketData(req)
	})
	resp.Err = func() error {
		return req.Err()
	}

	genericTickSB := strings.Builder{}
	for idx, gt := range opts.AdditionalGenericTicks {
		if idx > 0 {
			_, _ = genericTickSB.WriteRune(',')
		}
		_, _ = genericTickSB.WriteString(strconv.Itoa(int(gt)))
	}

	if c.isProtoBufAvailable(common.REQ_MKT_DATA) {
		pb := protobuf.MarketDataRequest{
			ReqId:              protofmt.Int32(req.ID()),
			Contract:           opts.Contract.Proto(nil),
			GenericTickList:    protofmt.String(genericTickSB.String()),
			Snapshot:           protofmt.Bool(opts.Snapshot),
			RegulatorySnapshot: protofmt.Bool(opts.RegulatorySnapshot),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MKT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 11
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_MKT_DATA).
			Int(version).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1)
		if opts.Contract.SecType == models.SecurityTypePair {
			msgEnc.Int(len(opts.Contract.ComboLegs))
			for _, comboLeg := range opts.Contract.ComboLegs {
				msgEnc.Marshal(comboLeg, 1)
			}
		}
		if opts.Contract.DeltaNeutralContract != nil {
			msgEnc.Bool(true).
				Marshal(opts.Contract.DeltaNeutralContract, 1)
		} else {
			msgEnc.Bool(false)
		}
		msgEnc.String(genericTickSB.String())
		msgEnc.Bool(opts.Snapshot)
		msgEnc.Bool(opts.RegulatorySnapshot).
			Marshal(&models.TagValueList{}, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}

	bindStreamRequestContext(ctx, req, resp.Cancel)
	return resp, nil
}

// RequestMarketDepthData subscribes to market depth updates for a contract.
func (c *Client) RequestMarketDepthData(ctx context.Context, opts models.MarketDepthDataRequestOptions) (*models.MarketDepthDataResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.RowsCount < 0 || opts.RowsCount > 1000 {
		return nil, errInvalidRowsCount
	}
	if opts.RowsCount == 0 {
		opts.RowsCount = 20
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newMarketDepthDataStreamResponse(opts.RowsCount)
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_MKT_DEPTH,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelMarketDepthData(req, opts.SmartDepth)
	})
	resp.Err = func() error {
		return req.Err()
	}

	if c.isProtoBufAvailable(common.REQ_MKT_DEPTH) {
		pb := protobuf.MarketDepthRequest{
			ReqId:        protofmt.Int32(req.ID()),
			Contract:     opts.Contract.Proto(nil),
			NumRows:      protofmt.Int32(int32(opts.RowsCount)),
			IsSmartDepth: protofmt.Bool(opts.SmartDepth),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MKT_DEPTH + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 5
		msgEnc = message.NewEncoder().Reserve(17).
			RawUInt32(common.REQ_MKT_DEPTH).
			Int(version).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			Int(opts.RowsCount).
			Bool(opts.SmartDepth).
			Marshal(&models.TagValueList{}, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}

	bindStreamRequestContext(ctx, req, resp.Cancel)
	return resp, nil
}

// RequestRealTimeBars subscribes to streaming real-time bars.
func (c *Client) RequestRealTimeBars(ctx context.Context, opts models.RealTimeBarsRequestOptions) (*models.RealTimeBarsResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.BarSize == 0 {
		opts.BarSize = 5
	}
	if opts.BarSize < 1 {
		return nil, errInvalidBarSize
	}
	switch opts.WhatToShow {
	case models.WhatToShowTrades:
	case models.WhatToShowMidPoint:
	case models.WhatToShowBid:
	case models.WhatToShowAsk:
	default:
		return nil, errInvalidWhatToShow
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newRealTimeBarsStreamResponse(int64(opts.BarSize))
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_REAL_TIME_BARS,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelRealTimeBars(req)
	})
	resp.Err = func() error {
		return req.Err()
	}

	if c.isProtoBufAvailable(common.REQ_REAL_TIME_BARS) {
		pb := protobuf.RealTimeBarsRequest{
			ReqId:      protofmt.Int32(req.ID()),
			Contract:   opts.Contract.Proto(nil),
			BarSize:    protofmt.Int32(int32(opts.BarSize)),
			WhatToShow: protofmt.String(opts.WhatToShow.String()),
			UseRTH:     protofmt.Bool(opts.OnlyRegularTradingHours),
		}
		if realTimeBarsOptions := tagValueListToMap(opts.Options); len(realTimeBarsOptions) > 0 {
			pb.RealTimeBarsOptions = realTimeBarsOptions
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_REAL_TIME_BARS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 3
		msgEnc = message.NewEncoder().Reserve(19).
			RawUInt32(common.REQ_REAL_TIME_BARS).
			Int(version).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			Int(opts.BarSize).
			String(opts.WhatToShow.String()).
			Bool(opts.OnlyRegularTradingHours).
			Marshal(&opts.Options, 1)
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

// RequestTickByTickData subscribes to tick-by-tick market data.
func (c *Client) RequestTickByTickData(ctx context.Context, opts models.TickByTickRequestOptions) (*models.TickByTickResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	switch opts.TickType {
	case models.TickByTickTypeLast:
	case models.TickByTickTypeAllLast:
	case models.TickByTickTypeBidAsk:
	case models.TickByTickTypeMidPoint:
	default:
		return nil, errInvalidTickByTickType
	}
	if opts.NumberOfTicks < 0 {
		return nil, errInvalidNumberOfTicks
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newTickByTickStreamResponse()
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_TICK_BY_TICK_DATA,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelTickByTickData(req)
	})
	resp.Err = func() error {
		return req.Err()
	}

	if c.isProtoBufAvailable(common.REQ_TICK_BY_TICK_DATA) {
		pb := protobuf.TickByTickRequest{
			ReqId:         protofmt.Int32(req.ID()),
			Contract:      opts.Contract.Proto(nil),
			TickType:      protofmt.String(opts.TickType.String()),
			NumberOfTicks: protofmt.Int32(int32(opts.NumberOfTicks)),
			IgnoreSize:    protofmt.Bool(opts.IgnoreSize),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_TICK_BY_TICK_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(17).
			RawUInt32(common.REQ_TICK_BY_TICK_DATA).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			String(opts.TickType.String()).
			Int(opts.NumberOfTicks).
			Bool(opts.IgnoreSize)
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

func (c *Client) cancelTopMarketData(req *Request) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_MKT_DATA) {
		pb := protobuf.CancelMarketData{
			ReqId: protofmt.Int32(req.ID()),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_MKT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 2
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.CANCEL_MKT_DATA).
			Int(version).
			RequestID(req.ID())
	}

	_ = c.sendMessage(msgEnc.Bytes())
	c.reqMgr.removeRequest(req, nil)
}

func (c *Client) cancelMarketDepthData(req *Request, isSmartDepth bool) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_MKT_DEPTH) {
		pb := protobuf.CancelMarketDepth{
			ReqId:        protofmt.Int32(req.ID()),
			IsSmartDepth: protofmt.Bool(isSmartDepth),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_MKT_DEPTH + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 1
		msgEnc = message.NewEncoder().Reserve(4).
			RawUInt32(common.CANCEL_MKT_DEPTH).
			Int(version).
			RequestID(req.ID()).
			Bool(isSmartDepth)
	}

	_ = c.sendMessage(msgEnc.Bytes())
	c.reqMgr.removeRequest(req, nil)
}

func (c *Client) cancelRealTimeBars(req *Request) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_REAL_TIME_BARS) {
		pb := protobuf.CancelRealTimeBars{
			ReqId: protofmt.Int32(req.ID()),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_REAL_TIME_BARS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 1
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.CANCEL_REAL_TIME_BARS).
			Int(version).
			RequestID(req.ID())
	}

	_ = c.sendMessage(msgEnc.Bytes())
	c.reqMgr.removeRequest(req, nil)
}

func (c *Client) cancelTickByTickData(req *Request) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_TICK_BY_TICK_DATA) {
		pb := protobuf.CancelTickByTick{
			ReqId: protofmt.Int32(req.ID()),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_TICK_BY_TICK_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(2).
			RawUInt32(common.CANCEL_TICK_BY_TICK_DATA).
			RequestID(req.ID())
	}

	_ = c.sendMessage(msgEnc.Bytes())
	c.reqMgr.removeRequest(req, nil)
}
