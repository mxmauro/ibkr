package models

import (
	"fmt"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/utils"
)

// -----------------------------------------------------------------------------

type DepthMktDataDescription struct {
	Exchange        string
	SecType         SecurityType
	ListingExchange string
	ServiceDataType string
	AggGroup        int64
}

// -----------------------------------------------------------------------------

func NewDepthMktDataDescription() DepthMktDataDescription {
	dmdd := DepthMktDataDescription{
		AggGroup: common.UNSET_INT,
	}
	return dmdd
}

func (d DepthMktDataDescription) String() string {
	return fmt.Sprintf("Exchange: %s, SecType: %s, ListingExchange: %s, ServiceDataType: %s, AggGroup: %s",
		d.Exchange, d.SecType, d.ListingExchange, d.ServiceDataType, utils.IntMaxString(d.AggGroup))
}
