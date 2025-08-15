package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderPriceCondition struct {
	OrderContractCondition
	Price         float64
	TriggerMethod TriggerMethod
}

// -----------------------------------------------------------------------------

func (pc *OrderPriceCondition) decode(msgDec *message.Decoder) {
	pc.OrderContractCondition.decode(msgDec)
	pc.Price = msgDec.Float()
	pc.TriggerMethod = TriggerMethod(msgDec.Int32())
}

func (pc *OrderPriceCondition) makeFields() []any {
	return append(pc.OrderContractCondition.makeFields(), pc.Price, pc.TriggerMethod)
}
