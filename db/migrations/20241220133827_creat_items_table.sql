-- +goose Up
-- +goose StatementBegin
CREATE TABLE items
(
    id    integer PRIMARY KEY,
    name  VARCHAR(255)   NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
