-- name: CreateFeed :one
INSERT INTO feeds (name, url)
VALUES ($1, $2)
RETURNING *;

-- -- name: GetFeed :one
-- SELECT * FROM users WHERE apikey=$1;