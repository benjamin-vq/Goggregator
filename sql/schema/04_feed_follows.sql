-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;