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

func newSoftDollarTier() SoftDollarTier {
	return SoftDollarTier{}
}

func (s SoftDollarTier) String() string {
	return fmt.Sprintf("Name: %s, Value: %s, DisplayName: %s",
		s.Name,
		s.Value,
		s.DisplayName)
}
