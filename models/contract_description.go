package models

import (
	"strings"
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
