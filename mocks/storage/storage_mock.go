// Code generated by MockGen. DO NOT EDIT.
// Source: ./util/util.go

// Package mock_util is a generated GoMock package.
package mock_util

import (
	http "net/http"
	reflect "reflect"

	util "github.com/chemax/url-shorter/util"
	gomock "github.com/golang/mock/gomock"
)

// MockDBInterface is a mock of DBInterface interface.
type MockDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDBInterfaceMockRecorder
}

// MockDBInterfaceMockRecorder is the mock recorder for MockDBInterface.
type MockDBInterfaceMockRecorder struct {
	mock *MockDBInterface
}

// NewMockDBInterface creates a new mock instance.
func NewMockDBInterface(ctrl *gomock.Controller) *MockDBInterface {
	mock := &MockDBInterface{ctrl: ctrl}
	mock.recorder = &MockDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBInterface) EXPECT() *MockDBInterfaceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockDBInterface) Get(code string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", code)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDBInterfaceMockRecorder) Get(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDBInterface)(nil).Get), code)
}

// Ping mocks base method.
func (m *MockDBInterface) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDBInterfaceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDBInterface)(nil).Ping))
}

// SaveURL mocks base method.
func (m *MockDBInterface) SaveURL(code string, URL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", code, URL)
	ret0, _ := ret[0].(error)
	return "", ret0
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockDBInterfaceMockRecorder) SaveURL(code, URL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockDBInterface)(nil).SaveURL), code, URL)
}

// Use mocks base method.
func (m *MockDBInterface) Use() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Use")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Use indicates an expected call of Use.
func (mr *MockDBInterfaceMockRecorder) Use() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Use", reflect.TypeOf((*MockDBInterface)(nil).Use))
}

// MockLoggerInterface is a mock of LoggerInterface interface.
type MockLoggerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerInterfaceMockRecorder
}

// MockLoggerInterfaceMockRecorder is the mock recorder for MockLoggerInterface.
type MockLoggerInterfaceMockRecorder struct {
	mock *MockLoggerInterface
}

// NewMockLoggerInterface creates a new mock instance.
func NewMockLoggerInterface(ctrl *gomock.Controller) *MockLoggerInterface {
	mock := &MockLoggerInterface{ctrl: ctrl}
	mock.recorder = &MockLoggerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerInterface) EXPECT() *MockLoggerInterfaceMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLoggerInterface) Debug(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerInterfaceMockRecorder) Debug(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLoggerInterface)(nil).Debug), args...)
}

// Error mocks base method.
func (m *MockLoggerInterface) Error(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerInterfaceMockRecorder) Error(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLoggerInterface)(nil).Error), args...)
}

// Info mocks base method.
func (m *MockLoggerInterface) Info(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerInterfaceMockRecorder) Info(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLoggerInterface)(nil).Info), args...)
}

// Middleware mocks base method.
func (m *MockLoggerInterface) Middleware(next http.Handler) http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Middleware", next)
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// Middleware indicates an expected call of Middleware.
func (mr *MockLoggerInterfaceMockRecorder) Middleware(next interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Middleware", reflect.TypeOf((*MockLoggerInterface)(nil).Middleware), next)
}

// Warn mocks base method.
func (m *MockLoggerInterface) Warn(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerInterfaceMockRecorder) Warn(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLoggerInterface)(nil).Warn), args...)
}

// MockStorageInterface is a mock of StorageInterface interface.
type MockStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStorageInterfaceMockRecorder
}

// MockStorageInterfaceMockRecorder is the mock recorder for MockStorageInterface.
type MockStorageInterfaceMockRecorder struct {
	mock *MockStorageInterface
}

// NewMockStorageInterface creates a new mock instance.
func NewMockStorageInterface(ctrl *gomock.Controller) *MockStorageInterface {
	mock := &MockStorageInterface{ctrl: ctrl}
	mock.recorder = &MockStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageInterface) EXPECT() *MockStorageInterfaceMockRecorder {
	return m.recorder
}

// AddNewURL mocks base method.
func (m *MockStorageInterface) AddNewURL(parsedURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewURL", parsedURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNewURL indicates an expected call of AddNewURL.
func (mr *MockStorageInterfaceMockRecorder) AddNewURL(parsedURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewURL", reflect.TypeOf((*MockStorageInterface)(nil).AddNewURL), parsedURL)
}

// BatchSave mocks base method.
func (m *MockStorageInterface) BatchSave(arr []*util.URLStructForBatch, httpPrefix string) ([]util.URLStructForBatchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSave", arr, httpPrefix)
	ret0, _ := ret[0].([]util.URLStructForBatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchSave indicates an expected call of BatchSave.
func (mr *MockStorageInterfaceMockRecorder) BatchSave(arr, httpPrefix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSave", reflect.TypeOf((*MockStorageInterface)(nil).BatchSave), arr, httpPrefix)
}

// GetURL mocks base method.
func (m *MockStorageInterface) GetURL(code string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", code)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockStorageInterfaceMockRecorder) GetURL(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockStorageInterface)(nil).GetURL), code)
}

// Ping mocks base method.
func (m *MockStorageInterface) Ping() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageInterfaceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorageInterface)(nil).Ping))
}

// MockConfigInterface is a mock of ConfigInterface interface.
type MockConfigInterface struct {
	ctrl     *gomock.Controller
	recorder *MockConfigInterfaceMockRecorder
}

// MockConfigInterfaceMockRecorder is the mock recorder for MockConfigInterface.
type MockConfigInterfaceMockRecorder struct {
	mock *MockConfigInterface
}

// NewMockConfigInterface creates a new mock instance.
func NewMockConfigInterface(ctrl *gomock.Controller) *MockConfigInterface {
	mock := &MockConfigInterface{ctrl: ctrl}
	mock.recorder = &MockConfigInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigInterface) EXPECT() *MockConfigInterfaceMockRecorder {
	return m.recorder
}

// GetHTTPAddr mocks base method.
func (m *MockConfigInterface) GetHTTPAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHTTPAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHTTPAddr indicates an expected call of GetHTTPAddr.
func (mr *MockConfigInterfaceMockRecorder) GetHTTPAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTTPAddr", reflect.TypeOf((*MockConfigInterface)(nil).GetHTTPAddr))
}

// GetNetAddr mocks base method.
func (m *MockConfigInterface) GetNetAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNetAddr indicates an expected call of GetNetAddr.
func (mr *MockConfigInterfaceMockRecorder) GetNetAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetAddr", reflect.TypeOf((*MockConfigInterface)(nil).GetNetAddr))
}
