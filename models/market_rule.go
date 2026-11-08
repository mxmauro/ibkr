package models

// -----------------------------------------------------------------------------

// MarketRuleRequestOptions identifies the market rule to retrieve.
type MarketRuleRequestOptions struct {
	MarketRuleID int64
}

// MarketRuleResponse contains the decoded price increments for a market rule.
type MarketRuleResponse struct {
	MarketRuleID    int64
	PriceIncrements []PriceIncrement
}
