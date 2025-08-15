package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderTimeCondition struct {
	OrderOperatorCondition
	Time string
}

// -----------------------------------------------------------------------------

func (tc *OrderTimeCondition) decode(msgDec *message.Decoder) {
	tc.OrderOperatorCondition.decode(msgDec)
	// tc.Time = decodeTime(fields[2], "20060102")
	tc.Time = msgDec.String()
}

func (tc *OrderTimeCondition) makeFields() []any {
	return append(tc.OrderOperatorCondition.makeFields(), tc.Time)
}
