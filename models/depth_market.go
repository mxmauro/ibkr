package models

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------------

type MarketDepthData interface {
	UpdateBook(book *MarketDepthBook)
}

type MarketDepthBook struct {
	Size int
	Bids []DepthMarketBookEntry
	Asks []DepthMarketBookEntry
}

type DepthMarketBookEntry struct {
	Price float64
	Size  Decimal
}

// -----------------------------------------------------------------------------

func NewMarketDepthBook(maxDepth int) MarketDepthBook {
	return MarketDepthBook{
		Size: maxDepth,
		Bids: make([]DepthMarketBookEntry, 0, maxDepth),
		Asks: make([]DepthMarketBookEntry, 0, maxDepth),
	}
}

func (book *MarketDepthBook) TableString() string {
	bidPriceStrings := make([]string, len(book.Bids))
	bidSizeStrings := make([]string, len(book.Bids))
	askPriceStrings := make([]string, len(book.Bids))
	askSizeStrings := make([]string, len(book.Bids))
	longestPriceString := 0
	longestSizeString := 0

	for idx, entry := range book.Bids {
		bidPriceStrings[idx] = fmt.Sprintf("%.2f", entry.Price)
		if len(bidPriceStrings[idx]) > longestPriceString {
			longestPriceString = len(bidPriceStrings[idx])
		}
		bidSizeStrings[idx] = entry.Size.StringMax()
		if len(bidSizeStrings[idx]) > longestSizeString {
			longestSizeString = len(bidSizeStrings[idx])
		}
	}

	for idx, entry := range book.Asks {
		askPriceStrings[idx] = fmt.Sprintf("%.2f", entry.Price)
		if len(askPriceStrings[idx]) > longestPriceString {
			longestPriceString = len(askPriceStrings[idx])
		}
		askSizeStrings[idx] = entry.Size.StringMax()
		if len(askSizeStrings[idx]) > longestSizeString {
			longestSizeString = len(askSizeStrings[idx])
		}
	}

	sb := strings.Builder{}
	n := longestPriceString + 3 + longestSizeString
	_, _ = sb.WriteString(strings.Repeat(" ", (n-3)/2+1))
	_, _ = sb.WriteString("BID")
	_, _ = sb.WriteString(strings.Repeat(" ", n-3-(n-3)/2))
	_, _ = sb.WriteString(" | ")
	_, _ = sb.WriteString(strings.Repeat(" ", (n-3)/2+1))
	_, _ = sb.WriteString("ASK")
	_, _ = sb.WriteString(strings.Repeat(" ", n-3-(n-3)/2))
	_, _ = sb.WriteString("\n")

	_, _ = sb.WriteString(strings.Repeat("-", n+1))
	_, _ = sb.WriteString("-+-")
	_, _ = sb.WriteString(strings.Repeat("-", n+1))
	_, _ = sb.WriteString("\n")

	c := len(book.Bids)
	if len(book.Asks) > c {
		c = len(book.Asks)
	}
	for idx := 0; idx < c; idx++ {
		if idx < len(book.Bids) {
			_, _ = sb.WriteString(" ")
			_, _ = sb.WriteString(strings.Repeat(" ", longestPriceString-len(bidPriceStrings[idx])))
			_, _ = sb.WriteString(bidPriceStrings[idx])
			_, _ = sb.WriteString(" (")
			_, _ = sb.WriteString(bidSizeStrings[idx])
			_, _ = sb.WriteString(")")
			_, _ = sb.WriteString(strings.Repeat(" ", longestSizeString-len(bidSizeStrings[idx])))
		} else {
			_, _ = sb.WriteString(strings.Repeat(" ", n+1))
		}
		_, _ = sb.WriteString(" | ")
		if idx < len(book.Asks) {
			_, _ = sb.WriteString(strings.Repeat(" ", longestPriceString-len(askPriceStrings[idx])))
			_, _ = sb.WriteString(askPriceStrings[idx])
			_, _ = sb.WriteString(" (")
			_, _ = sb.WriteString(askSizeStrings[idx])
			_, _ = sb.WriteString(")")
		}

		if idx < c {
			_, _ = sb.WriteString("\n")
		}
	}

	return sb.String()
}
