package usecases

import (
	"context"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"
)

// CreateTransaction creates a new record for BTC transaction.
// Only single transaction will create by this RPC for a specific User.
func (u *btcUsecase) CreateTransaction(ctx context.Context, params repository.CreateTransactionParams) (*rpc.Transaction, error) {
	return nil, nil
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (u *btcUsecase) ListTransaction(ctx context.Context, userID int64) (*rpc.ListTransactionResponse, error) {
	return nil, nil
}

// GetUserBalance get the latest balance for a specific User.
func (u *btcUsecase) GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error) {
	return nil, nil
}
