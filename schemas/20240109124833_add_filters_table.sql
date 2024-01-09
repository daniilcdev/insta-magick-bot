-- +goose Up
-- +goose StatementBegin
CREATE TABLE filters (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    receipt TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE filters;
-- +goose StatementEnd
