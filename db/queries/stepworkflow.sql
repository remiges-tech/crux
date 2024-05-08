-- name: GetWorkflowNameForStep :one
SELECT workflow,step FROM stepworkflow WHERE step = $1;