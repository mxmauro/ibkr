package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type SmartComponent struct {
	BitNumber      int64
	Exchange       string
	ExchangeLetter string
}

// -----------------------------------------------------------------------------

func newSmartComponent() SmartComponent {
	return SmartComponent{}
}

func (s SmartComponent) String() string {
	return fmt.Sprintf("BitNumber: %d, Exchange: %s, ExchangeLetter: %s", s.BitNumber, s.Exchange, s.ExchangeLetter)
}
