package datastore

import (
	"context"
	"fmt"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"

	"github.com/jackc/pgx/v5"
	"github.com/moemoe89/btc/internal/entities/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (r *btcRepo) CreateTransaction(ctx context.Context, params *repository.CreateTransactionParams) (*rpc.Transaction, error) {
	var id int64

	var err error

	err = r.db.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", params.UserID).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user id: %d not found", params.UserID)
	}

	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
	}()

	query := `INSERT INTO transactions (time, user_id, amount) VALUES ($1, $2, $3)`

	_, err = tx.Exec(ctx, query, params.Datetime, params.UserID, params.Amount)
	if err != nil {
		return nil, err
	}

	query = `UPDATE users SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(ctx, query, params.Amount, params.UserID)
	if err != nil {
		return nil, err
	}

	errCommit := tx.Commit(ctx)
	if errCommit != nil {
		return nil, fmt.Errorf("unable to commit transaction: %v", errCommit)
	}

	return &rpc.Transaction{
		UserId:   params.UserID,
		Datetime: timestamppb.New(params.Datetime),
		Amount:   params.Amount,
	}, nil
}

// ListTransaction get the list of records for BTC transaction.
// The record can be filtered by specific User.
func (r *btcRepo) ListTransaction(ctx context.Context, userID int64) ([]*rpc.Transaction, error) {
	query := `SELECT time, user_id, amount FROM transactions WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type data struct {
		Time   time.Time
		UserID int64
		Amount float32
	}

	var transactions []*rpc.Transaction

	for rows.Next() {
		var d data

		err = rows.Scan(&d.Time, &d.UserID, &d.Amount)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &rpc.Transaction{
			UserId:   d.UserID,
			Datetime: timestamppb.New(d.Time),
			Amount:   d.Amount,
		})
	}

	return transactions, nil
}

// GetUserBalance get the latest balance for a specific User.
func (r *btcRepo) GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error) {
	var balance float32

	err := r.db.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user id: %d not found", userID)
	}

	if err != nil {
		return nil, err
	}

	return &rpc.UserBalance{
		Balance: balance,
	}, nil
}
