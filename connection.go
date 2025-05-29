package ibkr

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	"sync/atomic"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/connection"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

func (c *Client) connectToServer(ctx context.Context, opts Options) error {
	var conn *connection.Connection
	var err error

	serverAddress := opts.Address

	for redirectionsCount := 0; ; redirectionsCount++ {
		var redirectedHost string

		conn, err = connection.New(ctx, serverAddress)
		if err != nil {
			return err
		}

		redirectedHost, err = c.initialHandshake(ctx, conn, opts)
		if err != nil {
			conn.Close()
			return err
		}
		if len(redirectedHost) == 0 {
			break
		}

		conn.Close()
		if redirectionsCount >= 10 {
			conn.Close()
			return fmt.Errorf("too many redirects")
		}

		serverAddress = redirectedHost
	}

	// Set up client id
	if opts.ClientID > 0 {
		c.clientID = opts.ClientID
	} else {
		c.clientID = 1 + int32(rand.Uint32()&0x7FFFFFFE)
	}

	// Start the API service
	err = c.startApi(conn, opts)
	if err != nil {
		conn.Close()
		return err
	}

	// Skip initial incoming messages until we receive the next valid request ID
	err = c.waitUntilNextReqID(ctx, conn)
	if err != nil {
		return err
	}

	// On success, save the connection link
	c.conn = conn

	// Done
	return nil
}

func (c *Client) initialHandshake(ctx context.Context, conn *connection.Connection, opts Options) (string, error) {
	// Initiate initialHandshake
	msg := fmt.Appendf(nil, "v%d..%d", common.MinClientVersion, common.MaxClientVersion)
	if len(opts.ConnectOptions) > 0 {
		msg = fmt.Appendf(msg, " %s", opts.ConnectOptions)
	}

	var tempBuf [8]byte
	tempBuf[0] = 'A'
	tempBuf[1] = 'P'
	tempBuf[2] = 'I'
	tempBuf[3] = 0
	binary.BigEndian.PutUint32(tempBuf[4:], uint32(len(msg)))
	err := conn.Send(tempBuf[:])
	if err != nil {
		return "", err
	}
	err = conn.Send(msg)
	if err != nil {
		return "", err
	}

	// Wait for server response
	msg, err = conn.WaitForNextMessage(ctx)
	if err != nil {
		return "", err
	}

	// Decode message
	msgDec := utils.NewMessageDecoder(msg)

	serverVersion := msgDec.Int64(false)
	connTimeOrNewServerHost := msgDec.String(false)
	if msgDec.Err() != nil {
		return "", err
	}

	if serverVersion < 0 {
		// redirect to the new server
		return connTimeOrNewServerHost, nil
	}
	if serverVersion < common.MinServerVersion {
		return "", fmt.Errorf("server version %d is not supported", serverVersion)
	}

	// Store server version
	c.serverVersion = int32(serverVersion)

	// Done
	return "", nil
}

func (c *Client) startApi(conn *connection.Connection, opts Options) error {
	const VERSION = 2

	msgEnc := utils.NewMessageEncoder().
		RawUInt32(uint32(common.START_API)).
		Int(VERSION).
		Int(int(c.clientID)).
		String(opts.OptionalCapabilities)

	return conn.Send(msgEnc.Bytes())
}

// This function ignores all the initial incoming messages until we receive the next valid request ID.
func (c *Client) waitUntilNextReqID(ctx context.Context, conn *connection.Connection) error {
	var msgID int64

	for {
		msg, err := conn.WaitForNextMessage(ctx)
		if err != nil {
			return err
		}

		msgID, _, err = c.getIncomingMessageID(msg)
		if err != nil {
			return err
		}
		if msgID == common.NEXT_VALID_ID {
			msgDec := utils.NewMessageDecoder(msg[4:])
			reqID := msgDec.Int64(false)
			if msgDec.Err() != nil {
				return msgDec.Err()
			}
			if reqID < 1 || reqID > math.MaxInt32 {
				reqID = 1
			}

			// Done
			atomic.StoreInt32(&c.nextValidReqID, int32(reqID))
			return nil
		}
	}
}

func (c *Client) connectionWorker() {
	var err error

	defer c.wg.Done()

	// Process messages and errors
	for {
		var msg []byte

		msg, err = c.conn.WaitForNextMessage(&c.rp)
		if err != nil {
			break
		}

		if !c.rp.Acquire() {
			break
		}
		err = c.processIncomingMessage(msg)
		c.rp.Release()
		if err != nil {
			break
		}
	}

	// If the connection dropped, mark as closed
	if err == nil || utils.IsConnectionDropError(err) {
		err = net.ErrClosed
	}

	// Close the connection
	c.connErrHolder.Store(err)
	c.isDisconnectedEv.Set()
	c.connMtx.Lock()
	conn := c.conn
	c.conn = nil
	c.connMtx.Unlock()
	conn.Close()

	// Cancel pending requests
	c.reqMgr.removeAndTryCancelAllRequests(err)

	// Raise the event if not closing and an event handler is present
	if c.rp.Acquire() {
		if c.eventsHandler != nil {
			c.eventsHandler.ConnectionClosed(err)
		}
		c.rp.Release()
	}
}

func (c *Client) sendMessage(msg []byte) error {
	c.connMtx.Lock()
	defer c.connMtx.Unlock()

	// Do we have a connection?
	if c.conn == nil {
		return net.ErrClosed
	}

	// Send the message
	return c.conn.Send(msg)
}

func (c *Client) sendRequest(msg []byte, req *Request) error {
	// Do we have a connection?
	c.connMtx.Lock()
	defer c.connMtx.Unlock()

	if c.conn == nil {
		return c.getConnError()
	}

	// Add to the active requests map
	c.reqMgr.addRequest(req)

	// Send the message
	err := c.conn.Send(msg)
	if err != nil {
		// On error remove the request from the active requests map
		c.reqMgr.removeRequest(req, err)
	}

	// Done
	return err
}

func (c *Client) getConnError() error {
	v := c.connErrHolder.Load()
	if v == nil {
		return nil
	}
	return v.(error)
}
