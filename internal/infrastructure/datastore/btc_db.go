package datastore

import (
	"context"
	"errors"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"
)

type btcRepo struct {
	*BaseRepo
}

// NewBTCRepo returns BTCRepo.
func NewBTCRepo(base *BaseRepo) repository.BTCRepo {
	return &btcRepo{
		BaseRepo: base,
	}
}

// CreateTransaction creates a new record for BTC transaction.
// Only single transaction will create by this RPC for a specific User.
func (r *btcRepo) CreateTransaction(ctx context.Context, params repository.CreateTransactionParams) (*rpc.Transaction, error) {
	return nil, errors.New("unimplemented")
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (r *btcRepo) ListTransaction(ctx context.Context, userID int64) (*rpc.ListTransactionResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetUserBalance get the latest balance for a specific User.
func (r *btcRepo) GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error) {
	return &rpc.UserBalance{}, nil
}