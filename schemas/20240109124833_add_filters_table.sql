-- +goose Up
-- +goose StatementBegin
CREATE TABLE filters (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    receipt TEXT NOT NULL
);
-- insert default filter
INSERT INTO filters (name, receipt)
VALUES ('Bright Summer',
'-adaptive-sharpen 10% -channel B -evaluate add 1.31 -channel G -evaluate add 1.37 +channel -modulate 120,142 -contrast-stretch -13%x-17% -enhance');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE filters;
-- +goose StatementEnd
