package models

import "time"

// -----------------------------------------------------------------------------

// OrderStatusUpdate contains an asynchronous order status update from IB.
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

// OrderErrorOrWarning captures an order-related error reported by IB.
type OrderErrorOrWarning struct {
	OrderID                 int32
	Timestamp               time.Time
	Code                    int
	Message                 string
	AdvancedOrderRejectJSON string
}

// OpenOrderData bundles the contract, order, and state reported by an open-order snapshot.
type OpenOrderData struct {
	Contract   *Contract
	Order      *Order
	OrderState *OrderState
}

// ExecutionData bundles an execution with the originating request and contract.
type ExecutionData struct {
	RequestID int32
	Contract  *Contract
	Execution *Execution
}

// CompletedOrderData bundles the contract, order, and state for a completed order.
type CompletedOrderData struct {
	Contract   *Contract
	Order      *Order
	OrderState *OrderState
}

// OpenOrdersResponse contains the open orders returned by an open-orders request.
type OpenOrdersResponse struct {
	Orders []OpenOrderData
}

// ExecutionsResponse contains the executions returned by an executions request.
type ExecutionsResponse struct {
	Executions []ExecutionData
}

// CompletedOrdersResponse contains the completed orders returned by IB.
type CompletedOrdersResponse struct {
	Orders []CompletedOrderData
}

// PlaceOrderRequestOptions configures a place-order request.
type PlaceOrderRequestOptions struct {
	Contract *Contract
	Order    *Order
}

// CancelOrderRequestOptions configures an order cancellation request.
type CancelOrderRequestOptions struct {
	OrderID     int32
	OrderCancel *OrderCancel
}

// AutoOpenOrdersRequestOptions configures whether new TWS orders auto-bind to the API client.
type AutoOpenOrdersRequestOptions struct {
	AutoBind bool
}

// ExecutionsRequestOptions configures an executions query.
type ExecutionsRequestOptions struct {
	Filter *ExecutionFilter
}

// CompletedOrdersRequestOptions configures a completed-orders query.
type CompletedOrdersRequestOptions struct {
	APIOnly bool
}
