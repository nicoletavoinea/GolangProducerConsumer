-- name: GetValueOfTasksByType :many
SELECT type, SUM(value) AS values_sum
FROM tasks
GROUP BY type