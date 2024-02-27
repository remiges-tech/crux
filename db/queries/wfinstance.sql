-- name: GetWFINstance :one
SELECT count(1)
FROM wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4;

-- name: AddWFNewInstances :many
INSERT INTO
    wfinstance (
        entityid, slice, app, class, workflow, step, loggedat, nextstep, parent
    )
VALUES (
        @entityid, @slice, @app, @class, @workflow, unnest(@step::text []), (NOW()::timestamp), @nextstep, @parent
    )
RETURNING *;

;
-- name: DeleteWfInstance :one
WITH deleted_parents AS (
   DELETE FROM wfinstance
   WHERE
       (id = sqlc.narg('id')::INTEGER OR entityid = sqlc.narg('entityid')::TEXT)
   RETURNING parent
),
deletion_count AS (
   SELECT COUNT(*) AS cnt FROM deleted_parents
),
delete_childrens AS (
    DELETE FROM wfinstance
    WHERE parent IN (SELECT parent FROM deleted_parents WHERE parent IS NOT NULL)
)
SELECT 
    CASE 
        WHEN (SELECT cnt FROM deletion_count) > 0 THEN 1
        ELSE -1 
    END AS result;


-- name: GetWFInstanceList :many
SELECT * FROM wfinstance
WHERE 
   (sqlc.narg('slice')::INTEGER is null OR slice = sqlc.narg('slice')::INTEGER)
   AND (sqlc.narg('entityid')::text is null OR entityid = sqlc.narg('entityid')::text)
   AND (sqlc.narg('app')::text is null OR app = sqlc.narg('app')::text)
   AND (sqlc.narg('workflow')::text is null OR workflow = sqlc.narg('workflow')::text)
   AND(sqlc.narg('parent')::INTEGER is null OR  parent = sqlc.narg('parent')::INTEGER);

    
-- name: GetWFInstanceListByParents :many
SELECT * FROM wfinstance
WHERE 
   (@id::INTEGER[] IS NOT NULL AND id = ANY(@id::INTEGER[]));




-- name: DeleteWfinstanceByID :many
  DELETE FROM wfinstance
   WHERE
       (id = sqlc.narg('id')::INTEGER OR entityid = sqlc.narg('entityid')::TEXT)
   RETURNING *;
    
-- name: DeleteWFInstanceListByParents :many
DELETE FROM wfinstance
WHERE 
   (@id::INTEGER[] IS NOT NULL AND id = ANY(@id::INTEGER[]) OR @parent::INTEGER[] IS NOT NULL AND parent = ANY(@parent::INTEGER[]))
    RETURNING *;