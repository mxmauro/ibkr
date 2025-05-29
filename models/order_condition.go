package models

import (
	"errors"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type OrderConditionType = int64

const (
	OrderConditionTypePrice         OrderConditionType = 1
	OrderConditionTypeTime          OrderConditionType = 3
	OrderConditionTypeMargin        OrderConditionType = 4
	OrderConditionTypeExecution     OrderConditionType = 5
	OrderConditionTypeVolume        OrderConditionType = 6
	OrderConditionTypePercentChange OrderConditionType = 7
)

// Trigger Methods.
const (
	TriggerMethodDefault = iota
	TriggerMethodDoubleBidAsk
	TriggerMethodLast
	TriggerMethodDoubleLast
	TriggerMethodBidAsk
	TriggerMethodLastBidAsk
	TriggerMethodMidPoint
)

type OrderCondition interface {
	Type() OrderConditionType
	String() string

	decode(msgDec *utils.MessageDecoder)
	makeFields() []any
}

type orderCondition struct {
	condType                OrderConditionType
	IsConjunctionConnection bool
}

// -----------------------------------------------------------------------------

func NewOrderCondition(condType OrderConditionType) (OrderCondition, error) {
	switch condType {
	case OrderConditionTypePrice:
		return &PriceCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: condType},
				},
			},
		}, nil

	case OrderConditionTypeTime:
		return &TimeCondition{
			OperatorCondition: &OperatorCondition{
				orderCondition: &orderCondition{condType: condType},
			},
		}, nil

	case OrderConditionTypeMargin:
		return &MarginCondition{
			OperatorCondition: &OperatorCondition{
				orderCondition: &orderCondition{condType: condType},
			},
		}, nil

	case OrderConditionTypeExecution:
		return &ExecutionCondition{
			orderCondition: &orderCondition{
				condType: condType,
			},
		}, nil

	case OrderConditionTypeVolume:
		return &VolumeCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: condType},
				},
			},
		}, nil

	case OrderConditionTypePercentChange:
		return &PercentChangeCondition{
			ContractCondition: &ContractCondition{
				OperatorCondition: &OperatorCondition{
					orderCondition: &orderCondition{condType: condType},
				},
			},
		}, nil
	}
	return nil, errors.New("invalid OrderConditionType")
}

func NewOrderConditionFromMessageDecoder(msgDec *utils.MessageDecoder) OrderCondition {
	condType := msgDec.Int64(false)
	if msgDec.Err() != nil {
		return nil
	}
	cond, err := NewOrderCondition(condType)
	if err != nil {
		msgDec.SetErr(err)
		return nil
	}
	cond.decode(msgDec)
	if msgDec.Err() != nil {
		return nil
	}
	return cond
}

// -----------------------------------------------------------------------------

func (oc *orderCondition) Type() OrderConditionType {
	return oc.condType
}

func (oc *orderCondition) decode(msgDec *utils.MessageDecoder) {
	connector := msgDec.String(false)
	oc.IsConjunctionConnection = connector == "a"
}

func (oc *orderCondition) makeFields() []any {
	if oc.IsConjunctionConnection {
		return []any{"a"}
	}
	return []any{"o"}
}

func (oc *orderCondition) String() string {
	if oc.IsConjunctionConnection {
		return "<AND>"
	}
	return "<OR>"
}
