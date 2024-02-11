-- name: CreateRequest :one
INSERT INTO requests (file, requester_id, filter_name)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetRequestsInStatus :many
SELECT file, requester_id FROM requests
WHERE status = $1;

-- name: UpdateRequestStatus :one
UPDATE requests
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteRequest :exec
DELETE FROM requests
WHERE id = $1;