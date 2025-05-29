package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type OperatorCondition struct {
	*orderCondition
	IsMore bool
}

// -----------------------------------------------------------------------------

func (oc *OperatorCondition) decode(msgDec *utils.MessageDecoder) {
	oc.orderCondition.decode(msgDec)
	oc.IsMore = msgDec.Bool()
}

func (oc *OperatorCondition) makeFields() []any {
	return append(oc.orderCondition.makeFields(), oc.IsMore)
}
