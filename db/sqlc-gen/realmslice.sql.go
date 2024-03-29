// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: realmslice.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const cloneRecordInConfigBySliceID = `-- name: CloneRecordInConfigBySliceID :execresult
INSERT INTO
    config (
        realm, slice, name, descr, val, ver, setby
    )
SELECT realm, $2, name, descr, val, ver, $3
FROM config
WHERE
    config.slice = $1
`

type CloneRecordInConfigBySliceIDParams struct {
	Slice   int32  `json:"slice"`
	Slice_2 int32  `json:"slice_2"`
	Setby   string `json:"setby"`
}

func (q *Queries) CloneRecordInConfigBySliceID(ctx context.Context, arg CloneRecordInConfigBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInConfigBySliceID, arg.Slice, arg.Slice_2, arg.Setby)
}

const cloneRecordInRealmSliceBySliceID = `-- name: CloneRecordInRealmSliceBySliceID :one
INSERT INTO
    realmslice (
        realm, descr, active, activateat, deactivateat,createdby
    )
SELECT
    realm,
    COALESCE(
        descr, $4::text
    ),
    true,
    activateat,
    deactivateat,
    $3
FROM realmslice
WHERE
    realmslice.id = $1
    AND realmslice.realm = $2
RETURNING
    realmslice.id
`

type CloneRecordInRealmSliceBySliceIDParams struct {
	ID        int32       `json:"id"`
	Realm     string      `json:"realm"`
	Createdby string      `json:"createdby"`
	Descr     pgtype.Text `json:"descr"`
}

func (q *Queries) CloneRecordInRealmSliceBySliceID(ctx context.Context, arg CloneRecordInRealmSliceBySliceIDParams) (int32, error) {
	row := q.db.QueryRow(ctx, cloneRecordInRealmSliceBySliceID,
		arg.ID,
		arg.Realm,
		arg.Createdby,
		arg.Descr,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const cloneRecordInRulesetBySliceID = `-- name: CloneRecordInRulesetBySliceID :execresult
INSERT INTO
    ruleset (
        realm, slice, app, brwf, class, setname, schemaid, is_active, is_internal, ruleset, createdby
    )
SELECT
    realm,
    $2,
    app,
    brwf,
    class,
    setname,
    schemaid,
    is_active,
    is_internal,
    ruleset,
    $3
FROM ruleset
WHERE
    ruleset.slice = $1
    AND (
        $4::text [] is null
        OR app = any ($4::text [])
    )
`

type CloneRecordInRulesetBySliceIDParams struct {
	Slice     int32    `json:"slice"`
	Slice_2   int32    `json:"slice_2"`
	Createdby string   `json:"createdby"`
	App       []string `json:"app"`
}

func (q *Queries) CloneRecordInRulesetBySliceID(ctx context.Context, arg CloneRecordInRulesetBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInRulesetBySliceID,
		arg.Slice,
		arg.Slice_2,
		arg.Createdby,
		arg.App,
	)
}

const cloneRecordInSchemaBySliceID = `-- name: CloneRecordInSchemaBySliceID :execresult
INSERT INTO
    schema (
        realm, slice, app, brwf, class, patternschema, actionschema, createdby
    )
SELECT
    realm,
    $2,
    app,
    brwf,
    class,
    patternschema,
    actionschema,
    $3
FROM schema
WHERE
    schema.slice = $1
    AND (
        $4::text [] is null
        OR app = any ($4::text [])
    )
`

type CloneRecordInSchemaBySliceIDParams struct {
	Slice     int32    `json:"slice"`
	Slice_2   int32    `json:"slice_2"`
	Createdby string   `json:"createdby"`
	App       []string `json:"app"`
}

func (q *Queries) CloneRecordInSchemaBySliceID(ctx context.Context, arg CloneRecordInSchemaBySliceIDParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, cloneRecordInSchemaBySliceID,
		arg.Slice,
		arg.Slice_2,
		arg.Createdby,
		arg.App,
	)
}

const getRealmSliceListByRealm = `-- name: GetRealmSliceListByRealm :many
SELECT
    id,descr,active,deactivateat,createdat,createdby,editedat,editedby
FROM
    realmslice
WHERE
    realm= $1
