-- +goose Up
CREATE TABLE sessions(
    Session_id     VARCHAR(255) NOT NULL PRIMARY KEY,
    User_id        VARCHAR(255) NOT NULL,
    Created_at  TIMESTAMP NOT NULL,
    Expires_at  TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE sessions;