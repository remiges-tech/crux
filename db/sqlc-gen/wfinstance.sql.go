// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: wfinstance.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addWFNewInstances = `-- name: AddWFNewInstances :many
INSERT INTO
    wfinstance (
        entityid, slice, app, class, workflow, step, loggedat, nextstep, parent
    )
VALUES (
        $1, $2, $3, $4, $5, unnest($6::text []), (NOW()::timestamp), $7, $8
    )
RETURNING id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent
`

type AddWFNewInstancesParams struct {
	Entityid string      `json:"entityid"`
	Slice    int32       `json:"slice"`
	App      string      `json:"app"`
	Class    string      `json:"class"`
	Workflow string      `json:"workflow"`
	Step     []string    `json:"step"`
	Nextstep string      `json:"nextstep"`
	Parent   pgtype.Int4 `json:"parent"`
}

func (q *Queries) AddWFNewInstances(ctx context.Context, arg AddWFNewInstancesParams) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, addWFNewInstances,
		arg.Entityid,
		arg.Slice,
		arg.App,
		arg.Class,
		arg.Workflow,
		arg.Step,
		arg.Nextstep,
		arg.Parent,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const deleteWFInstanceListByParents = `-- name: DeleteWFInstanceListByParents :many
DELETE FROM wfinstance
WHERE 
   ($1::INTEGER[] IS NOT NULL AND id = ANY($1::INTEGER[]) OR $2::INTEGER[] IS NOT NULL AND parent = ANY($2::INTEGER[]))
    RETURNING id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent
`

type DeleteWFInstanceListByParentsParams struct {
	ID     []int32 `json:"id"`
	Parent []int32 `json:"parent"`
}

func (q *Queries) DeleteWFInstanceListByParents(ctx context.Context, arg DeleteWFInstanceListByParentsParams) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, deleteWFInstanceListByParents, arg.ID, arg.Parent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const deleteWFInstances = `-- name: DeleteWFInstances :exec
DELETE FROM
    wfinstance
WHERE
     wfinstance.entityid IN (SELECT wfinstance.entityid FROM wfinstance WHERE wfinstance.id = $1)
    AND wfinstance.slice = $2
    AND wfinstance.app = $3
`

type DeleteWFInstancesParams struct {
	ID    int32  `json:"id"`
	Slice int32  `json:"slice"`
	App   string `json:"app"`
}

func (q *Queries) DeleteWFInstances(ctx context.Context, arg DeleteWFInstancesParams) error {
	_, err := q.db.Exec(ctx, deleteWFInstances, arg.ID, arg.Slice, arg.App)
	return err
}

const deleteWfinstanceByID = `-- name: DeleteWfinstanceByID :many
  DELETE FROM wfinstance
   WHERE
       (id = $1::INTEGER OR entityid = $2::TEXT)
   RETURNING id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent
`

type DeleteWfinstanceByIDParams struct {
	ID       pgtype.Int4 `json:"id"`
	Entityid pgtype.Text `json:"entityid"`
}

func (q *Queries) DeleteWfinstanceByID(ctx context.Context, arg DeleteWfinstanceByIDParams) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, deleteWfinstanceByID, arg.ID, arg.Entityid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const getWFINstance = `-- name: GetWFINstance :one
SELECT count(1)
FROM wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4
`

type GetWFINstanceParams struct {
	Slice    int32  `json:"slice"`
	App      string `json:"app"`
	Workflow string `json:"workflow"`
	Entityid string `json:"entityid"`
}

func (q *Queries) GetWFINstance(ctx context.Context, arg GetWFINstanceParams) (int64, error) {
	row := q.db.QueryRow(ctx, getWFINstance,
		arg.Slice,
		arg.App,
		arg.Workflow,
		arg.Entityid,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getWFInstanceCounts = `-- name: GetWFInstanceCounts :one
SELECT COUNT(*) 
FROM wfinstance
WHERE
    wfinstance.slice = $1
    AND wfinstance.app = $2
    AND wfinstance.workflow = $3
    AND wfinstance.entityid IN (SELECT wfinstance.entityid FROM wfinstance WHERE wfinstance.id = $4)
`

type GetWFInstanceCountsParams struct {
	Slice    int32  `json:"slice"`
	App      string `json:"app"`
	Workflow string `json:"workflow"`
	ID       int32  `json:"id"`
}

func (q *Queries) GetWFInstanceCounts(ctx context.Context, arg GetWFInstanceCountsParams) (int64, error) {
	row := q.db.QueryRow(ctx, getWFInstanceCounts,
		arg.Slice,
		arg.App,
		arg.Workflow,
		arg.ID,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getWFInstanceCurrent = `-- name: GetWFInstanceCurrent :one
 SELECT id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent FROM wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4
`

type GetWFInstanceCurrentParams struct {
	Slice    int32  `json:"slice"`
	App      string `json:"app"`
	Workflow string `json:"workflow"`
	Entityid string `json:"entityid"`
}

func (q *Queries) GetWFInstanceCurrent(ctx context.Context, arg GetWFInstanceCurrentParams) (Wfinstance, error) {
	row := q.db.QueryRow(ctx, getWFInstanceCurrent,
		arg.Slice,
		arg.App,
		arg.Workflow,
		arg.Entityid,
	)
	var i Wfinstance
	err := row.Scan(
		&i.ID,
		&i.Entityid,
		&i.Slice,
		&i.App,
		&i.Class,
		&i.Workflow,
		&i.Step,
		&i.Loggedat,
		&i.Doneat,
		&i.Nextstep,
		&i.Parent,
	)
	return i, err
}

const getWFInstanceFromId = `-- name: GetWFInstanceFromId :one
SELECT id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent FROM wfinstance 
WHERE 
    id = $1
