// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/snow/networking/handler (interfaces: Handler)
//
// Generated by this command:
//
//	mockgen -package=handlermock -destination=snow/networking/handler/handlermock/handler.go -mock_names=Handler=Handler github.com/MetalBlockchain/metalgo/snow/networking/handler Handler
//

// Package handlermock is a generated GoMock package.
package handlermock

import (
	context "context"
	reflect "reflect"
	time "time"

	ids "github.com/MetalBlockchain/metalgo/ids"
	snow "github.com/MetalBlockchain/metalgo/snow"
	handler "github.com/MetalBlockchain/metalgo/snow/networking/handler"
	gomock "go.uber.org/mock/gomock"
)

// Handler is a mock of Handler interface.
type Handler struct {
	ctrl     *gomock.Controller
	recorder *HandlerMockRecorder
}

// HandlerMockRecorder is the mock recorder for Handler.
type HandlerMockRecorder struct {
	mock *Handler
}

// NewHandler creates a new mock instance.
func NewHandler(ctrl *gomock.Controller) *Handler {
	mock := &Handler{ctrl: ctrl}
	mock.recorder = &HandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Handler) EXPECT() *HandlerMockRecorder {
	return m.recorder
}

// AwaitStopped mocks base method.
func (m *Handler) AwaitStopped(arg0 context.Context) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AwaitStopped", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AwaitStopped indicates an expected call of AwaitStopped.
func (mr *HandlerMockRecorder) AwaitStopped(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AwaitStopped", reflect.TypeOf((*Handler)(nil).AwaitStopped), arg0)
}

// Context mocks base method.
func (m *Handler) Context() *snow.ConsensusContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(*snow.ConsensusContext)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *HandlerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*Handler)(nil).Context))
}

// GetEngineManager mocks base method.
func (m *Handler) GetEngineManager() *handler.EngineManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEngineManager")
	ret0, _ := ret[0].(*handler.EngineManager)
	return ret0
}

// GetEngineManager indicates an expected call of GetEngineManager.
func (mr *HandlerMockRecorder) GetEngineManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEngineManager", reflect.TypeOf((*Handler)(nil).GetEngineManager))
}

// HealthCheck mocks base method.
func (m *Handler) HealthCheck(arg0 context.Context) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *HandlerMockRecorder) HealthCheck(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*Handler)(nil).HealthCheck), arg0)
}

// Len mocks base method.
func (m *Handler) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *HandlerMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*Handler)(nil).Len))
}

// Push mocks base method.
func (m *Handler) Push(arg0 context.Context, arg1 handler.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Push", arg0, arg1)
}

// Push indicates an expected call of Push.
func (mr *HandlerMockRecorder) Push(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*Handler)(nil).Push), arg0, arg1)
}

// SetEngineManager mocks base method.
func (m *Handler) SetEngineManager(arg0 *handler.EngineManager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetEngineManager", arg0)
}

// SetEngineManager indicates an expected call of SetEngineManager.
func (mr *HandlerMockRecorder) SetEngineManager(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEngineManager", reflect.TypeOf((*Handler)(nil).SetEngineManager), arg0)
}

// SetOnStopped mocks base method.
func (m *Handler) SetOnStopped(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnStopped", arg0)
}

// SetOnStopped indicates an expected call of SetOnStopped.
func (mr *HandlerMockRecorder) SetOnStopped(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnStopped", reflect.TypeOf((*Handler)(nil).SetOnStopped), arg0)
}

// ShouldHandle mocks base method.
func (m *Handler) ShouldHandle(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldHandle", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ShouldHandle indicates an expected call of ShouldHandle.
func (mr *HandlerMockRecorder) ShouldHandle(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldHandle", reflect.TypeOf((*Handler)(nil).ShouldHandle), arg0)
}

// Start mocks base method.
func (m *Handler) Start(arg0 context.Context, arg1 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0, arg1)
}

// Start indicates an expected call of Start.
func (mr *HandlerMockRecorder) Start(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*Handler)(nil).Start), arg0, arg1)
}

// Stop mocks base method.
func (m *Handler) Stop(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *HandlerMockRecorder) Stop(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*Handler)(nil).Stop), arg0)
}

// StopWithError mocks base method.
func (m *Handler) StopWithError(arg0 context.Context, arg1 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopWithError", arg0, arg1)
}

// StopWithError indicates an expected call of StopWithError.
func (mr *HandlerMockRecorder) StopWithError(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopWithError", reflect.TypeOf((*Handler)(nil).StopWithError), arg0, arg1)
}
