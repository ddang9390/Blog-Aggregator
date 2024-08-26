-- +goose Up
ALTER TABLE users
ADD CONSTRAINT name UNIQUE (name);


-- +goose Down
ALTER TABLE users
DROP CONSTRAINT name;