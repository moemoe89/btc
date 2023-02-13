package datastore_test

import (
	"context"
	"testing"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/di"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/internal/infrastructure/datastore"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBTCRepo_CreateTransaction(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *repository.CreateTransactionParams
	}

	type test struct {
		args       args
		want       *rpc.Transaction
		wantErr    error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}

	db := datastore.GetDatabase()

	tests := map[string]func(t *testing.T) test{
		"Given valid query of Create transaction, When query executed successfully, Return no error": func(t *testing.T) test {
			userID := int64(1988)
			amount := float32(100.5)
			datetime := time.Now().UTC()

			args := args{
				ctx: context.Background(),
				params: &repository.CreateTransactionParams{
					UserID:   userID,
					Datetime: datetime,
					Amount:   amount,
				},
			}

			want := &rpc.Transaction{
				UserId:   userID,
				Datetime: timestamppb.New(datetime),
				Amount:   amount,
			}

			return test{
				args:    args,
				want:    want,
				wantErr: nil,
				beforeFunc: func(t *testing.T) {
					t.Helper()

					// Remove existing data, if any.
					_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)

					// Insert test data.
					_, err = db.Exec(context.Background(), "INSERT INTO users (id, balance) VALUES ($1, $2)", userID, 0)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T) {
					t.Helper()

					// Check accumulated balance after insert transaction.
					var balance float32

					err := db.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
					assert.NoError(t, err)
					assert.Equal(t, amount, balance)

					// Clear data.
					_, err = db.Exec(context.Background(), "DELETE FROM transactions WHERE user_id = $1", userID)
					assert.NoError(t, err)

					_, err = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)
				},
			}
		},
		"Given valid query of Create transaction, When query executed successfully with no User found, Return an error": func(t *testing.T) test {
			userID := int64(999)
			amount := float32(100.5)
			datetime := time.Now().UTC()

			args := args{
				ctx: context.Background(),
				params: &repository.CreateTransactionParams{
					UserID:   userID,
					Datetime: datetime,
					Amount:   amount,
				},
			}

			return test{
				args:    args,
				want:    nil,
				wantErr: datastore.ErrNotFound,
				beforeFunc: func(t *testing.T) {
					t.Helper()

					_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)
				},
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t)
			}

			if tt.afterFunc != nil {
				defer tt.afterFunc(t)
			}

			sut := di.GetBTCRepo()

			got, err := sut.CreateTransaction(tt.args.ctx, tt.args.params)

			if !assert.ErrorIs(t, err, tt.wantErr) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBTCRepo_ListTransaction(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	type test struct {
		args       args
		want       []*rpc.Transaction
		wantErr    error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}

	db := datastore.GetDatabase()

	tests := map[string]func(t *testing.T) test{
		"Given valid query of Get List transactions, When query executed successfully, Return no error": func(t *testing.T) test {
			userID := int64(1988)
			balance1 := float32(100.5)
			balance2 := float32(900.6)
			datetime1 := &timestamppb.Timestamp{
				Seconds: 1676252796,
				Nanos:   0,
			}
			datetime2 := &timestamppb.Timestamp{
				Seconds: 1676339196,
				Nanos:   0,
			}

			args := args{
				ctx:    context.Background(),
				userID: userID,
			}

			want := []*rpc.Transaction{
				{
					UserId:   userID,
					Datetime: datetime1,
					Amount:   balance1,
				},
				{
					UserId:   userID,
					Datetime: datetime2,
					Amount:   balance2,
				},
			}

			return test{
				args:    args,
				want:    want,
				wantErr: nil,
				beforeFunc: func(t *testing.T) {
					t.Helper()

					// Remove existing data, if any.
					_, err := db.Exec(context.Background(), "DELETE FROM transactions WHERE user_id = $1", userID)
					assert.NoError(t, err)

					_, err = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)

					// Insert test data.
					_, err = db.Exec(context.Background(), "INSERT INTO users (id, balance) VALUES ($1, $2)", userID, 0)
					assert.NoError(t, err)

					_, err = db.Exec(context.Background(), "INSERT INTO transactions (time, user_id, amount) VALUES ($1, $2, $3)", datetime1.AsTime(), userID, balance1)
					assert.NoError(t, err)

					_, err = db.Exec(context.Background(), "INSERT INTO transactions (time, user_id, amount) VALUES ($1, $2, $3)", datetime2.AsTime(), userID, balance2)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T) {
					t.Helper()

					// Clear data.
					_, err := db.Exec(context.Background(), "DELETE FROM transactions WHERE user_id = $1", userID)
					assert.NoError(t, err)

					_, err = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)
				},
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t)
			}

			if tt.afterFunc != nil {
				defer tt.afterFunc(t)
			}

			sut := di.GetBTCRepo()

			got, err := sut.ListTransaction(tt.args.ctx, tt.args.userID)

			if !assert.ErrorIs(t, err, tt.wantErr) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBTCRepo_GetUserBalance(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	type test struct {
		args       args
		want       *rpc.UserBalance
		wantErr    error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}

	db := datastore.GetDatabase()

	tests := map[string]func(t *testing.T) test{
		"Given valid query of Get User balance, When query executed successfully, Return no error": func(t *testing.T) test {
			userID := int64(1989)
			balance := float32(100.5)

			args := args{
				ctx:    context.Background(),
				userID: userID,
			}

			want := &rpc.UserBalance{
				Balance: balance,
			}

			return test{
				args:    args,
				want:    want,
				wantErr: nil,
				beforeFunc: func(t *testing.T) {
					t.Helper()

					// Remove existing data, if any.
					_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)

					// Insert test data.
					_, err = db.Exec(context.Background(), "INSERT INTO users (id, balance) VALUES ($1, $2)", userID, balance)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T) {
					t.Helper()

					// Clear data.
					_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)
				},
			}
		},
		"Given valid query of Get User balance, When query executed successfully with no User found, Return an error": func(t *testing.T) test {
			userID := int64(999)

			args := args{
				ctx:    context.Background(),
				userID: userID,
			}

			return test{
				args:    args,
				want:    nil,
				wantErr: datastore.ErrNotFound,
				beforeFunc: func(t *testing.T) {
					t.Helper()

					// Remove existing data, if any.
					_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
					assert.NoError(t, err)
				},
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t)
			}

			if tt.afterFunc != nil {
				defer tt.afterFunc(t)
			}

			sut := di.GetBTCRepo()

			got, err := sut.GetUserBalance(tt.args.ctx, tt.args.userID)

			if !assert.ErrorIs(t, err, tt.wantErr) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
