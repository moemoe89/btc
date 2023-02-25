package grpchandler_test

import (
	"context"
	"errors"
	"testing"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/adapters/grpchandler"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/internal/usecases"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fields struct {
	uc usecases.BTCUsecase
}

func sut(f fields) grpchandler.BTCServiceServer {
	return grpchandler.NewBTCHandler(
		f.uc,
	)
}

func TestBTCServer_CreateTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.CreateTransactionRequest
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.Transaction
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Create Transaction, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.CreateTransactionRequest{
					UserId: 1,
					Datetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
					Amount: 100,
				},
			}

			want := &rpc.Transaction{
				UserId:   args.req.UserId,
				Datetime: args.req.Datetime,
				Amount:   args.req.Amount,
			}

			params := &repository.CreateTransactionParams{
				UserID:   args.req.UserId,
				Datetime: args.req.Datetime.AsTime(),
				Amount:   args.req.Amount,
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().CreateTransaction(args.ctx, params).Return(want, nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Create Transaction, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.CreateTransactionRequest{
					UserId: 1,
					Datetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
					Amount: 100,
				},
			}

			params := &repository.CreateTransactionParams{
				UserID:   args.req.UserId,
				Datetime: args.req.Datetime.AsTime(),
				Amount:   args.req.Amount,
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().CreateTransaction(args.ctx, params).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.CreateTransaction(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBTCServer_ListTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.ListTransactionRequest
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.ListTransactionResponse
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of List Transaction, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.ListTransactionRequest{
					UserId: 1,
					StartDatetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
					EndDatetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
				},
			}

			want := &rpc.ListTransactionResponse{
				Transactions: []*rpc.Transaction{
					{
						UserId:   args.req.UserId,
						Datetime: args.req.StartDatetime,
						Amount:   100,
					},
				},
			}

			params := &repository.ListTransactionParams{
				UserID:        args.req.UserId,
				StartDatetime: args.req.StartDatetime.AsTime(),
				EndDatetime:   args.req.EndDatetime.AsTime(),
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().ListTransaction(args.ctx, params).Return(want, nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List Transaction, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.ListTransactionRequest{
					UserId: 1,
					StartDatetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
					EndDatetime: &timestamppb.Timestamp{
						Seconds: 1676169338,
						Nanos:   0,
					},
				},
			}

			params := &repository.ListTransactionParams{
				UserID:        args.req.UserId,
				StartDatetime: args.req.StartDatetime.AsTime(),
				EndDatetime:   args.req.EndDatetime.AsTime(),
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().ListTransaction(args.ctx, params).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.ListTransaction(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBTCServer_GetUserBalance(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.GetUserBalanceRequest
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.UserBalance
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Get User Balance, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.GetUserBalanceRequest{
					UserId: 1,
				},
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().GetUserBalance(args.ctx, args.req.GetUserId()).Return(want, nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List Get User Balance, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.GetUserBalanceRequest{
					UserId: 1,
				},
			}

			ucMock := usecases.NewGoMockBTCUsecase(ctrl)
			ucMock.EXPECT().GetUserBalance(args.ctx, args.req.GetUserId()).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.GetUserBalance(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
