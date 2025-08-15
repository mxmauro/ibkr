package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type OrderAllocation struct {
	Account         string
	Position        *Decimal
	PositionDesired *Decimal
	PositionAfter   *Decimal
	DesiredAllocQty *Decimal
	AllowedAllocQty *Decimal
	IsMonetary      bool
}

// -----------------------------------------------------------------------------

func NewOrderAllocation() *OrderAllocation {
	oa := OrderAllocation{}
	return &oa
}

func (oa *OrderAllocation) String() string {
	return fmt.Sprint(
		"Account: ", oa.Account,
		", Position: ", oa.Position.StringMax(),
		", PositionDesired: ", oa.PositionDesired.StringMax(),
		", PositionAfter: ", oa.PositionAfter.StringMax(),
		", DesiredAllocQty: ", oa.DesiredAllocQty.StringMax(),
		", AllowedAllocQty: ", oa.AllowedAllocQty.StringMax(),
		", IsMonetary: ", oa.IsMonetary,
	)
}
