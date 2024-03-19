-- name: SchemaNew :one
INSERT INTO
    schema(
        realm, slice, app, brwf, class, patternschema, actionschema, createdat, createdby
    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, $8
    ) RETURNING id;

-- name: SchemaUpdate :exec
UPDATE schema
SET
    brwf = $5,
    patternschema = $6,
    actionschema = $7,
    editedat = CURRENT_TIMESTAMP,
    editedby = $8
WHERE
    realm = $1
    AND slice = $2
    AND class = $3
    AND app = $4;

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
    realm = $1
    AND slice = $2
    AND class = $3
    AND app = $4 FOR
UPDATE;

-- name: SchemaDelete :one
DELETE FROM schema WHERE id = $1 RETURNING id;



-- name: SchemaGet :many
SELECT schema.slice, schema.app, app.longname, schema.class, schema.createdby, schema.createdat, schema.editedby, schema.editedat
FROM schema
    JOIN app ON schema.app = app.shortname
    JOIN realmslice on schema.slice = realmslice.id
WHERE
    schema.realm = $1
    AND schema.slice = $2
    AND schema.class = $3
    AND schema.app = $4;

-- name: Wfschemaget :one
SELECT s.slice, s.app, s.class, rm.longname, s.patternschema, s.actionschema, s.createdat, s.createdby, s.editedat, s.editedby
FROM schema as s, realm as rm, realmslice as rs
WHERE
    s.realm = rm.shortname
    and rs.realm = rm.shortname
    and s.slice = rs.id
    and s.slice = $1
    and rm.shortname = @realm
    and s.class = $3
    AND s.app = $2;

-- name: Wfschemadelete :exec
DELETE from schema
where
    id in (
        select id
        from (
                SELECT schema.id
                FROM schema, realm, realmslice
                WHERE
                    schema.realm = realm.id
                    and schema.slice = realmslice.id
                    and schema.slice = $1
                    and realmslice.realm = realm.shortname
                    and schema.realm = @realm
                    and schema.class = $3
                    AND schema.app = $2
            ) as id
        where
            id not in(
                SELECT schemaid
                FROM ruleset
                where
                    realm = @realm
                    and slice = $1
                    and app = $2
                    and class = $3
            )
    );

-- name: WfPatternSchemaGet :one
SELECT patternschema
FROM public.schema
WHERE
    realm = $1
    AND slice = $2
    AND class = $3
    AND app = $4
    AND realm = $5
    AND brwf = 'W';

-- name: WfSchemaGet :one
SELECT *
FROM public.schema
WHERE
    realm = $1
    AND slice = $2
    AND class = $3
    AND app = $4;

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

-- name: AllSchemas :many
SELECT * FROM public.schema;