-- +goose Up
CREATE TABLE users(
    Id          VARCHAR(255) PRIMARY KEY,
    Created_at  TIMESTAMP,
    Updated_at  TIMESTAMP,
    Name        VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;