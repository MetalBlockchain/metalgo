// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/snow/uptime (interfaces: Calculator)
//
// Generated by this command:
//
//	mockgen -package=uptimemock -destination=snow/uptime/uptimemock/calculator.go -mock_names=Calculator=Calculator github.com/MetalBlockchain/metalgo/snow/uptime Calculator
//

// Package uptimemock is a generated GoMock package.
package uptimemock

import (
	reflect "reflect"
	time "time"

	ids "github.com/MetalBlockchain/metalgo/ids"
	gomock "go.uber.org/mock/gomock"
)

// Calculator is a mock of Calculator interface.
type Calculator struct {
	ctrl     *gomock.Controller
	recorder *CalculatorMockRecorder
}

// CalculatorMockRecorder is the mock recorder for Calculator.
type CalculatorMockRecorder struct {
	mock *Calculator
}

// NewCalculator creates a new mock instance.
func NewCalculator(ctrl *gomock.Controller) *Calculator {
	mock := &Calculator{ctrl: ctrl}
	mock.recorder = &CalculatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Calculator) EXPECT() *CalculatorMockRecorder {
	return m.recorder
}

// CalculateUptime mocks base method.
func (m *Calculator) CalculateUptime(arg0 ids.NodeID) (time.Duration, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateUptime", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CalculateUptime indicates an expected call of CalculateUptime.
func (mr *CalculatorMockRecorder) CalculateUptime(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateUptime", reflect.TypeOf((*Calculator)(nil).CalculateUptime), arg0)
}

// CalculateUptimePercent mocks base method.
func (m *Calculator) CalculateUptimePercent(arg0 ids.NodeID) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateUptimePercent", arg0)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateUptimePercent indicates an expected call of CalculateUptimePercent.
func (mr *CalculatorMockRecorder) CalculateUptimePercent(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateUptimePercent", reflect.TypeOf((*Calculator)(nil).CalculateUptimePercent), arg0)
}

// CalculateUptimePercentFrom mocks base method.
func (m *Calculator) CalculateUptimePercentFrom(arg0 ids.NodeID, arg1 time.Time) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateUptimePercentFrom", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateUptimePercentFrom indicates an expected call of CalculateUptimePercentFrom.
func (mr *CalculatorMockRecorder) CalculateUptimePercentFrom(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateUptimePercentFrom", reflect.TypeOf((*Calculator)(nil).CalculateUptimePercentFrom), arg0, arg1)
}
