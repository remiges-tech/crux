package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewProvider(connString string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return connPool, err
}
