package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/pkg/kvs"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
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
		params *repository.ListTransactionParams
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.ListTransactionResponse
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Get User balance, When repository executed successfully without cache, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			userID := int64(1)
			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.ListTransactionParams{
					UserID:        userID,
					StartDatetime: now,
					EndDatetime:   now,
				},
			}

			transactions := []*rpc.Transaction{
				{
					UserId:   userID,
					Datetime: timestamppb.New(now),
					Amount:   100,
				},
			}

			want := &rpc.ListTransactionResponse{
				Transactions: transactions,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.params).Return(transactions, nil)

			key := fmt.Sprintf("user:transactions:%d:%d:%d",
				userID,
				now.UnixNano(),
				now.UnixNano(),
			)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When repository executed successfully without cache and failed store cache, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			userID := int64(1)
			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.ListTransactionParams{
					UserID:        userID,
					StartDatetime: now,
					EndDatetime:   now,
				},
			}

			transactions := []*rpc.Transaction{
				{
					UserId:   userID,
					Datetime: timestamppb.New(now),
					Amount:   100,
				},
			}

			want := &rpc.ListTransactionResponse{
				Transactions: transactions,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.params).Return(transactions, nil)

			key := fmt.Sprintf("user:transactions:%d:%d:%d",
				userID,
				now.UnixNano(),
				now.UnixNano(),
			)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When cache exists, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			userID := int64(1)
			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.ListTransactionParams{
					UserID:        userID,
					StartDatetime: now,
					EndDatetime:   now,
				},
			}

			transactions := []*rpc.Transaction{
				{
					UserId:   userID,
					Datetime: timestamppb.New(now),
					Amount:   100,
				},
			}

			want := &rpc.ListTransactionResponse{
				Transactions: transactions,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			key := fmt.Sprintf("user:transactions:%d:%d:%d",
				userID,
				now.UnixNano(),
				now.UnixNano(),
			)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(string(b), nil)

			return test{
				fields: fields{
					redis: redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When failed getting cache but repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			userID := int64(1)
			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.ListTransactionParams{
					UserID:        userID,
					StartDatetime: now,
					EndDatetime:   now,
				},
			}

			transactions := []*rpc.Transaction{
				{
					UserId:   userID,
					Datetime: timestamppb.New(now),
					Amount:   100,
				},
			}

			want := &rpc.ListTransactionResponse{
				Transactions: transactions,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.params).Return(transactions, nil)

			key := fmt.Sprintf("user:transactions:%d:%d:%d",
				userID,
				now.UnixNano(),
				now.UnixNano(),
			)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, errors.New("error"))
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List transaction, When repository failed to executed with no cache, Return an error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			userID := int64(1)
			now := time.Now()

			args := args{
				ctx: ctx,
				params: &repository.ListTransactionParams{
					UserID:        userID,
					StartDatetime: now,
					EndDatetime:   now,
				},
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().ListTransaction(args.ctx, args.params).Return(nil, errInternal)

			key := fmt.Sprintf("user:transactions:%d:%d:%d",
				userID,
				now.UnixNano(),
				now.UnixNano(),
			)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
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

			got, err := sut.ListTransaction(tt.args.ctx, tt.args.params)
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
		"Given valid request of Get User balance, When repository executed successfully without cache, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(want, nil)

			key := fmt.Sprintf("user:balance:%d", args.userID)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When repository executed successfully without cache and failed store cache, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(want, nil)

			key := fmt.Sprintf("user:balance:%d", args.userID)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When failed getting cache but repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(want, nil)

			key := fmt.Sprintf("user:balance:%d", args.userID)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, errors.New("error"))
			redisKVS.EXPECT().Set(ctx, key, b, time.Second*1).Return(nil, nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When cache exists, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			want := &rpc.UserBalance{
				Balance: 100,
			}

			b, err := protojson.Marshal(want)
			assert.NoError(t, err)

			key := fmt.Sprintf("user:balance:%d", args.userID)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(string(b), nil)

			return test{
				fields: fields{
					redis: redisKVS,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Get User balance, When repository failed to executed with no cache, Return an error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				userID: 1,
			}

			mockJourneyRepo := repository.NewGoMockBTCRepo(ctrl)
			mockJourneyRepo.EXPECT().GetUserBalance(args.ctx, args.userID).Return(nil, errInternal)

			key := fmt.Sprintf("user:balance:%d", args.userID)

			redisKVS := kvs.NewGoMockClient(ctrl)
			redisKVS.EXPECT().Get(ctx, key).Return(nil, redis.Nil)

			return test{
				fields: fields{
					btcRepo: mockJourneyRepo,
					redis:   redisKVS,
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
