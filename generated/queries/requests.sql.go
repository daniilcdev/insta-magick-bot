// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: requests.sql

package queries

import (
	"context"
	"strings"
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

const deleteRequestsInStatus = `-- name: DeleteRequestsInStatus :exec
DELETE FROM requests
WHERE status = ?
`

func (q *Queries) DeleteRequestsInStatus(ctx context.Context, status string) error {
	_, err := q.db.ExecContext(ctx, deleteRequestsInStatus, status)
	return err
}

const getRequestsInStatus = `-- name: GetRequestsInStatus :many
SELECT file, requester_id FROM requests
WHERE status = ?
`

type GetRequestsInStatusRow struct {
	File        string
	RequesterID string
}

func (q *Queries) GetRequestsInStatus(ctx context.Context, status string) ([]GetRequestsInStatusRow, error) {
	rows, err := q.db.QueryContext(ctx, getRequestsInStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRequestsInStatusRow
	for rows.Next() {
		var i GetRequestsInStatusRow
		if err := rows.Scan(&i.File, &i.RequesterID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const schedulePending = `-- name: SchedulePending :many
UPDATE requests
SET status = "Processing"
WHERE id in (
        SELECT id
        FROM requests
        WHERE status = "Pending"
        LIMIT ?
    )
RETURNING file
`

// weird behaviour: naming parameter doesn't work wit sqlite for some reason
func (q *Queries) SchedulePending(ctx context.Context, limit int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, schedulePending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, err
		}
		items = append(items, file)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRequestsStatus = `-- name: UpdateRequestsStatus :exec
UPDATE requests
SET status = ?
WHERE file in (/*SLICE:filenames*/?)
`

type UpdateRequestsStatusParams struct {
	Status    string
	Filenames []string
}

func (q *Queries) UpdateRequestsStatus(ctx context.Context, arg UpdateRequestsStatusParams) error {
	query := updateRequestsStatus
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Status)
	if len(arg.Filenames) > 0 {
		for _, v := range arg.Filenames {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:filenames*/?", strings.Repeat(",?", len(arg.Filenames))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:filenames*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}
