-- name: SchemaNew :one 
-- :one
INSERT INTO schema (
    realm, slice, app, brwf, class, patternschema, actionschema, createdby, editedby
) VALUES (
    1, $1, $2, W, $3, $4, $5, $6, $7
) RETURNING *;

-- name: SchemaUpdate :one
-- :one
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
RETURNING *;

-- name: SchemaDelete :one
-- :one
DELETE FROM schema
WHERE
    id = $1
RETURNING *;

-- name: SchemaGet :one
-- :one
SELECT *
FROM schema
WHERE
    id = $1;

-- name: SchemaList :many
SELECT *
FROM schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3;
