package models

import (
	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderVolumeCondition struct {
	OrderContractCondition
	Volume int32
}

// -----------------------------------------------------------------------------

func (vc *OrderVolumeCondition) decode(msgDec *message.Decoder) {
	vc.OrderContractCondition.decode(msgDec)
	vc.Volume = msgDec.Int32()
}

func (vc *OrderVolumeCondition) makeFields() []any {
	return append(vc.OrderContractCondition.makeFields(), vc.Volume)
}
