package ibkr

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

func TestSubscribeNewsLegacyRoutesTopicsAndCancels(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufMarketData-1)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.NewsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.SubscribeNews(context.Background(), []models.NewsTopic{
			{Label: "first", Contract: testRequestContract(), ProviderCodes: []models.NewsProviderCode{models.NewsProviderBenzinga}},
			{Label: "second", Contract: testRequestContract()},
		})
		resultCh <- struct {
			resp *models.NewsResponse
			err  error
		}{resp, err}
	}()

	firstID := assertTextNewsMarketDataRequest(t, readFramedPayload(t, serverConn), "mdoff,292:BZ")
	secondID := assertTextNewsMarketDataRequest(t, readFramedPayload(t, serverConn), "mdoff,292")
	result := waitNewsSubscription(t, resultCh)
	if result.err != nil {
		t.Fatalf("SubscribeNews failed: %v", result.err)
	}

	if err := client.processIncomingMessage(buildTextTickNewsMessage(secondID, 1700000000123, "FLY", "article-2", "second headline", "extra")); err != nil {
		t.Fatalf("process second headline: %v", err)
	}
	if err := client.processIncomingMessage(buildTextTickNewsMessage(firstID, 1700000000456, "BZ", "article-1", "first headline", "")); err != nil {
		t.Fatalf("process first headline: %v", err)
	}
	assertNewsHeadline(t, waitNewsHeadline(t, result.resp.Channel), "second", 1700000000123, "FLY", "article-2", "second headline", "extra")
	assertNewsHeadline(t, waitNewsHeadline(t, result.resp.Channel), "first", 1700000000456, "BZ", "article-1", "first headline", "")

	result.resp.Cancel()
	assertTextCancelMarketDataRequest(t, readFramedPayload(t, serverConn), firstID)
	assertTextCancelMarketDataRequest(t, readFramedPayload(t, serverConn), secondID)
	waitNewsClosed(t, result.resp.Channel)
}

func TestSubscribeNewsProtobufAndChildError(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufMarketData)
	defer cleanup()

	resultCh := make(chan struct {
		resp *models.NewsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.SubscribeNews(context.Background(), []models.NewsTopic{{Label: "proto", Contract: testRequestContract()}})
		resultCh <- struct {
			resp *models.NewsResponse
			err  error
		}{resp, err}
	}()

	payload := readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_MKT_DATA+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected protobuf request id: %d", got)
	}
	pb := protobuf.MarketDataRequest{}
	if err := proto.Unmarshal(payload[4:], &pb); err != nil {
		t.Fatalf("unmarshal market data request: %v", err)
	}
	if pb.GetGenericTickList() != "mdoff,292" || pb.GetReqId() <= 0 {
		t.Fatalf("unexpected protobuf news request: %#v", pb)
	}
	result := waitNewsSubscription(t, resultCh)
	if result.err != nil {
		t.Fatalf("SubscribeNews failed: %v", result.err)
	}

	if err := client.processIncomingMessage(buildProtoMessage(common.TICK_NEWS, &protobuf.TickNews{
		ReqId: protofmt.Int32(pb.GetReqId()), Timestamp: protofmt.Int64(1700000000789), ProviderCode: protofmt.String("BZ"),
		ArticleId: protofmt.String("article"), Headline: protofmt.String("headline"), ExtraData: protofmt.String("extra"),
	})); err != nil {
		t.Fatalf("process protobuf headline: %v", err)
	}
	assertNewsHeadline(t, waitNewsHeadline(t, result.resp.Channel), "proto", 1700000000789, "BZ", "article", "headline", "extra")

	if err := client.processIncomingMessage(buildTextErrorMessage(pb.GetReqId(), 321, "news failed")); err != nil {
		t.Fatalf("process child error: %v", err)
	}
	waitNewsClosed(t, result.resp.Channel)
	if result.resp.Err() == nil {
		t.Fatal("stream error was not retained")
	}
	assertProtobufCancelMarketDataRequest(t, readFramedPayload(t, serverConn), pb.GetReqId())
}

