package models

import (
	"fmt"
	"strings"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type OrderState struct {
	Status string

	InitMarginBefore     string
	MaintMarginBefore    string
	EquityWithLoanBefore string
	InitMarginChange     string
	MaintMarginChange    string
	EquityWithLoanChange string
	InitMarginAfter      string
	MaintMarginAfter     string
	EquityWithLoanAfter  string

	CommissionAndFees              *float64
	MinCommissionAndFees           *float64
	MaxCommissionAndFees           *float64
	CommissionAndFeesCurrency      string
	MarginCurrency                 string
	InitMarginBeforeOutsideRTH     *float64
	MaintMarginBeforeOutsideRTH    *float64
	EquityWithLoanBeforeOutsideRTH *float64
	InitMarginChangeOutsideRTH     *float64
	MaintMarginChangeOutsideRTH    *float64
	EquityWithLoanChangeOutsideRTH *float64
	InitMarginAfterOutsideRTH      *float64
	MaintMarginAfterOutsideRTH     *float64
	EquityWithLoanAfterOutsideRTH  *float64
	SuggestedSize                  *Decimal
	RejectReason                   string
	OrderAllocations               []*OrderAllocation
	WarningText                    string

	CompletedTime   string
	CompletedStatus string
}

// -----------------------------------------------------------------------------

func NewOrderState() *OrderState {
	os := OrderState{}
	return &os
}

func (os *OrderState) String() string {
	s := fmt.Sprint(
		"Status: ", os.Status,
		", InitMarginBefore: ", os.InitMarginBefore,
		", MaintMarginBefore: ", os.MaintMarginBefore,
		", EquityWithLoanBefore: ", os.EquityWithLoanBefore,
		", InitMarginChange: ", os.InitMarginChange,
		", MaintMarginChange: ", os.MaintMarginChange,
		", EquityWithLoanChange: ", os.EquityWithLoanChange,
		", InitMarginAfter: ", os.InitMarginAfter,
		", MaintMarginAfter: ", os.MaintMarginAfter,
		", EquityWithLoanAfter: ", os.EquityWithLoanAfter,
		", CommissionAndFees: ", formatter.FloatMaxString(os.CommissionAndFees),
		", MinCommissionAndFees: ", formatter.FloatMaxString(os.MinCommissionAndFees),
		", MaxCommissionAndFees: ", formatter.FloatMaxString(os.MaxCommissionAndFees),
		", CommissionAndFeesCurrency: ", os.CommissionAndFeesCurrency,
		", MarginCurrency: ", os.MarginCurrency,
		", InitMarginBeforeOutsideRTH: ", formatter.FloatMaxString(os.InitMarginBeforeOutsideRTH),
		", MaintMarginBeforeOutsideRTH: ", formatter.FloatMaxString(os.MaintMarginBeforeOutsideRTH),
		", EquityWithLoanBeforeOutsideRTH: ", formatter.FloatMaxString(os.EquityWithLoanBeforeOutsideRTH),
		", InitMarginChangeOutsideRTH: ", formatter.FloatMaxString(os.InitMarginChangeOutsideRTH),
		", MaintMarginChangeOutsideRTH: ", formatter.FloatMaxString(os.MaintMarginChangeOutsideRTH),
		", EquityWithLoanChangeOutsideRTH: ", formatter.FloatMaxString(os.EquityWithLoanChangeOutsideRTH),
		", InitMarginAfterOutsideRTH: ", formatter.FloatMaxString(os.InitMarginAfterOutsideRTH),
		", MaintMarginAfterOutsideRTH: ", formatter.FloatMaxString(os.MaintMarginAfterOutsideRTH),
		", EquityWithLoanAfterOutsideRTH: ", formatter.FloatMaxString(os.EquityWithLoanAfterOutsideRTH),
		", SuggestedSize: ", os.SuggestedSize.StringMax(),
		", RejectReason: ", os.RejectReason,
		", WarningText: ", os.WarningText,
		", CompletedTime: ", os.CompletedTime,
		", CompletedStatus: ", os.CompletedStatus,
	)

	if os.OrderAllocations != nil {
		sb := strings.Builder{}
		_, _ = sb.WriteString(", OrderAllocations: [")
		for idx, oa := range os.OrderAllocations {
			if idx > 0 {
				_, _ = sb.WriteRune(',')
			}
			_, _ = sb.WriteString(oa.String())
		}
		_, _ = sb.WriteRune(']')
		s += sb.String()
	}

	return s
}
