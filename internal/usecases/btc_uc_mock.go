// Code generated by MockGen. DO NOT EDIT.
// Source: btc_uc.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "github.com/moemoe89/btc/api/go/grpc"
	repository "github.com/moemoe89/btc/internal/entities/repository"
)

// GoMockBTCUsecase is a mock of BTCUsecase interface.
type GoMockBTCUsecase struct {
	ctrl     *gomock.Controller
	recorder *GoMockBTCUsecaseMockRecorder
}

// GoMockBTCUsecaseMockRecorder is the mock recorder for GoMockBTCUsecase.
type GoMockBTCUsecaseMockRecorder struct {
	mock *GoMockBTCUsecase
}

// NewGoMockBTCUsecase creates a new mock instance.
func NewGoMockBTCUsecase(ctrl *gomock.Controller) *GoMockBTCUsecase {
	mock := &GoMockBTCUsecase{ctrl: ctrl}
	mock.recorder = &GoMockBTCUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *GoMockBTCUsecase) EXPECT() *GoMockBTCUsecaseMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *GoMockBTCUsecase) CreateTransaction(ctx context.Context, params *repository.CreateTransactionParams) (*grpc.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, params)
	ret0, _ := ret[0].(*grpc.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *GoMockBTCUsecaseMockRecorder) CreateTransaction(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*GoMockBTCUsecase)(nil).CreateTransaction), ctx, params)
}

// GetUserBalance mocks base method.
func (m *GoMockBTCUsecase) GetUserBalance(ctx context.Context, userID int64) (*grpc.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", ctx, userID)
	ret0, _ := ret[0].(*grpc.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *GoMockBTCUsecaseMockRecorder) GetUserBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*GoMockBTCUsecase)(nil).GetUserBalance), ctx, userID)
}

// ListTransaction mocks base method.
func (m *GoMockBTCUsecase) ListTransaction(ctx context.Context, params *repository.ListTransactionParams) (*grpc.ListTransactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransaction", ctx, params)
	ret0, _ := ret[0].(*grpc.ListTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransaction indicates an expected call of ListTransaction.
func (mr *GoMockBTCUsecaseMockRecorder) ListTransaction(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransaction", reflect.TypeOf((*GoMockBTCUsecase)(nil).ListTransaction), ctx, params)
}