func TestSubscribeNewsContextCancellation(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufMarketData-1)
	defer cleanup()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan struct {
		resp *models.NewsResponse
		err  error
	}, 1)
	go func() {
		resp, err := client.SubscribeNews(ctx, []models.NewsTopic{{Contract: testRequestContract()}})
		resultCh <- struct {
			resp *models.NewsResponse
			err  error
		}{resp, err}
	}()
	reqID := assertTextNewsMarketDataRequest(t, readFramedPayload(t, serverConn), "mdoff,292")
	result := waitNewsSubscription(t, resultCh)
	if result.err != nil {
		t.Fatalf("SubscribeNews failed: %v", result.err)
	}

	cancel()
	assertTextCancelMarketDataRequest(t, readFramedPayload(t, serverConn), reqID)
	waitNewsClosed(t, result.resp.Channel)
}

func TestRequestNewsProvidersAndArticle(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufNewsData-1)
	defer cleanup()

	providersCh := make(chan struct {
		providers []models.NewsProvider
		err       error
	}, 1)
	go func() {
		providers, err := client.RequestNewsProviders(context.Background())
		providersCh <- struct {
			providers []models.NewsProvider
			err       error
		}{providers, err}
	}()
	payload := readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_NEWS_PROVIDERS {
		t.Fatalf("unexpected providers request: %d", got)
	}
	if err := client.processIncomingMessage(message.NewRawEncoder().RawUInt32(common.NEWS_PROVIDERS).Int(2).String("BZ").String("Benzinga").String("FLY").String("Fly").Bytes()); err != nil {
		t.Fatalf("process providers: %v", err)
	}
	providersResult := waitNewsProviders(t, providersCh)
	if providersResult.err != nil || len(providersResult.providers) != 2 || providersResult.providers[0].Code != "BZ" {
		t.Fatalf("unexpected providers result: %#v, %v", providersResult.providers, providersResult.err)
	}

	articleCh := make(chan struct {
		article *models.NewsArticle
		err     error
	}, 1)
	go func() {
		article, err := client.RequestNewsArticle(context.Background(), models.NewsArticleRequestOptions{
			ProviderCode: models.NewsProviderBenzinga, ArticleID: "article", Options: models.TagValueList{{Tag: "x", Value: "y"}},
		})
		articleCh <- struct {
			article *models.NewsArticle
			err     error
		}{article, err}
	}()
	payload = readFramedPayload(t, serverConn)
	msgDec := message.NewDecoder(payload[4:])
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_NEWS_ARTICLE || msgDec.Int32() != 1 {
		t.Fatalf("unexpected article request")
	}
	reqID := msgDec.RequestID(false)
	if msgDec.String() != "BZ" || msgDec.String() != "article" || msgDec.String() != "x=y;" || msgDec.Err() != nil {
		t.Fatalf("unexpected article request fields: %v", msgDec.Err())
	}
	if err := client.processIncomingMessage(message.NewRawEncoder().RawUInt32(common.NEWS_ARTICLE).RequestID(reqID).Int32(1).String("SGVsbG8=").Bytes()); err != nil {
		t.Fatalf("process article: %v", err)
	}
	articleResult := waitNewsArticle(t, articleCh)
	if articleResult.err != nil || articleResult.article.ArticleType != 1 || articleResult.article.Text != "SGVsbG8=" {
		t.Fatalf("unexpected article result: %#v, %v", articleResult.article, articleResult.err)
	}
}

