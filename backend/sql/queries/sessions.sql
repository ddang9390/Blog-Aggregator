-- name: CreateSession :one
INSERT INTO sessions(session_id, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE session_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions 
WHERE session_id = $1;