package main

import (
	"database/sql"
	"time"
)

const (
	TABLE_NAME_CAMPAIGN = "campaigns"

	SQL_CREATE_TABLE_CAMPAIGN = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_CAMPAIGN + `(
			id						BIGSERIAL			PRIMARY KEY,
			title					VARCHAR(255)		NOT NULL,
			description				TEXT				NOT NULL,
			cover_picture_url		VARCHAR(511)		NOT NULL,
			thumbnail_picture_url	VARCHAR(511)		NOT NULL,
			amount					REAL				NOT NULL,
			deadline				TIMESTAMPTZ			NOT NULL,
			finished				BOOLEAN				NOT NULL,

			creator_id	BIGINT REFERENCES ` + TABLE_NAME_USER + `(id)	NOT NULL,
			claimer_id	BIGINT REFERENCES ` + TABLE_NAME_USER + `(id),

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// The Campaign model represents a funding effort with a clear goal and a deadline
type Campaign struct {
	Id                  int64     // The identifier of the campaign
	Title               string    // The title of the campaign
	Description         string    // The description of the campaign
	CoverPictureUrl     string    // The URL of this campaign's cover picture
	ThumbnailPictureUrl string    // The URL of this campaign's thumbnail picture
	Amount              float64   // The current amount that this campaign has raised
	Deadline            time.Time // When this campaign expires
	Finished            bool      // True if the campaign is over

	Creator       User           // The person who started this campaign; One-To-Many relationship (has one)
	CreatorId     int64          // The id of the creator; Foreign key for User (belongs to)
	Claimer       User           // The person who successfully claimed the Campaign; One-To-Many relationship (has one)
	ClaimerId     sql.NullInt64  // The id of the person who successfully claimed the Campaign; Foreign key for User (belongs to)
	Contributions []Contribution // All the contributions to this campaign; One-To-Many relationship (has many)
	Claims        []Claim        // All the claims for this campaign; One-To-Many relationship (has many)

	Active    bool      // True if this entity has not been soft deleted
	CreatedAt time.Time // The time when this campaign was created
	UpdatedAt time.Time // The time when this campaign was last updated
	DeletedAt time.Time // The time when this campaign was soft deleted
}

// Creates the Campaign table if it doesn't already exist
func CreateCampaignTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CAMPAIGN)
	return err
}