func TestNewsProtobufRequestsAndSDK1052Models(t *testing.T) {
	client, serverConn, cleanup := newRequestTestClientWithVersion(t, common.ServerVersionProtobufNewsData)
	defer cleanup()

	providersCh := make(chan error, 1)
	go func() {
		_, err := client.RequestNewsProviders(context.Background())
		providersCh <- err
	}()
	payload := readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_NEWS_PROVIDERS+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected protobuf providers request: %d", got)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.NEWS_PROVIDERS, &protobuf.NewsProviders{NewsProviders: []*protobuf.NewsProvider{{ProviderCode: protofmt.String("BRF"), ProviderName: protofmt.String("Briefing")}}})); err != nil {
		t.Fatalf("process protobuf providers: %v", err)
	}
	if err := <-providersCh; err != nil {
		t.Fatalf("RequestNewsProviders failed: %v", err)
	}

	articleCh := make(chan struct {
		article *models.NewsArticle
		err     error
	}, 1)
	go func() {
		article, err := client.RequestNewsArticle(context.Background(), models.NewsArticleRequestOptions{
			ProviderCode: models.NewsProviderBriefingTrader, ArticleID: "article",
		})
		articleCh <- struct {
			article *models.NewsArticle
			err     error
		}{article, err}
	}()
	payload = readFramedPayload(t, serverConn)
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_NEWS_ARTICLE+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected protobuf article request: %d", got)
	}
	articleRequest := protobuf.NewsArticleRequest{}
	if err := proto.Unmarshal(payload[4:], &articleRequest); err != nil || articleRequest.GetProviderCode() != "BRF" || articleRequest.GetArticleId() != "article" {
		t.Fatalf("unexpected protobuf article request: %#v, %v", articleRequest, err)
	}
	if err := client.processIncomingMessage(buildProtoMessage(common.NEWS_ARTICLE, &protobuf.NewsArticle{
		ReqId: articleRequest.ReqId, ArticleType: protofmt.Int32(0), ArticleText: protofmt.String("<p>article</p>"),
	})); err != nil {
		t.Fatalf("process protobuf article: %v", err)
	}
	articleResult := waitNewsArticle(t, articleCh)
	if articleResult.err != nil || articleResult.article.ArticleType != 0 || articleResult.article.Text != "<p>article</p>" {
		t.Fatalf("unexpected protobuf article result: %#v, %v", articleResult.article, articleResult.err)
	}

	order := models.NewOrder()
	order.ConditionsIncludeOvernight = true
	if !protofmt.NewDecoder(nil).Bool(order.Proto().ConditionsIncludeOvernight) {
		t.Fatal("protobuf order omitted conditions include overnight")
	}
	if !models.NewOrderFromProtobufDecoder(protofmt.NewDecoder(nil), order.Proto()).ConditionsIncludeOvernight {
		t.Fatal("protobuf order did not round trip conditions include overnight")
	}
	order.Action = "BUY"
	order.OrderType = "MKT"
	order.TotalQuantity = &models.DecimalOne
	legacyClient := &Client{serverVersion: common.ServerVersionConditionsIncludeOvernight - 1}
	if _, err := legacyClient.PlaceOrder(models.PlaceOrderRequestOptions{Contract: testRequestContract(), Order: order}); !errors.Is(err, errConditionsIncludeOvernightUnsupported) {
		t.Fatalf("unexpected old-server validation error: %v", err)
	}
	settlementMethod := "Cash"
	details := models.NewContractDetailsFromProtobufDecoder(protofmt.NewDecoder(nil), &protobuf.ContractDetails{SettlementMethod: &settlementMethod}, nil)
	if details.SettlementMethod != settlementMethod {
		t.Fatalf("unexpected settlement method: %q", details.SettlementMethod)
	}
	assertBroadTapeNewsContract(t, models.NewBriefingTraderBroadTapeNewsTopic("briefing").Contract, "BRF")
	assertBroadTapeNewsContract(t, models.NewBenzingaBroadTapeNewsTopic("benzinga").Contract, "BZ")
	assertBroadTapeNewsContract(t, models.NewFlyOnTheWallBroadTapeNewsTopic("fly").Contract, "FLY")
}

// -----------------------------------------------------------------------------

func buildTextTickNewsMessage(reqID int32, timestamp int64, providerCode, articleID, headline, extraData string) []byte {
	return message.NewRawEncoder().RawUInt32(common.TICK_NEWS).RequestID(reqID).Int64(timestamp).String(providerCode).String(articleID).String(headline).String(extraData).Bytes()
}

