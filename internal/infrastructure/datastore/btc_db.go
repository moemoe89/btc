package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/internal/entities/repository"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ErrNotFound is an error for indicates record not found.
	ErrNotFound = errors.New("error not found")
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

	err = r.dbSlave.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", params.UserID).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user id: %d not found: %w", params.UserID, ErrNotFound)
	}

	if err != nil {
		return nil, err
	}

	tx, err := r.dbMaster.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
	}()

	query := `INSERT INTO transactions (datetime, user_id, amount) VALUES ($1, $2, $3)`

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
func (r *btcRepo) ListTransaction(ctx context.Context, params *repository.ListTransactionParams) ([]*rpc.Transaction, error) {
	query := `SELECT datetime, user_id, amount FROM transactions WHERE user_id = $1 AND datetime >= $2::timestamptz AND datetime <= $3::timestamptz`

	rows, err := r.dbSlave.Query(ctx, query, params.UserID, params.StartDatetime, params.EndDatetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type data struct {
		Datetime time.Time
		UserID   int64
		Amount   float64
	}

	var transactions []*rpc.Transaction

	for rows.Next() {
		var d data

		err = rows.Scan(&d.Datetime, &d.UserID, &d.Amount)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &rpc.Transaction{
			UserId:   d.UserID,
			Datetime: timestamppb.New(d.Datetime),
			Amount:   d.Amount,
		})
	}

	return transactions, nil
}

// GetUserBalance get the latest balance for a specific User.
func (r *btcRepo) GetUserBalance(ctx context.Context, userID int64) (*rpc.UserBalance, error) {
	var balance float64

	err := r.dbSlave.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user id: %d not found: %w", userID, ErrNotFound)
	}

	if err != nil {
		return nil, err
	}

	return &rpc.UserBalance{
		Balance: balance,
	}, nil
}
