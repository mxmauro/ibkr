package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderOperatorCondition struct {
	orderConditionBase
	IsMore bool
}

// -----------------------------------------------------------------------------

func (oc *OrderOperatorCondition) decode(msgDec *message.Decoder) {
	oc.orderConditionBase.decode(msgDec)
	oc.IsMore = msgDec.Bool()
}

func (oc *OrderOperatorCondition) makeFields() []any {
	return append(oc.orderConditionBase.makeFields(), oc.IsMore)
}
