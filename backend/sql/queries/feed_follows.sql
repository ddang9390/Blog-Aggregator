-- name: CreateFeedFollows :one
INSERT INTO feed_follows (feed_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.feed_id = $1;

-- name: GetAllFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE feed_follows.user_id = $1;