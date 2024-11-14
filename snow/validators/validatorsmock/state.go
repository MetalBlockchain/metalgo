// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/snow/validators (interfaces: State)
//
// Generated by this command:
//
//	mockgen -package=validatorsmock -destination=snow/validators/validatorsmock/state.go -mock_names=State=State github.com/MetalBlockchain/metalgo/snow/validators State
//

// Package validatorsmock is a generated GoMock package.
package validatorsmock

import (
	context "context"
	reflect "reflect"

	ids "github.com/MetalBlockchain/metalgo/ids"
	validators "github.com/MetalBlockchain/metalgo/snow/validators"
	gomock "go.uber.org/mock/gomock"
)

// State is a mock of State interface.
type State struct {
	ctrl     *gomock.Controller
	recorder *StateMockRecorder
}

// StateMockRecorder is the mock recorder for State.
type StateMockRecorder struct {
	mock *State
}

// NewState creates a new mock instance.
func NewState(ctrl *gomock.Controller) *State {
	mock := &State{ctrl: ctrl}
	mock.recorder = &StateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *State) EXPECT() *StateMockRecorder {
	return m.recorder
}

// GetCurrentHeight mocks base method.
func (m *State) GetCurrentHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentHeight indicates an expected call of GetCurrentHeight.
func (mr *StateMockRecorder) GetCurrentHeight(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentHeight", reflect.TypeOf((*State)(nil).GetCurrentHeight), arg0)
}

// GetMinimumHeight mocks base method.
func (m *State) GetMinimumHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinimumHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinimumHeight indicates an expected call of GetMinimumHeight.
func (mr *StateMockRecorder) GetMinimumHeight(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinimumHeight", reflect.TypeOf((*State)(nil).GetMinimumHeight), arg0)
}

// GetSubnetID mocks base method.
func (m *State) GetSubnetID(arg0 context.Context, arg1 ids.ID) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubnetID", arg0, arg1)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubnetID indicates an expected call of GetSubnetID.
func (mr *StateMockRecorder) GetSubnetID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubnetID", reflect.TypeOf((*State)(nil).GetSubnetID), arg0, arg1)
}

// GetValidatorSet mocks base method.
func (m *State) GetValidatorSet(arg0 context.Context, arg1 uint64, arg2 ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorSet", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[ids.NodeID]*validators.GetValidatorOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorSet indicates an expected call of GetValidatorSet.
func (mr *StateMockRecorder) GetValidatorSet(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorSet", reflect.TypeOf((*State)(nil).GetValidatorSet), arg0, arg1, arg2)
}