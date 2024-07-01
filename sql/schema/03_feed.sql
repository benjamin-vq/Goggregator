-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE feeds;