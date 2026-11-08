package ibkr_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mxmauro/ibkr/models"
)

// -----------------------------------------------------------------------------

type orderStoreSpy struct {
	mtx               sync.Mutex
	createdOrders     []*models.Order
	openOrders        []*models.Order
	completedOrders   []*models.Order
	statuses          []models.OrderStatusUpdate
	executions        []*models.Execution
	permIDsByOrderID  map[int32]int64
	errorsByOrderID   map[int32][]string
	warningsByOrderID map[int32][]string
	commissionReports []*models.CommissionAndFeesReport
}

// -----------------------------------------------------------------------------

func newOrderStoreSpy() *orderStoreSpy {
	return &orderStoreSpy{
		permIDsByOrderID:  make(map[int32]int64),
		errorsByOrderID:   make(map[int32][]string),
		warningsByOrderID: make(map[int32][]string),
	}
}

func (s *orderStoreSpy) CreateOrder(_ *models.Contract, order *models.Order) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.createdOrders = append(s.createdOrders, cloneOrder(order))
}

func (s *orderStoreSpy) SetOrderPermanentID(orderID int32, _ int32, permID int64) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.permIDsByOrderID[orderID] = permID
}

func (s *orderStoreSpy) UpdateOrderStatus(update models.OrderStatusUpdate) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.statuses = append(s.statuses, update)
}

func (s *orderStoreSpy) RecordOrderError(orderID int32, _ time.Time, _ int, message string, _ string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.errorsByOrderID[orderID] = append(s.errorsByOrderID[orderID], message)
}

func (s *orderStoreSpy) RecordOrderWarning(orderID int32, _ time.Time, _ int, message string, _ string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.warningsByOrderID[orderID] = append(s.warningsByOrderID[orderID], message)
}

func (s *orderStoreSpy) RecordOpenOrder(_ *models.Contract, order *models.Order, _ *models.OrderState) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.openOrders = append(s.openOrders, cloneOrder(order))
}

func (s *orderStoreSpy) RecordExecution(_ *models.Contract, execution *models.Execution) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.executions = append(s.executions, cloneExecution(execution))
}

func (s *orderStoreSpy) RecordCommissionAndFeesReport(report *models.CommissionAndFeesReport) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.commissionReports = append(s.commissionReports, cloneCommissionAndFeesReport(report))
}

func (s *orderStoreSpy) BindOrder(permID int64, _ int32, orderID int32) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.permIDsByOrderID[orderID] = permID
}

func (s *orderStoreSpy) RecordCompletedOrder(_ *models.Contract, order *models.Order, _ *models.OrderState) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.completedOrders = append(s.completedOrders, cloneOrder(order))
}

func (s *orderStoreSpy) hasCreatedOrder(orderID int32, orderRef string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, order := range s.createdOrders {
		if order.OrderID == orderID && order.OrderRef == orderRef {
			return true
		}
	}
	return false
}

func (s *orderStoreSpy) hasOpenOrder(orderID int32) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, order := range s.openOrders {
		if order.OrderID == orderID {
			return true
		}
	}
	return false
}

func (s *orderStoreSpy) hasCompletedOrder(orderID int32) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, order := range s.completedOrders {
		if order.OrderID == orderID {
			return true
		}
	}
	return false
}

func (s *orderStoreSpy) hasStatus(orderID int32, statuses ...string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, update := range s.statuses {
		if update.OrderID == orderID && slices.Contains(statuses, update.Status) {
			return true
		}
	}
	return false
}

func (s *orderStoreSpy) permID(orderID int32) int64 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.permIDsByOrderID[orderID]
}

func (s *orderStoreSpy) errors(orderID int32) []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return append([]string(nil), s.errorsByOrderID[orderID]...)
}

func (s *orderStoreSpy) hasErrorContaining(orderID int32, substr string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, errMsg := range s.errorsByOrderID[orderID] {
		if strings.Contains(errMsg, substr) {
			return true
		}
	}
	return false
}

func (s *orderStoreSpy) warnings(orderID int32) []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return append([]string(nil), s.warningsByOrderID[orderID]...)
}

