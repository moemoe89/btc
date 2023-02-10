package grpchandler

import (
	"context"
	"errors"

	rpc "github.com/moemoe89/btc/api/go/grpc"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// BTCServiceServer is BTC service server contract.
type BTCServiceServer interface {
	rpc.BTCServiceServer
	health.HealthServer
}

// NewBTCHandler returns a new gRPC handler that implements BTCServiceServer interface.
func NewBTCHandler() BTCServiceServer {
	return &btcHandler{}
}

// btcHandler is a struct for handler.
type btcHandler struct {
	rpc.UnimplementedBTCServiceServer
}

// CreateTransaction creates a new record for BTC transaction.
// Only single transaction will create by this RPC for a specific User.
func (h *btcHandler) CreateTransaction(ctx context.Context, req *rpc.CreateTransactionRequest) (*rpc.Transaction, error) {
	return nil, errors.New("unimplemented")
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (h *btcHandler) ListTransaction(ctx context.Context, req *rpc.ListTransactionRequest) (*rpc.ListTransactionResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetUserBalance get the latest balance for a specific User.
func (h *btcHandler) GetUserBalance(ctx context.Context, req *rpc.GetUserBalanceRequest) (*rpc.UserBalance, error) {
	return nil, errors.New("unimplemented")
}
