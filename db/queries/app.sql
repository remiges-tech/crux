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


