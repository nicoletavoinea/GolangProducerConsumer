-- name: UpdateTask :one
UPDATE tasks
SET 
    status='PROCESSING', 
    lastupdatetime = strftime('%s','now')
WHERE 
    id=:param1
RETURNING *;