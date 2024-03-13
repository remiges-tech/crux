-- name: AppNew :many
INSERT INTO
    app (
        realm, shortname, shortnamelc, longname, setby
    )
VALUES (
        @realm, @shortname, @shortnamelc, @longname, @setby
    )
    RETURNING *;

-- name: GetAppName :one
select count(1) FROM app WHERE shortnamelc = $1 AND realm = $2;