-- +goose Up
ALTER TABLE requests
ADD status TEXT NOT NULL DEFAULT 'Pending';

-- +goose Down
ALTER TABLE requests DROP COLUMN status;