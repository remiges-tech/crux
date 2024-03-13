-- name: CloneRecordInRealmSliceBySliceID :one
INSERT INTO
    realmslice (
        realm, descr, active, activateat, deactivateat
    )
SELECT
    realm,
    COALESCE(descr, sqlc.narg('descr')::text),
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
    realmslice (
        realm, descr, active
    )
VALUES (
    $1, $2, true
)
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

