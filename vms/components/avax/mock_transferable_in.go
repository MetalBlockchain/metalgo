// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/vms/components/avax (interfaces: TransferableIn)
//
// Generated by this command:
//
//	mockgen -package=avax -destination=vms/components/avax/mock_transferable_in.go github.com/MetalBlockchain/metalgo/vms/components/avax TransferableIn
//

// Package avax is a generated GoMock package.
package avax

import (
	reflect "reflect"

	snow "github.com/MetalBlockchain/metalgo/snow"
	gomock "go.uber.org/mock/gomock"
)

// MockTransferableIn is a mock of TransferableIn interface.
type MockTransferableIn struct {
	ctrl     *gomock.Controller
	recorder *MockTransferableInMockRecorder
}

// MockTransferableInMockRecorder is the mock recorder for MockTransferableIn.
type MockTransferableInMockRecorder struct {
	mock *MockTransferableIn
}

// NewMockTransferableIn creates a new mock instance.
func NewMockTransferableIn(ctrl *gomock.Controller) *MockTransferableIn {
	mock := &MockTransferableIn{ctrl: ctrl}
	mock.recorder = &MockTransferableInMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferableIn) EXPECT() *MockTransferableInMockRecorder {
	return m.recorder
}

// Amount mocks base method.
func (m *MockTransferableIn) Amount() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Amount")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Amount indicates an expected call of Amount.
func (mr *MockTransferableInMockRecorder) Amount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Amount", reflect.TypeOf((*MockTransferableIn)(nil).Amount))
}

// Cost mocks base method.
func (m *MockTransferableIn) Cost() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cost")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cost indicates an expected call of Cost.
func (mr *MockTransferableInMockRecorder) Cost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cost", reflect.TypeOf((*MockTransferableIn)(nil).Cost))
}

// InitCtx mocks base method.
func (m *MockTransferableIn) InitCtx(arg0 *snow.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitCtx", arg0)
}

// InitCtx indicates an expected call of InitCtx.
func (mr *MockTransferableInMockRecorder) InitCtx(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitCtx", reflect.TypeOf((*MockTransferableIn)(nil).InitCtx), arg0)
}

// Verify mocks base method.
func (m *MockTransferableIn) Verify() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify")
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockTransferableInMockRecorder) Verify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockTransferableIn)(nil).Verify))
}
