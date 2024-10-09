-- name: GetNumberOfTasks :one
SELECT COUNT(*)
FROM tasks
WHERE state=$1::task_state;