package repository

import (
	"context"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
)

//go:generate rm -f ./btc_mock.go
//go:generate mockgen -destination btc_mock.go -package repository -mock_names BTCRepo=GoMockBTCRepo -source btc.go

// CreateTransactionParams parameter for creates a BTC transaction.
type CreateTransactionParams struct {
	UserID   int64     // required
	Datetime time.Time // required
	Amount   float64   // required
}

// ListTransactionParams parameter for lists a BTC transactions.
type ListTransactionParams struct {
	UserID        int64     // required
	StartDatetime time.Time // required
	EndDatetime   time.Time // required
}

// BTCRepo defines BTC repository.
type BTCRepo interface {
	// CreateTransaction creates a new record for BTC transaction.
	// Only single transaction will create by this RPC for a specific User.
	CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*rpc.Transaction, error)
	// ListTransaction get the list of records for BTC transaction.
	// The record can be filtered by specific User.
	ListTransaction(ctx context.Context, params *ListTransactionParams) ([]*rpc.Transaction, error)
	// GetUserBalance get the latest balance for a specific User.
	GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error)
}
