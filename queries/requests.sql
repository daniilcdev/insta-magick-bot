-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id, filter_name)
VALUES (?, ?, ?);

-- weird behaviour: naming parameter doesn't work wit sqlite for some reason
-- name: SchedulePending :many
UPDATE requests
SET status = "Processing"
WHERE id in (
        SELECT id
        FROM requests
        WHERE status = "Pending"
        LIMIT ?
    )
RETURNING file, filter_name;

-- name: GetRequestsInStatus :many
SELECT file, requester_id FROM requests
WHERE status = ?;

-- name: UpdateRequestsStatus :exec
UPDATE requests
SET status = ?
WHERE file in (sqlc.slice('filenames'));

-- name: DeleteRequestsInStatus :exec
DELETE FROM requests
WHERE status = ?;