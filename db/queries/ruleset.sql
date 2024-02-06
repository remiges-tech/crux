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
    and setname = $4;

-- name: WorkFlowNew :one
INSERT INTO
    ruleset (
        realm, slice, app, brwf, class, setname, is_active, is_internal, ruleset, createdat, createdby
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, $10
    )
RETURNING
    id;