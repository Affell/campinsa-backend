package auth

import (
	"context"

	"github.com/jackc/pgx/v4"
)

var (
	secretKey string = NewEncryptSecretKey()
	SQLConn   *pgx.ConnConfig
	SQLCtx    context.Context
)
