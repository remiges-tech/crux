-- name: GetWorkflow :many
SELECT workflow, step FROM stepworkflow WHERE step = $1;