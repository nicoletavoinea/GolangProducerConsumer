-- name: GetNumberOfTasks :one
SELECT COUNT(*)
FROM tasks
WHERE state=:param1