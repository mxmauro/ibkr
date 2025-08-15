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

func NewSmartComponent() SmartComponent {
	return SmartComponent{}
}

func (sc *SmartComponent) String() string {
	return fmt.Sprintf(
		"BitNumber: %d, Exchange: %s, ExchangeLetter: %s",
		sc.BitNumber,
		sc.Exchange,
		sc.ExchangeLetter,
	)
}
