package models

import (
	"strconv"

	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/quagmt/udecimal"
)

// -----------------------------------------------------------------------------

type Decimal struct {
	value udecimal.Decimal
}

// -----------------------------------------------------------------------------

var (
	DecimalZero = Decimal{value: udecimal.MustFromInt64(0, 0)}
	DecimalOne  = Decimal{value: udecimal.MustFromInt64(1, 0)}
)

// -----------------------------------------------------------------------------

func NewDecimalFromStringWithErr(s string) (Decimal, error) {
	if s == "" {
		return DecimalZero, nil
	}
	f, err := udecimal.Parse(s)
	if err != nil {
		return DecimalZero, err
	}
	return Decimal{value: f}, nil
}

func NewDecimalMaxFromStringWithErr(s string) (*Decimal, error) {
	if isMaxDecimal(s) {
		return nil, nil
	}
	f, err := udecimal.Parse(s)
	if err != nil {
		return nil, err
	}
	dec := Decimal{value: f}
	return &dec, nil
}

func NewDecimalFromMessageDecoder(msgDec *message.Decoder) Decimal {
	data := msgDec.String()

	// Done
	dec, err := NewDecimalFromStringWithErr(data)
	if err != nil {
		msgDec.SetErr(err)
		return DecimalZero
	}
	return dec
}

func NewDecimalMaxFromMessageDecoder(msgDec *message.Decoder) *Decimal {
	data := msgDec.String()

	// Done
	dec, err := NewDecimalMaxFromStringWithErr(data)
	if err != nil {
		msgDec.SetErr(err)
		return nil
	}
	return dec
}

func NewDecimalFromProtobufDecoder(msgDec *protofmt.Decoder, val *string) Decimal {
	// Handle empty value
	if val == nil {
		return DecimalZero
	}

	// Done
	dec, err := NewDecimalFromStringWithErr(*val)
	if err != nil {
		msgDec.SetErr(err)
		return DecimalZero
	}
	return dec
}

func NewDecimalMaxFromProtobufDecoder(msgDec *protofmt.Decoder, val *string) *Decimal {
	// Handle empty value
	if val == nil {
		return nil
	}

	// Done
	dec, err := NewDecimalMaxFromStringWithErr(*val)
	if err != nil {
		msgDec.SetErr(err)
		return nil
	}
	return dec
}

func (d *Decimal) String() string {
	if d == nil {
		return ""
	}
	return d.value.String()
}

func (d *Decimal) StringMax() string {
	if d == nil {
		return ""
	}
	return d.String()
}

func (d *Decimal) Int64() int64 {
	if d == nil {
		return 0
	}
	v, err := d.value.Int64()
	if err != nil {
		return 0
	}
	return v
}

func (d *Decimal) Float() float64 {
	if d == nil {
		return 0
	}
	rounded := d.value
	if rounded.PrecUint() > 15 {
		rounded = rounded.RoundBank(15)
	}
	return rounded.InexactFloat64()
}

func (d *Decimal) EncodeMessage(_ int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	msgEnc.String(d.String())
	return msgEnc.Bytes(), msgEnc.Err()
}

func isMaxDecimal(s string) bool {
	if s == "" || s == "2147483647" || s == "9223372036854775807" || s == "1.7976931348623157E308" || s == "-9223372036854775808" {
		return true
	}

	f, err := strconv.ParseFloat(s, 64)
	return err == nil && (f == 1.7976931348623157e+308 || f == -1.7976931348623157e+308)
}
