// Code generated by MockGen. DO NOT EDIT.
// Source: worker_queue.go

// Package lucy is a generated GoMock package.
package lucy

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockWorkerQueue is a mock of WorkerQueue interface
type MockWorkerQueue struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerQueueMockRecorder
}

// MockWorkerQueueMockRecorder is the mock recorder for MockWorkerQueue
type MockWorkerQueueMockRecorder struct {
	mock *MockWorkerQueue
}

// NewMockWorkerQueue creates a new mock instance
func NewMockWorkerQueue(ctrl *gomock.Controller) *MockWorkerQueue {
	mock := &MockWorkerQueue{ctrl: ctrl}
	mock.recorder = &MockWorkerQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkerQueue) EXPECT() *MockWorkerQueueMockRecorder {
	return m.recorder
}

// SubscribeRequests mocks base method
func (m *MockWorkerQueue) SubscribeRequests(ctx context.Context) (<-chan *Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeRequests", ctx)
	ret0, _ := ret[0].(<-chan *Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeRequests indicates an expected call of SubscribeRequests
func (mr *MockWorkerQueueMockRecorder) SubscribeRequests(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeRequests", reflect.TypeOf((*MockWorkerQueue)(nil).SubscribeRequests), ctx)
}

// RetryRequest mocks base method
func (m *MockWorkerQueue) RetryRequest(request *Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetryRequest", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// RetryRequest indicates an expected call of RetryRequest
func (mr *MockWorkerQueueMockRecorder) RetryRequest(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetryRequest", reflect.TypeOf((*MockWorkerQueue)(nil).RetryRequest), request)
}
