-- name: CreateRequest :exec
INSERT INTO requests (file, requester_id, filter_name)
VALUES ($1, $2, $3);

-- weird behaviour: naming parameter doesn't work wit sqlite for some reason
-- name: SchedulePending :many
UPDATE requests
SET status = 'Processing'
WHERE id in (
        SELECT id
        FROM requests
        WHERE status = 'Pending'
        LIMIT $1)
RETURNING file, filter_name;

-- name: GetRequestsInStatus :many
SELECT file, requester_id FROM requests
WHERE status = $1;

-- name: UpdateRequestsStatus :exec
UPDATE requests
SET status = $1
WHERE file in (sqlc.slice('filenames'));

-- name: DeleteRequestsInStatus :exec
DELETE FROM requests
WHERE status = $1;