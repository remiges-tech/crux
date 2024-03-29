// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: config.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const configGet = `-- name: ConfigGet :many
SELECT 
name AS attr,
val,ver,
setby AS by
FROM config 
where realm = $1
`

type ConfigGetRow struct {
	Attr string      `json:"attr"`
	Val  pgtype.Text `json:"val"`
	Ver  pgtype.Int4 `json:"ver"`
	By   string      `json:"by"`
}

func (q *Queries) ConfigGet(ctx context.Context, realm string) ([]ConfigGetRow, error) {
	rows, err := q.db.Query(ctx, configGet, realm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ConfigGetRow
	for rows.Next() {
		var i ConfigGetRow
		if err := rows.Scan(
			&i.Attr,
			&i.Val,
			&i.Ver,
			&i.By,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const configSet = `-- name: ConfigSet :exec
INSERT INTO
    config(
        realm, slice, name, descr, val, setby
    )
VALUES (
        $1, $2, $3, $4, $5, $6
    )
`

type ConfigSetParams struct {
	Realm string      `json:"realm"`
	Slice int32       `json:"slice"`
	Name  string      `json:"name"`
	Descr string      `json:"descr"`
	Val   pgtype.Text `json:"val"`
	Setby string      `json:"setby"`
}

func (q *Queries) ConfigSet(ctx context.Context, arg ConfigSetParams) error {
	_, err := q.db.Exec(ctx, configSet,
		arg.Realm,
		arg.Slice,
		arg.Name,
		arg.Descr,
		arg.Val,
		arg.Setby,
	)
	return err
}
