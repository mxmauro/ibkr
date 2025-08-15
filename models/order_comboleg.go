package models

import (
	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type OrderComboLeg struct {
	Price *float64
}

// -----------------------------------------------------------------------------

func NewOrderComboLeg() OrderComboLeg {
	ocl := OrderComboLeg{}
	return ocl
}

func (o OrderComboLeg) String() string {
	return formatter.FloatMaxString(o.Price)
}
