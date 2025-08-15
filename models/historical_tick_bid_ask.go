package models

import (
	"fmt"
	"time"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// HistoricalTickBidAsk is the historical tick's description.
// Used when requesting historical tick data with whatToShow = BID_ASK.
type HistoricalTickBidAsk struct {
	Time             time.Time // Epoch in seconds
	TickAttribBidAsk TickAttribBidAsk
	PriceBid         float64
	PriceAsk         float64
	SizeBid          *Decimal
	SizeAsk          *Decimal
}

// -----------------------------------------------------------------------------

func NewHistoricalTickBidAsk() HistoricalTickBidAsk {
	htba := HistoricalTickBidAsk{}
	return htba
}

func NewHistoricalTickBidAskFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.HistoricalTickBidAsk) HistoricalTickBidAsk {
	htba := NewHistoricalTickBidAsk()
	if pb == nil {
		return htba
	}
	htba.Time = msgDec.EpochTimestamp(pb.Time, false)
	htba.PriceBid = msgDec.Float(pb.PriceBid)
	htba.PriceAsk = msgDec.Float(pb.PriceAsk)
	htba.SizeBid = NewDecimalMaxFromProtobufDecoder(msgDec, pb.SizeBid)
	htba.SizeAsk = NewDecimalMaxFromProtobufDecoder(msgDec, pb.SizeAsk)
	return htba
}

func (htba HistoricalTickBidAsk) String() string {
	return fmt.Sprintf(
		"Time: %s, TickAttriBidAsk: %s, PriceBid: %f, PriceAsk: %f, SizeBid: %s, SizeAsk: %s",
		htba.Time.Format("2006/01/02 15:04:05"),
		htba.TickAttribBidAsk.String(),
		htba.PriceBid,
		htba.PriceAsk,
		htba.SizeBid.StringMax(),
		htba.SizeAsk.StringMax(),
	)
}
