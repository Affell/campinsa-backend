package config

import (
	"context"
	"os"
	"oui/auth"
	"oui/models/shotgun"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	URL      string
	Password string
}

func InitRedis() (redisConfig Redis) {
	if env := os.Getenv("REDIS_HOST"); env != "" {
		redisConfig.URL = env
	} else {
		redisConfig.URL = "localhost:6379"
	}

	ctx := context.Background()

	auth.UserTokenRedisConn = redis.NewClient(&redis.Options{
		Addr: redisConfig.URL,
		DB:   0,
	})
	auth.UserTokenRedisCtx = context.Background()
	auth.UserTokenRedisConn.FlushAll(ctx)

	shotgun.RedisConn = redis.NewClient(&redis.Options{
		Addr: redisConfig.URL,
		DB:   1,
	})
	shotgun.RedisCtx = context.Background()
	shotgun.RedisConn.FlushAll(ctx)
	return
}
