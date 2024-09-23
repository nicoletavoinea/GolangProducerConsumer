// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: updateTask.sql

package database

import (
	"context"
)

const updateTask = `-- name: UpdateTask :one
UPDATE tasks
SET 
    status='PROCESSING', 
    lastupdatetime = strftime('%s','now')
WHERE 
    id=?1
RETURNING id, type, value, state, creationtime, lastupdatetime
`

func (q *Queries) UpdateTask(ctx context.Context, param1 int64) (Task, error) {
	row := q.db.QueryRowContext(ctx, updateTask, param1)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Value,
		&i.State,
		&i.Creationtime,
		&i.Lastupdatetime,
	)
	return i, err
}
