package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type DeltaNeutralContract struct {
	ConID int32
	Delta float64
	Price float64
}

// -----------------------------------------------------------------------------

func NewDeltaNeutralContract() *DeltaNeutralContract {
	return &DeltaNeutralContract{}
}

func NewDeltaNeutralContractFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.DeltaNeutralContract) *DeltaNeutralContract {
	if pb == nil {
		return NewDeltaNeutralContract()
	}
	dnc := &DeltaNeutralContract{}
	dnc.ConID = msgDec.Int32(pb.ConId)
	dnc.Delta = msgDec.Float(pb.Delta)
	dnc.Price = msgDec.Float(pb.Price)
	return dnc
}

func (c *DeltaNeutralContract) EncodeMessage(_ int) ([]byte, error) {
	msgEnc := message.NewRawEncoder()
	msgEnc.Int32(c.ConID)
	msgEnc.Float(c.Delta)
	msgEnc.Float(c.Price)
	return msgEnc.Bytes(), msgEnc.Err()
}

func (c *DeltaNeutralContract) Proto() *protobuf.DeltaNeutralContract {
	pb := protobuf.DeltaNeutralContract{
		ConId: protofmt.Int32(c.ConID),
		Delta: protofmt.Float(c.Delta),
		Price: protofmt.Float(c.Price),
	}
	return &pb
}

func (c *DeltaNeutralContract) String() string {
	return fmt.Sprintf("%d, %f, %f", c.ConID, c.Delta, c.Price)
}
