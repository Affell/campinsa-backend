package config

import (
	"context"
	"os"
	"oui/models/shotgun"
	"oui/models/user"

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

	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		redisConfig.Password = env
	} else {
		redisConfig.Password = ""
	}

	ctx := context.Background()

	user.UserTokenRedisConn = redis.NewClient(&redis.Options{
		Addr:     redisConfig.URL,
		DB:       0,
		Password: redisConfig.Password,
	})
	user.UserTokenRedisCtx = context.Background()
	user.UserTokenRedisConn.FlushAll(ctx)

	shotgun.RedisConn = redis.NewClient(&redis.Options{
		Addr: redisConfig.URL,
		DB:   1,
	})
	shotgun.RedisCtx = context.Background()
	shotgun.RedisConn.FlushAll(ctx)
	return
}
