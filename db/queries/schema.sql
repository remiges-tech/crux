-- name: SchemaNew :one
INSERT INTO
    schema (
        realm, slice, app, brwf, class, patternschema, actionschema, createdby
    )
VALUES (
        1, $1, $2, 'W', $3, $4, $5, $6
    )
RETURNING
    id;

-- name: SchemaUpdate :one
UPDATE schema
SET
    brwf = 'W',
    patternschema = $4,
    actionschema = $5,
    editedat = CURRENT_TIMESTAMP,
    editedby = $6
WHERE
    slice = $1
    AND class = $2
    AND app = $3
RETURNING
    id;

-- name: SchemaDelete :one
DELETE FROM schema WHERE id = $1 RETURNING id;

-- name: SchemaList :many
SELECT schema.slice, realmslice.descr, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id;

-- name: SchemaListByApp :many
SELECT schema.slice,realmslice.descr, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    app = $1;

-- name: SchemaListByClass :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    class = $1;

-- name: SchemaListBySlice :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    slice = $1;

-- name: SchemaListByAppAndClass :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    app = $1
    AND class = $2;

-- name: SchemaListByAppAndSlice :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    app = $1
    AND slice = $2;

-- name: SchemaListByClassAndSlice :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    class = $1
    AND slice = $2;

-- name: SchemaGet :many
SELECT *
FROM schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3;


-- name: WfPatternSchemaGet :one
 SELECT patternschema
 FROM public.schema
 WHERE 
 slice = $1
    AND class = $2
    AND app = $3
    AND brwf = 'W';
