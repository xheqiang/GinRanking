package cache

import (
	"ginRanking/config"

	"context"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Rctx context.Context

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		DB:       config.REDIS_DB,
		Password: config.REDIS_PASSWORD,
	})

	Rctx = context.Background()
}
