-- name: GetWFINstance :one
SELECT count(1)
FROM public.wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4;

-- name: AddWFNewInstances :many
INSERT INTO
    public.wfinstance (
        entityid, slice, app, class, workflow, step, loggedat, nextstep, parent
    )
VALUES (
        @entityid, @slice, @app, @class, @workflow, unnest(@step::text []), (NOW()::timestamp), @nextstep, @parent
    )
RETURNING
    id,
    loggedat,
    step;

;
-- name: DeleteWfInstance :one
WITH deleted_parents AS (
   DELETE FROM public.wfinstance
   WHERE
       (id = sqlc.narg('id')::INTEGER OR entityid = sqlc.narg('entityid')::TEXT)
   RETURNING parent
),
deletion_count AS (
   SELECT COUNT(*) AS cnt FROM deleted_parents
),
delete_childrens AS (
    DELETE FROM public.wfinstance
    WHERE parent IN (SELECT parent FROM deleted_parents WHERE parent IS NOT NULL)
)
SELECT 
    CASE 
        WHEN (SELECT cnt FROM deletion_count) > 0 THEN 1
        ELSE -1 
    END AS result;

