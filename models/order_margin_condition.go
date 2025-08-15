package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderMarginCondition struct {
	OrderOperatorCondition
	Percent int32
}

// -----------------------------------------------------------------------------

func (mc *OrderMarginCondition) decode(msgDec *message.Decoder) {
	mc.OrderOperatorCondition.decode(msgDec)
	mc.Percent = msgDec.Int32()
}

func (mc *OrderMarginCondition) makeFields() []any {
	return append(mc.OrderOperatorCondition.makeFields(), mc.Percent)
}
