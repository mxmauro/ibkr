package connection

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/mxmauro/go-rundownprotection"
	"github.com/mxmauro/resetevent"
	"github.com/mxmauro/ringbuffer"
)

// -----------------------------------------------------------------------------

func TestCloseUnblocksReadMessageSend(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()

	c := &Connection{
		conn:          clientConn,
		msgCh:         make(chan []byte),
		errorCh:       make(chan error, 1),
		writeBuffer:   ringbuffer.New(512),
		writeBufferEv: resetevent.NewAutoResetEvent(),
	}
	c.rp = rundownprotection.RundownProtection{}
	c.rp.Initialize()

	c.wg.Add(1)
	go c.handleRead()

	var msgLen [4]byte
	binary.BigEndian.PutUint32(msgLen[:], 1)
	if _, err := serverConn.Write(msgLen[:]); err != nil {
		t.Fatal(err)
	}
	if _, err := serverConn.Write([]byte{0}); err != nil {
		t.Fatal(err)
	}

	doneCh := make(chan struct{})
	go func() {
		c.Close()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-time.After(time.Second):
		t.Fatal("connection close blocked while read goroutine was sending a message")
	}
}
