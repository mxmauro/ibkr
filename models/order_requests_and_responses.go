package models

import "time"

// -----------------------------------------------------------------------------

type OrderStatusUpdate struct {
	OrderID        int32
	Status         string
	Filled         Decimal
	Remaining      Decimal
	AvgFilledPrice float64
	PermID         int64
	ParentID       int32
	LastFillPrice  float64
	ClientID       int32
	WhyHeld        string
	MarketCapPrice float64
}

type OrderError struct {
	OrderID                 int32
	Timestamp               time.Time
	Code                    int
	Message                 string
	AdvancedOrderRejectJSON string
}

type OpenOrderData struct {
	Contract   *Contract
	Order      *Order
	OrderState *OrderState
}

type ExecutionData struct {
	RequestID int32
	Contract  *Contract
	Execution *Execution
}

type CompletedOrderData struct {
	Contract   *Contract
	Order      *Order
	OrderState *OrderState
}

type OpenOrdersResponse struct {
	Orders []OpenOrderData
}

type ExecutionsResponse struct {
	Executions []ExecutionData
}

type CompletedOrdersResponse struct {
	Orders []CompletedOrderData
}

type PlaceOrderRequestOptions struct {
	Contract *Contract
	Order    *Order
}

type CancelOrderRequestOptions struct {
	OrderID     int32
	OrderCancel *OrderCancel
}

type AutoOpenOrdersRequestOptions struct {
	AutoBind bool
}

type ExecutionsRequestOptions struct {
	Filter *ExecutionFilter
}

type CompletedOrdersRequestOptions struct {
	APIOnly bool
}
