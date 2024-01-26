// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: filters.sql

package queries

import (
	"context"
)

const createReceipt = `-- name: CreateReceipt :exec
INSERT INTO filters (id, name, receipt)
VALUES (DEFAULT, $1, $2)
`

type CreateReceiptParams struct {
	Name    string
	Receipt string
}

func (q *Queries) CreateReceipt(ctx context.Context, arg CreateReceiptParams) error {
	_, err := q.db.ExecContext(ctx, createReceipt, arg.Name, arg.Receipt)
	return err
}

const getDefaultReceipt = `-- name: GetDefaultReceipt :one
SELECT id, name, receipt
FROM filters
WHERE id = 1
LIMIT 1
`

func (q *Queries) GetDefaultReceipt(ctx context.Context) (Filter, error) {
	row := q.db.QueryRowContext(ctx, getDefaultReceipt)
	var i Filter
	err := row.Scan(&i.ID, &i.Name, &i.Receipt)
	return i, err
}

const getNames = `-- name: GetNames :many
SELECT
    name
FROM
    filters
`

func (q *Queries) GetNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReceiptWithName = `-- name: GetReceiptWithName :one
SELECT id, name, receipt
FROM filters
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetReceiptWithName(ctx context.Context, name string) (Filter, error) {
	row := q.db.QueryRowContext(ctx, getReceiptWithName, name)
	var i Filter
	err := row.Scan(&i.ID, &i.Name, &i.Receipt)
	return i, err
}
