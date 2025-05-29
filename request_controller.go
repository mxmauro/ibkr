package ibkr

import (
	"container/list"
	"context"
	"math"
	"net"
	"sync"
	"sync/atomic"
)

// -----------------------------------------------------------------------------

type RequestManager struct {
	mtx           sync.Mutex
	reqsWithID    map[int32]*Request
	reqsWithoutID map[int]*NonIdRequestList
}

type Request struct {
	_type       RequestType
	id          int32
	msgCode     int
	done        int32
	completedCh chan struct{}
	completeCB  RequestCompleteCallback
	errHolder   atomic.Value
	responseMtx sync.Mutex
	response    interface{}
}

type RequestOptions struct {
	Type       RequestType
	MsgCode    int
	Response   interface{}
	CompleteCB RequestCompleteCallback
}

type NonIdRequestList struct {
	list.List
}

type WithRequestWithIdCallback func(resp interface{}) (done bool, err error)
type WithRequestWithoutIdCallback func(resp interface{}) (err error)
type WithRequestWithTickerIdCallback func(resp interface{}) (done bool, err error)

type RequestCompleteCallback func(req *Request, err error)

type RequestType int

const (
	RequestTypeRequestWithID RequestType = iota
	RequestTypeRequestWithoutID
	RequestTypeRequestWithTickerID
)

// -----------------------------------------------------------------------------

func (c *Client) initRequestManager() {
	c.reqMgr = RequestManager{
		mtx:           sync.Mutex{},
		reqsWithID:    make(map[int32]*Request),
		reqsWithoutID: make(map[int]*NonIdRequestList),
	}
}

func (c *Client) createRequest(opts RequestOptions) *Request {
	req := &Request{
		_type:       opts.Type,
		id:          c.getNextRequestID(opts.Type),
		msgCode:     opts.MsgCode,
		completeCB:  opts.CompleteCB,
		responseMtx: sync.Mutex{},
		response:    opts.Response,
	}
	if opts.Type == RequestTypeRequestWithID || opts.Type == RequestTypeRequestWithoutID {
		req.completedCh = make(chan struct{})
	}
	return req
}

func (c *Client) getNextRequestID(_type RequestType) int32 {
	var nextReqID int32
	var addr *int32

	switch _type {
	case RequestTypeRequestWithID:
		fallthrough
	case RequestTypeRequestWithTickerID:
		addr = &c.nextValidReqID

	case RequestTypeRequestWithoutID:
		addr = &c.nextValidReqWithoutID
	}

	for {
		currentReqID := atomic.LoadInt32(addr)
		if currentReqID <= 0 {
			return 0
		}
		if currentReqID != math.MaxInt32 {
			nextReqID = currentReqID + 1
		} else {
			nextReqID = 1
		}
		if atomic.CompareAndSwapInt32(addr, currentReqID, nextReqID) {
			return currentReqID
		}
	}
}

func (c *Client) waitRequestCompletion(ctx context.Context, req *Request) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-c.isDisconnectedEv.WaitCh():
		return net.ErrClosed

	case <-req.CompleteCh():
		err := req.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (rm *RequestManager) addRequest(req *Request) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()

	switch req._type {
	case RequestTypeRequestWithID:
		fallthrough
	case RequestTypeRequestWithTickerID:
		rm.reqsWithID[req.id] = req

	case RequestTypeRequestWithoutID:
		var l *NonIdRequestList
		var ok bool

		if l, ok = rm.reqsWithoutID[req.msgCode]; !ok {
			l = &NonIdRequestList{
				List: list.List{},
			}
			rm.reqsWithoutID[req.msgCode] = l
		}
		l.queueRequest(req)
	}
}

func (rm *RequestManager) removeRequest(req *Request, err error) {
	rm.mtx.Lock()
	switch req._type {
	case RequestTypeRequestWithID:
		fallthrough
	case RequestTypeRequestWithTickerID:
		delete(rm.reqsWithID, req.id)

	case RequestTypeRequestWithoutID:
		if l, ok := rm.reqsWithoutID[req.msgCode]; ok {
			l.removeRequest(req.id)
		}
	}
	rm.mtx.Unlock()

	// Complete it
	req.complete(err)
}

func (rm *RequestManager) removeAndTryCancelAllRequests(err error) {
	rm.mtx.Lock()
	oldReqsWithID := rm.reqsWithID
	oldReqsWithoutID := rm.reqsWithoutID
	rm.reqsWithID = make(map[int32]*Request)
	rm.reqsWithoutID = make(map[int]*NonIdRequestList)
	rm.mtx.Unlock()

	for _, req := range oldReqsWithID {
		req.complete(err)
	}
	for _, l := range oldReqsWithoutID {
		for elem := l.List.Front(); elem != nil; elem = elem.Next() {
			req := elem.Value.(*Request)
			req.complete(err)
		}
	}
}

func (rm *RequestManager) withRequestWithID(reqID int32, cb WithRequestWithIdCallback) {
	rm.mtx.Lock()
	req, ok := rm.reqsWithID[reqID]
	rm.mtx.Unlock()

	if !ok || req.Err() != nil {
		return
	}

	req.responseMtx.Lock()
	if req.response == nil {
		req.responseMtx.Unlock()
		return
	}
	done, err := cb(req.response)
	req.responseMtx.Unlock()

	if done || err != nil {
		rm.mtx.Lock()
		delete(rm.reqsWithID, reqID)
		rm.mtx.Unlock()

		req.complete(err)
	}
}

func (rm *RequestManager) withRequestWithoutID(msgCode int, cb WithRequestWithoutIdCallback) {
	var req *Request

	rm.mtx.Lock()
	if l, ok := rm.reqsWithoutID[msgCode]; ok {
		req = l.popFirstRequest()
	}
	rm.mtx.Unlock()

	if req == nil || req.Err() != nil {
		return
	}

	req.responseMtx.Lock()
	if req.response == nil {
		req.responseMtx.Unlock()
		return
	}
	err := cb(req.response)
	req.responseMtx.Unlock()

	req.complete(err)
}

func (req *Request) Type() RequestType {
	return req._type
}

func (req *Request) ID() int32 {
	return req.id
}

func (req *Request) Err() error {
	v := req.errHolder.Load()
	if v == nil {
		return nil
	}
	return v.(error)
}

func (req *Request) CompleteCh() <-chan struct{} {
	return req.completedCh
}

func (req *Request) complete(err error) {
	if atomic.CompareAndSwapInt32(&req.done, 0, 1) {
		req.responseMtx.Lock()
		req.response = nil
		req.responseMtx.Unlock()

		if err != nil {
			req.errHolder.Store(err)
		}
		if req.completeCB != nil {
			req.completeCB(req, err)
		}
		if req.completedCh != nil {
			close(req.completedCh)
		}
	}
}

func (nirl *NonIdRequestList) queueRequest(req *Request) {
	_ = nirl.List.PushBack(req)
}

func (nirl *NonIdRequestList) popFirstRequest() *Request {
	elem := nirl.List.Front()
	if elem == nil {
		return nil
	}
	nirl.List.Remove(elem)
	return elem.Value.(*Request)
}

func (nirl *NonIdRequestList) removeRequest(id int32) {
	for elem := nirl.List.Front(); elem != nil; elem = elem.Next() {
		req := elem.Value.(*Request)
		if req.id == id {
			nirl.List.Remove(elem)
			return
		}
	}
}
