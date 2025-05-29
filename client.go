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
	"github.com/mxmauro/ibkr/utils"
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

// NewClient creates  new client object and establishes a connection to the given server.
func NewClient(ctx context.Context, opts Options) (*Client, error) {
	var tzOffset time.Duration
	var err error

	// Validate options
	if len(opts.Address) == 0 {
		return nil, errors.New("invalid host:port address")
	}
	if opts.EventsHandler == nil {
		return nil, errors.New("invalid incoming message handler")
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

// ConnectedCh returns a channel that is closed if the connection goes down.
func (c *Client) ConnectedCh() <-chan struct{} {
	return c.isDisconnectedEv.WaitCh()
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
	msgEnc := utils.NewMessageEncoder().
		RawUInt32(common.REQ_CURRENT_TIME_IN_MILLIS)

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
	const VERSION = 1
	msgEnc := utils.NewMessageEncoder().
		RawUInt32(common.REQ_MANAGED_ACCTS).
		Int(VERSION)

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
		Bars: make([]models.Bar, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_HISTORICAL_DATA,
		Response: resp,
	})

	// Build the message to send
	msgEnc := utils.NewMessageEncoder().Reserve(20).
		RawUInt32(common.REQ_HISTORICAL_DATA).
		RequestID(req.ID()).
		Marshal(opts.Contract).
		Bool(opts.Contract.IncludeExpired).
		String(opts.EndDate.Format("20060102-15:04:05")).
		String(opts.BarSize.String()).
		String(strconv.Itoa(opts.Duration) + " " + opts.DurationUnit.String()).
		Bool(opts.OnlyRegularTradingHours).
		String(opts.WhatToShow.String()).
		Int(2) // Return epoch timestamp

	if opts.Contract.SecType == models.SecurityTypePair {
		msgEnc.Int(len(opts.Contract.ComboLegs))
		for _, comboLeg := range opts.Contract.ComboLegs {
			msgEnc.Int64(comboLeg.ConID, false).
				Int64(comboLeg.Ratio, false).
				String(comboLeg.Action).
				String(comboLeg.Exchange)
		}
	}

	msgEnc.Bool(false). // KeepUpToDate
				Marshal(&models.TagValueList{})

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
	msgEnc := utils.NewMessageEncoder().Reserve(22).
		RawUInt32(common.REQ_HISTORICAL_TICKS).
		RequestID(req.ID()).
		Marshal(opts.Contract).
		Bool(opts.Contract.IncludeExpired).
		String(opts.StartDate.Format("20060102-15:04:05")).
		String(opts.EndDate.Format("20060102-15:04:05")).
		Int(opts.NumberOfTicks).
		String(opts.WhatToShow.String()).
		Bool(opts.OnlyRegularTradingHours).
		Bool(opts.IgnoreSize).
		Marshal(&models.TagValueList{})

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
		ContractDetails: make([]models.ContractDetails, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_CONTRACT_DATA,
		Response: resp,
	})

	// Build the message to send
	const VERSION = 8
	msgEnc := utils.NewMessageEncoder().Reserve(21).
		RawUInt32(common.REQ_CONTRACT_DATA).
		RequestID(req.ID()).
		Int(VERSION).
		Marshal(opts.Contract).
		Bool(opts.Contract.IncludeExpired).
		String(opts.Contract.SecIDType).
		String(opts.Contract.SecID).
		String(opts.Contract.IssuerID)

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
		ContractDescriptions: make([]models.ContractDescription, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_MATCHING_SYMBOLS,
		Response: resp,
	})

	// Build the message to send
	msgEnc := utils.NewMessageEncoder().Reserve(3).
		RawUInt32(common.REQ_MATCHING_SYMBOLS).
		RequestID(req.ID()).
		String(opts.Pattern)

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
	const VERSION = 1
	msgEnc := utils.NewMessageEncoder().Reserve(3).
		RawUInt32(common.REQ_MARKET_DATA_TYPE).
		Int(VERSION).
		Int(int(opts.Type))

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
	const VERSION = 11
	genericTickSB := strings.Builder{}
	for idx, gt := range opts.AdditionalGenericTicks {
		if idx > 0 {
			_, _ = genericTickSB.WriteRune(',')
		}
		_, _ = genericTickSB.WriteString(strconv.Itoa(int(gt)))
	}
	msgEnc := utils.NewMessageEncoder().Reserve(3).
		RawUInt32(common.REQ_MKT_DATA).
		Int(VERSION).
		RequestID(req.ID()).
		Marshal(opts.Contract)
	if opts.Contract.SecType == models.SecurityTypePair {
		comboLegsCount := len(opts.Contract.ComboLegs)
		msgEnc.Int(comboLegsCount)
		for _, comboLeg := range opts.Contract.ComboLegs {
			msgEnc.Int64(comboLeg.ConID, false).
				Int64(comboLeg.Ratio, false).
				String(comboLeg.Action).
				String(comboLeg.Exchange)
		}
	}
	if opts.Contract.DeltaNeutralContract != nil {
		msgEnc.Bool(true).
			Int64(opts.Contract.DeltaNeutralContract.ConID, false).
			Float64(opts.Contract.DeltaNeutralContract.Delta, false).
			Float64(opts.Contract.DeltaNeutralContract.Price, false)
	} else {
		msgEnc.Bool(false)
	}
	msgEnc.String(genericTickSB.String())
	msgEnc.Bool(opts.Snapshot)
	msgEnc.Bool(opts.RegulatorySnapshot).
		Marshal(&models.TagValueList{})

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
	const VERSION = 5
	msgEnc := utils.NewMessageEncoder().Reserve(17).
		RawUInt32(common.REQ_MKT_DEPTH).
		Int(VERSION).
		RequestID(req.ID()).
		Marshal(opts.Contract).
		Int(opts.RowsCount).
		Bool(opts.SmartDepth).
		Marshal(&models.TagValueList{})

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
	const VERSION = 2
	msgEnc := utils.NewMessageEncoder().Reserve(3).
		RawUInt32(common.CANCEL_MKT_DATA).
		Int(VERSION).
		RequestID(req.ID())

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
	const VERSION = 1
	msgEnc := utils.NewMessageEncoder().Reserve(4).
		RawUInt32(common.CANCEL_MKT_DEPTH).
		Int(VERSION).
		RequestID(req.ID()).
		Bool(isSmartDepth)

	// Send it
	_ = c.sendMessage(msgEnc.Bytes())

	// Remove the request from the manager
	c.reqMgr.removeRequest(req, nil)
}
