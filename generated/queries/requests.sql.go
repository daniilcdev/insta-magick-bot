// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: requests.sql

package queries

import (
	"context"
)

const createRequest = `-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id)
VALUES (?, ?)
`

type CreateRequestParams struct {
	File        string
	RequesterID string
}

func (q *Queries) CreateRequest(ctx context.Context, arg CreateRequestParams) error {
	_, err := q.db.ExecContext(ctx, createRequest, arg.File, arg.RequesterID)
	return err
}

const deleteRequest = `-- name: DeleteRequest :exec
DELETE FROM requests
WHERE file = ?
`

func (q *Queries) DeleteRequest(ctx context.Context, file string) error {
	_, err := q.db.ExecContext(ctx, deleteRequest, file)
	return err
}

const getRequest = `-- name: GetRequest :one
SELECT id, file, requester_id FROM requests
WHERE file = ? LIMIT 1
`

func (q *Queries) GetRequest(ctx context.Context, file string) (Request, error) {
	row := q.db.QueryRowContext(ctx, getRequest, file)
	var i Request
	err := row.Scan(&i.ID, &i.File, &i.RequesterID)
	return i, err
}
