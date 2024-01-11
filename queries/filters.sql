-- name: GetNames :many
SELECT name FROM filters;

-- name: GetReceiptOrDefault :one
SELECT id, name, receipt FROM filters
WHERE name = ? OR id = 1;

-- name: CreateReceipt :exec
INSERT INTO filters (name, receipt)
VALUES (?, ?);