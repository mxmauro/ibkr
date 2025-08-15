package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type OrderCancel struct {
	ManualOrderCancelTime string
	ExtOperator           string
	ManualOrderIndicator  *int32
}

// -----------------------------------------------------------------------------

func NewOrderCancel() OrderCancel {
	oc := OrderCancel{}
	return oc
}

func (o OrderCancel) String() string {
	return fmt.Sprintf(
		"ManualOrderCancelTime: %s, ManualOrderIndicator: %s",
		o.ManualOrderCancelTime,
		formatter.Int32MaxString(o.ManualOrderIndicator),
	)
}
