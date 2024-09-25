-- name: UpdateTask :one
UPDATE tasks
SET 
    state=:param2, 
    lastupdatetime = strftime('%s','now')
WHERE 
    id=:param1
RETURNING *;