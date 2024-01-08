-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id)
VALUES (?, ?);

-- name: SchedulePending :many
UPDATE requests
SET status = "Processing"
WHERE id in (
        SELECT id
        FROM requests
        WHERE status = "Pending"
        LIMIT ?
    )
RETURNING file;

-- name: ObtainCompleted :many
SELECT file, requester_id FROM requests
WHERE status = "Processed";

-- name: UpdateFilesStatus :exec
UPDATE requests
SET status = "Processed"
WHERE file in (sqlc.slice('filenames'));

-- name: DeleteCompletedRequests :exec
DELETE FROM requests
WHERE status = "Processed";