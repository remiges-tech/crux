-- name: GetWorkflow :one 
SELECT workflow FROM stepworkflow
WHERE step =$1;

