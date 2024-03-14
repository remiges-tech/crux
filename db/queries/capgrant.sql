-- name: DeleteCapGranForApp :exec

DELETE FROM capgrant where app = @app and realm = @realm;