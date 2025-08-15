package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/utils/formatter"
)

// -----------------------------------------------------------------------------

type DepthMktDataDescription struct {
	Exchange        string
	SecType         SecurityType
	ListingExchange string
	ServiceDataType string
	AggGroup        *int32
}

// -----------------------------------------------------------------------------

func NewDepthMktDataDescription() *DepthMktDataDescription {
	dmdd := DepthMktDataDescription{}
	return &dmdd
}

func (d *DepthMktDataDescription) String() string {
	return fmt.Sprintf(
		"Exchange: %s, SecType: %s, ListingExchange: %s, ServiceDataType: %s, AggGroup: %s",
		d.Exchange,
		d.SecType,
		d.ListingExchange,
		d.ServiceDataType,
		formatter.Int32MaxString(d.AggGroup),
	)
}
