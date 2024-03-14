-- name: CloneRecordInRealmSliceBySliceID :one
INSERT INTO
    realmslice (
        realm, descr, active, activateat, deactivateat
    )
SELECT
    realm,
    COALESCE(
        descr, sqlc.narg ('descr')::text
    ),
    true,
    activateat,
    deactivateat
FROM realmslice
WHERE
    realmslice.id = $1
    AND realmslice.realm = $2
RETURNING
    realmslice.id;

-- name: InsertNewRecordInRealmSlice :one
INSERT INTO
    realmslice (realm, descr, active)
VALUES ($1, $2, true)
RETURNING
    realmslice.id;

-- name: CloneRecordInConfigBySliceID :execresult
INSERT INTO
    config (
        realm, slice, name, descr, val, ver, setby
    )
SELECT realm, $2, name, descr, val, ver, $3
FROM config
WHERE
    config.slice = $1;

-- name: CloneRecordInSchemaBySliceID :execresult
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
        @app::text [] is null
        OR app = any (@app::text [])
    );

-- name: CloneRecordInRulesetBySliceID :execresult
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
        @app::text [] is null
        OR app = any (@app::text [])
    );

-- name: RealmSliceAppsList :many
SELECT a.shortname, a.longname
FROM realmslice
    JOIN app a ON realmslice.realm = a.realm
WHERE
    realmslice.id = $1;

-- name: RealmSlicePurge :execresult
WITH
    cte1 AS (
        DELETE FROM stepworkflow
    ),
    cte2 AS (
        DELETE FROM wfinstance
    ),
    cte3 AS (
        DELETE FROM ruleset
    ),
    cte4 AS (
        DELETE FROM schema
    ),
    cte5 AS (
        DELETE FROM config
    )
DELETE FROM realmslice
WHERE
    id IN (
        SELECT id
        FROM realmslice
        LIMIT 100
    );


-- name: RealmSliceActivate :one
UPDATE realmslice
SET
    active = @isactive,
    activateat = CASE
        WHEN (sqlc.narg('activateat')::TIMESTAMP) IS NULL
            THEN NOW()
        ELSE (sqlc.narg('activateat')::TIMESTAMP)
    END,
    deactivateat = NULL
WHERE
    id = @id
    RETURNING *;

-- name: RealmSliceDeactivate :one
UPDATE realmslice
SET
    active = @isactive,
    deactivateat = CASE
        WHEN (sqlc.narg('deactivateat')::TIMESTAMP) IS NULL
            THEN NOW()
        ELSE (sqlc.narg('deactivateat')::TIMESTAMP)
    END,
    activateat = NULL
WHERE
    id = @id
    RETURNING *;
