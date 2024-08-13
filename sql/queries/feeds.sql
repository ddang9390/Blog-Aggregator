-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- -- name: GetFeed :one
-- SELECT * FROM users WHERE apikey=$1;