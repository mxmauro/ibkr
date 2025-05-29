package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type ExecutionCondition struct {
	*orderCondition
	SecType  SecurityType
	Exchange string
	Symbol   string
}

// -----------------------------------------------------------------------------

func (ec *ExecutionCondition) decode(msgDec *utils.MessageDecoder) {
	ec.orderCondition.decode(msgDec)
	ec.SecType = NewSecurityTypeFromString(msgDec.String(false))
	ec.Exchange = msgDec.String(false)
	ec.Symbol = msgDec.String(false)
}

func (ec *ExecutionCondition) makeFields() []any {
	return append(ec.orderCondition.makeFields(), ec.SecType, ec.Exchange, ec.Symbol)
}

func (ec *ExecutionCondition) String() string {
	return fmt.Sprintf("trade occurs for %v symbol on %v exchange for %v security type", ec.Symbol, ec.Exchange, ec.SecType)
}
