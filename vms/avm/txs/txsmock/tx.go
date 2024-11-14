// Code generated by MockGen. DO NOT EDIT.
// Source: vms/avm/txs/tx.go
//
// Generated by this command:
//
//	mockgen -source=vms/avm/txs/tx.go -destination=vms/avm/txs/txsmock/tx.go -package=txsmock -exclude_interfaces= -mock_names=UnsignedTx=UnsignedTx
//

// Package txsmock is a generated GoMock package.
package txsmock

import (
	reflect "reflect"

	ids "github.com/MetalBlockchain/metalgo/ids"
	snow "github.com/MetalBlockchain/metalgo/snow"
	set "github.com/MetalBlockchain/metalgo/utils/set"
	txs "github.com/MetalBlockchain/metalgo/vms/avm/txs"
	avax "github.com/MetalBlockchain/metalgo/vms/components/avax"
	gomock "go.uber.org/mock/gomock"
)

// UnsignedTx is a mock of UnsignedTx interface.
type UnsignedTx struct {
	ctrl     *gomock.Controller
	recorder *UnsignedTxMockRecorder
}

// UnsignedTxMockRecorder is the mock recorder for UnsignedTx.
type UnsignedTxMockRecorder struct {
	mock *UnsignedTx
}

// NewUnsignedTx creates a new mock instance.
func NewUnsignedTx(ctrl *gomock.Controller) *UnsignedTx {
	mock := &UnsignedTx{ctrl: ctrl}
	mock.recorder = &UnsignedTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *UnsignedTx) EXPECT() *UnsignedTxMockRecorder {
	return m.recorder
}

// Bytes mocks base method.
func (m *UnsignedTx) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes.
func (mr *UnsignedTxMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*UnsignedTx)(nil).Bytes))
}

// InitCtx mocks base method.
func (m *UnsignedTx) InitCtx(ctx *snow.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitCtx", ctx)
}

// InitCtx indicates an expected call of InitCtx.
func (mr *UnsignedTxMockRecorder) InitCtx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitCtx", reflect.TypeOf((*UnsignedTx)(nil).InitCtx), ctx)
}

// InputIDs mocks base method.
func (m *UnsignedTx) InputIDs() set.Set[ids.ID] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputIDs")
	ret0, _ := ret[0].(set.Set[ids.ID])
	return ret0
}

// InputIDs indicates an expected call of InputIDs.
func (mr *UnsignedTxMockRecorder) InputIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputIDs", reflect.TypeOf((*UnsignedTx)(nil).InputIDs))
}

// InputUTXOs mocks base method.
func (m *UnsignedTx) InputUTXOs() []*avax.UTXOID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputUTXOs")
	ret0, _ := ret[0].([]*avax.UTXOID)
	return ret0
}

// InputUTXOs indicates an expected call of InputUTXOs.
func (mr *UnsignedTxMockRecorder) InputUTXOs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputUTXOs", reflect.TypeOf((*UnsignedTx)(nil).InputUTXOs))
}

// NumCredentials mocks base method.
func (m *UnsignedTx) NumCredentials() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumCredentials")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumCredentials indicates an expected call of NumCredentials.
func (mr *UnsignedTxMockRecorder) NumCredentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumCredentials", reflect.TypeOf((*UnsignedTx)(nil).NumCredentials))
}

// SetBytes mocks base method.
func (m *UnsignedTx) SetBytes(unsignedBytes []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBytes", unsignedBytes)
}

// SetBytes indicates an expected call of SetBytes.
func (mr *UnsignedTxMockRecorder) SetBytes(unsignedBytes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBytes", reflect.TypeOf((*UnsignedTx)(nil).SetBytes), unsignedBytes)
}

// Visit mocks base method.
func (m *UnsignedTx) Visit(visitor txs.Visitor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Visit", visitor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Visit indicates an expected call of Visit.
func (mr *UnsignedTxMockRecorder) Visit(visitor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Visit", reflect.TypeOf((*UnsignedTx)(nil).Visit), visitor)
}