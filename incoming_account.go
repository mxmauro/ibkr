package ibkr

import (
	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type updateAccountValueEventHandler interface {
	UpdateAccountValue(tag string, val string, currency string, accountName string)
}

type updatePortfolioEventHandler interface {
	UpdatePortfolio(contract *models.Contract, position models.Decimal, marketPrice float64, marketValue float64, averageCost float64, unrealizedPNL float64, realizedPNL float64, accountName string)
}

type updateAccountTimeEventHandler interface {
	UpdateAccountTime(timeStamp string)
}

type accountDownloadEndEventHandler interface {
	AccountDownloadEnd(accountName string)
}

type accountSummaryEventHandler interface {
	AccountSummary(reqID int64, account string, tag string, value string, currency string)
}

type accountSummaryEndEventHandler interface {
	AccountSummaryEnd(reqID int64)
}

// -----------------------------------------------------------------------------

func (c *Client) processAccountValueMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	data := models.AccountValueUpdate{
		Tag:         msgDec.String(),
		Value:       msgDec.String(),
		Currency:    msgDec.String(),
		AccountName: msgDec.String(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountValueMsg", msgDec.Err())
	}

	return c.processAccountValueCommon(data)
}

func (c *Client) processAccountValueProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.AccountValue{}
	msgDec.Unmarshal(&pb)
	data := models.AccountValueUpdate{
		Tag:         msgDec.String(pb.Key),
		Value:       msgDec.String(pb.Value),
		Currency:    msgDec.String(pb.Currency),
		AccountName: msgDec.String(pb.AccountName),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountValueProtobuf", msgDec.Err())
	}

	return c.processAccountValueCommon(data)
}

func (c *Client) processAccountValueCommon(data models.AccountValueUpdate) error {
	c.reqMgr.withActiveFirstRequestWithoutID(int(common.ACCT_VALUE), func(resp interface{}) error {
		streamResp := resp.(*accountUpdatesStreamResponse)
		streamResp.publish(data)
		return nil
	})

	if handler, ok := c.eventsHandler.(updateAccountValueEventHandler); ok {
		handler.UpdateAccountValue(data.Tag, data.Value, data.Currency, data.AccountName)
	}
	return nil
}

func (c *Client) processPortfolioValueMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	contract := decodeLegacyContract(msgDec, 1)
	data := models.PortfolioUpdate{
		Contract:      contract,
		Position:      models.NewDecimalFromMessageDecoder(msgDec),
		MarketPrice:   msgDec.Float(),
		MarketValue:   msgDec.Float(),
		AverageCost:   msgDec.Float(),
		UnrealizedPNL: msgDec.Float(),
		RealizedPNL:   msgDec.Float(),
		AccountName:   msgDec.String(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPortfolioValueMsg", msgDec.Err())
	}

	return c.processPortfolioValueCommon(data)
}

func (c *Client) processPortfolioValueProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.PortfolioValue{}
	msgDec.Unmarshal(&pb)
	data := models.PortfolioUpdate{
		Contract:      models.NewContractFromProtobufDecoder(msgDec, pb.Contract),
		Position:      models.NewDecimalFromProtobufDecoder(msgDec, pb.Position),
		MarketPrice:   msgDec.Float(pb.MarketPrice),
		MarketValue:   msgDec.Float(pb.MarketValue),
		AverageCost:   msgDec.Float(pb.AverageCost),
		UnrealizedPNL: msgDec.Float(pb.UnrealizedPNL),
		RealizedPNL:   msgDec.Float(pb.RealizedPNL),
		AccountName:   msgDec.String(pb.AccountName),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPortfolioValueProtobuf", msgDec.Err())
	}

	return c.processPortfolioValueCommon(data)
}

func (c *Client) processPortfolioValueCommon(data models.PortfolioUpdate) error {
	c.reqMgr.withActiveFirstRequestWithoutID(int(common.ACCT_VALUE), func(resp interface{}) error {
		streamResp := resp.(*accountUpdatesStreamResponse)
		streamResp.publish(data)
		return nil
	})

	if handler, ok := c.eventsHandler.(updatePortfolioEventHandler); ok {
		handler.UpdatePortfolio(
			data.Contract,
			data.Position,
			data.MarketPrice,
			data.MarketValue,
			data.AverageCost,
			data.UnrealizedPNL,
			data.RealizedPNL,
			data.AccountName,
		)
	}
	return nil
}

func (c *Client) processAccountUpdateTimeMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	data := models.AccountTimeUpdate{TimeStamp: msgDec.String()}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountUpdateTimeMsg", msgDec.Err())
	}

	return c.processAccountUpdateTimeCommon(data)
}

func (c *Client) processAccountUpdateTimeProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.AccountUpdateTime{}
	msgDec.Unmarshal(&pb)
	data := models.AccountTimeUpdate{TimeStamp: msgDec.String(pb.TimeStamp)}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountUpdateTimeProtobuf", msgDec.Err())
	}

	return c.processAccountUpdateTimeCommon(data)
}

func (c *Client) processAccountUpdateTimeCommon(data models.AccountTimeUpdate) error {
	c.reqMgr.withActiveFirstRequestWithoutID(int(common.ACCT_VALUE), func(resp interface{}) error {
		streamResp := resp.(*accountUpdatesStreamResponse)
		streamResp.publish(data)
		return nil
	})

	if handler, ok := c.eventsHandler.(updateAccountTimeEventHandler); ok {
		handler.UpdateAccountTime(data.TimeStamp)
	}
	return nil
}

func (c *Client) processAccountDownloadEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	data := models.AccountDownloadEndUpdate{AccountName: msgDec.String()}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountDownloadEndMsg", msgDec.Err())
	}

	return c.processAccountDownloadEndCommon(data)
}

func (c *Client) processAccountDownloadEndCommon(data models.AccountDownloadEndUpdate) error {
	c.reqMgr.withFirstRequestWithoutID(int(common.ACCT_VALUE), func(resp interface{}) (bool, error) {
		streamResp := resp.(*accountUpdatesStreamResponse)
		streamResp.publish(data)
		return true, nil
	})

	if handler, ok := c.eventsHandler.(accountDownloadEndEventHandler); ok {
		handler.AccountDownloadEnd(data.AccountName)
	}
	return nil
}

func (c *Client) processAccountSummaryMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	reqID := msgDec.RequestID(false)
	data := models.AccountSummaryValue{
		Account:  msgDec.String(),
		Tag:      msgDec.String(),
		Value:    msgDec.String(),
		Currency: msgDec.String(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountSummaryMsg", msgDec.Err())
	}

	return c.processAccountSummaryCommon(reqID, data)
}

func (c *Client) processAccountSummaryProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.AccountSummary{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	data := models.AccountSummaryValue{
		Account:  msgDec.String(pb.Account),
		Tag:      msgDec.String(pb.Tag),
		Value:    msgDec.String(pb.Value),
		Currency: msgDec.String(pb.Currency),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountSummaryProtobuf", msgDec.Err())
	}

	return c.processAccountSummaryCommon(reqID, data)
}

func (c *Client) processAccountSummaryCommon(reqID int32, data models.AccountSummaryValue) error {
	c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
		summaryResp := resp.(*models.AccountSummaryResponse)
		summaryResp.Values = append(summaryResp.Values, data)
		return false, nil
	})

	if handler, ok := c.eventsHandler.(accountSummaryEventHandler); ok {
		handler.AccountSummary(int64(reqID), data.Account, data.Tag, data.Value, data.Currency)
	}
	return nil
}

func (c *Client) processAccountSummaryEndMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	reqID := msgDec.RequestID(false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountSummaryEndMsg", msgDec.Err())
	}

	return c.processAccountSummaryEndCommon(reqID)
}

func (c *Client) processAccountSummaryEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.AccountSummaryEnd{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processAccountSummaryEndProtobuf", msgDec.Err())
	}

	return c.processAccountSummaryEndCommon(reqID)
}

func (c *Client) processAccountSummaryEndCommon(reqID int32) error {
	c.reqMgr.withRequestWithID(reqID, func(_ interface{}) (bool, error) {
		return true, nil
	})

	if handler, ok := c.eventsHandler.(accountSummaryEndEventHandler); ok {
		handler.AccountSummaryEnd(int64(reqID))
	}
	return nil
}
