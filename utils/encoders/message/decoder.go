package message

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mxmauro/ibkr/common"
)

// -----------------------------------------------------------------------------

type Decoder struct {
	fields       [][]byte
	nextFieldIdx int
	err          error
}

// -----------------------------------------------------------------------------

var errNoMoreFields = errors.New("no more fields")

// -----------------------------------------------------------------------------

func NewDecoder(msg []byte) *Decoder {
	d := Decoder{
		fields: bytes.Split(msg, []byte{common.MessageDelimiter}),
	}
	d.fields = d.fields[:len(d.fields)-1]

	// Done
	return &d
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) SetErr(err error) {
	if d.err == nil && err != nil {
		d.err = err
	}
}

func (d *Decoder) Skip() {
	_ = d.getNextField()
}

func (d *Decoder) RequestID(canBeOptional bool) int32 {
	reqID := d.Int32()
	if d.Err() != nil {
		return 0
	}
	if canBeOptional && reqID <= 0 {
		return -1
	}
	if reqID < 1 {
		d.SetErr(fmt.Errorf("received invalid request id: %d", reqID))
		return 0
	}
	return reqID
}

func (d *Decoder) Int32() int32 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return 0
	}

	// Parse value
	value, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		d.SetErr(errors.New("cannot decode int32 value"))
		return 0
	}

	// Done
	return int32(value)
}

func (d *Decoder) Int32Max() *int32 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return nil
	}

	// Parse value
	value, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		d.SetErr(errors.New("cannot decode int32 value"))
		return nil
	}

	// Done
	v := int32(value)
	return &v
}

func (d *Decoder) Int64() int64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return 0
	}

	// Parse value
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode int64 value"))
		return 0
	}

	// Done
	return value
}

func (d *Decoder) Int64Max() *int64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return nil
	}

	// Parse value
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode int64 value"))
		return nil
	}

	// Done
	return &value
}

func (d *Decoder) Float() float64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return 0.0
	}

	// Parse value
	value, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode float64 value"))
		return 0.0
	}

	// Done
	return value
}

func (d *Decoder) FloatMax() *float64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		return nil
	}

	// Parse value
	value, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode float64 value"))
		return nil
	}

	// Done
	return &value
}

func (d *Decoder) Bool() bool {
	data := d.getNextField()

	// Is nil, empty or false?
	if len(data) == 0 || bytes.Equal(data, []byte{'0'}) {
		return false
	}

	// Done
	return true
}

func (d *Decoder) String() string {
	data := d.getNextField()
	if len(data) == 0 {
		return ""
	}

	// Done
	return string(data)
}

func (d *Decoder) EscapedString() string {
	data := d.getNextField()
	if len(data) == 0 {
		return ""
	}

	// Decode unescaped?
	s, err := strconv.Unquote(`"` + string(data) + `"`)
	if err != nil {
		d.SetErr(errors.New("cannot decode string"))
		return ""
	}

	// Done
	return s
}

func (d *Decoder) EpochTimestamp(inMilliseconds bool) time.Time {
	ts := d.Int64()
	if ts < 0 {
		d.SetErr(errors.New("invalid epoch timestamp"))
		return time.Time{}
	}
	if inMilliseconds {
		return time.Unix(ts/1000, (ts%1000)*1000000).UTC()
	}
	return time.Unix(ts, 0).UTC()
}

func (d *Decoder) getNextField() []byte {
	if d.hasError() {
		return nil
	}
	if d.nextFieldIdx >= len(d.fields) {
		d.SetErr(errNoMoreFields)
		return nil
	}
	d.nextFieldIdx += 1
	return d.fields[d.nextFieldIdx-1]
}

func (d *Decoder) hasError() bool {
	return d.err != nil
}
