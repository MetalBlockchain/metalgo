// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MetalBlockchain/metalgo/snow/engine/avalanche/vertex (interfaces: LinearizableVM)
//
// Generated by this command:
//
//	mockgen -package=vertexmock -destination=snow/engine/avalanche/vertex/vertexmock/linearizable_vm.go -mock_names=LinearizableVM=LinearizableVM github.com/MetalBlockchain/metalgo/snow/engine/avalanche/vertex LinearizableVM
//

// Package vertexmock is a generated GoMock package.
package vertexmock

import (
	context "context"
	http "net/http"
	reflect "reflect"
	time "time"

	database "github.com/MetalBlockchain/metalgo/database"
	ids "github.com/MetalBlockchain/metalgo/ids"
	snow "github.com/MetalBlockchain/metalgo/snow"
	snowman "github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	snowstorm "github.com/MetalBlockchain/metalgo/snow/consensus/snowstorm"
	common "github.com/MetalBlockchain/metalgo/snow/engine/common"
	version "github.com/MetalBlockchain/metalgo/version"
	gomock "go.uber.org/mock/gomock"
)

// LinearizableVM is a mock of LinearizableVM interface.
type LinearizableVM struct {
	ctrl     *gomock.Controller
	recorder *LinearizableVMMockRecorder
}

// LinearizableVMMockRecorder is the mock recorder for LinearizableVM.
type LinearizableVMMockRecorder struct {
	mock *LinearizableVM
}

// NewLinearizableVM creates a new mock instance.
func NewLinearizableVM(ctrl *gomock.Controller) *LinearizableVM {
	mock := &LinearizableVM{ctrl: ctrl}
	mock.recorder = &LinearizableVMMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *LinearizableVM) EXPECT() *LinearizableVMMockRecorder {
	return m.recorder
}

// AppGossip mocks base method.
func (m *LinearizableVM) AppGossip(arg0 context.Context, arg1 ids.NodeID, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppGossip", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppGossip indicates an expected call of AppGossip.
func (mr *LinearizableVMMockRecorder) AppGossip(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppGossip", reflect.TypeOf((*LinearizableVM)(nil).AppGossip), arg0, arg1, arg2)
}

// AppRequest mocks base method.
func (m *LinearizableVM) AppRequest(arg0 context.Context, arg1 ids.NodeID, arg2 uint32, arg3 time.Time, arg4 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppRequest", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppRequest indicates an expected call of AppRequest.
func (mr *LinearizableVMMockRecorder) AppRequest(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppRequest", reflect.TypeOf((*LinearizableVM)(nil).AppRequest), arg0, arg1, arg2, arg3, arg4)
}

// AppRequestFailed mocks base method.
func (m *LinearizableVM) AppRequestFailed(arg0 context.Context, arg1 ids.NodeID, arg2 uint32, arg3 *common.AppError) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppRequestFailed", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppRequestFailed indicates an expected call of AppRequestFailed.
func (mr *LinearizableVMMockRecorder) AppRequestFailed(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppRequestFailed", reflect.TypeOf((*LinearizableVM)(nil).AppRequestFailed), arg0, arg1, arg2, arg3)
}

// AppResponse mocks base method.
func (m *LinearizableVM) AppResponse(arg0 context.Context, arg1 ids.NodeID, arg2 uint32, arg3 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppResponse", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppResponse indicates an expected call of AppResponse.
func (mr *LinearizableVMMockRecorder) AppResponse(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppResponse", reflect.TypeOf((*LinearizableVM)(nil).AppResponse), arg0, arg1, arg2, arg3)
}

// BuildBlock mocks base method.
func (m *LinearizableVM) BuildBlock(arg0 context.Context) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildBlock", arg0)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildBlock indicates an expected call of BuildBlock.
func (mr *LinearizableVMMockRecorder) BuildBlock(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildBlock", reflect.TypeOf((*LinearizableVM)(nil).BuildBlock), arg0)
}

// Connected mocks base method.
func (m *LinearizableVM) Connected(arg0 context.Context, arg1 ids.NodeID, arg2 *version.Application) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connected", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connected indicates an expected call of Connected.
func (mr *LinearizableVMMockRecorder) Connected(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connected", reflect.TypeOf((*LinearizableVM)(nil).Connected), arg0, arg1, arg2)
}

// CreateHandlers mocks base method.
func (m *LinearizableVM) CreateHandlers(arg0 context.Context) (map[string]http.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHandlers", arg0)
	ret0, _ := ret[0].(map[string]http.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHandlers indicates an expected call of CreateHandlers.
func (mr *LinearizableVMMockRecorder) CreateHandlers(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHandlers", reflect.TypeOf((*LinearizableVM)(nil).CreateHandlers), arg0)
}

// Disconnected mocks base method.
func (m *LinearizableVM) Disconnected(arg0 context.Context, arg1 ids.NodeID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnected", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnected indicates an expected call of Disconnected.
func (mr *LinearizableVMMockRecorder) Disconnected(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnected", reflect.TypeOf((*LinearizableVM)(nil).Disconnected), arg0, arg1)
}

// GetBlock mocks base method.
func (m *LinearizableVM) GetBlock(arg0 context.Context, arg1 ids.ID) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", arg0, arg1)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *LinearizableVMMockRecorder) GetBlock(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*LinearizableVM)(nil).GetBlock), arg0, arg1)
}

// GetBlockIDAtHeight mocks base method.
func (m *LinearizableVM) GetBlockIDAtHeight(arg0 context.Context, arg1 uint64) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockIDAtHeight", arg0, arg1)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockIDAtHeight indicates an expected call of GetBlockIDAtHeight.
func (mr *LinearizableVMMockRecorder) GetBlockIDAtHeight(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockIDAtHeight", reflect.TypeOf((*LinearizableVM)(nil).GetBlockIDAtHeight), arg0, arg1)
}

// HealthCheck mocks base method.
func (m *LinearizableVM) HealthCheck(arg0 context.Context) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *LinearizableVMMockRecorder) HealthCheck(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*LinearizableVM)(nil).HealthCheck), arg0)
}

// Initialize mocks base method.
func (m *LinearizableVM) Initialize(arg0 context.Context, arg1 *snow.Context, arg2 database.Database, arg3, arg4, arg5 []byte, arg6 chan<- common.Message, arg7 []*common.Fx, arg8 common.AppSender) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *LinearizableVMMockRecorder) Initialize(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*LinearizableVM)(nil).Initialize), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
}

