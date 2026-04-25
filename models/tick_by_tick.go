package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

type TickByTickType int32

const (
	TickByTickTypeUnset TickByTickType = iota
	TickByTickTypeLast
	TickByTickTypeAllLast
	TickByTickTypeBidAsk
	TickByTickTypeMidPoint
)

func (ttt TickByTickType) String() string {
	switch ttt {
	case TickByTickTypeLast:
		return "Last"
	case TickByTickTypeAllLast:
		return "AllLast"
	case TickByTickTypeBidAsk:
		return "BidAsk"
	case TickByTickTypeMidPoint:
		return "MidPoint"
	}
	return ""
}

type TickByTickData interface {
	TickByTickType() TickByTickType
	String() string
}

type TickByTickAllLast struct {
	tickType          TickByTickType
	Time              time.Time
	Price             float64
	Size              *Decimal
	TickAttribLast    TickAttribLast
	Exchange          string
	SpecialConditions string
}

type TickByTickBidAsk struct {
	Time             time.Time
	PriceBid         float64
	PriceAsk         float64
	SizeBid          *Decimal
	SizeAsk          *Decimal
	TickAttribBidAsk TickAttribBidAsk
}

type TickByTickMidPoint struct {
	Time     time.Time
	MidPoint float64
}

// -----------------------------------------------------------------------------

func NewTickByTickAllLast(tickType TickByTickType) *TickByTickAllLast {
	return &TickByTickAllLast{
		tickType: tickType,
	}
}

func NewTickByTickBidAsk() *TickByTickBidAsk {
	return &TickByTickBidAsk{}
}

func NewTickByTickMidPoint() *TickByTickMidPoint {
	return &TickByTickMidPoint{}
}

func (t *TickByTickAllLast) TickByTickType() TickByTickType {
	return t.tickType
}

func (t *TickByTickBidAsk) TickByTickType() TickByTickType {
	return TickByTickTypeBidAsk
}

func (t *TickByTickMidPoint) TickByTickType() TickByTickType {
	return TickByTickTypeMidPoint
}

func (t *TickByTickAllLast) String() string {
	return fmt.Sprintf(
		"Type=%s, Time=%s, Price=%f, Size=%s, Exchange=%s, SpecialConditions=%s, TickAttribLast=[%s]",
		t.tickType.String(),
		t.Time.Format("2006/01/02 15:04:05"),
		t.Price,
		t.Size.StringMax(),
		t.Exchange,
		t.SpecialConditions,
		t.TickAttribLast.String(),
	)
}

func (t *TickByTickBidAsk) String() string {
	return fmt.Sprintf(
		"Type=%s, Time=%s, Bid=%f, Ask=%f, BidSize=%s, AskSize=%s, TickAttribBidAsk=[%s]",
		TickByTickTypeBidAsk.String(),
		t.Time.Format("2006/01/02 15:04:05"),
		t.PriceBid,
		t.PriceAsk,
		t.SizeBid.StringMax(),
		t.SizeAsk.StringMax(),
		t.TickAttribBidAsk.String(),
	)
}

func (t *TickByTickMidPoint) String() string {
	return fmt.Sprintf(
		"Type=%s, Time=%s, MidPoint=%f",
		TickByTickTypeMidPoint.String(),
		t.Time.Format("2006/01/02 15:04:05"),
		t.MidPoint,
	)
}
