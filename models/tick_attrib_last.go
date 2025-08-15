package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type TickAttribLast struct {
	PastLimit  bool
	Unreported bool
}

// -----------------------------------------------------------------------------

func NewTickAttribLast() TickAttribLast {
	return TickAttribLast{}
}

func NewTickAttribLastFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.TickAttribLast) TickAttribLast {
	tal := TickAttribLast{}
	if pb == nil {
		return tal
	}
	tal.PastLimit = msgDec.Bool(pb.PastLimit)
	tal.Unreported = msgDec.Bool(pb.Unreported)
	return tal
}

func (tal *TickAttribLast) String() string {
	return fmt.Sprintf("PastLimit: %t, Unreported: %t", tal.PastLimit, tal.Unreported)
}
