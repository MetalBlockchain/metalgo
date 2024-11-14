// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/snow/networking/timeout (interfaces: Manager)
//
// Generated by this command:
//
//	mockgen -package=timeoutmock -destination=snow/networking/timeout/timeoutmock/manager.go -mock_names=Manager=Manager github.com/MetalBlockchain/metalgo/snow/networking/timeout Manager
//

// Package timeoutmock is a generated GoMock package.
package timeoutmock

import (
	reflect "reflect"
	time "time"

	ids "github.com/MetalBlockchain/metalgo/ids"
	message "github.com/MetalBlockchain/metalgo/message"
	snow "github.com/MetalBlockchain/metalgo/snow"
	gomock "go.uber.org/mock/gomock"
)

// Manager is a mock of Manager interface.
type Manager struct {
	ctrl     *gomock.Controller
	recorder *ManagerMockRecorder
}

// ManagerMockRecorder is the mock recorder for Manager.
type ManagerMockRecorder struct {
	mock *Manager
}

// NewManager creates a new mock instance.
func NewManager(ctrl *gomock.Controller) *Manager {
	mock := &Manager{ctrl: ctrl}
	mock.recorder = &ManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Manager) EXPECT() *ManagerMockRecorder {
	return m.recorder
}

// Dispatch mocks base method.
func (m *Manager) Dispatch() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch")
}

// Dispatch indicates an expected call of Dispatch.
func (mr *ManagerMockRecorder) Dispatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*Manager)(nil).Dispatch))
}

// IsBenched mocks base method.
func (m *Manager) IsBenched(arg0 ids.NodeID, arg1 ids.ID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBenched", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBenched indicates an expected call of IsBenched.
func (mr *ManagerMockRecorder) IsBenched(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBenched", reflect.TypeOf((*Manager)(nil).IsBenched), arg0, arg1)
}

// RegisterChain mocks base method.
func (m *Manager) RegisterChain(arg0 *snow.ConsensusContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterChain", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterChain indicates an expected call of RegisterChain.
func (mr *ManagerMockRecorder) RegisterChain(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChain", reflect.TypeOf((*Manager)(nil).RegisterChain), arg0)
}

// RegisterRequest mocks base method.
func (m *Manager) RegisterRequest(arg0 ids.NodeID, arg1 ids.ID, arg2 bool, arg3 ids.RequestID, arg4 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterRequest", arg0, arg1, arg2, arg3, arg4)
}

// RegisterRequest indicates an expected call of RegisterRequest.
func (mr *ManagerMockRecorder) RegisterRequest(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRequest", reflect.TypeOf((*Manager)(nil).RegisterRequest), arg0, arg1, arg2, arg3, arg4)
}

// RegisterRequestToUnreachableValidator mocks base method.
func (m *Manager) RegisterRequestToUnreachableValidator() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterRequestToUnreachableValidator")
}

// RegisterRequestToUnreachableValidator indicates an expected call of RegisterRequestToUnreachableValidator.
func (mr *ManagerMockRecorder) RegisterRequestToUnreachableValidator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRequestToUnreachableValidator", reflect.TypeOf((*Manager)(nil).RegisterRequestToUnreachableValidator))
}

// RegisterResponse mocks base method.
func (m *Manager) RegisterResponse(arg0 ids.NodeID, arg1 ids.ID, arg2 ids.RequestID, arg3 message.Op, arg4 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterResponse", arg0, arg1, arg2, arg3, arg4)
}

// RegisterResponse indicates an expected call of RegisterResponse.
func (mr *ManagerMockRecorder) RegisterResponse(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterResponse", reflect.TypeOf((*Manager)(nil).RegisterResponse), arg0, arg1, arg2, arg3, arg4)
}

// RemoveRequest mocks base method.
func (m *Manager) RemoveRequest(arg0 ids.RequestID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveRequest", arg0)
}

// RemoveRequest indicates an expected call of RemoveRequest.
func (mr *ManagerMockRecorder) RemoveRequest(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRequest", reflect.TypeOf((*Manager)(nil).RemoveRequest), arg0)
}

// Stop mocks base method.
func (m *Manager) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *ManagerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*Manager)(nil).Stop))
}

// TimeoutDuration mocks base method.
func (m *Manager) TimeoutDuration() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeoutDuration")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// TimeoutDuration indicates an expected call of TimeoutDuration.
func (mr *ManagerMockRecorder) TimeoutDuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeoutDuration", reflect.TypeOf((*Manager)(nil).TimeoutDuration))
}