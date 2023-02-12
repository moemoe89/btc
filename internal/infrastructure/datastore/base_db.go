package datastore

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
)

// NewBaseRepo returns a base repository.
func NewBaseRepo(db *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{db: db}
}

// BaseRepo is a base repository.
type BaseRepo struct {
	db *pgxpool.Pool
}

// GetDatabase returns postgresql Pool.
func GetDatabase() *pgxpool.Pool {
	if db != nil {
		return db
	}

	ctx := context.Background()

	var err error

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST")+":"+os.Getenv("POSTGRES_PORT"), // for lint purpose
		os.Getenv("POSTGRES_DB"),
	)

	db, err = pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to timescaleDB pool: %v", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}
