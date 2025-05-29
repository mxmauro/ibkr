package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type OrderCancel struct {
	ManualOrderCancelTime string
	ExtOperator           string
	ManualOrderIndicator  int64
}

// -----------------------------------------------------------------------------

func NewOrderCancel() OrderCancel {
	oc := OrderCancel{
		ManualOrderIndicator: common.UNSET_INT,
	}
	return oc
}

func (o OrderCancel) String() string {
	return fmt.Sprintf("ManualOrderCancelTime: %s, ManualOrderIndicator: %s",
		o.ManualOrderCancelTime,
		utils.IntMaxString(o.ManualOrderIndicator))
}
