package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/internal/usecases"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// BTCServiceServer is BTC service server contract.
type BTCServiceServer interface {
	rpc.BTCServiceServer
	health.HealthServer
}

// NewBTCHandler returns a new gRPC handler that implements BTCServiceServer interface.
func NewBTCHandler(uc usecases.BTCUsecase) BTCServiceServer {
	return &btcHandler{
		uc: uc,
	}
}

// btcHandler is a struct for handler.
type btcHandler struct {
	rpc.UnimplementedBTCServiceServer
	uc usecases.BTCUsecase
}

// CreateTransaction creates a new record for BTC transaction.
// Only single transaction will create by this RPC for a specific User.
func (h *btcHandler) CreateTransaction(ctx context.Context, req *rpc.CreateTransactionRequest) (*rpc.Transaction, error) {
	return h.uc.CreateTransaction(ctx, &repository.CreateTransactionParams{
		UserID:   req.GetUserId(),
		Datetime: req.GetDatetime().AsTime(),
		Amount:   req.GetAmount(),
	})
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (h *btcHandler) ListTransaction(ctx context.Context, req *rpc.ListTransactionRequest) (*rpc.ListTransactionResponse, error) {
	return h.uc.ListTransaction(ctx, req.GetUserId())
}

// GetUserBalance get the latest balance for a specific User.
func (h *btcHandler) GetUserBalance(ctx context.Context, req *rpc.GetUserBalanceRequest) (*rpc.UserBalance, error) {
	return h.uc.GetUserBalance(ctx, req.GetUserId())
}
