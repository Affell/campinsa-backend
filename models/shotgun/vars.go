package shotgun

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const K = "shotgun"

var (
	RedisConn *redis.Client
	RedisCtx  context.Context
)
