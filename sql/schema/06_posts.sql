-- +goose Up
CREATE TABLE posts (
    id          uuid PRIMARY KEY,
    title       VARCHAR(256) NOT NULL,
    url         TEXT UNIQUE  NOT NULL,
    description TEXT,
    publishedAt TIMESTAMP,
    feed_id     uuid         NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    createdAt   TIMESTAMP    NOT NULL,
    updatedAt   TIMESTAMP    NOT NULL
);

-- +goose Down
DROP TABLE posts;