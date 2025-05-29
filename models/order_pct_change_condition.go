package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type PercentChangeCondition struct {
	*ContractCondition
	ChangePercent float64
}

// -----------------------------------------------------------------------------

func (pcc *PercentChangeCondition) decode(msgDec *utils.MessageDecoder) {
	pcc.ContractCondition.decode(msgDec)
	pcc.ChangePercent = msgDec.Float64(false)
}

func (pcc *PercentChangeCondition) makeFields() []any {
	return append(pcc.ContractCondition.makeFields(), pcc.ChangePercent)
}
