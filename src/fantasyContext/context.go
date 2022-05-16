package fantasycontext

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	pool_value = "pool"
)

func WithDBPool(ctx context.Context, pool *pgxpool.Pool) context.Context {
	return context.WithValue(ctx, pool_value, pool)
}

func GetDBPool(ctx context.Context) *pgxpool.Pool {
	return ctx.Value(pool_value).(*pgxpool.Pool)
}
