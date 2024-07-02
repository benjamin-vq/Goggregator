-- +goose Up
ALTER TABLE feeds
ADD COLUMN lastFetchedAt TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN lastFetchedAt;