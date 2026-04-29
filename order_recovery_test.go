package ibkr

import (
	"testing"

	"github.com/mxmauro/ibkr/common"
	"github.com/mxmauro/ibkr/models"
	"github.com/mxmauro/resetevent"
)

// -----------------------------------------------------------------------------

func newRecoveryTestClient(store OrderStore) *Client {
	c := &Client{
		orderStore:       store,
		isDisconnectedEv: resetevent.NewManualResetEvent(),
		nextValidReqID:   1,
		serverVersion:    212,
	}
	c.initRequestManager()
	return c
}

func queueOpenOrdersSnapshotRequest(c *Client) (*Request, *models.OpenOrdersResponse) {
	resp := &models.OpenOrdersResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.OPEN_ORDER_END),
		Response: resp,
	})
	c.reqMgr.addRequest(req)
	return req, resp
}

func queueCompletedOrdersSnapshotRequest(c *Client) (*Request, *models.CompletedOrdersResponse) {
	resp := &models.CompletedOrdersResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithoutID,
		MsgCode:  int(common.COMPLETED_ORDERS_END),
		Response: resp,
	})
	c.reqMgr.addRequest(req)
	return req, resp
}

func queueExecutionsSnapshotRequest(c *Client) (*Request, *models.ExecutionsResponse) {
	resp := &models.ExecutionsResponse{}
	req := c.createRequest(RequestOptions{
		Type:     RequestTypeRequestWithID,
		MsgCode:  int(common.EXECUTION_DATA_END),
		Response: resp,
	})
	c.reqMgr.addRequest(req)
	return req, resp
}

func assertRequestCompleted(t *testing.T, req *Request, name string) {
	t.Helper()

	select {
	case <-req.CompleteCh():
	default:
		t.Fatalf("%s request was not completed", name)
	}
}

// -----------------------------------------------------------------------------

func TestOrderRecoveryRebuildsOpenOrdersAfterReconnect(t *testing.T) {
	store := NewInMemoryOrderStore()
	client := newRecoveryTestClient(store)
	req, resp := queueOpenOrdersSnapshotRequest(client)

	contract := &models.Contract{
		Symbol:   "BAC",
		SecType:  models.SecurityTypeStock,
		Exchange: "SMART",
		Currency: "USD",
	}
	order := models.NewOrder()
	order.OrderID = 41
	order.ClientID = 7
	order.PermID = 9001
	order.OrderRef = "reconnect-open-41"
	order.Action = models.OrderActionBuy
	order.OrderType = models.OrderTypeLimit
	order.TotalQuantity = &models.DecimalOne
	orderState := models.NewOrderState()
	orderState.Status = "Submitted"

	if err := client.processOpenOrderCommon(contract, order, orderState); err != nil {
		t.Fatal(err)
	}
	if err := client.processOpenOrdersEndCommon(); err != nil {
		t.Fatal(err)
	}

	assertRequestCompleted(t, req, "open orders snapshot")
	if len(resp.Orders) != 1 {
		t.Fatalf("expected 1 open order in snapshot, got %d", len(resp.Orders))
	}

	recovered, ok := store.GetByOrderRef(order.OrderRef)
	if !ok {
		t.Fatalf("order %q was not rebuilt in the recovery store", order.OrderRef)
	}
	if !recovered.IsOpen {
		t.Fatalf("order %q was not marked as open", order.OrderRef)
	}
	if recovered.Order == nil || recovered.Order.PermID != order.PermID {
		t.Fatalf("order %q perm id was not rebuilt", order.OrderRef)
	}
	if recovered.OrderState == nil || recovered.OrderState.Status != "Submitted" {
		t.Fatalf("order %q status was not rebuilt from the open order snapshot", order.OrderRef)
	}
}

func TestOrderRecoveryRebuildsCompletedOrdersAndExecutionsAfterReconnect(t *testing.T) {
	store := NewInMemoryOrderStore()
	client := newRecoveryTestClient(store)
	completedReq, completedResp := queueCompletedOrdersSnapshotRequest(client)
	execReq, execResp := queueExecutionsSnapshotRequest(client)

	contract := &models.Contract{
		Symbol:   "BAC",
		SecType:  models.SecurityTypeStock,
		Exchange: "SMART",
		Currency: "USD",
	}
	order := models.NewOrder()
	order.OrderID = 58
	order.ClientID = 11
	order.PermID = 9058
	order.OrderRef = "reconnect-completed-58"
	order.Action = models.OrderActionBuy
	order.OrderType = models.OrderTypeLimit
	order.TotalQuantity = &models.DecimalOne
	orderState := models.NewOrderState()
	orderState.Status = "Cancelled"
	orderState.CompletedStatus = "Cancelled"
	orderState.CompletedTime = "20260429 18:00:00"

	execution := models.NewExecution()
	execution.OrderID = order.OrderID
	execution.PermID = order.PermID
	execution.OrderRef = order.OrderRef
	execution.ExecID = "E-58-1"
	execution.Side = "BOT"
	execution.AcctNumber = "DU0000000"
	execution.Shares = &models.DecimalOne
	execution.CumQty = &models.DecimalOne
	execution.Price = 42.15

	if err := client.processCompletedOrderCommon(contract, order, orderState); err != nil {
		t.Fatal(err)
	}
	if err := client.processCompletedOrdersEndCommon(); err != nil {
		t.Fatal(err)
	}
	if err := client.processExecutionDetailsCommon(execReq.ID(), contract, execution); err != nil {
		t.Fatal(err)
	}
	if err := client.processExecutionDetailsEndCommon(execReq.ID()); err != nil {
		t.Fatal(err)
	}

	assertRequestCompleted(t, completedReq, "completed orders snapshot")
	assertRequestCompleted(t, execReq, "executions snapshot")
	if len(completedResp.Orders) != 1 {
		t.Fatalf("expected 1 completed order in snapshot, got %d", len(completedResp.Orders))
	}
	if len(execResp.Executions) != 1 {
		t.Fatalf("expected 1 execution in snapshot, got %d", len(execResp.Executions))
	}

	recovered, ok := store.GetByOrderRef(order.OrderRef)
	if !ok {
		t.Fatalf("order %q was not rebuilt in the recovery store", order.OrderRef)
	}
	if !recovered.IsCompleted {
		t.Fatalf("order %q was not marked as completed", order.OrderRef)
	}
	if recovered.Order == nil || recovered.Order.PermID != order.PermID {
		t.Fatalf("order %q perm id was not rebuilt", order.OrderRef)
	}
	if recovered.OrderState == nil || recovered.OrderState.CompletedStatus != "Cancelled" {
		t.Fatalf("order %q completed status was not rebuilt", order.OrderRef)
	}
	if len(recovered.Executions) != 1 || recovered.Executions[0].ExecID != execution.ExecID {
		t.Fatalf("order %q executions were not rebuilt", order.OrderRef)
	}
}
