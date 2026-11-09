package ibkr

import (
	"context"
	"net"
	"strings"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// SubscribeNews subscribes to live headlines for each topic and merges them into one stream.
func (c *Client) SubscribeNews(ctx context.Context, topics []models.NewsTopic) (*models.NewsResponse, error) {
	if len(topics) == 0 {
		return nil, errInvalidNewsTopics
	}
	for _, topic := range topics {
		if topic.Contract == nil {
			return nil, errInvalidContract
		}
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	streamResp := newNewsStreamResponse(c)
	resp := streamResp.resp
	resp.Cancel = newRequestCancelFunc(func() {
		streamResp.close(nil)
	})
	resp.Err = streamResp.err

	for _, topic := range topics {
		req, msgEnc := c.newNewsTopicRequest(streamResp, topic)
		if msgEnc.Err() != nil {
			streamResp.close(msgEnc.Err())
			return nil, msgEnc.Err()
		}

		if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
			streamResp.close(err)
			return nil, err
		}
		streamResp.addRequest(req)
	}

	if ctx != nil {
		go func() {
			select {
			case <-ctx.Done():
				resp.Cancel()
			case <-streamResp.closedCh:
			}
		}()
	}

	return resp, nil
}

// RequestNewsProviders retrieves the API news providers enabled for the account.
func (c *Client) RequestNewsProviders(ctx context.Context) ([]models.NewsProvider, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := make([]models.NewsProvider, 0)
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.NEWS_PROVIDERS),
		Response: &resp,
	})
	if c.isProtoBufAvailable(common.REQ_NEWS_PROVIDERS) {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_NEWS_PROVIDERS + common.PROTOBUF_MSG_ID).
			Proto(&protobuf.NewsProvidersRequest{})
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_NEWS_PROVIDERS).
			Int(1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// RequestNewsArticle retrieves an article body. Text is returned unchanged, including Base64 payloads.
func (c *Client) RequestNewsArticle(ctx context.Context, opts models.NewsArticleRequestOptions) (*models.NewsArticle, error) {
	var msgEnc *message.Encoder

	if len(opts.ProviderCode) == 0 || len(opts.ArticleID) == 0 {
		return nil, errInvalidNewsArticle
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.NewsArticle{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_NEWS_ARTICLE,
		Response: resp,
	})
	if c.isProtoBufAvailable(common.REQ_NEWS_ARTICLE) {
		pb := protobuf.NewsArticleRequest{
			ReqId:        protofmt.Int32(req.ID()),
			ProviderCode: protofmt.String(string(opts.ProviderCode)),
			ArticleId:    protofmt.String(opts.ArticleID),
		}
		if options := tagValueListToMap(opts.Options); len(options) > 0 {
			pb.NewsArticleOptions = options
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_NEWS_ARTICLE + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_NEWS_ARTICLE).
			Int(1).
			RequestID(req.ID()).
			String(string(opts.ProviderCode)).
			String(opts.ArticleID).
			Marshal(&opts.Options, 1)
	}
	if msgEnc.Err() != nil {
		return nil, msgEnc.Err()
	}

	if err := c.sendRequest(msgEnc.Bytes(), req); err != nil {
		return nil, err
	}
	defer c.reqMgr.removeRequest(req, context.Canceled)

	if err := c.waitRequestCompletion(ctx, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------

func (c *Client) newNewsTopicRequest(streamResp *newsStreamResponse, topic models.NewsTopic) (*Request, *message.Encoder) {
	providerCodes := make([]string, len(topic.ProviderCodes))
	for i, providerCode := range topic.ProviderCodes {
		providerCodes[i] = string(providerCode)
	}
	genericTicks := "mdoff,292"
	if len(providerCodes) > 0 {
		genericTicks += ":" + strings.Join(providerCodes, "+")
	}
	req := c.createRequest(RequestOptions{
		Type:    RequestTypeRequestWithTickerID,
		MsgCode: common.REQ_MKT_DATA,
		Response: &newsTopicStreamResponse{
			stream: streamResp,
			label:  topic.Label,
		},
		CompleteCB: func(_ *Request, err error) {
			if err != nil {
				streamResp.close(err)
			}
		},
	})

	if c.isProtoBufAvailable(common.REQ_MKT_DATA) {
		pb := protobuf.MarketDataRequest{
			ReqId:              protofmt.Int32(req.ID()),
			Contract:           topic.Contract.Proto(nil),
			GenericTickList:    protofmt.String(genericTicks),
			Snapshot:           protofmt.Bool(false),
			RegulatorySnapshot: protofmt.Bool(false),
		}
		return req, message.NewEncoder().
			RawUInt32(common.REQ_MKT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	}

	msgEnc := message.NewEncoder().Reserve(3).
		RawUInt32(common.REQ_MKT_DATA).
		Int(11).
		RequestID(req.ID()).
		Marshal(topic.Contract, 1)
	if topic.Contract.SecType == models.SecurityTypePair {
		msgEnc.Int(len(topic.Contract.ComboLegs))
		for _, comboLeg := range topic.Contract.ComboLegs {
			msgEnc.Marshal(comboLeg, 1)
		}
	}
	if topic.Contract.DeltaNeutralContract != nil {
		msgEnc.Bool(true).Marshal(topic.Contract.DeltaNeutralContract, 1)
	} else {
		msgEnc.Bool(false)
	}

	return req, msgEnc.String(genericTicks).
		Bool(false).
		Bool(false).
		Marshal(&models.TagValueList{}, 1)
}
