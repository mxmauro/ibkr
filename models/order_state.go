package models

import (
	"fmt"
	"strings"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
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

	CommissionAndFees              float64 // UNSET_FLOAT
	MinCommissionAndFees           float64 // UNSET_FLOAT
	MaxCommissionAndFees           float64 // UNSET_FLOAT
	CommissionAndFeesCurrency      string
	MarginCurrency                 string
	InitMarginBeforeOutsideRTH     float64 // UNSET_FLOAT
	MaintMarginBeforeOutsideRTH    float64 // UNSET_FLOAT
	EquityWithLoanBeforeOutsideRTH float64 // UNSET_FLOAT
	InitMarginChangeOutsideRTH     float64 // UNSET_FLOAT
	MaintMarginChangeOutsideRTH    float64 // UNSET_FLOAT
	EquityWithLoanChangeOutsideRTH float64 // UNSET_FLOAT
	InitMarginAfterOutsideRTH      float64 // UNSET_FLOAT
	MaintMarginAfterOutsideRTH     float64 // UNSET_FLOAT
	EquityWithLoanAfterOutsideRTH  float64 // UNSET_FLOAT
	SuggestedSize                  Decimal // UNSET_DECIMAL
	RejectReason                   string
	OrderAllocations               []*OrderAllocation
	WarningText                    string

	CompletedTime   string
	CompletedStatus string
}

// -----------------------------------------------------------------------------

func NewOrderState() *OrderState {
	os := &OrderState{
		CommissionAndFees:              common.UNSET_FLOAT,
		MinCommissionAndFees:           common.UNSET_FLOAT,
		MaxCommissionAndFees:           common.UNSET_FLOAT,
		InitMarginBeforeOutsideRTH:     common.UNSET_FLOAT,
		MaintMarginBeforeOutsideRTH:    common.UNSET_FLOAT,
		EquityWithLoanBeforeOutsideRTH: common.UNSET_FLOAT,
		InitMarginChangeOutsideRTH:     common.UNSET_FLOAT,
		MaintMarginChangeOutsideRTH:    common.UNSET_FLOAT,
		EquityWithLoanChangeOutsideRTH: common.UNSET_FLOAT,
		InitMarginAfterOutsideRTH:      common.UNSET_FLOAT,
		MaintMarginAfterOutsideRTH:     common.UNSET_FLOAT,
		EquityWithLoanAfterOutsideRTH:  common.UNSET_FLOAT,
		SuggestedSize:                  UNSET_DECIMAL,
	}
	return os
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
		", CommissionAndFees: ", utils.FloatMaxString(os.CommissionAndFees),
		", MinCommissionAndFees: ", utils.FloatMaxString(os.MinCommissionAndFees),
		", MaxCommissionAndFees: ", utils.FloatMaxString(os.MaxCommissionAndFees),
		", CommissionAndFeesCurrency: ", os.CommissionAndFeesCurrency,
		", MarginCurrency: ", os.MarginCurrency,
		", InitMarginBeforeOutsideRTH: ", utils.FloatMaxString(os.InitMarginBeforeOutsideRTH),
		", MaintMarginBeforeOutsideRTH: ", utils.FloatMaxString(os.MaintMarginBeforeOutsideRTH),
		", EquityWithLoanBeforeOutsideRTH: ", utils.FloatMaxString(os.EquityWithLoanBeforeOutsideRTH),
		", InitMarginChangeOutsideRTH: ", utils.FloatMaxString(os.InitMarginChangeOutsideRTH),
		", MaintMarginChangeOutsideRTH: ", utils.FloatMaxString(os.MaintMarginChangeOutsideRTH),
		", EquityWithLoanChangeOutsideRTH: ", utils.FloatMaxString(os.EquityWithLoanChangeOutsideRTH),
		", InitMarginAfterOutsideRTH: ", utils.FloatMaxString(os.InitMarginAfterOutsideRTH),
		", MaintMarginAfterOutsideRTH: ", utils.FloatMaxString(os.MaintMarginAfterOutsideRTH),
		", EquityWithLoanAfterOutsideRTH: ", utils.FloatMaxString(os.EquityWithLoanAfterOutsideRTH),
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
