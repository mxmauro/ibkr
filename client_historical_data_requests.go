package ibkr

import (
	"context"
	"net"
	"strconv"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// RequestHistoricalData retrieves historical market data.
func (c *Client) RequestHistoricalData(ctx context.Context, opts models.HistoricalDataRequestOptions) (*models.HistoricalDataResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.Duration < 1 {
		return nil, errInvalidDuration
	}
	if len(opts.DurationUnit.String()) == 0 {
		return nil, errInvalidDurationUnits
	}
	if !opts.KeepUpToDate && (opts.EndDate.IsZero() || opts.EndDate.Before(utils.EpochDate)) {
		return nil, errInvalidEndDate
	}
	if len(opts.BarSize.String()) == 0 {
		return nil, errInvalidBarSize
	}
	switch opts.WhatToShow {
	case models.WhatToShowTrades:
	case models.WhatToShowMidPoint:
	case models.WhatToShowBid:
	case models.WhatToShowAsk:
	case models.WhatToShowBidAsk:
	case models.WhatToShowHistoricalVolatility:
	case models.WhatToShowOptionImpliedVolatility:
	default:
		return nil, errInvalidWhatToShow
	}
	if len(opts.WhatToShow.String()) == 0 {
		return nil, errInvalidWhatToShow
	}
	endDateTime := ""
	if !opts.KeepUpToDate {
		endDateTime = opts.EndDate.Format("20060102-15:04:05")
	}

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newHistoricalDataStreamResponse(opts.KeepUpToDate, opts.BarSize.Seconds())
	resp := streamResp.resp
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_DATA,
		Response: streamResp,
		CompleteCB: func(req *Request, err error) {
			streamResp.close()
		},
	})
	resp.Cancel = newRequestCancelFunc(func() {
		c.cancelHistoricalData(req)
	})
	resp.Err = func() error {
		return req.Err()
	}

	if c.isProtoBufAvailable(common.REQ_HISTORICAL_DATA) {
		pb := protobuf.HistoricalDataRequest{
			ReqId:          protofmt.Int32(req.ID()),
			Contract:       opts.Contract.Proto(nil),
			EndDateTime:    protofmt.String(endDateTime),
			BarSizeSetting: protofmt.String(opts.BarSize.String()),
			Duration:       protofmt.String(strconv.Itoa(opts.Duration) + " " + opts.DurationUnit.String()),
			UseRTH:         protofmt.Bool(opts.OnlyRegularTradingHours),
			WhatToShow:     protofmt.String(opts.WhatToShow.String()),
			FormatDate:     protofmt.Int32(2),
			KeepUpToDate:   protofmt.Bool(opts.KeepUpToDate),
		}
		if chartOptions := tagValueListToMap(opts.ChartOptions); len(chartOptions) > 0 {
			pb.ChartOptions = chartOptions
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HISTORICAL_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(20).
			RawUInt32(common.REQ_HISTORICAL_DATA).
			RequestID(req.ID()).
			Marshal(opts.Contract, 2).
			String(endDateTime).
			String(opts.BarSize.String()).
			String(strconv.Itoa(opts.Duration) + " " + opts.DurationUnit.String()).
			Bool(opts.OnlyRegularTradingHours).
			String(opts.WhatToShow.String()).
			Int(2)
		if opts.Contract.SecType == models.SecurityTypePair {
			msgEnc.Int(len(opts.Contract.ComboLegs))
			for _, comboLeg := range opts.Contract.ComboLegs {
				msgEnc.Marshal(comboLeg, 1)
			}
		}
		msgEnc.Bool(opts.KeepUpToDate).
			Marshal(&opts.ChartOptions, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	if opts.KeepUpToDate {
		bindStreamRequestContext(ctx, req, resp.Cancel)
		return resp, nil
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RequestHistoricalTicks retrieves historical market ticks.
func (c *Client) RequestHistoricalTicks(ctx context.Context, opts models.HistoricalTicksRequestOptions) (*models.HistoricalTicksResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if opts.StartDate.IsZero() || opts.StartDate.Before(utils.EpochDate) {
		return nil, errInvalidStartDate
	}
	if opts.EndDate.IsZero() || opts.EndDate.Before(opts.StartDate) {
		return nil, errInvalidEndDate
	}
	if opts.NumberOfTicks < 1 || opts.NumberOfTicks > 1000 {
		return nil, errInvalidNumberOfTicks
	}

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.HistoricalTicksResponse{}
	switch opts.WhatToShow {
	case models.WhatToShowMidPoint:
		resp.Ticks = make([]models.HistoricalTick, 0)
	case models.WhatToShowBidAsk:
		resp.TicksBidAsk = make([]models.HistoricalTickBidAsk, 0)
	case models.WhatToShowTrades:
		resp.TicksLast = make([]models.HistoricalTickLast, 0)
	default:
		return nil, errInvalidWhatToShow
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_TICKS,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_HISTORICAL_TICKS) {
		pb := protobuf.HistoricalTicksRequest{
			ReqId:         protofmt.Int32(req.ID()),
			Contract:      opts.Contract.Proto(nil),
			StartDateTime: protofmt.String(opts.StartDate.Format("20060102-15:04:05")),
			EndDateTime:   protofmt.String(opts.EndDate.Format("20060102-15:04:05")),
			NumberOfTicks: protofmt.Int32(int32(opts.NumberOfTicks)),
			WhatToShow:    protofmt.String(opts.WhatToShow.String()),
			UseRTH:        protofmt.Bool(opts.OnlyRegularTradingHours),
			IgnoreSize:    protofmt.Bool(opts.IgnoreSize),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HISTORICAL_TICKS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(22).
			RawUInt32(common.REQ_HISTORICAL_TICKS).
			RequestID(req.ID()).
			Marshal(opts.Contract, 2).
			String(opts.StartDate.Format("20060102-15:04:05")).
			String(opts.EndDate.Format("20060102-15:04:05")).
			Int(opts.NumberOfTicks).
			String(opts.WhatToShow.String()).
			Bool(opts.OnlyRegularTradingHours).
			Bool(opts.IgnoreSize).
			Marshal(&models.TagValueList{}, 1)
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

func (c *Client) cancelHistoricalData(req *Request) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	if c.isProtoBufAvailable(common.CANCEL_HISTORICAL_DATA) {
		pb := protobuf.CancelHistoricalData{
			ReqId: protofmt.Int32(req.ID()),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_HISTORICAL_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 1
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.CANCEL_HISTORICAL_DATA).
			Int(version).
			RequestID(req.ID())
	}

	_ = c.sendMessage(msgEnc.Bytes())
	c.reqMgr.removeRequest(req, nil)
}