// LastAccepted mocks base method.
func (m *LinearizableVM) LastAccepted(arg0 context.Context) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastAccepted", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastAccepted indicates an expected call of LastAccepted.
func (mr *LinearizableVMMockRecorder) LastAccepted(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastAccepted", reflect.TypeOf((*LinearizableVM)(nil).LastAccepted), arg0)
}

// Linearize mocks base method.
func (m *LinearizableVM) Linearize(arg0 context.Context, arg1 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Linearize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Linearize indicates an expected call of Linearize.
func (mr *LinearizableVMMockRecorder) Linearize(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Linearize", reflect.TypeOf((*LinearizableVM)(nil).Linearize), arg0, arg1)
}

// ParseBlock mocks base method.
func (m *LinearizableVM) ParseBlock(arg0 context.Context, arg1 []byte) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBlock", arg0, arg1)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseBlock indicates an expected call of ParseBlock.
func (mr *LinearizableVMMockRecorder) ParseBlock(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBlock", reflect.TypeOf((*LinearizableVM)(nil).ParseBlock), arg0, arg1)
}

// ParseTx mocks base method.
func (m *LinearizableVM) ParseTx(arg0 context.Context, arg1 []byte) (snowstorm.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseTx", arg0, arg1)
	ret0, _ := ret[0].(snowstorm.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseTx indicates an expected call of ParseTx.
func (mr *LinearizableVMMockRecorder) ParseTx(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseTx", reflect.TypeOf((*LinearizableVM)(nil).ParseTx), arg0, arg1)
}

// SetPreference mocks base method.
func (m *LinearizableVM) SetPreference(arg0 context.Context, arg1 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPreference", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPreference indicates an expected call of SetPreference.
func (mr *LinearizableVMMockRecorder) SetPreference(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPreference", reflect.TypeOf((*LinearizableVM)(nil).SetPreference), arg0, arg1)
}

// SetState mocks base method.
func (m *LinearizableVM) SetState(arg0 context.Context, arg1 snow.State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState.
func (mr *LinearizableVMMockRecorder) SetState(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*LinearizableVM)(nil).SetState), arg0, arg1)
}

// Shutdown mocks base method.
func (m *LinearizableVM) Shutdown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *LinearizableVMMockRecorder) Shutdown(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*LinearizableVM)(nil).Shutdown), arg0)
}

// Version mocks base method.
func (m *LinearizableVM) Version(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Version indicates an expected call of Version.
func (mr *LinearizableVMMockRecorder) Version(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*LinearizableVM)(nil).Version), arg0)
}
