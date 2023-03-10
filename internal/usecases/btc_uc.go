package usecases

//go:generate rm -f ./btc_uc_mock.go
//go:generate mockgen -destination btc_uc_mock.go -package usecases -mock_names BTCUsecase=GoMockBTCUsecase -source btc_uc.go

import (
	"context"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/pkg/kvs"
	"github.com/moemoe89/btc/pkg/logging"
	"github.com/moemoe89/btc/pkg/trace"
)

// BTCUsecase defines BTC transactions related domain functionality.
type BTCUsecase interface {
	// CreateTransaction creates a new record for BTC transaction.
	// Only single transaction will create by this RPC for a specific User.
	CreateTransaction(ctx context.Context, params *repository.CreateTransactionParams) (*rpc.Transaction, error)
	// ListTransaction get the list of records for BTC transaction.
	// The record can be filtered by specific User.
	ListTransaction(ctx context.Context, params *repository.ListTransactionParams) (*rpc.ListTransactionResponse, error)
	// GetUserBalance get the latest balance for a specific User.
	GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error)
}

// compile time interface implementation check.
var _ BTCUsecase = (*btcUsecase)(nil)

// NewBTCUsecase returns BTCUsecase.
func NewBTCUsecase(
	btcRepo repository.BTCRepo,
	trace trace.Tracer,
	logger logging.Logger,
	redis kvs.Client,
) BTCUsecase {
	return &btcUsecase{
		btcRepo: btcRepo,
		trace:   trace,
		logger:  logger,
		redis:   redis,
	}
}

// btcUsecase is a struct for usecase.
type btcUsecase struct {
	btcRepo repository.BTCRepo
	trace   trace.Tracer
	logger  logging.Logger
	redis   kvs.Client
}
