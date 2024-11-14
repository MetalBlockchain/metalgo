// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/vms/avm/txs/mempool (interfaces: Mempool)
//
// Generated by this command:
//
//	mockgen -package=mempoolmock -destination=vms/avm/txs/mempool/mempoolmock/mempool.go -mock_names=Mempool=Mempool github.com/MetalBlockchain/metalgo/vms/avm/txs/mempool Mempool
//

// Package mempoolmock is a generated GoMock package.
package mempoolmock

import (
	reflect "reflect"

	ids "github.com/MetalBlockchain/metalgo/ids"
	txs "github.com/MetalBlockchain/metalgo/vms/avm/txs"
	gomock "go.uber.org/mock/gomock"
)

// Mempool is a mock of Mempool interface.
type Mempool struct {
	ctrl     *gomock.Controller
	recorder *MempoolMockRecorder
}

// MempoolMockRecorder is the mock recorder for Mempool.
type MempoolMockRecorder struct {
	mock *Mempool
}

// NewMempool creates a new mock instance.
func NewMempool(ctrl *gomock.Controller) *Mempool {
	mock := &Mempool{ctrl: ctrl}
	mock.recorder = &MempoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mempool) EXPECT() *MempoolMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *Mempool) Add(arg0 *txs.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MempoolMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*Mempool)(nil).Add), arg0)
}

// Get mocks base method.
func (m *Mempool) Get(arg0 ids.ID) (*txs.Tx, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MempoolMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Mempool)(nil).Get), arg0)
}

// GetDropReason mocks base method.
func (m *Mempool) GetDropReason(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDropReason", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDropReason indicates an expected call of GetDropReason.
func (mr *MempoolMockRecorder) GetDropReason(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDropReason", reflect.TypeOf((*Mempool)(nil).GetDropReason), arg0)
}

// Iterate mocks base method.
func (m *Mempool) Iterate(arg0 func(*txs.Tx) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Iterate", arg0)
}

// Iterate indicates an expected call of Iterate.
func (mr *MempoolMockRecorder) Iterate(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iterate", reflect.TypeOf((*Mempool)(nil).Iterate), arg0)
}

// Len mocks base method.
func (m *Mempool) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MempoolMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*Mempool)(nil).Len))
}

// MarkDropped mocks base method.
func (m *Mempool) MarkDropped(arg0 ids.ID, arg1 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkDropped", arg0, arg1)
}

// MarkDropped indicates an expected call of MarkDropped.
func (mr *MempoolMockRecorder) MarkDropped(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDropped", reflect.TypeOf((*Mempool)(nil).MarkDropped), arg0, arg1)
}

// Peek mocks base method.
func (m *Mempool) Peek() (*txs.Tx, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peek")
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Peek indicates an expected call of Peek.
func (mr *MempoolMockRecorder) Peek() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peek", reflect.TypeOf((*Mempool)(nil).Peek))
}

// Remove mocks base method.
func (m *Mempool) Remove(arg0 ...*txs.Tx) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Remove", varargs...)
}

// Remove indicates an expected call of Remove.
func (mr *MempoolMockRecorder) Remove(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*Mempool)(nil).Remove), arg0...)
}

// RequestBuildBlock mocks base method.
func (m *Mempool) RequestBuildBlock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestBuildBlock")
}

// RequestBuildBlock indicates an expected call of RequestBuildBlock.
func (mr *MempoolMockRecorder) RequestBuildBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestBuildBlock", reflect.TypeOf((*Mempool)(nil).RequestBuildBlock))
}