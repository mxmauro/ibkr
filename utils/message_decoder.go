package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mxmauro/ibkr/common"
)

// -----------------------------------------------------------------------------

type MessageDecoder struct {
	fields       [][]byte
	nextFieldIdx int
	err          error
}

// -----------------------------------------------------------------------------

const (
	delim = '\x00'
)

var errNoMoreFields = errors.New("no more fields")

// -----------------------------------------------------------------------------

func NewMessageDecoder(msg []byte) *MessageDecoder {
	d := MessageDecoder{
		fields: bytes.Split(msg, []byte{delim}),
	}
	d.fields = d.fields[:len(d.fields)-1]

	// Done
	return &d
}

func (d *MessageDecoder) Err() error {
	return d.err
}

func (d *MessageDecoder) SetErr(err error) {
	if d.err == nil && err != nil {
		d.err = err
	}
}

func (d *MessageDecoder) Skip() {
	_ = d.getNextField()
}

func (d *MessageDecoder) Int64(emptyIsUnset bool) int64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		if emptyIsUnset {
			return common.UNSET_INT
		}
		return 0
	}

	// Parse value
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode int64 value"))
		if emptyIsUnset {
			return common.UNSET_INT
		}
		return 0
	}

	// Done
	return value
}

func (d *MessageDecoder) Float64(emptyIsUnset bool) float64 {
	data := d.getNextField()

	// Handle nil or empty value
	if len(data) == 0 {
		if emptyIsUnset {
			return common.UNSET_FLOAT
		}
		return 0.0
	}

	// Parse value
	value, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode float64 value"))
		if emptyIsUnset {
			return common.UNSET_FLOAT
		}
		return 0.0
	}

	// Done
	return value
}

func (d *MessageDecoder) Bool() bool {
	data := d.getNextField()

	// Is nil, empty or false?
	if len(data) == 0 || bytes.Equal(data, []byte{'0'}) {
		return false
	}

	// Done
	return true
}

func (d *MessageDecoder) String(unescape bool) string {
	data := d.getNextField()
	if data == nil {
		return ""
	}

	// Decode unescaped?
	if unescape {
		s, err := strconv.Unquote(fmt.Sprint("\"", data, "\""))
		if err != nil {
			d.SetErr(errors.New("cannot decode string"))
			return ""
		}
		return s
	}

	// Done
	return string(data)
}

func (d *MessageDecoder) EpochTimestamp(inMilliseconds bool) time.Time {
	ts := d.Int64(false)
	if ts < 0 {
		d.SetErr(errors.New("invalid epoch timestamp"))
		return time.Time{}
	}
	if inMilliseconds {
		return time.Unix(ts/1000, (ts%1000)*1000000).UTC()
	}
	return time.Unix(ts, 0).UTC()
}

func (d *MessageDecoder) getNextField() []byte {
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

func (d *MessageDecoder) hasError() bool {
	return d.err != nil
}
