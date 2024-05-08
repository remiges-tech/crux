-- name: AppNew :many
INSERT INTO
    app (
        realm, shortname, shortnamelc, longname, setby
    )
VALUES (
        @realm, @shortname, @shortnamelc, @longname, @setby
    )
RETURNING
    *;

-- name: GetAppName :many
select * FROM app WHERE shortnamelc = $1 AND realm = $2;

-- name: AppUpdate :exec
UPDATE app
set
    longname = @longname,
    setat = NOW(),
    setby = @setby
WHERE
    shortnamelc = @shortnamelc
    AND realm = @realm;

-- name: AppExist :one
SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM schema WHERE schema.app = @app) OR
             EXISTS (SELECT 1 FROM ruleset WHERE ruleset.app = @app)
        THEN 1
        ELSE 0
    END AS value_exists ;


-- name: AppDelete :exec
DELETE FROM app
WHERE shortnamelc = @shortnamelc AND realm = @realm;

-- name: GetAppList :many
SELECT
    a.shortnamelc AS name,
    a.longname AS descr,
    a.setat AS createdat,
    a.setby AS createdby,
    ( SELECT COUNT(DISTINCT "user")
        FROM capgrant
        WHERE app = a.shortnamelc
    ) AS nusers,
    ( SELECT COUNT(*)
        FROM ruleset
        WHERE app = a.shortnamelc AND brwf = 'B'
    ) AS nrulesetsbre,
    ( SELECT COUNT(*)
        FROM ruleset
        WHERE app = a.shortnamelc AND brwf = 'W'
    ) AS nrulesetswfe
FROM
    app a
WHERE
a.realm= @realm;


-- name: GetAppNames :many
SELECT shortnamelc FROM app WHERE realm = @realm;