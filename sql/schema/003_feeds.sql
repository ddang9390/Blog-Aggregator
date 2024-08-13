-- +goose Up
CREATE TABLE feeds(
    Name        VARCHAR(255),
    Url         VARCHAR(255) UNIQUE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;