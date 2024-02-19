-- name: SchemaNew :exec
INSERT INTO
    schema(
        realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, $8
    );

-- name: SchemaUpdate :exec
UPDATE schema
SET
    brwf = $4,
    patternschema = $5,
    actionschema = $6,
    editedat = CURRENT_TIMESTAMP,
    editedby = $7
WHERE
    slice = $1
    AND class = $2
    AND app = $3;

-- name: GetSchemaWithLock :one
SELECT
    id,
    brwf,
    patternschema,
    actionschema,
    editedat = CURRENT_TIMESTAMP,
    editedby
FROM schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3 FOR
UPDATE;

-- name: SchemaDelete :one
DELETE FROM schema WHERE id = $1 RETURNING id;

DELETE FROM schema WHERE id = $1 RETURNING id;

-- name: SchemaList :many
SELECT schema.slice, realmslice.descr, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id;

-- name: SchemaListByApp :many
SELECT schema.slice, realmslice.descr, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
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
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    slice = $1
    AND class = $2
    AND app = $3;

-- name: Wfschemaget :one
SELECT s.slice, s.app, s.class, rm.longname, s.patternschema, s.actionschema, s.createdat, s.createdby, s.editedat, s.editedby
FROM schema as s, realm as rm, realmslice as rs
WHERE
    s.realm = rm.id
    and rs.realm = rm.shortname
    and s.slice = rs.id
    and s.slice = $1
    and rm.id = @realm
    and s.class = $3
    AND s.app = $2;

-- name: Wfschemadelete :exec
DELETE from schema
where
    id in (
        SELECT schema.id
        FROM schema, realm, realmslice
        WHERE
            schema.realm = realm.id
            and schema.slice = realmslice.id
            and schema.slice = $1
            and realmslice.realm = realm.shortname
            and realm.id = @realm
            and schema.class = $3
            AND schema.app = $2
    );

-- name: WfPatternSchemaGet :one
SELECT patternschema
FROM public.schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3
    AND realm = $4;

-- name: WfSchemaGet :one
SELECT *
FROM public.schema
WHERE
    slice = $1
    AND class = $2
    AND app = $3;

-- name: WfSchemaList :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema, app, realmslice
where
schema.app = app.shortname
and schema.slice = realmslice.id
    AND schema.realm =  @relam
    AND ((sqlc.narg('slice')::INTEGER is null) OR (schema.slice = @slice::INTEGER))
    AND ((sqlc.narg('app')::text is null) OR (schema.app = @app::text))
    AND (sqlc.narg('class')::text is null OR schema.class = sqlc.narg('class')::text);
