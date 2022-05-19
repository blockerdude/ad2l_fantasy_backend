package fantasycontext

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	pool_value  = "pool"
	authn_value = "authn"
)

func WithDBPool(ctx context.Context, pool *pgxpool.Pool) context.Context {
	return context.WithValue(ctx, pool_value, pool)
}

func GetDBPool(ctx context.Context) *pgxpool.Pool {
	val := ctx.Value(pool_value)
	return val.(*pgxpool.Pool)
}

func WithAuthn(ctx context.Context, authn *model.Authn) context.Context {
	return context.WithValue(ctx, authn_value, authn)
}

func GetAuthn(ctx context.Context) *model.Authn {
	val := ctx.Value(authn_value)
	return val.(*model.Authn)
}
