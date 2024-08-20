-- +goose Up
ALTER TABLE feeds
ADD COLUMN created_at TIMESTAMP NOT NULL,
ADD COLUMN updated_at TIMESTAMP NOT NULL;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN created_at,
DROP COLUMN updated_at;