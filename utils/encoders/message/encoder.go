package message

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/mxmauro/ibkr/common"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------

// Encoder builds legacy delimited or protobuf-framed outbound messages.
type Encoder struct {
	buf bytes.Buffer
	raw bool
	err error
}

// EncoderMarshaller is implemented by values that can encode themselves into the legacy wire format.
type EncoderMarshaller interface {
	EncodeMessage(mode int) ([]byte, error)
}

// -----------------------------------------------------------------------------

// NewEncoder allocates an encoder for framed outbound messages.
func NewEncoder() *Encoder {
	e := Encoder{
		buf: bytes.Buffer{},
	}

	// Reserve 4 bytes for the message size header
	_, e.err = e.buf.Write([]byte{0, 0, 0, 0})

	// Done
	return &e
}

// NewRawEncoder allocates an encoder that omits the 4-byte message length header.
func NewRawEncoder() *Encoder {
	e := Encoder{
		buf: bytes.Buffer{},
		raw: true,
	}

	// Done
	return &e
}

// Reserve grows the backing buffer using an estimated field count.
func (e *Encoder) Reserve(estimatedFieldsCount int) *Encoder {
	e.buf.Grow(4 + 8*estimatedFieldsCount)
	// Done
	return e
}

// RequestID encodes an IB request identifier as a delimited string field.
func (e *Encoder) RequestID(v int32) *Encoder {
	return e.String(strconv.Itoa(int(v)))
}

// Int encodes an int as a delimited string field.
func (e *Encoder) Int(v int) *Encoder {
	return e.String(strconv.FormatInt(int64(v), 10))
}

// Int32 encodes an int32 as a delimited string field.
func (e *Encoder) Int32(v int32) *Encoder {
	return e.String(strconv.FormatInt(int64(v), 10))
}

// Int32Max encodes an optional int32 as a delimited string field.
func (e *Encoder) Int32Max(v *int32) *Encoder {
	if v == nil {
		if e.err == nil {
			e.err = e.buf.WriteByte(common.MessageDelimiter)
		}
		return e
	}
	return e.String(strconv.FormatInt(int64(*v), 10))
}

// Int64 encodes an int64 as a delimited string field.
func (e *Encoder) Int64(v int64) *Encoder {
	return e.String(strconv.FormatInt(v, 10))
}

// Int64Max encodes an optional int64 as a delimited string field.
func (e *Encoder) Int64Max(v *int64) *Encoder {
	if v == nil {
		if e.err == nil {
			e.err = e.buf.WriteByte(common.MessageDelimiter)
		}
		return e
	}
	return e.String(strconv.FormatInt(*v, 10))
}

// Float encodes a float64 as a delimited string field.
func (e *Encoder) Float(v float64) *Encoder {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return e.String("")
	}
	return e.String(strconv.FormatFloat(v, 'f', -1, 64))
}

// FloatMax encodes an optional float64 as a delimited string field.
func (e *Encoder) FloatMax(v *float64) *Encoder {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return e.String("")
	}
	return e.String(strconv.FormatFloat(*v, 'f', -1, 64))
}

// String encodes a string field and appends the protocol delimiter.
func (e *Encoder) String(v string) *Encoder {
	if e.err != nil {
		return e
	}
	_, e.err = e.buf.WriteString(v)
	if e.err == nil {
		e.err = e.buf.WriteByte(common.MessageDelimiter)
	}
	return e
}

// Bool encodes a bool as 1 or 0 followed by the protocol delimiter.
func (e *Encoder) Bool(v bool) *Encoder {
	if e.err != nil {
		return e
	}
	if v {
		e.err = e.buf.WriteByte('1')
	} else {
		e.err = e.buf.WriteByte('0')
	}
	if e.err == nil {
		_ = e.buf.WriteByte(common.MessageDelimiter)
	}
	return e
}

// AddDelim appends an empty delimited field.
func (e *Encoder) AddDelim() *Encoder {
	if e.err != nil {
		return e
	}
	e.err = e.buf.WriteByte(common.MessageDelimiter)
	return e
}

// Marshal appends a legacy wire representation produced by v.
func (e *Encoder) Marshal(v EncoderMarshaller, mode int) *Encoder {
	var encoded []byte

	if e.err != nil {
		return e
	}
	encoded, e.err = v.EncodeMessage(mode)
	if e.err == nil {
		_, e.err = e.buf.Write(encoded)
	}
	return e
}

// Proto appends a marshaled protobuf payload to the current message.
func (e *Encoder) Proto(v proto.Message) *Encoder {
	var encoded []byte

	if e.err != nil {
		return e
	}
	encoded, e.err = proto.Marshal(v)
	if e.err == nil {
		_, e.err = e.buf.Write(encoded)
	}
	return e
}

// Raw appends raw bytes to the current message buffer.
func (e *Encoder) Raw(v []byte) *Encoder {
	if e.err != nil {
		return e
	}
	_, e.err = e.buf.Write(v)
	return e
}

// RawUInt32 appends an unsigned 32-bit integer using big-endian byte order.
func (e *Encoder) RawUInt32(v uint32) *Encoder {
	var data [4]byte

	if e.err != nil {
		return e
	}
	binary.BigEndian.PutUint32(data[:], v)
	_, e.err = e.buf.Write(data[:])
	return e
}

// RawInt64 appends a signed 64-bit integer using big-endian byte order.
func (e *Encoder) RawInt64(v int64) *Encoder {
	var data [8]byte

	if e.err != nil {
		return e
	}
	binary.BigEndian.PutUint64(data[:], uint64(v))
	_, e.err = e.buf.Write(data[:])
	return e
}

// Err returns the first encoding error encountered by the encoder.
func (e *Encoder) Err() error {
	return e.err
}

// Bytes finalize the message by writing the size header and returning the complete message
func (e *Encoder) Bytes() []byte {
	if e.err != nil {
		return nil
	}

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
