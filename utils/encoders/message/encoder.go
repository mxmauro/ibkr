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

type Encoder struct {
	buf bytes.Buffer
	raw bool
	err error
}

type EncoderMarshaller interface {
	EncodeMessage(mode int) ([]byte, error)
}

// -----------------------------------------------------------------------------

func NewEncoder() *Encoder {
	e := Encoder{
		buf: bytes.Buffer{},
	}

	// Reserve 4 bytes for the message size header
	_, e.err = e.buf.Write([]byte{0, 0, 0, 0})

	// Done
	return &e
}

func NewRawEncoder() *Encoder {
	e := Encoder{
		buf: bytes.Buffer{},
		raw: true,
	}

	// Done
	return &e
}

func (e *Encoder) Reserve(estimatedFieldsCount int) *Encoder {
	e.buf.Grow(4 + 8*estimatedFieldsCount)
	// Done
	return e
}

func (e *Encoder) RequestID(v int32) *Encoder {
	return e.String(strconv.Itoa(int(v)))
}

func (e *Encoder) Int(v int) *Encoder {
	return e.String(strconv.FormatInt(int64(v), 10))
}

func (e *Encoder) Int32(v int32) *Encoder {
	return e.String(strconv.FormatInt(int64(v), 10))
}

func (e *Encoder) Int32Max(v *int32) *Encoder {
	if v == nil {
		if e.err == nil {
			e.err = e.buf.WriteByte(common.MessageDelimiter)
		}
		return e
	}
	return e.String(strconv.FormatInt(int64(*v), 10))
}

func (e *Encoder) Int64(v int64) *Encoder {
	return e.String(strconv.FormatInt(v, 10))
}

func (e *Encoder) Int64Max(v *int64) *Encoder {
	if v == nil {
		if e.err == nil {
			e.err = e.buf.WriteByte(common.MessageDelimiter)
		}
		return e
	}
	return e.String(strconv.FormatInt(*v, 10))
}

func (e *Encoder) Float(v float64) *Encoder {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return e.String("")
	}
	return e.String(strconv.FormatFloat(v, 'f', -1, 64))
}

func (e *Encoder) FloatMax(v *float64) *Encoder {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return e.String("")
	}
	return e.String(strconv.FormatFloat(*v, 'f', -1, 64))
}

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

func (e *Encoder) AddDelim() *Encoder {
	if e.err != nil {
		return e
	}
	e.err = e.buf.WriteByte(common.MessageDelimiter)
	return e
}

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

func (e *Encoder) Raw(v []byte) *Encoder {
	if e.err != nil {
		return e
	}
	_, e.err = e.buf.Write(v)
	return e
}

func (e *Encoder) RawUInt32(v uint32) *Encoder {
	var data [4]byte

	if e.err != nil {
		return e
	}
	binary.BigEndian.PutUint32(data[:], v)
	_, e.err = e.buf.Write(data[:])
	return e
}

func (e *Encoder) RawInt64(v int64) *Encoder {
	var data [8]byte

	if e.err != nil {
		return e
	}
	binary.BigEndian.PutUint64(data[:], uint64(v))
	_, e.err = e.buf.Write(data[:])
	return e
}

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
