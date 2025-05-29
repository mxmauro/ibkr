package utils

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/mxmauro/ibkr/common"
)

// -----------------------------------------------------------------------------

type MessageEncoder struct {
	buf bytes.Buffer
	raw bool
}

type MessageEncoderMarshaller interface {
	EncodeMessage() []byte
}

// -----------------------------------------------------------------------------

func Delimiter() byte {
	return delim
}

func NewMessageEncoder() *MessageEncoder {
	e := MessageEncoder{
		buf: bytes.Buffer{},
	}

	// Reserve 4 bytes for the message size header
	_, _ = e.buf.Write([]byte{0, 0, 0, 0})

	// Done
	return &e
}

func NewRawMessageEncoder() *MessageEncoder {
	e := MessageEncoder{
		buf: bytes.Buffer{},
		raw: true,
	}

	// Done
	return &e
}

func (e *MessageEncoder) Reserve(estimatedFieldsCount int) *MessageEncoder {
	e.buf.Grow(4 + 8*estimatedFieldsCount)
	// Done
	return e
}

func (e *MessageEncoder) RequestID(v int32) *MessageEncoder {
	return e.String(strconv.Itoa(int(v)))
}

func (e *MessageEncoder) Int(v int) *MessageEncoder {
	return e.String(strconv.Itoa(v))
}

func (e *MessageEncoder) Int64(v int64, emptyIfUnset bool) *MessageEncoder {
	if emptyIfUnset && v == common.UNSET_INT {
		_ = e.buf.WriteByte(delim)
		return e
	}
	return e.String(strconv.FormatInt(v, 10))
}

func (e *MessageEncoder) Float64(v float64, emptyIfUnset bool) *MessageEncoder {
	if emptyIfUnset && v == common.UNSET_FLOAT {
		_ = e.buf.WriteByte(delim)
		return e
	}
	return e.String(strconv.FormatFloat(v, 'f', -1, 64))
}

func (e *MessageEncoder) String(v string) *MessageEncoder {
	_, _ = e.buf.WriteString(v)
	_ = e.buf.WriteByte(delim)
	return e
}

func (e *MessageEncoder) Bool(v bool) *MessageEncoder {
	if v {
		e.buf.WriteByte('1')
	} else {
		e.buf.WriteByte('0')
	}
	_ = e.buf.WriteByte(delim)
	return e
}

func (e *MessageEncoder) AddDelim() *MessageEncoder {
	_ = e.buf.WriteByte(delim)
	return e
}

func (e *MessageEncoder) Marshal(v MessageEncoderMarshaller) *MessageEncoder {
	_, _ = e.buf.Write(v.EncodeMessage())
	return e
}

func (e *MessageEncoder) Raw(v []byte) *MessageEncoder {
	_, _ = e.buf.Write(v)
	return e
}

func (e *MessageEncoder) RawUInt32(v uint32) *MessageEncoder {
	var data [4]byte

	binary.BigEndian.PutUint32(data[:], v)
	_, _ = e.buf.Write(data[:])
	return e
}

func (e *MessageEncoder) RawInt64(v int64) *MessageEncoder {
	var data [8]byte

	binary.BigEndian.PutUint64(data[:], uint64(v))
	_, _ = e.buf.Write(data[:])
	return e
}

// Bytes finalizes the message by writing the size header and returning the complete message
func (e *MessageEncoder) Bytes() []byte {
	data := e.buf.Bytes()
	if !e.raw {
		// Calculate message size (excluding the 4-byte header)
		msgSize := len(data) - 4

		// Write the size back into the header
		binary.BigEndian.PutUint32(data[:4], uint32(msgSize))
	}

	// Done
	return data
}
