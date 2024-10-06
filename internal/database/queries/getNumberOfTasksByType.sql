-- name: GetNumberOfTasksByType :many
SELECT type, COUNT(*) AS task_count
FROM tasks
GROUP BY type;
