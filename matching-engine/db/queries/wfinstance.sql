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
SELECT COUNT(*) AS instance_count
FROM public.wfinstance
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
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4
    AND parent = $5;
-- name: GetWFInstanceCurrent :one
 SELECT * FROM wfinstance
WHERE 
    slice = $1
    AND app = $2
    AND workflow = $3
    AND entityid = $4;


