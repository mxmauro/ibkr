package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type HistogramData struct {
	Price float64
	Size  *Decimal
}

// -----------------------------------------------------------------------------

func NewHistogramData() HistogramData {
	hd := HistogramData{}
	return hd
}

func (hd HistogramData) String() string {
	return fmt.Sprintf("Price: %v, Size: %v", hd.Price, hd.Size.String())
}
