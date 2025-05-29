package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type PriceIncrement struct {
	LowEdge   float64
	Increment float64
}

// -----------------------------------------------------------------------------

func newPriceIncrement() PriceIncrement {
	return PriceIncrement{}
}

func (p PriceIncrement) String() string {
	return fmt.Sprintf("LowEdge: %f, Increment: %f", p.LowEdge, p.Increment)
}
