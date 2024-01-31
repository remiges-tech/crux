package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

func NewProvider(connString string) (sqlc.Querier, error) {
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return sqlc.New(connPool), nil
}
