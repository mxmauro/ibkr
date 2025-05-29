package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type ContractCondition struct {
	*OperatorCondition
	ConID    int64
	Exchange string
}

// -----------------------------------------------------------------------------

func (cc *ContractCondition) decode(msgDec *utils.MessageDecoder) {
	cc.OperatorCondition.decode(msgDec)
	cc.ConID = msgDec.Int64(false)
	cc.Exchange = msgDec.String(false)
}

func (cc *ContractCondition) makeFields() []any {
	return append(cc.OperatorCondition.makeFields(), cc.ConID, cc.Exchange)
}
