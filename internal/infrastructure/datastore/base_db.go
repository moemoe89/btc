package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbOnce sync.Once
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

func getConnString() string {
	if os.Getenv("APP_ENV") == "test" || os.Getenv("APP_ENV") == "" {
		return "postgres://test:test@localhost:5432/test?sslmode=disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST")+":"+os.Getenv("POSTGRES_PORT"), // for lint purpose
		os.Getenv("POSTGRES_DB"),
	)
}

// GetDatabase returns postgresql Pool.
func GetDatabase() *pgxpool.Pool {
	dbOnce.Do(func() {
		ctx := context.Background()

		var err error

		connString := getConnString()

		// Use default config.
		db, err = pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatalf("failed to connect to timescaleDB pool: %v", err)
		}

		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}
	})

	return db
}
