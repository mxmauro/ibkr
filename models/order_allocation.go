package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type OrderAllocation struct {
	Account         string
	Position        Decimal // UNSET_DECIMAL
	PositionDesired Decimal // UNSET_DECIMAL
	PositionAfter   Decimal // UNSET_DECIMAL
	DesiredAllocQty Decimal // UNSET_DECIMAL
	AllowedAllocQty Decimal // UNSET_DECIMAL
	IsMonetary      bool
}

// -----------------------------------------------------------------------------

func NewOrderAllocation() *OrderAllocation {
	oa := &OrderAllocation{
		Position:        UNSET_DECIMAL,
		PositionDesired: UNSET_DECIMAL,
		PositionAfter:   UNSET_DECIMAL,
		DesiredAllocQty: UNSET_DECIMAL,
		AllowedAllocQty: UNSET_DECIMAL,
	}
	return oa
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
