-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;