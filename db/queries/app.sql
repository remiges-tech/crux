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