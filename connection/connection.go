package connection

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mxmauro/go-rundownprotection"
	"github.com/mxmauro/resetevent"
	"github.com/mxmauro/ringbuffer"
)

// -----------------------------------------------------------------------------

// Connection manages a framed TCP connection to TWS or IB Gateway.
type Connection struct {
	mtx sync.Mutex
	wg  sync.WaitGroup

	rp rundownprotection.RundownProtection

	conn net.Conn

	msgCh     chan []byte
	errorSent int32
	errorCh   chan error

	writeBuffer   *ringbuffer.RingBuffer
	writeBufferEv *resetevent.AutoResetEvent

	closed uint32

	closeOnce sync.Once
}

// -----------------------------------------------------------------------------

// New opens a TCP connection and starts the background read and write workers.
func New(ctx context.Context, address string) (*Connection, error) {
	var err error

	// Create the connection object
	c := Connection{
		mtx: sync.Mutex{},
		wg:  sync.WaitGroup{},

		msgCh:   make(chan []byte),
		errorCh: make(chan error, 1),

		writeBuffer:   ringbuffer.New(512),
		writeBufferEv: resetevent.NewAutoResetEvent(),
	}
	c.rp.Initialize()

	// Set up a limited context
	dialCtx, dialCancelCtx := context.WithTimeout(ctx, 10*time.Second)
	defer dialCancelCtx()

	// Dial
	dialer := net.Dialer{}
	c.conn, err = dialer.DialContext(dialCtx, "tcp", address)
	if err != nil {
		return nil, err
	}

	// Handle read and write
	c.wg.Add(2)
	go c.handleRead()
	go c.handleWrite()

	// Done
	return &c, nil
}

// Close shuts down the connection and waits for the background workers to stop.
func (c *Connection) Close() {
	c.closeOnce.Do(func() {
		c.rp.Wait()

		_ = c.conn.Close()
		c.wg.Wait()

		close(c.msgCh)
		close(c.errorCh)
	})
}

// Send queues a fully encoded message for transmission.
func (c *Connection) Send(msg []byte) error {
	// Check if the connection is closed
	if c.rp.Acquire() == false {
		return net.ErrClosed
	}
	defer c.rp.Release()

	_, _ = c.writeBuffer.Write(msg)
	c.writeBufferEv.Set()

	// Done
	return nil
}

// WaitForNextMessage waits for the next inbound message or connection error.
func (c *Connection) WaitForNextMessage(ctx context.Context) ([]byte, error) {
	select {
	case msg := <-c.msgCh:
		return msg, nil

	case err := <-c.errorCh:
		return nil, err

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// MsgRcvCh returns the inbound message channel produced by the reader worker.
func (c *Connection) MsgRcvCh() <-chan []byte {
	return c.msgCh
}

// ErrorCh returns the asynchronous connection error channel.
func (c *Connection) ErrorCh() <-chan error {
	return c.errorCh
}

func (c *Connection) handleRead() {
	var msgLenBytes [4]byte
	var ofs int

	defer c.wg.Done()

	for {
		// Read the message length
		for ofs = 0; ofs < 4; {
			read, err := c.conn.Read(msgLenBytes[ofs:])
			if err != nil {
				c.handleError(err)
				return
			}
			ofs += read
		}

		msgLen := int(uint(binary.BigEndian.Uint32(msgLenBytes[:])))
		if msgLen == 0 {
			c.handleError(errors.New("invalid message length"))
			return
		}

		// Read the message
		msg := make([]byte, msgLen)
		for ofs = 0; ofs < msgLen; {
			read, err := c.conn.Read(msg[ofs:])
			if err != nil {
				c.handleError(err)
				return
			}
			ofs += read
		}

		// Forward the complete frame unless shutdown won the race first.
		if c.rp.Acquire() == false {
			return
		}
		select {
		case c.msgCh <- msg:
		case <-c.rp.Done():
			c.rp.Release()
			return
		}
		c.rp.Release()
	}
}

func (c *Connection) handleWrite() {
	defer c.wg.Done()

	toWriteBuf := make([]byte, 1024)

	for {
		select {
		case <-c.rp.Done():
			return

		case <-c.writeBufferEv.WaitCh():
			for {
				n, err := c.writeBuffer.Read(toWriteBuf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					c.handleError(err)
					return
				}

				if n == 0 {
					break
				}

				// Refresh the write deadline for each chunk so a stuck socket does not block shutdown forever.
				err = c.conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
				if err == nil {
					var written int

					for ofs := 0; ofs < n; {
						written, err = c.conn.Write(toWriteBuf[ofs:n])
						if err != nil {
							break
						}

						// c.debugDump(toWriteBuf[ofs:ofs+written], "Sent data:")

						ofs += written
					}
				}
				if err != nil {
					c.handleError(err)
					return
				}
			}
		}
	}
}

func (c *Connection) handleError(err error) {
	// Cancellation errors are not important
	if errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
		return
	}

	if c.rp.Acquire() == false {
		return
	}
	defer c.rp.Release()

	// We only send one error
	if atomic.CompareAndSwapInt32(&c.errorSent, 0, 1) {
		c.errorCh <- err
	}
}

func (c *Connection) parseFields(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, io.EOF
	}

	if len(data) < 4 {
		return 0, nil, nil // will try to read more data
	}

	totalSize := int(binary.BigEndian.Uint32(data[:4])) + 4

	if totalSize > len(data) {
		return 0, nil, nil
	}

	// msgBytes := make([]byte, totalSize-4, totalSize-4)
	// copy(msgBytes, data[4:totalSize])
	// not copy here, copied by callee more reasonable
	return totalSize, data[4:totalSize], nil
}

// debugDump prints a hexadecimal dump of a raw protocol frame.
func (c *Connection) debugDump(data []byte, header string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	fmt.Println(header)
	fmt.Println(hex.Dump(data))
}
