package models

import (
	"fmt"
	"strconv"

	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

// Execution is the information of an order`s execution.
type Execution struct {
	ExecID               string
	Time                 string
	AcctNumber           string
	Exchange             string
	Side                 string
	Shares               Decimal //UNSET_DECIMAL
	Price                float64
	PermID               int64
	ClientID             int64
	OrderID              int64
	Liquidation          int64
	CumQty               Decimal // UNSET_DECIMAL
	AvgPrice             float64
	OrderRef             string
	EVRule               string
	EVMultiplier         float64
	ModelCode            string
	LastLiquidity        int64
	PendingPriceRevision bool
	Submitter            string
	OptExerciseOrLapse   OptionExercise
}

// -----------------------------------------------------------------------------

func NewExecution() *Execution {
	e := &Execution{
		Shares:             UNSET_DECIMAL,
		CumQty:             UNSET_DECIMAL,
		OptExerciseOrLapse: OptionExerciseNone,
	}
	return e
}

func (e Execution) String() string {
	return fmt.Sprintf("ExecId: %s, Time: %s, Account: %s, Exchange: %s, Side: %s, Shares: %s, Price: %s, PermId: %s, ClientId: %s, OrderId: %s, Liquidation: %s, CumQty: %s, AvgPrice: %s, OrderRef: %s, EvRule: %s, EvMultiplier: %s, ModelCode: %s, LastLiquidity: %s,  PendingPriceRevision: %s, Submitter: %s, OptionExerciseType: %s",
		e.ExecID, e.Time, e.AcctNumber, e.Exchange, e.Side, e.Shares.StringMax(), utils.FloatMaxString(e.Price),
		utils.IntMaxString(e.PermID), utils.IntMaxString(e.ClientID), utils.IntMaxString(e.OrderID),
		utils.IntMaxString(e.Liquidation), e.CumQty.StringMax(), utils.FloatMaxString(e.AvgPrice), e.OrderRef,
		e.EVRule, utils.FloatMaxString(e.EVMultiplier), e.ModelCode, utils.IntMaxString(e.LastLiquidity),
		strconv.FormatBool(e.PendingPriceRevision), e.Submitter, e.OptExerciseOrLapse.String())
}
