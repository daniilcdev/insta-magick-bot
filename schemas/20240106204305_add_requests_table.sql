-- +goose Up
CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    file TEXT NOT NULL,
    requester_id TEXT NOT NULL
);

-- +goose Down
DROP TABLE requests;