-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    usernames TEXT[] NOT NULL
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE chat;
-- +goose StatementEnd
