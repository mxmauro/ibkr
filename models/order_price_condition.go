package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type PriceCondition struct {
	*ContractCondition
	Price         float64
	TriggerMethod int64
}

// -----------------------------------------------------------------------------

func (pc *PriceCondition) decode(msgDec *utils.MessageDecoder) {
	pc.ContractCondition.decode(msgDec)
	pc.Price = msgDec.Float64(false)
	pc.TriggerMethod = msgDec.Int64(false)
}

func (pc *PriceCondition) makeFields() []any {
	return append(pc.ContractCondition.makeFields(), pc.Price, pc.TriggerMethod)
}
