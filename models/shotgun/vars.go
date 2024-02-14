package shotgun

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	RedisConn *redis.Client
	RedisCtx  context.Context
)
