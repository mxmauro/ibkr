package models

import (
	"fmt"
	"time"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// HistoricalTick is the historical tick's description.
// Used when requesting historical tick data with whatToShow = MIDPOINT.
type HistoricalTick struct {
	Time  time.Time
	Price float64
	Size  *Decimal
}

// -----------------------------------------------------------------------------

func NewHistoricalTick() HistoricalTick {
	ht := HistoricalTick{}
	return ht
}

func NewHistoricalTickFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.HistoricalTick) HistoricalTick {
	ht := NewHistoricalTick()
	if pb == nil {
		return ht
	}
	ht.Time = msgDec.EpochTimestamp(pb.Time, false)
	ht.Price = msgDec.Float(pb.Price)
	ht.Size = NewDecimalMaxFromProtobufDecoder(msgDec, pb.Size)
	return ht
}

func (h HistoricalTick) String() string {
	return fmt.Sprintf(
		"Time: %s, Price: %f, Size: %s",
		h.Time.Format("2006/01/02 15:04:05"),
		h.Price,
		h.Size.StringMax(),
	)
}
