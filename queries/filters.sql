-- name: GetNames :many
SELECT name FROM filters;

-- name: GetReceipt :one
SELECT receipt FROM filters
WHERE name = ?;

-- name: CreateReceipt :exec
INSERT INTO filters (name, receipt)
VALUES (?, ?);