package models

import (
	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type OrderComboLeg struct {
	Price float64 `default:"UNSET_FLOAT"`
}

// -----------------------------------------------------------------------------

func NewOrderComboLeg() OrderComboLeg {
	ocl := OrderComboLeg{
		Price: common.UNSET_FLOAT,
	}
	return ocl
}

func (o OrderComboLeg) String() string {
	return utils.FloatMaxString(o.Price)
}
