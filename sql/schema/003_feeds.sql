-- +goose Up
CREATE TABLE feeds(
    Name        VARCHAR(255),
    Url         VARCHAR(255) UNIQUE,
    user_id     VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;