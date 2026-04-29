package ibkr

import (
	"sync"
	"time"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

// StoredOrder holds the reconstructed in-memory state for a single order identity.
type StoredOrder struct {
	Contract              *models.Contract
	Order                 *models.Order
	OrderState            *models.OrderState
	LastStatus            *models.OrderStatusUpdate
	Executions            []*models.Execution
	Errors                []models.OrderError
	CommissionAndFeeItems []*models.CommissionAndFeesReport
	IsOpen                bool
	IsCompleted           bool
}

// InMemoryOrderStore stores order lifecycle data in process for tests and lightweight integrations.
type InMemoryOrderStore struct {
	mtx                   sync.Mutex
	byOrderID             map[int32]*StoredOrder
	byPermID              map[int64]*StoredOrder
	byOrderRef            map[string]*StoredOrder
	commissionAndFeeItems []*models.CommissionAndFeesReport
}

// -----------------------------------------------------------------------------

// NewInMemoryOrderStore creates an empty in-memory order store.
func NewInMemoryOrderStore() *InMemoryOrderStore {
	return &InMemoryOrderStore{
		byOrderID:  make(map[int32]*StoredOrder),
		byPermID:   make(map[int64]*StoredOrder),
		byOrderRef: make(map[string]*StoredOrder),
	}
}

// CreateOrder records a locally created order before IB acknowledges it.
func (s *InMemoryOrderStore) CreateOrder(contract *models.Contract, order *models.Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(order.OrderID, order.OrderRef, order.PermID)
	state.Contract = contract
	state.Order = order
}

// SetOrderPermanentID links a local order ID to the permanent ID assigned by IB.
func (s *InMemoryOrderStore) SetOrderPermanentID(orderID int32, _ int32, permID int64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(orderID, "", permID)
	if state.Order != nil {
		state.Order.PermID = permID
	}
}

// UpdateOrderStatus stores the latest status update seen for an order.
func (s *InMemoryOrderStore) UpdateOrderStatus(update models.OrderStatusUpdate) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(update.OrderID, "", update.PermID)
	updateCopy := update
	state.LastStatus = &updateCopy
}

// RecordOrderError appends an order-scoped API error to the stored state.
func (s *InMemoryOrderStore) RecordOrderError(
	orderID int32, ts time.Time, code int, message string, advancedOrderRejectJSON string,
) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(orderID, "", 0)
	state.Errors = append(state.Errors, models.OrderError{
		OrderID:                 orderID,
		Timestamp:               ts,
		Code:                    code,
		Message:                 message,
		AdvancedOrderRejectJSON: advancedOrderRejectJSON,
	})
}

// RecordOpenOrder merges an open-order snapshot or callback into the stored state.
func (s *InMemoryOrderStore) RecordOpenOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(order.OrderID, order.OrderRef, order.PermID)
	state.Contract = contract
	state.Order = order
	state.OrderState = orderState
	state.IsOpen = true
}

// RecordExecution appends an execution to the stored order state.
func (s *InMemoryOrderStore) RecordExecution(_ *models.Contract, execution *models.Execution) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(execution.OrderID, execution.OrderRef, execution.PermID)
	state.Executions = append(state.Executions, execution)
}

// RecordCommissionAndFeesReport stores commission reports independently from order records.
func (s *InMemoryOrderStore) RecordCommissionAndFeesReport(report *models.CommissionAndFeesReport) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if report == nil {
		return
	}
	s.commissionAndFeeItems = append(s.commissionAndFeeItems, report)
}

// BindOrder links a permanent ID to a local order ID.
func (s *InMemoryOrderStore) BindOrder(permID int64, _ int32, orderID int32) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_ = s.ensureOrderLocked(orderID, "", permID)
}

// RecordCompletedOrder merges a completed-order snapshot or callback into the stored state.
func (s *InMemoryOrderStore) RecordCompletedOrder(contract *models.Contract, order *models.Order, orderState *models.OrderState) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state := s.ensureOrderLocked(order.OrderID, order.OrderRef, order.PermID)
	state.Contract = contract
	state.Order = order
	state.OrderState = orderState
	state.IsCompleted = true
}

// GetByOrderID returns a copy of the stored order state for a local order ID.
func (s *InMemoryOrderStore) GetByOrderID(orderID int32) (*StoredOrder, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state, ok := s.byOrderID[orderID]
	if !ok {
		return nil, false
	}
	return cloneStoredOrder(state), true
}

// GetByPermID returns a copy of the stored order state for a permanent ID.
func (s *InMemoryOrderStore) GetByPermID(permID int64) (*StoredOrder, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state, ok := s.byPermID[permID]
	if !ok {
		return nil, false
	}
	return cloneStoredOrder(state), true
}

// GetByOrderRef returns a copy of the stored order state for an order reference.
func (s *InMemoryOrderStore) GetByOrderRef(orderRef string) (*StoredOrder, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	state, ok := s.byOrderRef[orderRef]
	if !ok {
		return nil, false
	}
	return cloneStoredOrder(state), true
}

// CommissionAndFeesReports returns the recorded commission reports.
func (s *InMemoryOrderStore) CommissionAndFeesReports() []*models.CommissionAndFeesReport {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return append([]*models.CommissionAndFeesReport(nil), s.commissionAndFeeItems...)
}

func (s *InMemoryOrderStore) ensureOrderLocked(orderID int32, orderRef string, permID int64) *StoredOrder {
	var state *StoredOrder

	if permID != 0 {
		state = s.byPermID[permID]
	}
	if state == nil && orderID != 0 {
		state = s.byOrderID[orderID]
	}
	if state == nil && orderRef != "" {
		state = s.byOrderRef[orderRef]
	}
	if state == nil {
		state = &StoredOrder{}
	}

	if permID != 0 {
		s.byPermID[permID] = state
	}
	if orderID != 0 {
		s.byOrderID[orderID] = state
	}
	if orderRef != "" {
		s.byOrderRef[orderRef] = state
	}
	return state
}

func cloneStoredOrder(state *StoredOrder) *StoredOrder {
	if state == nil {
		return nil
	}

	clone := *state
	if state.LastStatus != nil {
		statusCopy := *state.LastStatus
		clone.LastStatus = &statusCopy
	}
	if len(state.Executions) > 0 {
		clone.Executions = append([]*models.Execution(nil), state.Executions...)
	}
	if len(state.Errors) > 0 {
		clone.Errors = append([]models.OrderError(nil), state.Errors...)
	}
	if len(state.CommissionAndFeeItems) > 0 {
		clone.CommissionAndFeeItems = append([]*models.CommissionAndFeesReport(nil), state.CommissionAndFeeItems...)
	}
	return &clone
}
