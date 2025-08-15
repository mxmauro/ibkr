package ibkr

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mxmauro/go-rundownprotection"
	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/connection"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/mxmauro/resetevent"
)

// -----------------------------------------------------------------------------

type Client struct {
	rp            rundownprotection.RundownProtection
	eventsHandler Events
	tzOffset      time.Duration

	wg          sync.WaitGroup
	destroyOnce sync.Once

	connMtx          sync.Mutex
	conn             *connection.Connection
	connErrHolder    atomic.Value
	isDisconnectedEv *resetevent.ManualResetEvent
	clientID         int32
	serverVersion    int32

	nextValidReqID        int32
	nextValidReqWithoutID int32

	reqMgr RequestManager
}

type Options struct {
	Address       string
	EventsHandler Events
	TzOffset      interface{}

	ConnectOptions       string
	OptionalCapabilities string
	ClientID             int32
}

// -----------------------------------------------------------------------------

// NewClient creates a new client object and establishes a connection to the given server.
func NewClient(ctx context.Context, opts Options) (*Client, error) {
	var tzOffset time.Duration
	var err error

	// Validate options
	if len(opts.Address) == 0 {
		return nil, errors.New("invalid host:port address")
	}
	if len(opts.ConnectOptions) > 0 && !utils.IsPrintableAsciiString(opts.ConnectOptions) {
		return nil, errors.New("invalid optional capabilities")
	}
	if len(opts.OptionalCapabilities) > 0 && !utils.IsPrintableAsciiString(opts.OptionalCapabilities) {
		return nil, errors.New("invalid optional capabilities")
	}
	if opts.TzOffset == nil {
		_, offset := time.Now().Local().Zone()
		tzOffset = time.Duration(offset) * time.Second
	} else {
		switch t := opts.TzOffset.(type) {
		case int:
			tzOffset = time.Duration(t) * time.Second
		case int64:
			tzOffset = time.Duration(t) * time.Second
		case time.Duration:
			tzOffset = t
		case string:
			tzOffset, err = time.ParseDuration(t)
			if err != nil {
				return nil, err
			}
		}
	}
	if opts.ClientID < 0 {
		return nil, errors.New("invalid client id")
	}

	// Create the client object
	c := Client{
		eventsHandler: opts.EventsHandler,
		tzOffset:      tzOffset,

		wg:          sync.WaitGroup{},
		destroyOnce: sync.Once{},

		connMtx:          sync.Mutex{},
		isDisconnectedEv: resetevent.NewManualResetEvent(),
	}
	c.rp.Initialize()
	c.initRequestManager()
	atomic.StoreInt32(&c.nextValidReqWithoutID, 1)

	// Try to connect to the server
	err = c.connectToServer(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Initiate the connection background worker
	c.wg.Add(1)
	go c.connectionWorker()

	// Done
	return &c, nil
}

// Destroy shuts down the current connection, cancels all pending requests and destroys the client object.
func (c *Client) Destroy() {
	c.destroyOnce.Do(func() {
		c.rp.Wait()
		c.wg.Wait()
	})
}

// IsConnected returns true if the connection to the server is alive.
func (c *Client) IsConnected() bool {
	connected := false
	if c.rp.Acquire() {
		select {
		case <-c.isDisconnectedEv.WaitCh():
		default:
			connected = true
		}
		c.rp.Release()
	}
	return connected
}

// ConnectedCh returns a channel closed if the connection goes down.
func (c *Client) ConnectedCh() <-chan struct{} {
	return c.isDisconnectedEv.WaitCh()
}

// ServerVersion returns the version of the server.
func (c *Client) ServerVersion() int {
	version := 0
	if c.rp.Acquire() {
		select {
		case <-c.isDisconnectedEv.WaitCh():
		default:
			version = int(c.serverVersion)
		}
		c.rp.Release()
	}
	return version
}

// RequestCurrentTime asks the current system time on the server side.
func (c *Client) RequestCurrentTime(ctx context.Context) (time.Time, error) {
	// Rundown protect
	if !c.rp.Acquire() {
		return time.Time{}, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.CurrentTimeResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  common.REQ_CURRENT_TIME_IN_MILLIS,
		Response: resp,
	})

	// Build the message to send
	msgEnc := message.NewEncoder().
		RawUInt32(common.REQ_CURRENT_TIME_IN_MILLIS)
	if msgEnc.Err() != nil {
		return time.Time{}, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return time.Time{}, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return time.Time{}, err
	}

	// Done
	return resp.CurrentTime, nil
}

