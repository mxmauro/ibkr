package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type ComboLeg struct {
	ConID     int64
	Ratio     int64
	Action    string // BUY/SELL/SSHORT
	Exchange  string
	OpenClose int64
	// for stock legs when doing short sale
	ShortSaleSlot      int64 // 1 = clearing broker, 2 = third party
	DesignatedLocation string
	ExemptCode         int64
}

// -----------------------------------------------------------------------------

// NewComboLeg creates a default ComboLeg.
func NewComboLeg() ComboLeg {
	cl := ComboLeg{}
	cl.ExemptCode = -1
	return cl
}

func (c ComboLeg) String() string {
	return fmt.Sprintf("%d, %d, %s, %s, %d, %d, %s, %d",
		c.ConID, c.Ratio, c.Action, c.Exchange, c.OpenClose, c.ShortSaleSlot, c.DesignatedLocation, c.ExemptCode)
}
