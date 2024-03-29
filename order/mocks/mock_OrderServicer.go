// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/enrico5b1b4/order-service/order (interfaces: OrderServicer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	errors "github.com/enrico5b1b4/order-service/errors"
	order "github.com/enrico5b1b4/order-service/order"
	gomock "github.com/golang/mock/gomock"
	go_uuid "github.com/satori/go.uuid"
)

// MockOrderServicer is a mock of OrderServicer interface
type MockOrderServicer struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServicerMockRecorder
}

// MockOrderServicerMockRecorder is the mock recorder for MockOrderServicer
type MockOrderServicerMockRecorder struct {
	mock *MockOrderServicer
}

// NewMockOrderServicer creates a new mock instance
func NewMockOrderServicer(ctrl *gomock.Controller) *MockOrderServicer {
	mock := &MockOrderServicer{ctrl: ctrl}
	mock.recorder = &MockOrderServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrderServicer) EXPECT() *MockOrderServicerMockRecorder {
	return m.recorder
}

// CompleteOrder mocks base method
func (m *MockOrderServicer) CompleteOrder(arg0 go_uuid.UUID, arg1 string) *errors.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteOrder", arg0, arg1)
	ret0, _ := ret[0].(*errors.Error)
	return ret0
}

// CompleteOrder indicates an expected call of CompleteOrder
func (mr *MockOrderServicerMockRecorder) CompleteOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteOrder", reflect.TypeOf((*MockOrderServicer)(nil).CompleteOrder), arg0, arg1)
}

// CreateOrder mocks base method
func (m *MockOrderServicer) CreateOrder(arg0 *order.Order) (*order.Order, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0)
	ret0, _ := ret[0].(*order.Order)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder
func (mr *MockOrderServicerMockRecorder) CreateOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderServicer)(nil).CreateOrder), arg0)
}

// GetOrderByID mocks base method
func (m *MockOrderServicer) GetOrderByID(arg0 go_uuid.UUID) (*order.Order, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", arg0)
	ret0, _ := ret[0].(*order.Order)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID
func (mr *MockOrderServicerMockRecorder) GetOrderByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockOrderServicer)(nil).GetOrderByID), arg0)
}

// GetOrders mocks base method
func (m *MockOrderServicer) GetOrders(arg0 string) ([]*order.Order, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", arg0)
	ret0, _ := ret[0].([]*order.Order)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders
func (mr *MockOrderServicerMockRecorder) GetOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockOrderServicer)(nil).GetOrders), arg0)
}

// ProcessOrder mocks base method
func (m *MockOrderServicer) ProcessOrder(arg0 go_uuid.UUID) *errors.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessOrder", arg0)
	ret0, _ := ret[0].(*errors.Error)
	return ret0
}

// ProcessOrder indicates an expected call of ProcessOrder
func (mr *MockOrderServicerMockRecorder) ProcessOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessOrder", reflect.TypeOf((*MockOrderServicer)(nil).ProcessOrder), arg0)
}
