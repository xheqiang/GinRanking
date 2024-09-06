package cache

import (
	"ginRanking/config"

	"context"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Rctx context.Context

func init() {

	redisAddr := config.AppConf.RedisConfig.Host + ":" + config.AppConf.RedisConfig.Port
	redisDb := config.AppConf.RedisConfig.DB
	redisPassword := config.AppConf.RedisConfig.Password

	Redis = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       redisDb,
		Password: redisPassword,
	})

	Rctx = context.Background()
}
