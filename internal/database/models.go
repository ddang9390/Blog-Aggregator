// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
)

type Feed struct {
	Name   sql.NullString
	Url    sql.NullString
	UserID string
}

type User struct {
	ID        string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Name      string
	Apikey    string
}
