-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, createdat, updatedat)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, user_id, feed_id, createdat, updatedat)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2
RETURNING *;

-- name: GetFeedFollowsByUser :many
SELECT * FROM feed_follows
WHERE user_id = $1;