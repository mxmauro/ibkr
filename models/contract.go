package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

// Contract describes an instrument's definition.
type Contract struct {
	ConID                        int32
	Symbol                       string
	SecType                      SecurityType
	LastTradeDateOrContractMonth string
	LastTradeDate                string
	Strike                       *float64
	Right                        string   // Either Put or Call (i.e.: Options). Valid values are P, PUT, C, CALL.
	Multiplier                   *float64 // Originally it was a string
	Exchange                     string
	PrimaryExchange              string // Pick an actual (i.e.: non-aggregate) exchange that the contract trades on.  DO NOT SET TO SMART.
	Currency                     string
	LocalSymbol                  string
	TradingClass                 string
	IncludeExpired               bool
	SecIDType                    string // CUSIP;SEDOL;ISIN;RIC
	SecID                        string
	Description                  string
	IssuerID                     string
	ComboLegsDescription         string // received in open order 14 and up for all combos
	ComboLegs                    []*ComboLeg
	DeltaNeutralContract         *DeltaNeutralContract
}

// -----------------------------------------------------------------------------

func NewContract() *Contract {
	return &Contract{}
}

func NewContractFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.Contract) *Contract {
	c := NewContract()
	if pb == nil {
		return c
	}
	c.ConID = msgDec.Int32(pb.ConId)
	c.Symbol = msgDec.String(pb.Symbol)
	c.SecType = SecurityType(msgDec.String(pb.SecType))
	c.LastTradeDateOrContractMonth = msgDec.String(pb.LastTradeDateOrContractMonth)
	c.LastTradeDate = msgDec.String(pb.LastTradeDate)
	c.Strike = msgDec.FloatMax(pb.Strike)
	c.Right = msgDec.String(pb.Right)
	c.Multiplier = msgDec.FloatMax(pb.Multiplier)
	c.Exchange = msgDec.String(pb.Exchange)
	c.PrimaryExchange = msgDec.String(pb.PrimaryExch)
	c.Currency = msgDec.String(pb.Currency)
	c.LocalSymbol = msgDec.String(pb.LocalSymbol)
	c.TradingClass = msgDec.String(pb.TradingClass)
	c.IncludeExpired = msgDec.Bool(pb.IncludeExpired)
	c.SecIDType = msgDec.String(pb.SecIdType)
	c.SecID = msgDec.String(pb.SecId)
	c.Description = msgDec.String(pb.Description)
	c.IssuerID = msgDec.String(pb.IssuerId)
	c.ComboLegsDescription = msgDec.String(pb.ComboLegsDescrip)
	c.ComboLegs = make([]*ComboLeg, len(pb.ComboLegs))
	for idx, pbcl := range pb.ComboLegs {
		c.ComboLegs[idx] = NewComboLegFromProtobufDecoder(msgDec, pbcl)
	}
	c.DeltaNeutralContract = NewDeltaNeutralContractFromProtobufDecoder(msgDec, pb.DeltaNeutralContract)
	return c
}

func (c *Contract) Equal(other *Contract) bool {
	if c.ConID != 0 && other.ConID != 0 {
		return c.ConID == other.ConID
	}
	if len(c.SecIDType) > 0 && len(other.SecIDType) > 0 && c.SecIDType == other.SecIDType {
		return c.SecID == other.SecID
	}
	return c.Symbol == other.Symbol &&
		c.SecType == other.SecType &&
		c.Exchange == other.Exchange &&
		c.Currency == other.Currency &&
		c.LastTradeDate == other.LastTradeDate &&
		c.Strike == other.Strike &&
		c.Right == other.Right
}

func (c *Contract) EncodeMessage(mode int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	msgEnc.Int32(c.ConID)
	msgEnc.String(c.Symbol)
	msgEnc.String(string(c.SecType))
	msgEnc.String(c.LastTradeDateOrContractMonth)
	msgEnc.FloatMax(c.Strike)
	msgEnc.String(c.Right)
	msgEnc.FloatMax(c.Multiplier)
	msgEnc.String(c.Exchange)
	msgEnc.String(c.PrimaryExchange)
	msgEnc.String(c.Currency)
	msgEnc.String(c.LocalSymbol)
	msgEnc.String(c.TradingClass)
	if mode == 2 || mode == 3 {
		msgEnc.Bool(c.IncludeExpired)
	}
	if mode == 3 {
		msgEnc.String(c.SecIDType)
		msgEnc.String(c.SecID)
		msgEnc.String(c.IssuerID)
	}
	return msgEnc.Bytes(), msgEnc.Err()
}

func (c *Contract) Proto(o *Order) *protobuf.Contract {
	pb := protobuf.Contract{
		ConId:                        protofmt.Int32(c.ConID),
		Symbol:                       protofmt.String(c.Symbol),
		SecType:                      protofmt.String(string(c.SecType)),
		LastTradeDateOrContractMonth: protofmt.String(c.LastTradeDateOrContractMonth),
		Strike:                       protofmt.FloatMax(c.Strike),
		Right:                        protofmt.String(c.Right),
		Multiplier:                   protofmt.FloatMax(c.Multiplier),
		Exchange:                     protofmt.String(c.Exchange),
		PrimaryExch:                  protofmt.String(c.PrimaryExchange),
		Currency:                     protofmt.String(c.Currency),
		LocalSymbol:                  protofmt.String(c.LocalSymbol),
		TradingClass:                 protofmt.String(c.TradingClass),
		SecIdType:                    protofmt.String(c.SecIDType),
		SecId:                        protofmt.String(c.SecID),
		Description:                  protofmt.String(c.Description),
		IssuerId:                     protofmt.String(c.IssuerID),
		IncludeExpired:               protofmt.Bool(c.IncludeExpired),
		ComboLegsDescrip:             protofmt.String(c.ComboLegsDescription),
		LastTradeDate:                protofmt.String(c.LastTradeDate),
	}
	if c.DeltaNeutralContract != nil {
		pb.DeltaNeutralContract = c.DeltaNeutralContract.Proto()
	}
	for idx := range c.ComboLegs {
		var perLegPrice *float64
		if o != nil && idx < len(o.OrderComboLegs) {
			perLegPrice = o.OrderComboLegs[idx].Price
		}
		pbcl := c.ComboLegs[idx].Proto(perLegPrice)
		pb.ComboLegs = append(pb.ComboLegs, pbcl)
	}
	return &pb
}

func (c *Contract) String() string {
	s := fmt.Sprintf("%d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		c.ConID,
		c.Symbol,
		c.SecType,
		c.LastTradeDateOrContractMonth,
		c.LastTradeDate,
		formatter.FloatMaxString(c.Strike),
		c.Right,
		formatter.FloatMaxString(c.Multiplier),
		c.Exchange,
		c.PrimaryExchange,
		c.Currency,
		c.LocalSymbol,
		c.TradingClass,
		c.IncludeExpired,
		c.SecIDType,
		c.SecID,
		c.Description,
		c.IssuerID,
	)
	if len(c.ComboLegs) > 1 {
		s += ", combo:" + c.ComboLegsDescription
		for _, leg := range c.ComboLegs {
			s += fmt.Sprintf("; %s", leg)
		}
	}

	if c.DeltaNeutralContract != nil {
		s += fmt.Sprintf("; %s", c.DeltaNeutralContract)
	}

	return s
}
