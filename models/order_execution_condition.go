package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderExecutionCondition struct {
	orderConditionBase
	SecType  SecurityType
	Exchange string
	Symbol   string
}

// -----------------------------------------------------------------------------

func (ec *OrderExecutionCondition) decode(msgDec *message.Decoder) {
	ec.orderConditionBase.decode(msgDec)
	ec.SecType = NewSecurityTypeFromString(msgDec.String())
	ec.Exchange = msgDec.String()
	ec.Symbol = msgDec.String()
}

func (ec *OrderExecutionCondition) makeFields() []any {
	return append(ec.orderConditionBase.makeFields(), ec.SecType, ec.Exchange, ec.Symbol)
}

func (ec *OrderExecutionCondition) String() string {
	return fmt.Sprintf("trade occurs for %v symbol on %v exchange for %v security type", ec.Symbol, ec.Exchange, ec.SecType)
}
