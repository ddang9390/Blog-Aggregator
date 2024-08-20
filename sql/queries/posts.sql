-- name: CreatePost :one
INSERT INTO posts(created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULLIF($1,''), $2, NULLIF($3,''), NULLIF($4,''), $5)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.*
FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;