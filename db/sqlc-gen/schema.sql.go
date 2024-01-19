// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: schema.sql

package sqlc

import (
	"context"
	"encoding/json"
)

const schemaDelete = `-- name: SchemaDelete :one
DELETE FROM schema
WHERE
    id = $1
RETURNING id, realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby, editedat, editedby
`

// :one
func (q *Queries) SchemaDelete(ctx context.Context, id int32) (Schema, error) {
	row := q.db.QueryRowContext(ctx, schemaDelete, id)
	var i Schema
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Slice,
		&i.App,
		&i.Brwf,
		&i.Class,
		&i.Patternschema,
		&i.Actionschema,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}

const schemaGet = `-- name: SchemaGet :one
SELECT id, realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby, editedat, editedby
FROM schema
WHERE
    id = $1
`

// :one
func (q *Queries) SchemaGet(ctx context.Context, id int32) (Schema, error) {
	row := q.db.QueryRowContext(ctx, schemaGet, id)
	var i Schema
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Slice,
		&i.App,
		&i.Brwf,
		&i.Class,
		&i.Patternschema,
		&i.Actionschema,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}

const schemaList = `-- name: SchemaList :many
SELECT id, realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby, editedat, editedby
FROM schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3
`

type SchemaListParams struct {
	Slice int32  `json:"slice"`
	Class string `json:"class"`
	App   string `json:"app"`
}

func (q *Queries) SchemaList(ctx context.Context, arg SchemaListParams) ([]Schema, error) {
	rows, err := q.db.QueryContext(ctx, schemaList, arg.Slice, arg.Class, arg.App)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Schema
	for rows.Next() {
		var i Schema
		if err := rows.Scan(
			&i.ID,
			&i.Realm,
			&i.Slice,
			&i.App,
			&i.Brwf,
			&i.Class,
			&i.Patternschema,
			&i.Actionschema,
			&i.Createdat,
			&i.Createdby,
			&i.Editedat,
			&i.Editedby,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const schemaNew = `-- name: SchemaNew :one
INSERT INTO schema (
    realm, slice, app, brwf, class, patternschema, actionschema, createdby, editedby
) VALUES (
    1, $1, $2, W, $3, $4, $5, $6, $7
) RETURNING id, realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby, editedat, editedby
`

type SchemaNewParams struct {
	Slice         int32           `json:"slice"`
	App           string          `json:"app"`
	Class         string          `json:"class"`
	Patternschema json.RawMessage `json:"patternschema"`
	Actionschema  json.RawMessage `json:"actionschema"`
	Createdby     string          `json:"createdby"`
	Editedby      string          `json:"editedby"`
}

// :one
func (q *Queries) SchemaNew(ctx context.Context, arg SchemaNewParams) (Schema, error) {
	row := q.db.QueryRowContext(ctx, schemaNew,
		arg.Slice,
		arg.App,
		arg.Class,
		arg.Patternschema,
		arg.Actionschema,
		arg.Createdby,
		arg.Editedby,
	)
	var i Schema
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Slice,
		&i.App,
		&i.Brwf,
		&i.Class,
		&i.Patternschema,
		&i.Actionschema,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}

const schemaUpdate = `-- name: SchemaUpdate :one
UPDATE schema
SET
    app = $2,
    brwf = $3,
    class = $4,
    patternschema = $5,
    actionschema = $6,
    editedat = CURRENT_TIMESTAMP,
    editedby = $7
WHERE
    id = $1
RETURNING id, realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby, editedat, editedby
`

type SchemaUpdateParams struct {
	ID            int32           `json:"id"`
	App           string          `json:"app"`
	Brwf          string          `json:"brwf"`
	Class         string          `json:"class"`
	Patternschema json.RawMessage `json:"patternschema"`
	Actionschema  json.RawMessage `json:"actionschema"`
	Editedby      string          `json:"editedby"`
}

// :one
func (q *Queries) SchemaUpdate(ctx context.Context, arg SchemaUpdateParams) (Schema, error) {
	row := q.db.QueryRowContext(ctx, schemaUpdate,
		arg.ID,
		arg.App,
		arg.Brwf,
		arg.Class,
		arg.Patternschema,
		arg.Actionschema,
		arg.Editedby,
	)
	var i Schema
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Slice,
		&i.App,
		&i.Brwf,
		&i.Class,
		&i.Patternschema,
		&i.Actionschema,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}
