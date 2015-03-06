package main

import (
	"database/sql"
	"errors"
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
	SQL_CREATE_NEW_USER = `
		INSERT INTO ` + TABLE_NAME_USER + `
		(first_name, last_name, email, hashed_password, stripe_id, picture_url, active, created_at, updated_at) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	ERR_USER_CREATION_FAILED = "Could not create new user: "
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

// Fills user with data from a db row
func (u User) populateFromRow(row *sql.Row) error {
	// Scan for member fields
	return row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.HashedPassword, &u.StripeId, &u.PictureUrl, &u.Active, &u.CreatedAt, &u.UpdatedAt)
}

// Creates the User table if it doesn't already exist
func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_USER)
	return err
}

// Creates a new User in the database
func CreateNewUser(
	db *sql.DB, // The database
	FirstName string, // The first name of the user
	LastName string, // The last name of the user
	Email string, // The email address of the user (indexed)
	HashedPassword string, // The bcrypted password of the user
	StripeId string, // The id of the user with Stripe's API
	PictureUrl string, // The URL to user's picture
) (*User, error) {
	var newUser User
	newRow := db.QueryRow(SQL_CREATE_NEW_USER, FirstName, LastName, Email, HashedPassword, StripeId, PictureUrl, true, time.Now(), time.Now(), time.Now())
	if err := newUser.populateFromRow(newRow); err != nil {
		return nil, errors.New(ERR_USER_CREATION_FAILED + err.Error())
	} else {
		return &newUser, nil
	}
}
