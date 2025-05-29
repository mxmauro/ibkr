package models

import (
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type TimeCondition struct {
	*OperatorCondition
	Time string
}

// -----------------------------------------------------------------------------

func (tc *TimeCondition) decode(msgDec *utils.MessageDecoder) {
	tc.OperatorCondition.decode(msgDec)
	// tc.Time = decodeTime(fields[2], "20060102")
	tc.Time = msgDec.String(false)
}

func (tc *TimeCondition) makeFields() []any {
	return append(tc.OperatorCondition.makeFields(), tc.Time)
}
