package di

import (
	"log"
	"os"

	"github.com/moemoe89/btc/pkg/di"
	"github.com/moemoe89/btc/pkg/kvs"
	"github.com/moemoe89/btc/pkg/kvs/redis"
)

// GetRedis get the Redis KVS client.
func GetRedis() kvs.Client {
	r, err := redis.New(redis.WithAddr(os.Getenv("REDIS_HOST")))
	if err != nil {
		log.Fatal(err)
	}

	di.RegisterCloser("RedisConnection", r)

	return r
}
