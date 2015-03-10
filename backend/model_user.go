package main

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strings"
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
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;
	`
	SQL_SELECT_USER_BY_ID = `
		SELECT * FROM ` + TABLE_NAME_USER + ` WHERE (id = $1);
	`
	SQL_SELECT_USER_BY_EMAIL = `
		SELECT * FROM ` + TABLE_NAME_USER + ` WHERE (email = $1);
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

	Active    bool        // True if this entity has not been soft deleted
	CreatedAt time.Time   // The time when this user was created
	UpdatedAt time.Time   // The time when this user was last updated
	DeletedAt pq.NullTime // The time when this user was soft deleted
}

// Fills user with data from a db row
func (u User) populateFromRow(row *sql.Row) error {
	// Scan for member fields
	Debug("Populate from row ", *row)
	return row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.HashedPassword, &u.StripeId, &u.PictureUrl, &u.Active, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
}

// Creates the User table if it doesn't already exist
func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_USER)
	return err
}

// Gets a User from the database by id
func GetUser(
	db *sql.DB,
	id int64,
) (*User, error) {
	rows, err := db.Query(SQL_SELECT_USER_BY_ID, id)
	if err != nil {
		return nil, PUBERR_ENTITY_NOT_FOUND
	}
	// Read the rows
	var (
		newId             int64
		newFirstName      string
		newLastName       string
		newEmail          string
		newHashedPassword string
		newStripeId       string
		newPictureUrl     string
		newActive         bool
		newCreatedAt      time.Time
		newUpdatedAt      time.Time
		newDeletedAt      pq.NullTime
	)
	for rows.Next() {
		err = rows.Scan(&newId, &newFirstName, &newLastName, &newEmail, &newHashedPassword, &newStripeId, &newPictureUrl, &newActive, &newCreatedAt, &newUpdatedAt, &newDeletedAt)
		if err != nil {
			return nil, err
		} else {
			return &User{
				Id:             newId,
				FirstName:      newFirstName,
				LastName:       newLastName,
				Email:          newEmail,
				HashedPassword: newHashedPassword,
				StripeId:       newStripeId,
				PictureUrl:     newPictureUrl,
				Active:         newActive,
				CreatedAt:      newCreatedAt,
				UpdatedAt:      newUpdatedAt,
				DeletedAt:      newDeletedAt,
			}, nil
		}
	}
	return nil, PUBERR_ENTITY_NOT_FOUND

}

// Finds a User by email
func FindUserByEmail(
	db *sql.DB,
	email string,
) (*User, error) {
	var newUser User
	if err := newUser.populateFromRow(db.QueryRow(SQL_SELECT_USER_BY_EMAIL, email)); err != nil {
		return nil, PUBERR_ENTITY_NOT_FOUND
	} else {
		return &newUser, nil
	}
}

// Creates a new User in the database; returns the id of the new user
func CreateNewUser(
	db *sql.DB, // The database
	FirstName string, // The first name of the user
	LastName string, // The last name of the user
	Email string, // The email address of the user (indexed)
	HashedPassword string, // The bcrypted password of the user
	StripeId string, // The id of the user with Stripe's API
	PictureUrl string, // The URL to user's picture
) (int64, error) {
	var (
		id  int64
		now = time.Now()
	)

	err := db.QueryRow(SQL_CREATE_NEW_USER, FirstName, LastName, Email, HashedPassword, StripeId, PictureUrl, true, now, now).Scan(&id)
	if err != nil {
		// Check if the issue is email related
		if strings.Contains(err.Error(), "violates unique constraint \"users_email_key\"") {
			return -1, PUBERR_USER_CREATION_FAILED_EMAIL_TAKEN
		} else {
			return -1, errors.New(ERR_USER_CREATION_FAILED + err.Error())
		}
	} else {
		return id, nil
	}
}
