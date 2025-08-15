package models

import (
	"time"
)

// -----------------------------------------------------------------------------

type ExecutionFilter struct {
	ClientID      int32
	AcctCode      string
	Time          string
	Symbol        string
	SecType       SecurityType
	Exchange      string
	Side          string
	LastNDays     *int32
	SpecificDates []time.Time
}

// -----------------------------------------------------------------------------

func NewExecutionFilter() *ExecutionFilter {
	ef := ExecutionFilter{}
	return &ef
}
