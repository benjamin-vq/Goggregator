-- name: CreatePost :exec
INSERT INTO posts (id, title, url, description, publishedat, feed_id, createdat, updatedat)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
JOIN feeds f on f.id = posts.feed_id
JOIN users u on f.user_id = u.id
WHERE u.id = $1
ORDER BY publishedat DESC
LIMIT $2;