package sqlc

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBQuerier defines all functions to execute db queries and transactions
type DBQuerier interface {
	PerformDBOperationWithTX(ctx context.Context, queries func(Querier) (any, error)) (any, error)
	Querier
}

// QuerierTX provides all functions to execute SQL queries and transactions
type QuerierTX struct {
	connPool *pgxpool.Pool
	*Queries
}

/*
NewQuerierWithTX creates a new QuerierTX object
and returns QuerierWithTX interface
*/
func NewQuerierWithTX(connPool *pgxpool.Pool) DBQuerier {
	return &QuerierTX{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

/*
PerformTX performs transaction operations for the queries callbackFunction
methods returns a result and error.PerformTx consist execTX methods which deals
with the transaction commit and rollback logic
*/
func (tx *QuerierTX) PerformDBOperationWithTX(ctx context.Context, queries func(Querier) (any, error)) (any, error) {
	return tx.execTX(ctx, queries)
}

/*
execTX performs transaction commit and rollBack for the seeded
queryCallBack function arguments.execTX will return a result and error Based
on the transaction behavior
*/
func (qtx *QuerierTX) execTX(ctx context.Context, queryCallBack func(Querier) (any, error)) (any, error) {
	tx, err := qtx.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	q := New(tx)
	result, err := queryCallBack(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return nil, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return result, err
	}
	return result, tx.Commit(ctx)
}
