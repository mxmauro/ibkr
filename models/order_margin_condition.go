package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type MarginCondition struct {
	*OperatorCondition
	Percent int64
}

// -----------------------------------------------------------------------------

func (mc *MarginCondition) decode(msgDec *utils.MessageDecoder) {
	mc.OperatorCondition.decode(msgDec)
	mc.Percent = msgDec.Int64(false)
}

func (mc *MarginCondition) makeFields() []any {
	return append(mc.OperatorCondition.makeFields(), mc.Percent)
}
