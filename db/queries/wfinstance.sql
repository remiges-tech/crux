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