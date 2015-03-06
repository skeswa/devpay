package main

import (
	"database/sql"
	"time"
)

const (
	TABLE_NAME_CLAIM = "claims"

	SQL_CREATE_TABLE_CLAIM = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_CLAIM + `(
			id						BIGSERIAL			PRIMARY KEY,
			description				TEXT				NOT NULL,

			claimer_id	BIGINT REFERENCES ` + TABLE_NAME_USER + `(id)		NOT NULL,
			campaign_id	BIGINT REFERENCES ` + TABLE_NAME_CAMPAIGN + `(id)	NOT NULL,

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// The Claim model represents a claim to the proceeds of a Campaign
type Claim struct {
	Id          int64  // The identifier of the contribution
	Description string // The description of the claim

	ClaimerId  int64           // The id of the claimer; Foreign key for User (belongs to)
	CampaignId int64           // The id of the campaign; Foreign key for the Campaign (belongs to)
	Evidence   []ClaimEvidence // The evidence of the claim; One-To-Many relationship (has many)
	Votes      []ClaimVote     // The votes concerning this claim; One-To-Many relationship (has many)

	Active    bool      // True if this entity has not been soft deleted
	CreatedAt time.Time // The time when this claim was created
	UpdatedAt time.Time // The time when this claim was last updated
	DeletedAt time.Time // The time when this claim was soft deleted
}

// Creates the Claim table if it doesn't already exist
func CreateClaimTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CLAIM)
	return err
}
