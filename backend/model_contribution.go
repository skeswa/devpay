package main

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// The Contribution model represents an amount paid by a user to a Campaign
type Contribution struct {
	Id     int64   // The identifier of the contribution
	Amount float64 // The amount of the contribution

	ContributorId int64 // The id of the contributor; Foreign key for User (belongs to)
	CampaignId    int64 // The id of the campaign; Foreign key for the Campaign (belongs to)

	Active    bool        // True if this entity has not been soft deleted
	CreatedAt time.Time   // The time when this contribution was created
	UpdatedAt time.Time   // The time when this contribution was last updated
	DeletedAt pq.NullTime // The time when this user was soft deleted
}

const (
	TABLE_NAME_CONTRIBUTION = "contributions"

	SQL_CREATE_TABLE_CONTRIBUTION = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_CONTRIBUTION + `(
			id			BIGSERIAL		PRIMARY KEY,
			amount		REAL			NOT NULL,
			stripe_id	VARCHAR(255)	NOT NULL,

			contributor_id	BIGINT REFERENCES ` + TABLE_NAME_USER + `(id)		NOT NULL,
			campaign_id		BIGINT REFERENCES ` + TABLE_NAME_CAMPAIGN + `(id)	NOT NULL,

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// Creates the Contribution table if it doesn't already exist
func CreateContributionTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CONTRIBUTION)
	return err
}
