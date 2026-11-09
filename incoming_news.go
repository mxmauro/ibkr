package ibkr

import (
	"time"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func (c *Client) processTickNewsMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	timestamp := time.UnixMilli(msgDec.Int64()).UTC()
	headline := models.NewsHeadline{
		Timestamp:    timestamp,
		ProviderCode: models.NewsProviderCode(msgDec.String()),
		ArticleID:    msgDec.String(),
		Headline:     msgDec.String(),
		ExtraData:    msgDec.String(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickNewsMsg", msgDec.Err())
	}
	return c.processTickNewsCommon(reqID, headline)
}

func (c *Client) processTickNewsProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.TickNews{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	headline := models.NewsHeadline{
		Timestamp:    time.UnixMilli(msgDec.Int64(pb.Timestamp)).UTC(),
		ProviderCode: models.NewsProviderCode(msgDec.String(pb.ProviderCode)),
		ArticleID:    msgDec.String(pb.ArticleId),
		Headline:     msgDec.String(pb.Headline),
		ExtraData:    msgDec.String(pb.ExtraData),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processTickNewsProtobuf", msgDec.Err())
	}
	return c.processTickNewsCommon(reqID, headline)
}

func (c *Client) processTickNewsCommon(reqID int32, headline models.NewsHeadline) error {
	c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
		topicResp, ok := resp.(*newsTopicStreamResponse)
		if !ok {
			return false, nil
		}
		headline.TopicLabel = topicResp.label
		topicResp.stream.enqueue(headline)
		return false, nil
	})
	return nil
}

func (c *Client) processNewsProvidersMsg(msgDec *message.Decoder) error {
	count := int(msgDec.Int32())
	if count < 0 {
		msgDec.SetErr(errInvalidReceivedMessage)
		return wrapDecoderErr("processNewsProvidersMsg", msgDec.Err())
	}
	providers := make([]models.NewsProvider, 0, count)
	for i := 0; i < count; i++ {
		providers = append(providers, models.NewsProvider{
			Code: models.NewsProviderCode(msgDec.String()),
			Name: msgDec.String(),
		})
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processNewsProvidersMsg", msgDec.Err())
	}
	return c.processNewsProvidersCommon(providers)
}

func (c *Client) processNewsProvidersProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.NewsProviders{}
	msgDec.Unmarshal(&pb)
	providers := make([]models.NewsProvider, 0, len(pb.NewsProviders))
	for _, provider := range pb.NewsProviders {
		providers = append(providers, models.NewsProvider{
			Code: models.NewsProviderCode(msgDec.String(provider.ProviderCode)),
			Name: msgDec.String(provider.ProviderName),
		})
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processNewsProvidersProtobuf", msgDec.Err())
	}
	return c.processNewsProvidersCommon(providers)
}

func (c *Client) processNewsProvidersCommon(providers []models.NewsProvider) error {
	c.reqMgr.withRequestWithoutID(int(common.NEWS_PROVIDERS), func(resp interface{}) error {
		result := resp.(*[]models.NewsProvider)
		*result = append((*result)[:0], providers...)
		return nil
	})
	return nil
}

func (c *Client) processNewsArticleMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	article := models.NewsArticle{
		ArticleType: msgDec.Int32(),
		Text:        msgDec.String(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processNewsArticleMsg", msgDec.Err())
	}
	return c.processNewsArticleCommon(reqID, article)
}

func (c *Client) processNewsArticleProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.NewsArticle{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	article := models.NewsArticle{
		ArticleType: msgDec.Int32(pb.ArticleType),
		Text:        msgDec.String(pb.ArticleText),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processNewsArticleProtobuf", msgDec.Err())
	}
	return c.processNewsArticleCommon(reqID, article)
}

func (c *Client) processNewsArticleCommon(reqID int32, article models.NewsArticle) error {
	c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
		result := resp.(*models.NewsArticle)
		*result = article
		return true, nil
	})
	return nil
}
