-- name: CreateNewSliceBY :one
INSERT INTO
    realmslice (
        realm, descr, active, activateat, deactivateat
    )
SELECT
    realm,
    (
        @ descr::VARCHAR is null
        OR descr = @ descr::VARCHAR
    ),
    -- descr,
    active,
    activateat,
    deactivateat
FROM realmslice
WHERE
    realmslice.id = $1
    AND realmslice.realm = $2
RETURNING
    realmslice.id;

-- name: CopyConfig :execresult
INSERT INTO
    config (
        realm, slice, name, descr, val, ver, setby
    )
SELECT realm, $2, name, descr, val, ver, $3
FROM config
WHERE
    config.slice = $1;

-- name: CopySchema :execresult
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
        @ app::text [] is null
        OR app = any (@ app::text [])
    );

-- name: CopyRuleset :execresult
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
        @ app::text [] is null
        OR app = any (@ app::text [])
    );
