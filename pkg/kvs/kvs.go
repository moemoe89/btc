package kvs

//go:generate rm -f ./kvs_mock.go
//go:generate mockgen -destination kvs_mock.go -package kvs -mock_names Client=GoMockClient -source kvs.go

import (
	"context"
	"time"
)

// Client is an interface for KVS cache.
type Client interface {
	// Set sets the value with expiration time.
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (interface{}, error)
	// Get gets the value by the given key.
	Get(ctx context.Context, key string) (interface{}, error)
	// Close closes the connection of KVS client.
	Close() error
}
