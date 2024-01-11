-- name: GetNames :many
SELECT
    name
FROM
    filters;

-- name: GetReceiptWithName :one
SELECT
    id,
    name,
    receipt
FROM
    filters
WHERE
    name = ?
LIMIT
    1;

-- name: GetDefaultReceipt :one
SELECT
    id,
    name,
    receipt
FROM
    filters
WHERE
    id = 1
LIMIT
    1;

-- name: CreateReceipt :exec
INSERT INTO
    filters (name, receipt)
VALUES
    (?, ?);