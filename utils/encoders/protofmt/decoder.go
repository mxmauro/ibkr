package protofmt

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------

type Decoder struct {
	msg []byte
	err error
}

var protoUnmarshaller = proto.UnmarshalOptions{
	DiscardUnknown: true,
	RecursionLimit: protowire.DefaultRecursionLimit,
	NoLazyDecoding: true,
}

// -----------------------------------------------------------------------------

func NewDecoder(msg []byte) *Decoder {
	d := Decoder{
		msg: msg,
	}

	// Done
	return &d
}

func (d *Decoder) Unmarshal(m proto.Message) {
	err := protoUnmarshaller.Unmarshal(d.msg, m)
	if err != nil {
		d.SetErr(err)
	}
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) SetErr(err error) {
	if d.err == nil && err != nil {
		d.err = err
	}
}

func (d *Decoder) RequestID(val *int32, canBeOptional bool) int32 {
	if val == nil {
		if canBeOptional {
			return -1
		}
		d.SetErr(errors.New("received invalid request id"))
		return 0
	}
	if canBeOptional && *val <= 0 {
		return -1
	}
	if *val < 1 {
		d.SetErr(fmt.Errorf("received invalid request id: %d", *val))
		return 0
	}
	return *val
}

func (d *Decoder) Bool(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}

func (d *Decoder) Int32(val *int32) int32 {
	if val == nil {
		return 0
	}
	return *val
}

func (d *Decoder) Int32Max(val *int32) *int32 {
	return val
}

func (d *Decoder) Int64(val *int64) int64 {
	if val == nil {
		return 0
	}
	return *val
}

func (d *Decoder) Int64Max(val *int64) *int64 {
	return val
}

func (d *Decoder) Float(val *float64) float64 {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return 0
	}
	return *val
}

func (d *Decoder) FloatMax(val *float64) *float64 {
	if val == nil || math.IsNaN(*val) || math.IsInf(*val, 0) {
		return nil
	}
	return val
}

func (d *Decoder) FloatFromString(val *string) float64 {
	if val == nil || len(*val) == 0 {
		return 0
	}

	// Parse value
	value, err := strconv.ParseFloat(*val, 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode int64 value"))
		return 0
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}

func (d *Decoder) FloatMaxFromString(val *string) *float64 {
	if val == nil || len(*val) == 0 {
		return nil
	}

	// Parse value
	value, err := strconv.ParseFloat(*val, 64)
	if err != nil {
		d.SetErr(errors.New("cannot decode int64 value"))
		return nil
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return nil
	}
	return &value
}

func (d *Decoder) EpochTimestamp(val *int64, inMilliseconds bool) time.Time {
	if val == nil {
		return time.Now().UTC()
	}
	if *val < 0 {
		d.SetErr(errors.New("invalid epoch timestamp"))
		return time.Time{}
	}
	if inMilliseconds {
		return time.Unix((*val)/1000, ((*val)%1000)*1000000).UTC()
	}
	return time.Unix(*val, 0).UTC()
}

func (d *Decoder) EpochTimestampFromString(val *string, inMilliseconds bool) time.Time {
	if val == nil {
		return time.Now().UTC()
	}
	value, err := strconv.ParseInt(*val, 10, 64)
	if err != nil || value < 0 {
		d.SetErr(errors.New("invalid epoch timestamp"))
		return time.Time{}
	}
	if inMilliseconds {
		return time.Unix(value/1000, (value%1000)*1000000).UTC()
	}
	return time.Unix(value, 0).UTC()
}

func (d *Decoder) String(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
