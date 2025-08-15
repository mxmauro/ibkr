package models

import (
	"fmt"
	"strconv"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

// Execution is the information of an order`s execution.
type Execution struct {
	ExecID               string
	Time                 string
	AcctNumber           string
	Exchange             string
	Side                 string
	Shares               *Decimal
	Price                float64
	PermID               int64
	ClientID             int32
	OrderID              int32
	Liquidation          int32
	CumQty               *Decimal
	AvgPrice             float64
	OrderRef             string
	EVRule               string
	EVMultiplier         float64
	ModelCode            string
	LastLiquidity        Liquidities
	PendingPriceRevision bool
	Submitter            string
	OptExerciseOrLapse   OptionExercise
}

// -----------------------------------------------------------------------------

func NewExecution() *Execution {
	e := Execution{
		OptExerciseOrLapse: OptionExerciseNone,
	}
	return &e
}

func (e *Execution) String() string {
	return fmt.Sprintf(
		"ExecId: %s, Time: %s, Account: %s, Exchange: %s, Side: %s, Shares: %s, Price: %s, PermId: %s, ClientId: %s, OrderId: %s, Liquidation: %s, CumQty: %s, AvgPrice: %s, OrderRef: %s, EvRule: %s, EvMultiplier: %s, ModelCode: %s, LastLiquidity: %s,  PendingPriceRevision: %s, Submitter: %s, OptionExerciseType: %s",
		e.ExecID,
		e.Time,
		e.AcctNumber,
		e.Exchange,
		e.Side,
		e.Shares.StringMax(),
		formatter.FloatString(e.Price),
		formatter.Int64String(e.PermID),
		formatter.Int32String(e.ClientID),
		formatter.Int32String(e.OrderID),
		formatter.Int32String(e.Liquidation),
		e.CumQty.StringMax(),
		formatter.FloatString(e.AvgPrice),
		e.OrderRef,
		e.EVRule,
		formatter.FloatString(e.EVMultiplier),
		e.ModelCode, e.LastLiquidity.String(),
		strconv.FormatBool(e.PendingPriceRevision),
		e.Submitter,
		e.OptExerciseOrLapse.String(),
	)
}
