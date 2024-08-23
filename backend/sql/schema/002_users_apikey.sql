-- +goose Up
ALTER TABLE users
ADD apikey VARCHAR(64) NOT NULL UNIQUE DEFAULT ENCODE(sha256(CAST(CAST(random() AS text) AS bytea)), 'hex');

-- +goose Down
ALTER TABLE users
DROP COLUMN apikey; 