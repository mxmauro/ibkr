package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type DeltaNeutralContract struct {
	ConID int64
	Delta float64
	Price float64
}

// -----------------------------------------------------------------------------

func NewDeltaNeutralContract() DeltaNeutralContract {
	return DeltaNeutralContract{}
}

func (c DeltaNeutralContract) String() string {
	return fmt.Sprintf("%d, %f, %f", c.ConID, c.Delta, c.Price)
}
