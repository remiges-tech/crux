-- name: DeleteCapGranForApp :exec

DELETE FROM capgrant WHERE app = @app AND realm = @realm AND "user" = @userId;

-- name: GetCapGrantForApp :many

SELECT * FROM capgrant WHERE app = @app AND realm = @realm AND "user" = @userId;

-- name: UserActivate :one
UPDATE capgrant
SET
    "from" = CASE
        WHEN (
            sqlc.narg ('activateat')::TIMESTAMP
        ) IS NULL THEN NOW()
        ELSE (
            sqlc.narg ('activateat')::TIMESTAMP
        )
    END,
    "to" = NULL
WHERE
    "user" = @userid
    and realm = @realm
RETURNING *;

-- name: UserDeactivate :one
UPDATE capgrant
SET
    "to" = CASE
        WHEN (
            sqlc.narg ('deactivateat')::TIMESTAMP
        ) IS NULL THEN NOW()
        ELSE (
            sqlc.narg ('deactivateat')::TIMESTAMP
        )
    END,
    "from" = NULL
WHERE
    "user" = @userid
    and realm = @realm
RETURNING *;

-- name: CapGet :many
SELECT app,cap,setby,setat,"from","to" from capgrant WHERE realm = @realm and "user" = @userId;

-- name: CapList :many
SELECT "user",app,cap,"from","to",setat,setby from capgrant
WHERE realm = @realm
and ((@app::text[] is null) OR ( app = any(@app::text[])))
and ((@cap::text[] is null) OR ( cap = any(@cap::text[])));

-- name: UpdateCapGranForUser :exec
UPDATE capgrant set cap = NULL WHERE "user" = @userId;

-- name: GetUserRealm :many
SELECT  realm  FROM capgrant  WHERE "user" = @userId;

-- name: GrantRealmCapability :exec
INSERT INTO capgrant (realm,"user",cap,"from","to",setat,setby)
VALUES(@realm, @userId,unnest(@cap::text []), sqlc.narg ('from') ,sqlc.narg('to'),(NOW() AT TIME ZONE 'UTC'),@setby);

-- name: GrantAppCapability :exec
INSERT INTO capgrant (realm, "user", app, cap, "from", "to", setat, setby)
SELECT
    @realm,
    @userId,
    app,
    cap,
    sqlc.narg('from'),
    sqlc.narg('to'),
    NOW() AT TIME ZONE 'UTC',
    @setby
FROM
    unnest(@app::text[]) AS app,
    unnest(@cap::text[]) AS cap;


-- name: CapRevoke :execresult
DELETE FROM capgrant 
WHERE ((sqlc.narg('cap')::text [] is null) OR (capgrant.cap = any (@cap::text [])))
AND ((sqlc.narg('app')::text [] is null) OR (capgrant.app = any (@app::text [])))
AND capgrant.user = $1;


-- name: CountOfRootCapUser :one
SELECT count(1)
FROM capgrant
WHERE cap = 'root';

-- name: UserExists :one
SELECT count(*)
FROM capgrant
WHERE "user" = $1;

-- name: AppExists :one
SELECT count(*)
FROM capgrant
WHERE @app::text [];

-- name: CapExists :one
SELECT count(*)
FROM capgrant
WHERE @cap::text [];