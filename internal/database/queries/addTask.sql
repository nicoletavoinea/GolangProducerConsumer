-- name: AddTask :one
INSERT INTO tasks (type, value,state,creationtime,lastupdatetime)
VALUES (:param1, :param2, 'RECEIVED',strftime('%s','now'),strftime('%s','now')) RETURNING id, creationtime;