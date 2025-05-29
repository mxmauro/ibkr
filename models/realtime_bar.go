package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type RealTimeBar struct {
	Time    int64
	EndTime int64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volume  Decimal
	Wap     Decimal
	Count   int64
}

// -----------------------------------------------------------------------------

func NewRealTimeBar() RealTimeBar {
	rtb := RealTimeBar{}
	rtb.Volume = UNSET_DECIMAL
	rtb.Wap = UNSET_DECIMAL
	return rtb
}

func (rb RealTimeBar) String() string {
	return fmt.Sprintf("Time: %d, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, Wap: %s, Count: %d",
		rb.Time, rb.Open, rb.High, rb.Low, rb.Close, rb.Volume.StringMax(), rb.Wap.StringMax(), rb.Count)
}
