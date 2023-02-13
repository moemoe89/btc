package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errInternal = errors.New("error")

func TestBTCUC_CreateTransaction(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *repository.CreateTransactionParams
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.Transaction
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Create transaction, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.CreateTransactionParams{
					UserID:   1,
					Datetime: now,
					Amount:   100,
				},
			}

			want := &rpc.Transaction{
				UserId:   1,
				Datetime: timestamppb.New(now),
				Amount:   100,
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().CreateTransaction(args.ctx, args.params).Return(want, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Create transaction, When repository failed to executed, Return an error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.CreateTransactionParams{
					UserID:   1,
					Datetime: now,
					Amount:   100,
				},
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().CreateTransaction(args.ctx, args.params).Return(nil, errInternal)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    nil,
				wantErr: errInternal,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.CreateTransaction(tt.args.ctx, tt.args.params)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBTCUC_ListTransaction(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	type test struct {
		fields  fields
		args    args
		want    []*rpc.Transaction
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of List transactions, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			now := time.Now()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := []*rpc.Transaction{
				{
					UserId:   1,
					Datetime: timestamppb.New(now),
					Amount:   100,
				},
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.userID).Return(want, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List transaction, When repository failed to executed, Return an error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.userID).Return(nil, errInternal)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    nil,
				wantErr: errInternal,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.ListTransaction(tt.args.ctx, tt.args.userID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBTCUC_GetUserBalance(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.UserBalance
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Get User balance, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(want, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When repository failed to executed, Return an error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(nil, errInternal)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
				},
				args:    args,
				want:    nil,
				wantErr: errInternal,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.GetUserBalance(tt.args.ctx, tt.args.userID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