// RequestManagedAccounts requests the list of managed accounts.
func (c *Client) RequestManagedAccounts(ctx context.Context) ([]string, error) {
	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.ManagedAccountsResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  common.REQ_MANAGED_ACCTS,
		Response: resp,
	})

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_MANAGED_ACCTS) {
		pb := protobuf.ManagedAccountsRequest{}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MANAGED_ACCTS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const VERSION = 1
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MANAGED_ACCTS).
			Int(VERSION)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp.Accounts, nil
}

// RequestHistoricalData retrieves historical market data.
func (c *Client) RequestHistoricalData(ctx context.Context, opts models.HistoricalDataRequestOptions) (*models.HistoricalDataResponse, error) {
	// Validate options
	if opts.Contract == nil {
		return nil, errors.New("invalid contract")
	}
	if opts.Duration < 1 {
		return nil, errors.New("invalid duration")
	}
	if len(opts.DurationUnit.String()) == 0 {
		return nil, errors.New("invalid duration units")
	}
	if opts.EndDate.IsZero() || opts.EndDate.Before(utils.EpochDate) {
		return nil, errors.New("invalid end date")
	}
	if len(opts.BarSize.String()) == 0 {
		return nil, errors.New("invalid bar size")
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
		return nil, errors.New("invalid what to show")
	}
	if len(opts.WhatToShow.String()) == 0 {
		return nil, errors.New("invalid what to show")
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.HistoricalDataResponse{
		Bars: make([]models.HistoricalDataBar, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_DATA,
		Response: resp,
	})

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_HISTORICAL_DATA) {
		pb := protobuf.HistoricalDataRequest{
			ReqId:          protofmt.Int32(req.ID()),
			Contract:       opts.Contract.Proto(nil),
			EndDateTime:    protofmt.String(opts.EndDate.Format("20060102-15:04:05")),
			BarSizeSetting: protofmt.String(opts.BarSize.String()),
			Duration:       protofmt.String(strconv.Itoa(opts.Duration) + " " + opts.DurationUnit.String()),
			UseRTH:         protofmt.Bool(opts.OnlyRegularTradingHours),
			WhatToShow:     protofmt.String(opts.WhatToShow.String()),
			FormatDate:     protofmt.Int32(2), // Return epoch timestamp
			KeepUpToDate:   protofmt.Bool(false),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_HISTORICAL_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(20).
			RawUInt32(common.REQ_HISTORICAL_DATA).
			RequestID(req.ID()).
			Marshal(opts.Contract, 2).
			String(opts.EndDate.Format("20060102-15:04:05")).
			String(opts.BarSize.String()).
			String(strconv.Itoa(opts.Duration) + " " + opts.DurationUnit.String()).
			Bool(opts.OnlyRegularTradingHours).
			String(opts.WhatToShow.String()).
			Int(2) // Return epoch timestamp
		if opts.Contract.SecType == models.SecurityTypePair {
			msgEnc.Int(len(opts.Contract.ComboLegs))
			for _, comboLeg := range opts.Contract.ComboLegs {
				msgEnc.Marshal(comboLeg, 1)
			}
		}
		msgEnc.Bool(false). // KeepUpToDate
					Marshal(&models.TagValueList{}, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

// RequestHistoricalTicks retrieves historical market ticks.
func (c *Client) RequestHistoricalTicks(ctx context.Context, opts models.HistoricalTicksRequestOptions) (*models.HistoricalTicksResponse, error) {
	// Validate options
	if opts.Contract == nil {
		return nil, errors.New("invalid contract")
	}
	if opts.StartDate.IsZero() || opts.StartDate.Before(utils.EpochDate) {
		return nil, errors.New("invalid start date")
	}
	if opts.EndDate.IsZero() || opts.EndDate.Before(opts.StartDate) {
		return nil, errors.New("invalid end date")
	}
	if opts.NumberOfTicks < 1 || opts.NumberOfTicks > 1000 {
		return nil, errors.New("invalid number of ticks")
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.HistoricalTicksResponse{}
	switch opts.WhatToShow {
	case models.WhatToShowMidPoint:
		resp.Ticks = make([]models.HistoricalTick, 0)
	case models.WhatToShowBidAsk:
		resp.TicksBidAsk = make([]models.HistoricalTickBidAsk, 0)
	case models.WhatToShowTrades:
		resp.TicksLast = make([]models.HistoricalTickLast, 0)
	default:
		return nil, errors.New("invalid what to show")
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_TICKS,
		Response: resp,
	})

	// Build the message to send
	var msgEnc *message.Encoder
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

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

// RequestContractDetails retrieves all details for contracts matching the provided input.
func (c *Client) RequestContractDetails(ctx context.Context, opts models.ContractDetailsRequestOptions) (*models.ContractDetailsResponse, error) {
	// Validate options
	if opts.Contract == nil {
		return nil, errors.New("invalid contract")
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.ContractDetailsResponse{
		ContractDetails: make([]*models.ContractDetails, 0),
	}

	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_CONTRACT_DATA,
		Response: resp,
	})

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_CONTRACT_DATA) {
		pb := protobuf.ContractDataRequest{
			ReqId:    protofmt.Int32(req.ID()),
			Contract: opts.Contract.Proto(nil),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_CONTRACT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const VERSION = 8
		msgEnc = message.NewEncoder().Reserve(21).
			RawUInt32(common.REQ_CONTRACT_DATA).
			RequestID(req.ID()).
			Int(VERSION).
			Marshal(opts.Contract, 3)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

// RequestMatchingSymbols retrieves all details for contracts matching the provided input.
func (c *Client) RequestMatchingSymbols(ctx context.Context, opts models.MatchingSymbolsRequestOptions) (*models.MatchingSymbolsResponse, error) {
	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.MatchingSymbolsResponse{
		ContractDescriptions: make([]*models.ContractDescription, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_MATCHING_SYMBOLS,
		Response: resp,
	})

	// Build the message to send
	msgEnc := message.NewEncoder().Reserve(3).
		RawUInt32(common.REQ_MATCHING_SYMBOLS).
		RequestID(req.ID()).
		String(opts.Pattern)
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	// Wait until the response is fulfilled
	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

func (c *Client) RequestMarketDataType(_ context.Context, opts models.MarketDataTypeRequestOptions) error {
	// Validate options
	if len(opts.Type.String()) == 0 {
		return errors.New("invalid market data type")
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return net.ErrClosed
	}
	defer c.rp.Release()

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_MARKET_DATA_TYPE) {
		pb := protobuf.MarketDataTypeRequest{
			MarketDataType: protofmt.Int32(int32(opts.Type)),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MARKET_DATA_TYPE + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const VERSION = 1
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_MARKET_DATA_TYPE).
			Int(VERSION).
			Int(int(opts.Type))
	}
	if msgEnc.Err() != nil {
		return msgEnc.Err()
	}

	// Send it
	err := c.sendMessage(msgEnc.Bytes())
	if err != nil {
		return err
	}

	// Done
	return nil
}

func (c *Client) RequestTopMarketData(_ context.Context, opts models.TopMarketDataRequestOptions) (*models.TopMarketDataResponse, error) {
	// Validate options
	if opts.Contract == nil {
		return nil, errors.New("invalid contract")
	}
	if opts.Snapshot && len(opts.AdditionalGenericTicks) > 0 {
		return nil, errors.New("generic ticks cannot be used when requesting a snapshot")
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.TopMarketDataResponse{
		Channel: make(chan models.TopMarketData, 4),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_MKT_DATA,
		Response: resp,
		CompleteCB: func(req *Request, err error) {
			close(resp.Channel)
		},
	})
	resp.Cancel = func() {
		c.cancelTopMarketData(req)
	}
	resp.Err = func() error {
		return req.Err()
	}

	// Build the message to send
	var msgEnc *message.Encoder

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
		const VERSION = 11
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_MKT_DATA).
			Int(VERSION).
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

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

func (c *Client) RequestMarketDepthData(_ context.Context, opts models.MarketDepthDataRequestOptions) (*models.MarketDepthDataResponse, error) {
	// Validate options
	if opts.Contract == nil {
		return nil, errors.New("invalid contract")
	}
	if opts.RowsCount < 0 || opts.RowsCount > 1000 {
		return nil, errors.New("invalid rows count")
	}
	if opts.RowsCount == 0 {
		opts.RowsCount = 20
	}

	// Rundown protect
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	// Create the new request and response holder
	resp := &models.MarketDepthDataResponse{
		Book:    models.NewMarketDepthBook(opts.RowsCount),
		Channel: make(chan models.MarketDepthData, 4),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithTickerID,
		MsgCode:  common.REQ_MKT_DEPTH,
		Response: resp,
		CompleteCB: func(req *Request, err error) {
			close(resp.Channel)
		},
	})
	resp.Cancel = func() {
		c.cancelMarketDepthData(req, opts.SmartDepth)
	}
	resp.Err = func() error {
		return req.Err()
	}

	// Build the message to send
	var msgEnc *message.Encoder
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
		const VERSION = 5
		msgEnc = message.NewEncoder().Reserve(17).
			RawUInt32(common.REQ_MKT_DEPTH).
			Int(VERSION).
			RequestID(req.ID()).
			Marshal(opts.Contract, 1).
			Int(opts.RowsCount).
			Bool(opts.SmartDepth).
			Marshal(&models.TagValueList{}, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	// Send it
	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}

	// Done
	return resp, nil
}

func (c *Client) cancelTopMarketData(req *Request) {
	// Rundown protect
	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.CANCEL_MKT_DATA) {
		pb := protobuf.CancelMarketData{
			ReqId: protofmt.Int32(req.ID()),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_MKT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const VERSION = 2
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.CANCEL_MKT_DATA).
			Int(VERSION).
			RequestID(req.ID())
	}

	// Send it
	_ = c.sendMessage(msgEnc.Bytes())

	// Remove the request from the manager
	c.reqMgr.removeRequest(req, nil)
}

func (c *Client) cancelMarketDepthData(req *Request, isSmartDepth bool) {
	// Rundown protect
	if !c.rp.Acquire() {
		return
	}
	defer c.rp.Release()

	// Build the message to send
	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.CANCEL_MKT_DEPTH) {
		pb := protobuf.CancelMarketDepth{
			ReqId:        protofmt.Int32(req.ID()),
			IsSmartDepth: protofmt.Bool(isSmartDepth),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.CANCEL_MKT_DEPTH + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const VERSION = 1
		msgEnc = message.NewEncoder().Reserve(4).
			RawUInt32(common.CANCEL_MKT_DEPTH).
			Int(VERSION).
			RequestID(req.ID()).
			Bool(isSmartDepth)
	}

	// Send it
	_ = c.sendMessage(msgEnc.Bytes())

	// Remove the request from the manager
	c.reqMgr.removeRequest(req, nil)
}

func (c *Client) isProtoBufAvailable(msgType uint32) bool {
	minServerVer, ok := common.PROTOBUF_MSG_IDS[msgType]
	return ok && c.serverVersion >= minServerVer
}
