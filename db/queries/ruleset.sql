-- name: GetApp :one
SELECT app
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND brwf = 'W';

-- name: GetClass :one
SELECT class
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND brwf = 'W';

-- name: GetWFActiveStatus :one
SELECT is_active
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND brwf = 'W'
    AND setname = $4;

-- name: GetWFInternalStatus :one
SELECT is_internal
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND brwf = 'W'
    AND setname = $4;

-- name: Workflowget :one
select
    id,
    slice,
    app,
    class,
    setname as name,
    is_active,
    is_internal,
    ruleset as flowrules,
    createdat,
    createdby,
    editedat,
    editedby
from ruleset
where
    slice = $1
    and app = $2
    and class = $3
    and setname = $4
    AND brwf = 'W';

-- name: WorkFlowNew :one
INSERT INTO
    ruleset (
        realm, slice, app, brwf, class, setname, schemaid, is_active, is_internal, ruleset, createdat, createdby
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, $11
    )
RETURNING
    id;

-- name: WorkflowList :many
select
    id,
    slice,
    app,
    class,
    setname as name,
    is_active,
    is_internal,
    createdat,
    createdby,
    editedat,
    editedby
from ruleset
where
    brwf = 'W'
    AND (sqlc.narg('slice')::INTEGER is null OR slice = sqlc.narg('slice')::INTEGER)
    AND (sqlc.narg('app')::VARCHAR(20) is null OR app = sqlc.narg('app')::VARCHAR(20))
    AND (sqlc.narg('class')::VARCHAR(20) is null OR class = sqlc.narg('class')::VARCHAR(20))
    AND (sqlc.narg('setname')::VARCHAR(20) is null OR setname = sqlc.narg('setname')::VARCHAR(20))
    AND (sqlc.narg('is_active')::BOOLEAN is null OR is_active = sqlc.narg('is_active')::BOOLEAN)
    AND (sqlc.narg('is_internal')::BOOLEAN is null OR is_internal = sqlc.narg('is_internal')::BOOLEAN);
