package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderContractCondition struct {
	OrderOperatorCondition
	ConID    int32
	Exchange string
}

// -----------------------------------------------------------------------------

func (cc *OrderContractCondition) decode(msgDec *message.Decoder) {
	cc.OrderOperatorCondition.decode(msgDec)
	cc.ConID = msgDec.Int32()
	cc.Exchange = msgDec.String()
}

func (cc *OrderContractCondition) makeFields() []any {
	return append(cc.OrderOperatorCondition.makeFields(), cc.ConID, cc.Exchange)
}
