package models

import (
	"fmt"
)

// -----------------------------------------------------------------------------

type HistoricalSession struct {
	StartDateTime string
	EndDateTime   string
	RefDate       string
}

// -----------------------------------------------------------------------------

func NewHistoricalSession() HistoricalSession {
	return HistoricalSession{}
}

func (hs HistoricalSession) String() string {
	return fmt.Sprintf("Start: %s, End: %s, Ref Date: %s", hs.StartDateTime, hs.EndDateTime, hs.RefDate)
}
