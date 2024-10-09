-- name: AddTask :one
INSERT INTO tasks (type, value,state,creationtime,lastupdatetime)
VALUES ($1, $2, 'RECEIVED',EXTRACT(EPOCH FROM NOW()),EXTRACT(EPOCH FROM NOW())) RETURNING id, creationtime;