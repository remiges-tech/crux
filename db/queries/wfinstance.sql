-- name: GetWFINstance :many

SELECT * 
FROM wfinstance
WHERE slice = $1 
AND app = $2
AND workflow = $3
AND entityid = $4;