func buildTextErrorMessage(reqID int32, code int, text string) []byte {
	return message.NewRawEncoder().RawUInt32(common.ERR_MSG).RequestID(reqID).Int(code).String(text).String("").String("").Bytes()
}

func assertTextNewsMarketDataRequest(t *testing.T, payload []byte, genericTicks string) int32 {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.REQ_MKT_DATA {
		t.Fatalf("unexpected request id: %d", got)
	}
	msgDec := message.NewDecoder(payload[4:])
	if msgDec.Int32() != 11 {
		t.Fatal("unexpected market data version")
	}
	reqID := msgDec.RequestID(false)
	assertDecodedContract(t, msgDec, testRequestContract())
	if msgDec.Bool() || msgDec.String() != genericTicks || msgDec.Bool() || msgDec.Bool() || msgDec.Int32() != 0 || msgDec.Err() != nil {
		t.Fatalf("unexpected news market data request: %v", msgDec.Err())
	}
	return reqID
}

func assertTextCancelMarketDataRequest(t *testing.T, payload []byte, reqID int32) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.CANCEL_MKT_DATA {
		t.Fatalf("unexpected cancel id: %d", got)
	}
	msgDec := message.NewDecoder(payload[4:])
	if msgDec.Int32() != 2 || msgDec.RequestID(false) != reqID || msgDec.Err() != nil {
		t.Fatalf("unexpected cancel request: %v", msgDec.Err())
	}
}

func assertProtobufCancelMarketDataRequest(t *testing.T, payload []byte, reqID int32) {
	t.Helper()
	if got := binary.BigEndian.Uint32(payload[:4]); got != common.CANCEL_MKT_DATA+common.PROTOBUF_MSG_ID {
		t.Fatalf("unexpected protobuf cancel id: %d", got)
	}
	pb := protobuf.CancelMarketData{}
	if err := proto.Unmarshal(payload[4:], &pb); err != nil || pb.GetReqId() != reqID {
		t.Fatalf("unexpected protobuf cancel request: %#v, %v", pb, err)
	}
}

func waitNewsSubscription(t *testing.T, ch <-chan struct {
	resp *models.NewsResponse
	err  error
}) struct {
	resp *models.NewsResponse
	err  error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for news subscription")
		return struct {
			resp *models.NewsResponse
			err  error
		}{}
	}
}

func waitNewsProviders(t *testing.T, ch <-chan struct {
	providers []models.NewsProvider
	err       error
}) struct {
	providers []models.NewsProvider
	err       error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for providers")
		return struct {
			providers []models.NewsProvider
			err       error
		}{}
	}
}

func waitNewsArticle(t *testing.T, ch <-chan struct {
	article *models.NewsArticle
	err     error
}) struct {
	article *models.NewsArticle
	err     error
} {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for article")
		return struct {
			article *models.NewsArticle
			err     error
		}{}
	}
}

func waitNewsHeadline(t *testing.T, ch <-chan models.NewsHeadline) models.NewsHeadline {
	t.Helper()
	select {
	case headline := <-ch:
		return headline
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for headline")
		return models.NewsHeadline{}
	}
}

func waitNewsClosed(t *testing.T, ch <-chan models.NewsHeadline) {
	t.Helper()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("news stream did not close")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for news stream to close")
	}
}

func assertNewsHeadline(t *testing.T, got models.NewsHeadline, label string, timestamp int64, providerCode, articleID, headline, extraData string) {
	t.Helper()
	if got.TopicLabel != label || !got.Timestamp.Equal(time.UnixMilli(timestamp).UTC()) || string(got.ProviderCode) != providerCode || got.ArticleID != articleID || got.Headline != headline || got.ExtraData != extraData {
		t.Fatalf("unexpected headline: %#v", got)
	}
}

func assertBroadTapeNewsContract(t *testing.T, contract *models.Contract, provider string) {
	t.Helper()
	if contract.Symbol != provider+":"+provider+"_ALL" || contract.SecType != models.SecurityTypeNews || contract.Exchange != provider {
		t.Fatalf("unexpected broad-tape contract: %#v", contract)
	}
}
