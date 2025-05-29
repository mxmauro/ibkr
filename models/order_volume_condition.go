package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type VolumeCondition struct {
	*ContractCondition
	Volume int64
}

// -----------------------------------------------------------------------------

func (vc *VolumeCondition) decode(msgDec *utils.MessageDecoder) {
	vc.ContractCondition.decode(msgDec)
	vc.Volume = msgDec.Int64(false)
}

func (vc *VolumeCondition) makeFields() []any {
	return append(vc.ContractCondition.makeFields(), vc.Volume)
}
