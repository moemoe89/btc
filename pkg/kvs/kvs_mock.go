// Code generated by MockGen. DO NOT EDIT.
// Source: kvs.go

// Package kvs is a generated GoMock package.
package kvs

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// GoMockClient is a mock of Client interface.
type GoMockClient struct {
	ctrl     *gomock.Controller
	recorder *GoMockClientMockRecorder
}

// GoMockClientMockRecorder is the mock recorder for GoMockClient.
type GoMockClientMockRecorder struct {
	mock *GoMockClient
}

// NewGoMockClient creates a new mock instance.
func NewGoMockClient(ctrl *gomock.Controller) *GoMockClient {
	mock := &GoMockClient{ctrl: ctrl}
	mock.recorder = &GoMockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *GoMockClient) EXPECT() *GoMockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *GoMockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *GoMockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*GoMockClient)(nil).Close))
}

// Get mocks base method.
func (m *GoMockClient) Get(ctx context.Context, key string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *GoMockClientMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*GoMockClient)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *GoMockClient) Set(ctx context.Context, key string, value interface{}, expire time.Duration) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expire)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *GoMockClientMockRecorder) Set(ctx, key, value, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*GoMockClient)(nil).Set), ctx, key, value, expire)
}
