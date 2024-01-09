-- name: GetNames :many
SELECT name FROM filters;

-- name: GetReceipt :one
SELECT id, name, receipt FROM filters
WHERE name = ?;

-- name: CreateReceipt :exec
INSERT INTO filters (name, receipt)
VALUES (?, ?);