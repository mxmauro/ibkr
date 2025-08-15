package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

// SoftDollarTier stores the Soft Dollar Tier information.
type SoftDollarTier struct {
	Name        string
	Value       string
	DisplayName string
}

// -----------------------------------------------------------------------------

func NewSoftDollarTier() SoftDollarTier {
	return SoftDollarTier{}
}

func (sdt SoftDollarTier) String() string {
	return fmt.Sprintf("Name: %s, Value: %s, DisplayName: %s",
		sdt.Name,
		sdt.Value,
		sdt.DisplayName)
}
