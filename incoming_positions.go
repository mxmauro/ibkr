package ibkr

import (
	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

type positionEventHandler interface {
	Position(account string, contract *models.Contract, position models.Decimal, avgCost float64)
}

type positionEndEventHandler interface {
	PositionEnd()
}

// -----------------------------------------------------------------------------

func (c *Client) processPositionDataMsg(msgDec *message.Decoder) error {
	msgDec.Skip()
	data := models.Position{
		Account:  msgDec.String(),
		Contract: decodeLegacyContract(msgDec, 1),
		Position: models.NewDecimalFromMessageDecoder(msgDec),
		AvgCost:  msgDec.Float(),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPositionDataMsg", msgDec.Err())
	}

	return c.processPositionDataCommon(data)
}

func (c *Client) processPositionDataProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.Position{}
	msgDec.Unmarshal(&pb)
	data := models.Position{
		Account:  msgDec.String(pb.Account),
		Contract: models.NewContractFromProtobufDecoder(msgDec, pb.Contract),
		Position: models.NewDecimalFromProtobufDecoder(msgDec, pb.Position),
		AvgCost:  msgDec.Float(pb.AvgCost),
	}
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPositionDataProtobuf", msgDec.Err())
	}

	return c.processPositionDataCommon(data)
}

func (c *Client) processPositionDataCommon(data models.Position) error {
	c.reqMgr.withActiveFirstRequestWithoutID(int(common.POSITION_DATA), func(resp interface{}) error {
		positionsResp := resp.(*models.PositionsResponse)
		positionsResp.Positions = append(positionsResp.Positions, data)
		return nil
	})

	if handler, ok := c.eventsHandler.(positionEventHandler); ok {
		handler.Position(data.Account, data.Contract, data.Position, data.AvgCost)
	}
	return nil
}

func (c *Client) processPositionEndMsg(msgDec *message.Decoder) error {
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPositionEndMsg", msgDec.Err())
	}
	return c.processPositionEndCommon()
}

func (c *Client) processPositionEndProtobuf(msgDec *protofmt.Decoder) error {
	pb := protobuf.PositionEnd{}
	msgDec.Unmarshal(&pb)
	if msgDec.Err() != nil {
		return wrapDecoderErr("processPositionEndProtobuf", msgDec.Err())
	}
	_ = pb
	return c.processPositionEndCommon()
}

func (c *Client) processPositionEndCommon() error {
	c.reqMgr.withFirstRequestWithoutID(int(common.POSITION_DATA), func(_ interface{}) (bool, error) {
		return true, nil
	})

	if handler, ok := c.eventsHandler.(positionEndEventHandler); ok {
		handler.PositionEnd()
	}
	return nil
}
