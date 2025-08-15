package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

type RealTimeBar struct {
	Time    time.Time
	EndTime time.Time
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volume  *Decimal
	Wap     *Decimal
	Count   int32
}

// -----------------------------------------------------------------------------

func NewRealTimeBar() RealTimeBar {
	rtb := RealTimeBar{}
	return rtb
}

func (rb RealTimeBar) String() string {
	return fmt.Sprintf(
		"Time: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, Wap: %s, Count: %d",
		rb.Time.Format("2006/01/02 15:04:05"),
		rb.Open,
		rb.High,
		rb.Low,
		rb.Close,
		rb.Volume.StringMax(),
		rb.Wap.StringMax(),
		rb.Count,
	)
}