func (s *orderStoreSpy) hasWarningContaining(orderID int32, substr string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, errMsg := range s.warningsByOrderID[orderID] {
		if strings.Contains(errMsg, substr) {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------

func testPaperOrders(t *testing.T) {
	skipIfLiveTrading(t, "order creation and manipulation tests are paper-only")
	if !usStockRegularTradingHoursOpen(time.Now()) {
		testSkip(t, "US regular session is closed; paper order lifecycle assertions are not reliable")
	}

	store := newOrderStoreSpy()
	client := connectWithOrderStore(t, store)
	t.Cleanup(func() {
		client.Destroy()
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelCtx()

	orderRef := fmt.Sprintf("codex-paper-%d", time.Now().UnixNano())
	order := models.NewOrder()
	order.Action = models.OrderActionBuy
	order.OrderType = models.OrderTypeLimit
	order.TotalQuantity = &models.DecimalOne
	limitPrice := 1.00
	order.LmtPrice = &limitPrice
	order.TIF = models.TimeInForceDay
	order.OrderRef = orderRef
	order.Transmit = true

	orderID, err := client.PlaceOrder(models.PlaceOrderRequestOptions{
		Contract: getContract("BAC", "SMART"),
		Order:    order,
	})
	if err != nil {
		testFatal(t, err)
	}
	testLogf(t, "  placed paper order %d with ref %s", orderID, orderRef)

	waitForCondition(t, 20*time.Second, "order store create callback", func() bool {
		return store.hasCreatedOrder(orderID, orderRef)
	})
	waitForConditionOr(t, 10*time.Second, "open order callback or read-only rejection", func() (bool, string) {
		if store.hasOpenOrder(orderID) {
			return true, ""
		}
		if store.hasErrorContaining(orderID, "Read-Only mode") {
			return true, "gateway API is in Read-Only mode"
		}
		return false, ""
	})
	if store.hasErrorContaining(orderID, "Read-Only mode") {
		testSkip(t, "paper gateway API is in Read-Only mode; order lifecycle tests require write access")
	}
	waitForCondition(t, 30*time.Second, "open order callback", func() bool {
		return store.hasOpenOrder(orderID)
	})
	waitForCondition(t, 30*time.Second, "order status callback", func() bool {
		return store.hasStatus(orderID, "PendingSubmit", "PreSubmitted", "Submitted", "ApiPending", "ApiCancelled")
	})

	if errs := store.errors(orderID); len(errs) > 0 {
		testFatalf(t, "unexpected order errors for %d: %v", orderID, errs)
	}

	openOrders, err := client.RequestOpenOrders(ctx)
	if err != nil {
		testFatal(t, err)
	}
	if !containsOpenOrder(openOrders, orderID, orderRef) {
		testFatalf(t, "order %d was not returned by RequestOpenOrders", orderID)
	}

	if permID := store.permID(orderID); permID != 0 {
		testLogf(t, "  permanent id: %d", permID)
	}

	if err := client.CancelOrder(models.CancelOrderRequestOptions{
		OrderID: orderID,
	}); err != nil {
		testFatal(t, err)
	}
	testLogf(t, "  cancel requested for order %d", orderID)

	waitForCondition(t, 30*time.Second, "cancelled status", func() bool {
		return store.hasStatus(orderID, "PendingCancel", "Cancelled", "ApiCancelled")
	})
	waitForCondition(t, 45*time.Second, "completed order snapshot after cancel", func() bool {
		completedOrders, err := client.RequestCompletedOrders(ctx, models.CompletedOrdersRequestOptions{
			APIOnly: true,
		})
		if err != nil {
			testFatalf(t, "RequestCompletedOrders failed: %v", err)
		}
		if containsCompletedOrder(completedOrders, orderID, orderRef) {
			return true
		}
		return store.hasCompletedOrder(orderID)
	})

	executions, err := client.RequestExecutions(ctx, models.ExecutionsRequestOptions{
		Filter: &models.ExecutionFilter{
			Symbol: "BAC",
		},
	})
	if err != nil {
		testFatal(t, err)
	}
	testLogf(t, "  executions returned: %d", len(executions.Executions))
}

// -----------------------------------------------------------------------------

func waitForCondition(t *testing.T, timeout time.Duration, description string, cond func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	testFatalf(t, "timeout waiting for %s", description)
}

func waitForConditionOr(t *testing.T, timeout time.Duration, description string, cond func() (bool, string)) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ok, reason := cond()
		if ok {
			if reason != "" {
				testLogf(t, "  %s: %s", description, reason)
			}
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	testFatalf(t, "timeout waiting for %s", description)
}

func containsOpenOrder(resp *models.OpenOrdersResponse, orderID int32, orderRef string) bool {
	for _, orderData := range resp.Orders {
		if orderData.Order != nil && orderData.Order.OrderID == orderID && orderData.Order.OrderRef == orderRef {
			return true
		}
	}
	return false
}

func containsCompletedOrder(resp *models.CompletedOrdersResponse, orderID int32, orderRef string) bool {
	for _, orderData := range resp.Orders {
		if orderData.Order != nil && orderData.Order.OrderID == orderID && orderData.Order.OrderRef == orderRef {
			return true
		}
	}
	return false
}

func cloneOrder(order *models.Order) *models.Order {
	if order == nil {
		return nil
	}
	clone := *order
	return &clone
}

func cloneExecution(execution *models.Execution) *models.Execution {
	if execution == nil {
		return nil
	}
	clone := *execution
	return &clone
}

func cloneCommissionAndFeesReport(report *models.CommissionAndFeesReport) *models.CommissionAndFeesReport {
	if report == nil {
		return nil
	}
	clone := *report
	return &clone
}
