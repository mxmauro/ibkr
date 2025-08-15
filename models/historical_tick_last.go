package models

import (
	"fmt"
	"time"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// HistoricalTickLast is the historical last tick's description.
// Used when requesting historical tick data with whatToShow = TRADES.
type HistoricalTickLast struct {
	Time              time.Time // Epoch in seconds
	TickAttribLast    TickAttribLast
	Price             float64
	Size              *Decimal
	Exchange          string
	SpecialConditions string
}

// -----------------------------------------------------------------------------

func NewHistoricalTickLast() HistoricalTickLast {
	htl := HistoricalTickLast{}
	return htl
}

func NewHistoricalTickLastFromProtobufDecoder(
	msgDec *protofmt.Decoder, pb *protobuf.HistoricalTickLast,
) HistoricalTickLast {
	htl := NewHistoricalTickLast()
	if pb == nil {
		return htl
	}
	htl.Time = msgDec.EpochTimestamp(pb.Time, false)
	htl.TickAttribLast = NewTickAttribLastFromProtobufDecoder(msgDec, pb.TickAttribLast)
	htl.Price = msgDec.Float(pb.Price)
	htl.Size = NewDecimalMaxFromProtobufDecoder(msgDec, pb.Size)
	htl.Exchange = msgDec.String(pb.Exchange)
	htl.SpecialConditions = msgDec.String(pb.SpecialConditions)
	return htl
}

func (htl HistoricalTickLast) String() string {
	return fmt.Sprintf(
		"Time: %s, Price: %f, Size: %s, Exchange: %s, SpecialConditions: %s, TickAttribLast [%s]",
		htl.Time.Format("2006/01/02 15:04:05"),
		htl.Price,
		htl.Size.StringMax(),
		htl.Exchange,
		htl.SpecialConditions,
		htl.TickAttribLast.String(),
	)
}
