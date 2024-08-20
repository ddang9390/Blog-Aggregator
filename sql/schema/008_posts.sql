-- +goose Up
CREATE TABLE posts(
    id                      UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at              TIMESTAMP,
    updated_at              TIMESTAMP,
    title                   VARCHAR(255),
    url                     TEXT UNIQUE,
    description             TEXT,
    published_at            TIMESTAMP,
    feed_id                 VARCHAR(255) NOT NULL,
    FOREIGN KEY(feed_id)    REFERENCES feed_follows(feed_id)
);

-- +goose Down
DROP TABLE posts;