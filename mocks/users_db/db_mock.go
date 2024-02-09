// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/users/users.go

// Package mock_users is a generated GoMock package.
package mock_users

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockdatabaseInterface is a mock of databaseInterface interface.
type MockdatabaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockdatabaseInterfaceMockRecorder
}

// MockdatabaseInterfaceMockRecorder is the mock recorder for MockdatabaseInterface.
type MockdatabaseInterfaceMockRecorder struct {
	mock *MockdatabaseInterface
}

// NewMockdatabaseInterface creates a new mock instance.
func NewMockdatabaseInterface(ctrl *gomock.Controller) *MockdatabaseInterface {
	mock := &MockdatabaseInterface{ctrl: ctrl}
	mock.recorder = &MockdatabaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdatabaseInterface) EXPECT() *MockdatabaseInterfaceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockdatabaseInterface) CreateUser() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockdatabaseInterfaceMockRecorder) CreateUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockdatabaseInterface)(nil).CreateUser))
}

// Use mocks base method.
func (m *MockdatabaseInterface) Use() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Use")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Use indicates an expected call of Use.
func (mr *MockdatabaseInterfaceMockRecorder) Use() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Use", reflect.TypeOf((*MockdatabaseInterface)(nil).Use))
}
