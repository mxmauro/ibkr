package models

import (
	"strings"

	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// ContractDescription includes contract and DerivativeSecTypes.
type ContractDescription struct {
	Contract           *Contract
	DerivativeSecTypes []SecurityType
}

// -----------------------------------------------------------------------------

func NewContractDescription() *ContractDescription {
	cd := ContractDescription{
		Contract: NewContract(),
	}
	return &cd
}

func NewContractDescriptionFromProtobufDecoder(msgDec *protofmt.Decoder, pb *protobuf.ContractDescription) *ContractDescription {
	cd := NewContractDescription()
	if pb == nil {
		return cd
	}
	cd.Contract = NewContractFromProtobufDecoder(msgDec, pb.Contract)
	cd.DerivativeSecTypes = make([]SecurityType, len(pb.DerivativeSecTypes))
	for idx, dst := range pb.DerivativeSecTypes {
		cd.DerivativeSecTypes[idx] = NewSecurityTypeFromString(dst)
	}
	return cd
}

func (c *ContractDescription) String() string {
	sb := strings.Builder{}
	_, _ = sb.WriteString(c.Contract.String())
	_, _ = sb.WriteString(", [")
	for idx, st := range c.DerivativeSecTypes {
		if idx > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(string(st))
	}
	_, _ = sb.WriteString("]")
	return sb.String()
}
