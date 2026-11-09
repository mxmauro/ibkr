package ibkr

import (
	"context"
	"net"
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
	"github.com/mxmauro/resetevent"
)

// -----------------------------------------------------------------------------

// Client manages the connection, request routing, and event processing for IBKR.
type Client struct {
	rp            rundownprotection.RundownProtection
	eventsHandler Events
	orderStore    OrderStore
	tzOffset      time.Duration

	wg          sync.WaitGroup
	destroyOnce sync.Once

	connMtx          sync.Mutex
	conn             *connection.Connection
	connErrHolder    atomic.Value
	isDisconnectedEv *resetevent.ManualResetEvent
	clientID         int32
	serverVersion    int32
	marketDataType   int32

	nextValidReqID        int32
	nextValidReqWithoutID int32

	reqMgr RequestManager
}

// Options configures how a client connects to TWS or IB Gateway.
type Options struct {
	Address       string
	EventsHandler Events
	OrderStore    OrderStore
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

	if len(opts.Address) == 0 {
		return nil, errInvalidHostPortAddress
	}
	if len(opts.ConnectOptions) > 0 && !utils.IsPrintableAsciiString(opts.ConnectOptions) {
		return nil, errInvalidOptionalCapabilities
	}
	if len(opts.OptionalCapabilities) > 0 && !utils.IsPrintableAsciiString(opts.OptionalCapabilities) {
		return nil, errInvalidOptionalCapabilities
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
		default:
			return nil, errInvalidTzOffset
		}
	}
	if opts.ClientID < 0 {
		return nil, errInvalidClientID
	}

	c := Client{
		eventsHandler: opts.EventsHandler,
		orderStore:    opts.OrderStore,
		tzOffset:      tzOffset,

		wg:          sync.WaitGroup{},
		destroyOnce: sync.Once{},

		connMtx:          sync.Mutex{},
		isDisconnectedEv: resetevent.NewManualResetEvent(),
	}
	if c.orderStore == nil {
		c.orderStore = noopOrderStore{}
	}
	c.rp.Initialize()
	c.initRequestManager()
	atomic.StoreInt32(&c.nextValidReqWithoutID, 1)
	atomic.StoreInt32(&c.marketDataType, int32(models.MarketDataTypeDelayed))

	err = c.connectToServer(ctx, opts)
	if err != nil {
		return nil, err
	}

	c.wg.Add(1)
	go c.connectionWorker()

	err = c.RequestMarketDataType(ctx, models.MarketDataTypeRequestOptions{
		Type: models.MarketDataTypeRealtime,
	})
	if err != nil {
		c.connMtx.Lock()
		if c.conn != nil {
			c.conn.Close()
		}
		c.connMtx.Unlock()
		c.wg.Wait()
		return nil, err
	}

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
	if !c.rp.Acquire() {
		return time.Time{}, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.CurrentTimeResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  common.REQ_CURRENT_TIME_IN_MILLIS,
		Response: resp,
	})

	var msgEnc *message.Encoder
	if c.isProtoBufAvailable(common.REQ_CURRENT_TIME_IN_MILLIS) {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_CURRENT_TIME_IN_MILLIS + common.PROTOBUF_MSG_ID).
			Proto(&protobuf.CurrentTimeInMillisRequest{})
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_CURRENT_TIME_IN_MILLIS)
	}
	if msgEnc.Err() != nil {
		return time.Time{}, msgEnc.Err()
	}

	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return time.Time{}, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return time.Time{}, err
	}

	return resp.CurrentTime, nil
}

// RequestManagedAccounts requests the list of managed accounts.
func (c *Client) RequestManagedAccounts(ctx context.Context) ([]string, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.ManagedAccountsResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  common.REQ_MANAGED_ACCTS,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_MANAGED_ACCTS) {
		pb := protobuf.ManagedAccountsRequest{}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MANAGED_ACCTS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 1
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MANAGED_ACCTS).
			Int(version)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	err := c.sendRequest(msgEnc.Bytes(), req)
	if err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	err = c.waitRequestCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Accounts, nil
}

// -----------------------------------------------------------------------------

func (c *Client) isProtoBufAvailable(msgType uint32) bool {
	minServerVer, ok := common.PROTOBUF_MSG_IDS[msgType]
	return ok && c.serverVersion >= minServerVer
}

func (c *Client) currentMarketDataType() models.MarketDataType {
	return models.MarketDataType(atomic.LoadInt32(&c.marketDataType))
}

func tagValueListToMap(values models.TagValueList) map[string]string {
	if len(values) == 0 {
		return nil
	}
	m := make(map[string]string, len(values))
	for _, tv := range values {
		m[tv.Tag] = tv.Value
	}
	return m
}