`

func (q *Queries) GetWFInstanceFromId(ctx context.Context, id int32) (Wfinstance, error) {
	row := q.db.QueryRow(ctx, getWFInstanceFromId, id)
	var i Wfinstance
	err := row.Scan(
		&i.ID,
		&i.Entityid,
		&i.Slice,
		&i.App,
		&i.Class,
		&i.Workflow,
		&i.Step,
		&i.Loggedat,
		&i.Doneat,
		&i.Nextstep,
		&i.Parent,
	)
	return i, err
}

const getWFInstanceList = `-- name: GetWFInstanceList :many
SELECT id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent FROM wfinstance
WHERE 
   ($1::INTEGER is null OR slice = $1::INTEGER)
   AND ($2::text is null OR entityid = $2::text)
   AND ($3::text is null OR app = $3::text)
   AND ($4::text is null OR workflow = $4::text)
   AND($5::INTEGER is null OR  parent = $5::INTEGER)
`

type GetWFInstanceListParams struct {
	Slice    pgtype.Int4 `json:"slice"`
	Entityid pgtype.Text `json:"entityid"`
	App      pgtype.Text `json:"app"`
	Workflow pgtype.Text `json:"workflow"`
	Parent   pgtype.Int4 `json:"parent"`
}

func (q *Queries) GetWFInstanceList(ctx context.Context, arg GetWFInstanceListParams) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, getWFInstanceList,
		arg.Slice,
		arg.Entityid,
		arg.App,
		arg.Workflow,
		arg.Parent,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const getWFInstanceListByParents = `-- name: GetWFInstanceListByParents :many
SELECT id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent FROM wfinstance
WHERE 
   ($1::INTEGER[] IS NOT NULL AND id = ANY($1::INTEGER[]))
`

func (q *Queries) GetWFInstanceListByParents(ctx context.Context, id []int32) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, getWFInstanceListByParents, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const getWFInstanceListForMarkDone = `-- name: GetWFInstanceListForMarkDone :many
SELECT id, entityid, slice, app, class, workflow, step, loggedat, doneat, nextstep, parent FROM wfinstance 
WHERE
    wfinstance.slice = $1
    AND wfinstance.app = $2
    AND wfinstance.workflow = $3
    AND wfinstance.entityid IN (SELECT wfinstance.entityid FROM wfinstance WHERE wfinstance.id = $4)
`

type GetWFInstanceListForMarkDoneParams struct {
	Slice    int32  `json:"slice"`
	App      string `json:"app"`
	Workflow string `json:"workflow"`
	ID       int32  `json:"id"`
}

func (q *Queries) GetWFInstanceListForMarkDone(ctx context.Context, arg GetWFInstanceListForMarkDoneParams) ([]Wfinstance, error) {
	rows, err := q.db.Query(ctx, getWFInstanceListForMarkDone,
		arg.Slice,
		arg.App,
		arg.Workflow,
		arg.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wfinstance
	for rows.Next() {
		var i Wfinstance
		if err := rows.Scan(
			&i.ID,
			&i.Entityid,
			&i.Slice,
			&i.App,
			&i.Class,
			&i.Workflow,
			&i.Step,
			&i.Loggedat,
			&i.Doneat,
			&i.Nextstep,
			&i.Parent,
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

const updateWFInstanceDoneat = `-- name: UpdateWFInstanceDoneat :exec
UPDATE public.wfinstance
SET 
    doneat = $1 -- Set doneat to the provided timestamp
WHERE
    id = $2
    AND slice = $3
    AND app = $4
    AND workflow = $5
`

type UpdateWFInstanceDoneatParams struct {
	Doneat   pgtype.Timestamp `json:"doneat"`
	ID       int32            `json:"id"`
	Slice    int32            `json:"slice"`
	App      string           `json:"app"`
	Workflow string           `json:"workflow"`
}

func (q *Queries) UpdateWFInstanceDoneat(ctx context.Context, arg UpdateWFInstanceDoneatParams) error {
	_, err := q.db.Exec(ctx, updateWFInstanceDoneat,
		arg.Doneat,
		arg.ID,
		arg.Slice,
		arg.App,
		arg.Workflow,
	)
	return err
}

const updateWFInstanceStep = `-- name: UpdateWFInstanceStep :exec
UPDATE public.wfinstance
SET step = $1,
doneat = $6
WHERE
    id = $2
    AND slice = $3
    AND app = $4
    AND workflow = $5
`

type UpdateWFInstanceStepParams struct {
	Step     string           `json:"step"`
	ID       int32            `json:"id"`
	Slice    int32            `json:"slice"`
	App      string           `json:"app"`
	Workflow string           `json:"workflow"`
	Doneat   pgtype.Timestamp `json:"doneat"`
}

func (q *Queries) UpdateWFInstanceStep(ctx context.Context, arg UpdateWFInstanceStepParams) error {
	_, err := q.db.Exec(ctx, updateWFInstanceStep,
		arg.Step,
		arg.ID,
		arg.Slice,
		arg.App,
		arg.Workflow,
		arg.Doneat,
	)
	return err
}
