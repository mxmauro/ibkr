package models

import (
	"errors"

	"github.com/mxmauro/ibkr/utils"
	"github.com/robaho/fixed"
)

// -----------------------------------------------------------------------------

type Decimal fixed.Fixed

// -----------------------------------------------------------------------------

var (
	UNSET_DECIMAL = Decimal(fixed.NaN)
	DECIMAL_ZERO  = Decimal(fixed.ZERO)
	DECIMAL_ONE   = Decimal(fixed.NewF(1))
)

// -----------------------------------------------------------------------------

func NewDecimalFromString(s string) Decimal {
	d, _ := NewDecimalFromStringWithErr(s)
	return d
}

func NewDecimalFromStringWithErr(s string) (Decimal, error) {
	if s == "" || s == "2147483647" || s == "9223372036854775807" || s == "1.7976931348623157E308" || s == "-9223372036854775808" {
		return UNSET_DECIMAL, errors.New("unset decimal")
	}
	f, err := fixed.NewSErr(s)
	if err != nil {
		return UNSET_DECIMAL, err
	}
	return Decimal(f), nil
}

func NewDecimalFromMessageDecoder(msgDec *utils.MessageDecoder, emptyIsUnset bool) Decimal {
	data := msgDec.String(false)

	// Handle empty value
	if len(data) == 0 {
		if emptyIsUnset {
			return UNSET_DECIMAL
		}
		return Decimal{}
	}

	// Done
	dec, err := NewDecimalFromStringWithErr(data)
	if err != nil {
		msgDec.SetErr(err)
		if emptyIsUnset {
			return UNSET_DECIMAL
		}
		return Decimal{}
	}
	return dec
}

func (d *Decimal) String() string {
	return fixed.Fixed(*d).String()
}

func (d *Decimal) StringMax() string {
	if *d == UNSET_DECIMAL {
		return ""
	}
	return d.String()
}

func (d *Decimal) Int() int64 {
	return fixed.Fixed(*d).Int()
}

func (d *Decimal) Float() float64 {
	return fixed.Fixed(*d).Float()
}

func (d *Decimal) EncodeMessage() []byte {
	msgEnc := utils.NewRawMessageEncoder().
		String(d.String())
	return msgEnc.Bytes()
}

/*
func (d *Decimal) MarshalBinary() ([]byte, error) {
	return fixed.Fixed(*d).MarshalBinary()
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	var f fixed.Fixed
	err := f.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	*d = Decimal(f)
	return nil
}
*/
