-- name: SchemaNew :one
-- :one
INSERT INTO
    schema(
        realm, slice, app, brwf, class, patternschema, actionschema, createdby, editedby
    )
VALUES (
        1, $1, $2, W, $3, $4, $5, $6, $7
    ) RETURNING id;

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
    id = $1 RETURNING id;

-- name: SchemaDelete :one
-- :one
DELETE FROM schema WHERE id = $1 RETURNING id;

-- name: SchemaList :many
SELECT slice,app,class,createdby,createdat,editedby,editedat FROM schema;
-- name: SchemaListByApp :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE app = $1;

-- name: SchemaListByClass :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE class = $1;

-- name: SchemaListBySlice :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE slice = $1;

-- name: SchemaListByAppAndClass :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE app = $1 AND class = $2;

-- name: SchemaListByAppAndSlice :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE app = $1 AND slice = $2;

-- name: SchemaListByClassAndSlice :many
SELECT slice,app,class,createdby,createdat,editedby,editedat  FROM schema WHERE class = $1 AND slice = $2;

-- name: SchemaGet :many
SELECT *
FROM schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3;