`

type GetRealmSliceListByRealmRow struct {
	ID           int32            `json:"id"`
	Descr        string           `json:"descr"`
	Active       bool             `json:"active"`
	Deactivateat pgtype.Timestamp `json:"deactivateat"`
	Createdat    pgtype.Timestamp `json:"createdat"`
	Createdby    string           `json:"createdby"`
	Editedat     pgtype.Timestamp `json:"editedat"`
	Editedby     pgtype.Text      `json:"editedby"`
}

func (q *Queries) GetRealmSliceListByRealm(ctx context.Context, realm string) ([]GetRealmSliceListByRealmRow, error) {
	rows, err := q.db.Query(ctx, getRealmSliceListByRealm, realm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRealmSliceListByRealmRow
	for rows.Next() {
		var i GetRealmSliceListByRealmRow
		if err := rows.Scan(
			&i.ID,
			&i.Descr,
			&i.Active,
			&i.Deactivateat,
			&i.Createdat,
			&i.Createdby,
			&i.Editedat,
			&i.Editedby,
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

const insertNewRecordInRealmSlice = `-- name: InsertNewRecordInRealmSlice :one
INSERT INTO
    realmslice (
        realm, descr, active, createdby
    )
VALUES ($1, $2, true, $3) RETURNING realmslice.id
`

type InsertNewRecordInRealmSliceParams struct {
	Realm     string `json:"realm"`
	Descr     string `json:"descr"`
	Createdby string `json:"createdby"`
}

func (q *Queries) InsertNewRecordInRealmSlice(ctx context.Context, arg InsertNewRecordInRealmSliceParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertNewRecordInRealmSlice, arg.Realm, arg.Descr, arg.Createdby)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const realmSliceActivate = `-- name: RealmSliceActivate :one
UPDATE realmslice
SET
    active = $1,
    activateat = CASE
        WHEN (
            $2::TIMESTAMP
        ) IS NULL THEN NOW()
        ELSE (
            $2::TIMESTAMP
        )
    END,
    deactivateat = NULL
WHERE
    id = $3
RETURNING
    id, realm, descr, active, activateat, deactivateat, createdat, createdby, editedat, editedby
`

type RealmSliceActivateParams struct {
	Isactive   bool             `json:"isactive"`
	Activateat pgtype.Timestamp `json:"activateat"`
	ID         int32            `json:"id"`
}

func (q *Queries) RealmSliceActivate(ctx context.Context, arg RealmSliceActivateParams) (Realmslice, error) {
	row := q.db.QueryRow(ctx, realmSliceActivate, arg.Isactive, arg.Activateat, arg.ID)
	var i Realmslice
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Descr,
		&i.Active,
		&i.Activateat,
		&i.Deactivateat,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}

const realmSliceAppsList = `-- name: RealmSliceAppsList :many
SELECT a.shortname, a.longname
FROM realmslice
    JOIN app a ON realmslice.realm = a.realm
WHERE
    realmslice.id = $1
`

type RealmSliceAppsListRow struct {
	Shortname string `json:"shortname"`
	Longname  string `json:"longname"`
}

func (q *Queries) RealmSliceAppsList(ctx context.Context, id int32) ([]RealmSliceAppsListRow, error) {
	rows, err := q.db.Query(ctx, realmSliceAppsList, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RealmSliceAppsListRow
	for rows.Next() {
		var i RealmSliceAppsListRow
		if err := rows.Scan(&i.Shortname, &i.Longname); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const realmSliceDeactivate = `-- name: RealmSliceDeactivate :one
UPDATE realmslice
SET
    active = $1,
    deactivateat = CASE
        WHEN (
            $2::TIMESTAMP
        ) IS NULL THEN NOW()
        ELSE (
            $2::TIMESTAMP
        )
    END,
    activateat = NULL
WHERE
    id = $3
RETURNING
    id, realm, descr, active, activateat, deactivateat, createdat, createdby, editedat, editedby
`

type RealmSliceDeactivateParams struct {
	Isactive     bool             `json:"isactive"`
	Deactivateat pgtype.Timestamp `json:"deactivateat"`
	ID           int32            `json:"id"`
}

func (q *Queries) RealmSliceDeactivate(ctx context.Context, arg RealmSliceDeactivateParams) (Realmslice, error) {
	row := q.db.QueryRow(ctx, realmSliceDeactivate, arg.Isactive, arg.Deactivateat, arg.ID)
	var i Realmslice
	err := row.Scan(
		&i.ID,
		&i.Realm,
		&i.Descr,
		&i.Active,
		&i.Activateat,
		&i.Deactivateat,
		&i.Createdat,
		&i.Createdby,
		&i.Editedat,
		&i.Editedby,
	)
	return i, err
}

const realmSlicePurge = `-- name: RealmSlicePurge :execresult
WITH
    del_stepworkflow AS (
        DELETE FROM stepworkflow
    ),
    del_wfinstance AS (
        DELETE FROM wfinstance
    ),
    del_ruleset AS (
        DELETE FROM ruleset
    ),
    del_schema AS (
        DELETE FROM schema
    ),
    del_config AS (
        DELETE FROM config
    )
DELETE FROM realmslice
WHERE
    realmslice.realm = $1
`

func (q *Queries) RealmSlicePurge(ctx context.Context, realm string) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, realmSlicePurge, realm)
}
