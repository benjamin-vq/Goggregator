-- +goose Up
ALTER TABLE users
ADD COLUMN apiKey VARCHAR(64) NOT NULL UNIQUE
default encode(sha256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users
DROP COLUMN apiKey;