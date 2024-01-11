-- +goose Up
-- +goose StatementBegin
ALTER TABLE requests
ADD filter_name TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE requests DROP COLUMN filter_name;
-- +goose StatementEnd