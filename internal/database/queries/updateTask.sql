-- name: UpdateTask :one
UPDATE tasks
SET 
    state=$2::task_state, 
    lastupdatetime = EXTRACT(EPOCH FROM NOW())
WHERE 
    id=$1
RETURNING *;