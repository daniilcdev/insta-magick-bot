-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id, status)
VALUES (?, ?, "Pending");

-- name: GetRequest :one
SELECT *
FROM requests
WHERE file = ?
LIMIT 1;

-- name: DeleteRequest :exec
DELETE FROM requests
WHERE file = ?;

-- name: ObtainPendingFiles :many
UPDATE requests
SET status = "Processing"
WHERE id in (
        SELECT id
        FROM requests
        WHERE status = "Pending"
        LIMIT ?
    )
RETURNING file;

-- name: GetRequestersByFilenames :many
SELECT file, requester_id
FROM requests
WHERE file in (sqlc.slice('filenames'));

-- name: UpdateFilesStatus :exec
UPDATE requests
SET status = "Processed"
WHERE file in (sqlc.slice('filenames'));