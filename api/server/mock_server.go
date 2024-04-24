// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/api/server (interfaces: Server)
//
// Generated by this command:
//
//	mockgen -package=server -destination=api/server/mock_server.go github.com/MetalBlockchain/metalgo/api/server Server
//

// Package server is a generated GoMock package.
package server

import (
	http "net/http"
	reflect "reflect"

	snow "github.com/MetalBlockchain/metalgo/snow"
	common "github.com/MetalBlockchain/metalgo/snow/engine/common"
	gomock "go.uber.org/mock/gomock"
)

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// AddAliases mocks base method.
func (m *MockServer) AddAliases(arg0 string, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddAliases", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAliases indicates an expected call of AddAliases.
func (mr *MockServerMockRecorder) AddAliases(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAliases", reflect.TypeOf((*MockServer)(nil).AddAliases), varargs...)
}

// AddAliasesWithReadLock mocks base method.
func (m *MockServer) AddAliasesWithReadLock(arg0 string, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddAliasesWithReadLock", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAliasesWithReadLock indicates an expected call of AddAliasesWithReadLock.
func (mr *MockServerMockRecorder) AddAliasesWithReadLock(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAliasesWithReadLock", reflect.TypeOf((*MockServer)(nil).AddAliasesWithReadLock), varargs...)
}

// AddRoute mocks base method.
func (m *MockServer) AddRoute(arg0 http.Handler, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRoute", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRoute indicates an expected call of AddRoute.
func (mr *MockServerMockRecorder) AddRoute(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRoute", reflect.TypeOf((*MockServer)(nil).AddRoute), arg0, arg1, arg2)
}

// AddRouteWithReadLock mocks base method.
func (m *MockServer) AddRouteWithReadLock(arg0 http.Handler, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRouteWithReadLock", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRouteWithReadLock indicates an expected call of AddRouteWithReadLock.
func (mr *MockServerMockRecorder) AddRouteWithReadLock(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRouteWithReadLock", reflect.TypeOf((*MockServer)(nil).AddRouteWithReadLock), arg0, arg1, arg2)
}

// Dispatch mocks base method.
func (m *MockServer) Dispatch() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dispatch")
	ret0, _ := ret[0].(error)
	return ret0
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockServerMockRecorder) Dispatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockServer)(nil).Dispatch))
}

// RegisterChain mocks base method.
func (m *MockServer) RegisterChain(arg0 string, arg1 *snow.ConsensusContext, arg2 common.VM) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterChain", arg0, arg1, arg2)
}

// RegisterChain indicates an expected call of RegisterChain.
func (mr *MockServerMockRecorder) RegisterChain(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChain", reflect.TypeOf((*MockServer)(nil).RegisterChain), arg0, arg1, arg2)
}

// Shutdown mocks base method.
func (m *MockServer) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockServerMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockServer)(nil).Shutdown))
}
