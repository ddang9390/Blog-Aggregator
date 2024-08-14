-- +goose Up
CREATE TABLE feed_follows(
    feed_id                 VARCHAR(255) NOT NULL UNIQUE,
    user_id                 VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id)    REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;