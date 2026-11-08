package ibkr

import (
	"context"
	"net"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/ibkr/proto/protobuf"
	"github.com/mxmauro/ibkr/utils/encoders/message"
	"github.com/mxmauro/ibkr/utils/encoders/protofmt"
)

// -----------------------------------------------------------------------------

// RequestContractDetails retrieves all details for contracts matching the provided input.
func (c *Client) RequestContractDetails(ctx context.Context, opts models.ContractDetailsRequestOptions) (*models.ContractDetailsResponse, error) {
	var msgEnc *message.Encoder

	if opts.Contract == nil {
		return nil, errInvalidContract
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.ContractDetailsResponse{
		ContractDetails: make([]*models.ContractDetails, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_CONTRACT_DATA,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_CONTRACT_DATA) {
		pb := protobuf.ContractDataRequest{
			ReqId:    protofmt.Int32(req.ID()),
			Contract: opts.Contract.Proto(nil),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_CONTRACT_DATA + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		const version = 8
		msgEnc = message.NewEncoder().Reserve(21).
			RawUInt32(common.REQ_CONTRACT_DATA).
			RequestID(req.ID()).
			Int(version).
			Marshal(opts.Contract, 3)
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

// RequestMatchingSymbols retrieves all details for contracts matching the provided input.
func (c *Client) RequestMatchingSymbols(ctx context.Context, opts models.MatchingSymbolsRequestOptions) (*models.MatchingSymbolsResponse, error) {
	var msgEnc *message.Encoder

	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.MatchingSymbolsResponse{
		ContractDescriptions: make([]*models.ContractDescription, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_MATCHING_SYMBOLS,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_MATCHING_SYMBOLS) {
		pb := protobuf.MatchingSymbolsRequest{
			ReqId:   protofmt.Int32(req.ID()),
			Pattern: protofmt.String(opts.Pattern),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_MATCHING_SYMBOLS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(3).
			RawUInt32(common.REQ_MATCHING_SYMBOLS).
			RequestID(req.ID()).
			String(opts.Pattern)
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

// RequestSecurityDefinitionOptionParameters retrieves option chain metadata for an underlying.
func (c *Client) RequestSecurityDefinitionOptionParameters(
	ctx context.Context, opts models.SecurityDefinitionOptionParametersRequestOptions,
) (*models.SecurityDefinitionOptionParametersResponse, error) {
	var msgEnc *message.Encoder

	if len(opts.UnderlyingSymbol) == 0 {
		return nil, errInvalidUnderlyingSymbol
	}
	if len(opts.UnderlyingSecurityType) == 0 {
		return nil, errInvalidUnderlyingSecurityType
	}
	if opts.UnderlyingContractID <= 0 {
		return nil, errInvalidUnderlyingContractID
	}
	if !c.rp.Acquire() {
		return nil, net.ErrClosed
	}
	defer c.rp.Release()

	resp := &models.SecurityDefinitionOptionParametersResponse{
		Parameters: make([]models.SecurityDefinitionOptionParameter, 0),
	}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  common.REQ_SEC_DEF_OPT_PARAMS,
		Response: resp,
	})

	if c.isProtoBufAvailable(common.REQ_SEC_DEF_OPT_PARAMS) {
		pb := protobuf.SecDefOptParamsRequest{
			ReqId:             protofmt.Int32(req.ID()),
			UnderlyingSymbol:  protofmt.String(opts.UnderlyingSymbol),
			FutFopExchange:    protofmt.String(opts.FutFopExchange),
			UnderlyingSecType: protofmt.String(string(opts.UnderlyingSecurityType)),
			UnderlyingConId:   protofmt.Int32(opts.UnderlyingContractID),
		}
		msgEnc = message.NewEncoder().
			RawUInt32(common.REQ_SEC_DEF_OPT_PARAMS + common.PROTOBUF_MSG_ID).
			Proto(&pb)
	} else {
		msgEnc = message.NewEncoder().Reserve(6).
			RawUInt32(common.REQ_SEC_DEF_OPT_PARAMS).
			RequestID(req.ID()).
			String(opts.UnderlyingSymbol).
			String(opts.FutFopExchange).
			String(string(opts.UnderlyingSecurityType)).
			Int32(opts.UnderlyingContractID)
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
