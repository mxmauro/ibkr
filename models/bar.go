package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

type Bar struct {
	Date     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   Decimal
	Wap      Decimal
	BarCount int64
}

// -----------------------------------------------------------------------------

func NewBar() Bar {
	b := Bar{}
	b.Volume = UNSET_DECIMAL
	b.Wap = UNSET_DECIMAL
	return b
}

func (b Bar) String() string {
	return fmt.Sprintf("Time: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, WAP: %s, BarCount: %d",
		b.Date.Format("2006/01/02 15:04:05"), b.Open, b.High, b.Low, b.Close, b.Volume.StringMax(), b.Wap.StringMax(), b.BarCount)
}
