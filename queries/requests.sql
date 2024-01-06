-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id)
VALUES (?, ?);

-- name: GetRequest :one
SELECT * FROM requests
WHERE file = ? LIMIT 1;

-- name: DeleteRequest :exec
DELETE FROM requests
WHERE file = ?;
