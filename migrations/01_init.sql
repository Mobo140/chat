-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    usernames TEXT[] NOT NULL
);

CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL, 
    from_user VARCHAR(255) NOT NULL, 
    text TEXT NOT NULL, 
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(), 
    FOREIGN KEY (chat_id) REFERENCES chat(id) ON DELETE CASCADE 
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE chat;
DROP TABLE message;
-- +goose StatementEnd
