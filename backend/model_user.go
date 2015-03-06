package main

import (
	"database/sql"
	"time"
)

const (
	TABLE_NAME_USER = "users"

	SQL_CREATE_TABLE_USER = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_USER + `(
			id				BIGSERIAL			PRIMARY KEY,
			first_name		VARCHAR(100)		NOT NULL,
			last_name		VARCHAR(100)		NOT NULL,
			email			VARCHAR(255)		UNIQUE NOT NULL,
			hashed_password	VARCHAR(255)		NOT NULL,
			stripe_id		VARCHAR(255)		NOT NULL,
			picture_url		VARCHAR(511),

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// The User model represents people who have accounts
type User struct {
	Id             int64  // The identifier of the user
	FirstName      string // The first name of the user
	LastName       string // The last name of the user
	Email          string // The email address of the user (indexed)
	HashedPassword string // The bcrypted password of the user
	StripeId       string // The id of the user with Stripe's API
	PictureUrl     string // The URL to user's picture

	Active    bool      // True if this entity has not been soft deleted
	CreatedAt time.Time // The time when this user was created
	UpdatedAt time.Time // The time when this user was last updated
	DeletedAt time.Time // The time when this user was soft deleted
}

// Creates the User table if it doesn't already exist
func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_USER)
	return err
}
