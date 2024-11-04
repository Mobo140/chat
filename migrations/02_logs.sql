-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    activity VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
