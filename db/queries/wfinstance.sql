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

-- name: UpdateWFInstanceStep :exec
UPDATE public.wfinstance
SET step = $1
WHERE
    entityid = $2
    AND slice = $3
    AND app = $4
    AND workflow = $5;
-- name: UpdateWFInstanceDoneat :exec

UPDATE public.wfinstance
SET 
    doneat = $1 -- Set doneat to the provided timestamp
WHERE
    entityid = $2
    AND slice = $3
    AND app = $4
    AND workflow = $5;


-- name: GetWFInstanceCounts :one
SELECT COUNT(*) 
FROM wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4;

-- name: DeleteWFInstances :exec
DELETE FROM
    public.wfinstance
WHERE
    entityid = $1
    AND slice = $2
    AND app = $3;

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


-- name: GetWFInstanceCurrent :one
 SELECT * FROM wfinstance
WHERE
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4;
