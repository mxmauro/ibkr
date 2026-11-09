package models

import "time"

// -----------------------------------------------------------------------------

const (
	NewsProviderBriefingTrader NewsProviderCode = "BRF"
	NewsProviderBenzinga       NewsProviderCode = "BZ"
	NewsProviderFlyOnTheWall   NewsProviderCode = "FLY"
)

// -----------------------------------------------------------------------------

// NewsProviderCode identifies an API news source.
type NewsProviderCode string

// NewsTopic subscribes to news for one exact contract. Label is copied to each delivered headline.
type NewsTopic struct {
	Label         string
	Contract      *Contract
	ProviderCodes []NewsProviderCode
}

// NewsHeadline is one live news headline.
type NewsHeadline struct {
	TopicLabel   string
	Timestamp    time.Time
	ProviderCode NewsProviderCode
	ArticleID    string
	Headline     string
	ExtraData    string
}

// NewsResponse exposes a merged live-news stream and controls.
type NewsResponse struct {
	Channel chan NewsHeadline
	Cancel  CancelFunc
	Err     ErrFunc
}

// NewsProvider describes an API news provider enabled for the account.
type NewsProvider struct {
	Code NewsProviderCode
	Name string
}

// NewsArticleRequestOptions configures a news article lookup.
type NewsArticleRequestOptions struct {
	ProviderCode NewsProviderCode
	ArticleID    string
	Options      TagValueList
}

// NewsArticle contains an article body exactly as returned by TWS.
// ArticleType is 0 for text or HTML and 1 for Base64-encoded binary data.
type NewsArticle struct {
	ArticleType int32
	Text        string
}

// -----------------------------------------------------------------------------

// NewBriefingTraderBroadTapeNewsTopic creates a Briefing Trader broad-tape topic.
func NewBriefingTraderBroadTapeNewsTopic(label string) NewsTopic {
	return NewBroadTapeNewsTopic(label, NewsProviderBriefingTrader)
}

// NewBenzingaBroadTapeNewsTopic creates a Benzinga broad-tape topic.
func NewBenzingaBroadTapeNewsTopic(label string) NewsTopic {
	return NewBroadTapeNewsTopic(label, NewsProviderBenzinga)
}

// NewFlyOnTheWallBroadTapeNewsTopic creates a Fly on the Wall broad-tape topic.
func NewFlyOnTheWallBroadTapeNewsTopic(label string) NewsTopic {
	return NewBroadTapeNewsTopic(label, NewsProviderFlyOnTheWall)
}

// NewBroadTapeNewsTopic creates a broad-tape topic for a provider.
func NewBroadTapeNewsTopic(label string, provider NewsProviderCode) NewsTopic {
	return NewsTopic{Label: label, Contract: NewBroadTapeNewsContract(provider)}
}

// NewBroadTapeNewsContract creates the NEWS contract used for a provider's broad tape.
func NewBroadTapeNewsContract(provider NewsProviderCode) *Contract {
	code := string(provider)
	return &Contract{
		Symbol:   code + ":" + code + "_ALL",
		SecType:  SecurityTypeNews,
		Exchange: code,
	}
}
