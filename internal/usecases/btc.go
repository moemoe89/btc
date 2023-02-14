package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateTransaction creates a new record for BTC transaction.
// Only single transaction will create by this RPC for a specific User.
func (u *btcUsecase) CreateTransaction(ctx context.Context, params *repository.CreateTransactionParams) (*rpc.Transaction, error) {
	ctx, span := u.trace.StartSpan(ctx, "UC.CreateTransaction", nil)
	defer span.End()

	return u.btcRepo.CreateTransaction(ctx, params)
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (u *btcUsecase) ListTransaction( //nolint: funlen
	ctx context.Context, params *repository.ListTransactionParams,
) (*rpc.ListTransactionResponse, error) {
	ctx, span := u.trace.StartSpan(ctx, "UC.ListTransaction", nil)
	defer span.End()

	var (
		list         = new(rpc.ListTransactionResponse)
		transactions []*rpc.Transaction
		err          error
	)

	// Create key for user balance cache based on User ID.
	key := fmt.Sprintf("user:transactions:%d:%d:%d",
		params.UserID,
		params.StartDatetime.UnixNano(),
		params.EndDatetime.UnixNano(),
	)

	// Gets the cache from Redis.
	val, err := u.redis.Get(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		// Ignore the error and just log it, because it supposed to not blocking.
		u.logger.Warn("failed gets user transactions from redis", zap.Error(err))
	}

	// If no error, assume the cache exists.
	if err == nil {
		// Unmarshal the cache []byte to proto message.
		// This is also should not become blocker,
		// thus if error happen we should do the normal query.
		err = protojson.Unmarshal([]byte(val.(string)), list)
		if err == nil {
			return list, err
		}

		u.logger.Warn("failed to unmarshal balance proto", zap.Error(err))
	}

	transactions, err = u.btcRepo.ListTransaction(ctx, params)
	if err != nil {
		return nil, err
	}

	list = &rpc.ListTransactionResponse{
		Transactions: transactions,
	}

	// Marshal the proto message to []byte
	b, err := protojson.Marshal(list)
	if err != nil {
		u.logger.Warn("failed to marshal transactions proto", zap.Error(err))
	}

	// Store the cache with expires 1 second.
	// The number can be changed later based on the traffic.
	_, err = u.redis.Set(ctx, key, b, time.Second*1)
	if err != nil {
		u.logger.Warn("failed stores user transactions to redis", zap.Error(err))
	}

	return list, nil
}

// GetUserBalance get the latest balance for a specific User.
func (u *btcUsecase) GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error) {
	ctx, span := u.trace.StartSpan(ctx, "UC.GetUserBalance", nil)
	defer span.End()

	var (
		balance = &rpc.UserBalance{}
		err     error
	)

	// Create key for user balance cache based on User ID.
	key := fmt.Sprintf("user:balance:%d", userID)

	// Gets the cache from Redis.
	val, err := u.redis.Get(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		// Ignore the error and just log it, because it supposed to not blocking.
		u.logger.Warn("failed gets user balance from redis", zap.Error(err))
	}

	// If no error, assume the cache exists.
	if err == nil {
		// Unmarshal the cache []byte to proto message.
		// This is also should not become blocker,
		// thus if error happen we should do the normal query.
		err = protojson.Unmarshal([]byte(val.(string)), balance)
		if err == nil {
			return balance, err
		}

		u.logger.Warn("failed to unmarshal balance proto", zap.Error(err))
	}

	balance, err = u.btcRepo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Marshal the proto message to []byte
	b, err := protojson.Marshal(balance)
	if err != nil {
		u.logger.Warn("failed to marshal balance proto", zap.Error(err))
	}

	// Store the cache with expires 1 second.
	// The number can be changed later based on the traffic.
	_, err = u.redis.Set(ctx, key, b, time.Second*1)
	if err != nil {
		u.logger.Warn("failed stores user balance to redis", zap.Error(err))
	}

	return balance, nil
}
