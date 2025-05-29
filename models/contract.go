package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

// Contract describes an instrument's definition.
type Contract struct {
	ConID                        int64
	Symbol                       string
	SecType                      SecurityType
	LastTradeDateOrContractMonth string
	LastTradeDate                string
	Strike                       float64 // UNSET_FLOAT
	Right                        string  // Either Put or Call (i.e. Options). Valid values are P, PUT, C, CALL.
	Multiplier                   string
	Exchange                     string
	PrimaryExchange              string // pick an actual (ie non-aggregate) exchange that the contract trades on.  DO NOT SET TO SMART.
	Currency                     string
	LocalSymbol                  string
	TradingClass                 string
	IncludeExpired               bool
	SecIDType                    string // CUSIP;SEDOL;ISIN;RIC
	SecID                        string
	Description                  string
	IssuerID                     string

	// combo legs
	ComboLegsDescription string // received in open order 14 and up for all combos
	ComboLegs            []ComboLeg

	// delta neutral contract
	DeltaNeutralContract *DeltaNeutralContract
}

// -----------------------------------------------------------------------------

func NewContract() *Contract {
	return &Contract{
		Strike: common.UNSET_FLOAT,
	}
}

func (c *Contract) Equal(other *Contract) bool {
	if c.ConID != 0 && other.ConID != 0 {
		return c.ConID == other.ConID
	}
	if len(c.SecIDType) > 0 && len(other.SecIDType) > 0 && c.SecIDType == other.SecIDType {
		return c.SecID == other.SecID
	}
	return c.Symbol == other.Symbol &&
		c.SecType == other.SecType &&
		c.Exchange == other.Exchange &&
		c.Currency == other.Currency &&
		c.LastTradeDate == other.LastTradeDate &&
		c.Strike == other.Strike &&
		c.Right == other.Right
}

func (c *Contract) EncodeMessage() []byte {
	msgEnc := utils.NewRawMessageEncoder().
		Int64(c.ConID, false).
		String(c.Symbol).
		String(string(c.SecType)).
		String(c.LastTradeDateOrContractMonth).
		Float64(c.Strike, true).
		String(c.Right).
		String(c.Multiplier).
		String(c.Exchange).
		String(c.PrimaryExchange).
		String(c.Currency).
		String(c.LocalSymbol).
		String(c.TradingClass)
	return msgEnc.Bytes()
}

func (c *Contract) String() string {
	s := fmt.Sprintf("%d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %t, %s, %s, %s, %s",
		c.ConID,
		c.Symbol,
		c.SecType,
		c.LastTradeDateOrContractMonth,
		c.LastTradeDate,
		utils.FloatMaxString(c.Strike),
		c.Right,
		c.Multiplier,
		c.Exchange,
		c.PrimaryExchange,
		c.Currency,
		c.LocalSymbol,
		c.TradingClass,
		c.IncludeExpired,
		c.SecIDType,
		c.SecID,
		c.Description,
		c.IssuerID,
	)
	if len(c.ComboLegs) > 1 {
		s += ", combo:" + c.ComboLegsDescription
		for _, leg := range c.ComboLegs {
			s += fmt.Sprintf("; %s", leg)
		}
	}

	if c.DeltaNeutralContract != nil {
		s += fmt.Sprintf("; %s", c.DeltaNeutralContract)
	}

	return s
}
