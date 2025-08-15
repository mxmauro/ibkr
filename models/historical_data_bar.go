package models

import (
	"fmt"
	"time"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type HistoricalDataBar struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume *Decimal
	Wap    *Decimal
	Count  int32
}

// -----------------------------------------------------------------------------

func NewHistoricalDataBar() HistoricalDataBar {
	bar := HistoricalDataBar{}
	return bar
}

func NewHistoricalDataBarFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.HistoricalDataBar) HistoricalDataBar {
	bar := NewHistoricalDataBar()
	if pb == nil {
		return bar
	}
	// Epoch because the request of historical data has the date format equal to 2.
	bar.Date = msgDec.EpochTimestampFromString(pb.Date, false)
	bar.Open = msgDec.Float(pb.Open)
	bar.High = msgDec.Float(pb.High)
	bar.Low = msgDec.Float(pb.Low)
	bar.Close = msgDec.Float(pb.Close)
	bar.Volume = NewDecimalMaxFromProtobufDecoder(msgDec, pb.Volume)
	bar.Wap = NewDecimalMaxFromProtobufDecoder(msgDec, pb.WAP)
	bar.Count = msgDec.Int32(pb.BarCount)
	return bar
}

func (bar HistoricalDataBar) String() string {
	return fmt.Sprintf(
		"Time: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, WAP: %s, Count: %d",
		bar.Date.Format("2006-01-02 15:04:05"),
		bar.Open,
		bar.High,
		bar.Low,
		bar.Close,
		bar.Volume.StringMax(),
		bar.Wap.StringMax(),
		bar.Count,
	)
}
