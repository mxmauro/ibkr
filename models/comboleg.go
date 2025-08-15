package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type ComboLeg struct {
	ConID     int32
	Ratio     int32
	Action    string // BUY/SELL/SSHORT
	Exchange  string
	OpenClose int32
	// for stock legs when doing short sale
	ShortSalesSlot     int32 // 1 = clearing broker, 2 = third party
	DesignatedLocation string
	ExemptCode         int32
}

// -----------------------------------------------------------------------------

// NewComboLeg creates a default ComboLeg.
func NewComboLeg() *ComboLeg {
	cl := ComboLeg{}
	cl.ExemptCode = -1
	return &cl
}

func NewComboLegFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.ComboLeg) *ComboLeg {
	if pb == nil {
		return NewComboLeg()
	}
	cl := NewComboLeg()
	cl.ConID = msgDec.Int32(pb.ConId)
	cl.Ratio = msgDec.Int32(pb.Ratio)
	cl.Action = msgDec.String(pb.Action)
	cl.Exchange = msgDec.String(pb.Exchange)
	cl.OpenClose = msgDec.Int32(pb.OpenClose)
	cl.ShortSalesSlot = msgDec.Int32(pb.ShortSalesSlot)
	cl.DesignatedLocation = msgDec.String(pb.DesignatedLocation)
	cl.ExemptCode = msgDec.Int32(pb.ExemptCode)
	return cl
}

func (c *ComboLeg) EncodeMessage(_ int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	msgEnc.Int32(c.ConID)
	msgEnc.Int32(c.Ratio)
	msgEnc.String(c.Action)
	msgEnc.String(c.Exchange)
	return msgEnc.Bytes(), msgEnc.Err()
}

func (c *ComboLeg) Proto(perLegPrice *float64) *protobuf.ComboLeg {
	pb := protobuf.ComboLeg{
		ConId:              protofmt.Int32(c.ConID),
		Ratio:              protofmt.Int32(c.Ratio),
		Action:             protofmt.String(c.Action),
		Exchange:           protofmt.String(c.Exchange),
		OpenClose:          protofmt.Int32(c.OpenClose),
		ShortSalesSlot:     protofmt.Int32(c.ShortSalesSlot),
		DesignatedLocation: protofmt.String(c.DesignatedLocation),
		ExemptCode:         protofmt.Int32(c.ExemptCode),
		PerLegPrice:        protofmt.FloatMax(perLegPrice),
	}
	return &pb
}

func (c *ComboLeg) String() string {
	return fmt.Sprintf(
		"%d, %d, %s, %s, %d, %d, %s, %d",
		c.ConID,
		c.Ratio,
		c.Action,
		c.Exchange,
		c.OpenClose,
		c.ShortSalesSlot,
		c.DesignatedLocation,
		c.ExemptCode,
	)
}
