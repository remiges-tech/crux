-- name: DeleteCapGranForApp :exec

DELETE FROM capgrant WHERE app = @app AND realm = @realm;

-- name: GetCapGrantForApp :many

SELECT * FROM capgrant WHERE app = @app AND realm = @realm;

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
RETURNING *;