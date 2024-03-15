-- name: DeleteCapGranForApp :exec

DELETE FROM capgrant WHERE app = @app AND realm = @realm;

-- name: GetCapGrantForApp :many

SELECT * FROM capgrant WHERE app = @app AND realm = @realm;