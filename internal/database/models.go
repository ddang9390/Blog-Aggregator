// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            string
	Name          string
	Url           string
	UserID        string
	LastFetchedAt sql.NullTime
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FeedFollow struct {
	FeedID string
	UserID string
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	Title       sql.NullString
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      string
}

type User struct {
	ID        string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Name      string
	Apikey    string
}
