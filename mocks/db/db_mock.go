// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/db/db.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

// MockLoggerer is a mock of Loggerer interface.
type MockLoggerer struct {
	ctrl     *gomock.Controller
	recorder *MockLoggererMockRecorder
}

// MockLoggererMockRecorder is the mock recorder for MockLoggerer.
type MockLoggererMockRecorder struct {
	mock *MockLoggerer
}

// NewMockLoggerer creates a new mock instance.
func NewMockLoggerer(ctrl *gomock.Controller) *MockLoggerer {
	mock := &MockLoggerer{ctrl: ctrl}
	mock.recorder = &MockLoggererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerer) EXPECT() *MockLoggererMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockLoggerer) Error(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggererMockRecorder) Error(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLoggerer)(nil).Error), args...)
}

// Warnln mocks base method.
func (m *MockLoggerer) Warnln(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnln", varargs...)
}

// Warnln indicates an expected call of Warnln.
func (mr *MockLoggererMockRecorder) Warnln(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnln", reflect.TypeOf((*MockLoggerer)(nil).Warnln), args...)
}

// MockPgxIface is a mock of PgxIface interface.
type MockPgxIface struct {
	ctrl     *gomock.Controller
	recorder *MockPgxIfaceMockRecorder
}

// MockPgxIfaceMockRecorder is the mock recorder for MockPgxIface.
type MockPgxIfaceMockRecorder struct {
	mock *MockPgxIface
}

// NewMockPgxIface creates a new mock instance.
func NewMockPgxIface(ctrl *gomock.Controller) *MockPgxIface {
	mock := &MockPgxIface{ctrl: ctrl}
	mock.recorder = &MockPgxIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgxIface) EXPECT() *MockPgxIfaceMockRecorder {
	return m.recorder
}

// Acquire mocks base method.
func (m *MockPgxIface) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx)
	ret0, _ := ret[0].(*pgxpool.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Acquire indicates an expected call of Acquire.
func (mr *MockPgxIfaceMockRecorder) Acquire(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockPgxIface)(nil).Acquire), ctx)
}

// Close mocks base method.
func (m *MockPgxIface) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockPgxIfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPgxIface)(nil).Close))
}

// Config mocks base method.
func (m *MockPgxIface) Config() *pgxpool.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*pgxpool.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockPgxIfaceMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockPgxIface)(nil).Config))
}

// Exec mocks base method.
func (m *MockPgxIface) Exec(arg0 context.Context, arg1 string, arg2 ...interface{}) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockPgxIfaceMockRecorder) Exec(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockPgxIface)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *MockPgxIface) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPgxIfaceMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPgxIface)(nil).Ping), arg0)
}

// Query mocks base method.
func (m *MockPgxIface) Query(arg0 context.Context, arg1 string, arg2 ...interface{}) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPgxIfaceMockRecorder) Query(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPgxIface)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockPgxIface) QueryRow(arg0 context.Context, arg1 string, arg2 ...interface{}) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockPgxIfaceMockRecorder) QueryRow(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockPgxIface)(nil).QueryRow), varargs...)
}

// Reset mocks base method.
func (m *MockPgxIface) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockPgxIfaceMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockPgxIface)(nil).Reset))
}

// Stat mocks base method.
func (m *MockPgxIface) Stat() *pgxpool.Stat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(*pgxpool.Stat)
	return ret0
}

// Stat indicates an expected call of Stat.
func (mr *MockPgxIfaceMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockPgxIface)(nil).Stat))
}
