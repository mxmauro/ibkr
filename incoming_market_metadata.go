package ibkr

import (
	"fmt"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type mktDepthExchangesEventHandler interface {
	MktDepthExchanges(depthMktDataDescriptions []models.DepthMktDataDescription)
}

// -----------------------------------------------------------------------------

func (c *Client) processHistogramDataMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	count := int(msgDec.Int64())
	if count < 0 {
		msgDec.SetErr(fmt.Errorf("negative histogram size: %d", count))
		return wrapDecoderErr("processHistogramDataMsg", msgDec.Err())
	}

	data := make([]models.HistogramData, 0, count)
	for idx := 0; idx < count; idx++ {
		item := models.NewHistogramData()
		item.Price = msgDec.Float()
		item.Size = models.NewDecimalMaxFromMessageDecoder(msgDec)
		data = append(data, item)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistogramDataMsg", msgDec.Err())
	}

	return c.processHistogramDataCommon(reqID, data)
}

func (c *Client) processHistogramDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.HistogramData{}
	msgDec.Unmarshal(&pb)
	reqID := msgDec.RequestID(pb.ReqId, false)
	data := make([]models.HistogramData, 0, len(pb.HistogramDataEntries))
	for _, entry := range pb.HistogramDataEntries {
		item := models.NewHistogramData()
		item.Price = msgDec.Float(entry.Price)
		item.Size = models.NewDecimalMaxFromProtobufDecoder(msgDec, entry.Size)
		data = append(data, item)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processHistogramDataProtobuf", msgDec.Err())
	}

	return c.processHistogramDataCommon(reqID, data)
}

func (c *Client) processHistogramDataCommon(reqID int32, data []models.HistogramData) error {
	c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
		histResp := resp.(*models.HistogramDataResponse)
		histResp.Data = append(histResp.Data[:0], data...)
		return true, nil
	})
	return nil
}

func (c *Client) processMktDepthExchangesMsg(msgDec *message.Decoder) error {
	count := int(msgDec.Int64())
	if count < 0 {
		msgDec.SetErr(fmt.Errorf("negative depth exchanges size: %d", count))
		return wrapDecoderErr("processMktDepthExchangesMsg", msgDec.Err())
	}

	data := make([]models.DepthMktDataDescription, 0, count)
	for idx := 0; idx < count; idx++ {
		entry := models.DepthMktDataDescription{
			Exchange:        msgDec.String(),
			SecType:         models.SecurityType(msgDec.String()),
			ListingExchange: msgDec.String(),
			ServiceDataType: msgDec.String(),
			AggGroup:        msgDec.Int32Max(),
		}
		data = append(data, entry)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processMktDepthExchangesMsg", msgDec.Err())
	}

	return c.processMktDepthExchangesCommon(data)
}

func (c *Client) processMktDepthExchangesCommon(data []models.DepthMktDataDescription) error {
	c.reqMgr.withRequestWithoutID(int(common.MKT_DEPTH_EXCHANGES), func(resp interface{}) error {
		depthResp := resp.(*models.MktDepthExchangesResponse)
		depthResp.Exchanges = append(depthResp.Exchanges[:0], data...)
		return nil
	})

	if handler, ok := c.eventsHandler.(mktDepthExchangesEventHandler); ok {
		handler.MktDepthExchanges(data)
	}
	return nil
}

func (c *Client) processSmartComponentsMsg(msgDec *message.Decoder) error {
	reqID := msgDec.RequestID(false)
	count := int(msgDec.Int64())
	if count < 0 {
		msgDec.SetErr(fmt.Errorf("negative smart components size: %d", count))
		return wrapDecoderErr("processSmartComponentsMsg", msgDec.Err())
	}

	data := make([]models.SmartComponent, 0, count)
	for idx := 0; idx < count; idx++ {
		item := models.NewSmartComponent()
		item.BitNumber = msgDec.Int64()
		item.Exchange = msgDec.String()
		item.ExchangeLetter = msgDec.String()
		data = append(data, item)
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processSmartComponentsMsg", msgDec.Err())
	}

	return c.processSmartComponentsCommon(reqID, data)
}

func (c *Client) processSmartComponentsCommon(reqID int32, data []models.SmartComponent) error {
	c.reqMgr.withRequestWithID(reqID, func(resp interface{}) (bool, error) {
		smartResp := resp.(*models.SmartComponentsResponse)
		smartResp.ReqID = reqID
		smartResp.Components = append(smartResp.Components[:0], data...)
		return true, nil
	})
	return nil
}

func decodeLegacyContract(msgDec *message.Decoder, mode int) *models.Contract {
	contract := models.NewContract()
	contract.ConID = msgDec.Int32()
	contract.Symbol = msgDec.String()
	contract.SecType = models.SecurityType(msgDec.String())
	contract.LastTradeDateOrContractMonth = msgDec.String()
	contract.Strike = msgDec.FloatMax()
	contract.Right = msgDec.String()
	contract.Multiplier = msgDec.FloatMax()
	contract.Exchange = msgDec.String()
	contract.PrimaryExchange = msgDec.String()
	contract.Currency = msgDec.String()
	contract.LocalSymbol = msgDec.String()
	contract.TradingClass = msgDec.String()
	if mode >= 2 {
		contract.IncludeExpired = msgDec.Bool()
	}
	if mode >= 3 {
		contract.SecIDType = msgDec.String()
		contract.SecID = msgDec.String()
		contract.IssuerID = msgDec.String()
	}
	return contract
}
