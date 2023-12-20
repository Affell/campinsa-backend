package auth

import (
	"context"

	"github.com/jackc/pgx/v4"
)

var (
	SQLConn *pgx.ConnConfig
	SQLCtx  context.Context
)
