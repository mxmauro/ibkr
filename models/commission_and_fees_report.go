package models

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------

type CommissionAndFeesReport struct {
	ExecID              string
	CommissionAndFees   float64
	Currency            string
	RealizedPNL         float64
	Yield               float64
	YieldRedemptionDate time.Time // YYYYMMDD format
}

// -----------------------------------------------------------------------------

func NewCommissionAndFeesReport() *CommissionAndFeesReport {
	return &CommissionAndFeesReport{}
}

func (cr *CommissionAndFeesReport) String() string {
	return fmt.Sprintf(
		"ExecId: %s, CommissionAndFees: %f, Currency: %s, RealizedPnL: %f, Yield: %f, YieldRedemptionDate: %s",
		cr.ExecID,
		cr.CommissionAndFees,
		cr.Currency,
		cr.RealizedPNL,
		cr.Yield,
		cr.YieldRedemptionDate.Format("2006-01-02"),
	)
}
