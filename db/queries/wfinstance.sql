-- name: GetWFINstance :many

SELECT * 
FROM wfinstance
WHERE slice = $1 
AND app = $2
AND workflow = $3
AND entityid = $4;


-- name: AddWFNewInstace :one
INSERT INTO 
wfinstance
(entityid,slice, app, class, workflow, step,loggedat, nextstep)
VALUES ($1,$2,$3,$4,$5,$6,(NOW() :: timestamp),$7) 
RETURNING id;

