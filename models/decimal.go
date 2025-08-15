package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/robaho/fixed"
)

// -----------------------------------------------------------------------------

type Decimal fixed.Fixed

// -----------------------------------------------------------------------------

var (
	DecimalZero = Decimal(fixed.ZERO)
	DecimalOne  = Decimal(fixed.NewF(1))
)

// -----------------------------------------------------------------------------

func NewDecimalFromStringWithErr(s string) (Decimal, error) {
	if s == "" {
		return DecimalZero, nil
	}
	f, err := fixed.NewSErr(s)
	if err != nil {
		return DecimalZero, err
	}
	return Decimal(f), nil
}

func NewDecimalMaxFromStringWithErr(s string) (*Decimal, error) {
	if isMaxDecimal(s) {
		return nil, nil
	}
	f, err := fixed.NewSErr(s)
	if err != nil {
		return nil, err
	}
	dec := Decimal(f)
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
	return fixed.Fixed(*d).String()
}

func (d *Decimal) StringMax() string {
	if d == nil {
		return ""
	}
	return d.String()
}

func (d *Decimal) Int64() int64 {
	return fixed.Fixed(*d).Int()
}

func (d *Decimal) Float() float64 {
	return fixed.Fixed(*d).Float()
}

func (d *Decimal) EncodeMessage(_ int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	msgEnc.String(d.String())
	return msgEnc.Bytes(), msgEnc.Err()
}

func isMaxDecimal(s string) bool {
	return s == "" || s == "2147483647" || s == "9223372036854775807" || s == "1.7976931348623157E308" || s == "-9223372036854775808"
}
