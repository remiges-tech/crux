// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
<<<<<<< HEAD

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
=======
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
>>>>>>> 95154679e693ee1f5f61feec17444188e01ce8d6
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

<<<<<<< HEAD
func (q *Queries) WithTx(tx pgx.Tx) *Queries {
=======
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
>>>>>>> 95154679e693ee1f5f61feec17444188e01ce8d6
	return &Queries{
		db: tx,
	}
}