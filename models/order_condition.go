package models

import (
	"errors"
	"strings"

	"github.com/mxmauro/ibkr/utils/encoders/message"
)

// -----------------------------------------------------------------------------

type OrderConditionType int32

const (
	OrderConditionTypePrice         OrderConditionType = 1
	OrderConditionTypeTime          OrderConditionType = 3
	OrderConditionTypeMargin        OrderConditionType = 4
	OrderConditionTypeExecution     OrderConditionType = 5
	OrderConditionTypeVolume        OrderConditionType = 6
	OrderConditionTypePercentChange OrderConditionType = 7
)

type OrderCondition interface {
	Type() OrderConditionType
	String() string

	decode(msgDec *message.Decoder)
	makeFields() []any
}

type orderConditionBase struct {
	condType                OrderConditionType
	IsConjunctionConnection bool
}

// -----------------------------------------------------------------------------

func NewOrderCondition(condType OrderConditionType) (OrderCondition, error) {
	switch condType {
	case OrderConditionTypePrice:
		return &OrderPriceCondition{
			OrderContractCondition: OrderContractCondition{
				OrderOperatorCondition: OrderOperatorCondition{
					orderConditionBase: orderConditionBase{
						condType: condType,
					},
				},
			},
		}, nil

	case OrderConditionTypeTime:
		return &OrderTimeCondition{
			OrderOperatorCondition: OrderOperatorCondition{
				orderConditionBase: orderConditionBase{
					condType: condType,
				},
			},
		}, nil

	case OrderConditionTypeMargin:
		return &OrderMarginCondition{
			OrderOperatorCondition: OrderOperatorCondition{
				orderConditionBase: orderConditionBase{
					condType: condType,
				},
			},
		}, nil

	case OrderConditionTypeExecution:
		return &OrderExecutionCondition{
			orderConditionBase: orderConditionBase{
				condType: condType,
			},
		}, nil

	case OrderConditionTypeVolume:
		return &OrderVolumeCondition{
			OrderContractCondition: OrderContractCondition{
				OrderOperatorCondition: OrderOperatorCondition{
					orderConditionBase: orderConditionBase{
						condType: condType,
					},
				},
			},
		}, nil

	case OrderConditionTypePercentChange:
		return &OrderPercentChangeCondition{
			OrderContractCondition: OrderContractCondition{
				OrderOperatorCondition: OrderOperatorCondition{
					orderConditionBase: orderConditionBase{
						condType: condType,
					},
				},
			},
		}, nil
	}
	return nil, errors.New("invalid OrderConditionType")
}

func NewOrderConditionFromMessageDecoder(msgDec *message.Decoder) OrderCondition {
	condType := OrderConditionType(msgDec.Int32())
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

func (oc *orderConditionBase) Type() OrderConditionType {
	return oc.condType
}

func (oc *orderConditionBase) decode(msgDec *message.Decoder) {
	connector := msgDec.String()
	oc.IsConjunctionConnection = strings.Compare(connector, "a") == 0
}

func (oc *orderConditionBase) makeFields() []any {
	if oc.IsConjunctionConnection {
		return []any{"a"}
	}
	return []any{"o"}
}

func (oc *orderConditionBase) String() string {
	if oc.IsConjunctionConnection {
		return "<AND>"
	}
	return "<OR>"
}
