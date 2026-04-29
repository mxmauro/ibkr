package ibkr

import (
	"time"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

// OrderStore persists order lifecycle events emitted by the client.
type OrderStore interface {
	CreateOrder(contract *models.Contract, order *models.Order)
	SetOrderPermanentID(orderID int32, clientID int32, permID int64)
	UpdateOrderStatus(update models.OrderStatusUpdate)
	RecordOrderError(orderID int32, ts time.Time, code int, message string, advancedOrderRejectJSON string)
	RecordOpenOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState)
	RecordExecution(contract *models.Contract, execution *models.Execution)
	RecordCommissionAndFeesReport(report *models.CommissionAndFeesReport)
	BindOrder(permID int64, clientID int32, orderID int32)
	RecordCompletedOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState)
}

type noopOrderStore struct{}

// -----------------------------------------------------------------------------

func (noopOrderStore) CreateOrder(*models.Contract, *models.Order) {}

func (noopOrderStore) SetOrderPermanentID(int32, int32, int64) {}

func (noopOrderStore) UpdateOrderStatus(models.OrderStatusUpdate) {}

func (noopOrderStore) RecordOrderError(int32, time.Time, int, string, string) {}

func (noopOrderStore) RecordOpenOrder(*models.Contract, *models.Order, *models.OrderState) {}

func (noopOrderStore) RecordExecution(*models.Contract, *models.Execution) {}

func (noopOrderStore) RecordCommissionAndFeesReport(*models.CommissionAndFeesReport) {}

func (noopOrderStore) BindOrder(int64, int32, int32) {}

func (noopOrderStore) RecordCompletedOrder(*models.Contract, *models.Order, *models.OrderState) {}
