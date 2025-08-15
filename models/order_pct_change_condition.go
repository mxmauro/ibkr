package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderPercentChangeCondition struct {
	OrderContractCondition
	ChangePercent float64
}

// -----------------------------------------------------------------------------

func (pcc *OrderPercentChangeCondition) decode(msgDec *message.Decoder) {
	pcc.OrderContractCondition.decode(msgDec)
	pcc.ChangePercent = msgDec.Float()
}

func (pcc *OrderPercentChangeCondition) makeFields() []any {
	return append(pcc.OrderContractCondition.makeFields(), pcc.ChangePercent)
}
