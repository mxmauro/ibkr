package models

import (
	"github.com/mxmauro/ibkr/common"
)

// -----------------------------------------------------------------------------

type ExecutionFilter struct {
	ClientID      int64
	AcctCode      string
	Time          string
	Symbol        string
	SecType       SecurityType
	Exchange      string
	Side          string
	LastNDays     int64
	SpecificDates []int64
}

// -----------------------------------------------------------------------------

func NewExecutionFilter() *ExecutionFilter {
	ef := &ExecutionFilter{
		LastNDays: common.UNSET_INT,
	}
	return ef
}
