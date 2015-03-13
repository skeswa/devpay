package main

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// The Claim model represents a claim to the proceeds of a Campaign
type Claim struct {
	Id          int64  `json:"id"`          // The identifier of the contribution
	Description string `json:"description"` // The description of the claim

	Claimer    *User            `json:"claimer,omitempty"`  // The person who made this claim
	ClaimerId  int64            `json:"-"`                  // The id of the claimer; Foreign key for User (belongs to)
	Campaign   *Campaign        `json:"campaign,omitempty"` // The campaign this claim was made for
	CampaignId int64            `json:"-"`                  // The id of the campaign; Foreign key for the Campaign (belongs to)
	Evidence   []*ClaimEvidence `json:"evidence"`           // The evidence of the claim; One-To-Many relationship (has many)
	Votes      []*ClaimVote     `json:"votes"`              // The votes concerning this claim; One-To-Many relationship (has many)

	Active    bool        `json:"active"`    // True if this entity has not been soft deleted
	CreatedAt time.Time   `json:"createdAt"` // The time when this contribution was created
	UpdatedAt time.Time   `json:"updatedAt"` // The time when this contribution was last updated
	DeletedAt pq.NullTime `json:"-"`         // The time when this user was soft deleted
}

const (
	TABLE_NAME_CLAIM = "claims"

	FIELD_CLAIM_CLAIMER_ID  = "claimer_id"
	FIELD_CLAIM_CAMPAIGN_ID = "campaign_id"

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
	SQL_SELECT_CLAIM_BY_CAMPAIGN_ID = `
		SELECT * FROM ` + TABLE_NAME_CLAIM + `
			LEFT JOIN ` + TABLE_NAME_USER + ` as claimers ON ` + TABLE_NAME_CLAIM + `.` + FIELD_CLAIM_CLAIMER_ID + `=claimers.id
		WHERE (` + FIELD_CONTRIBUTION_CAMPAIGN_ID + ` = $1);
	`
)

// Creates the Claim table if it doesn't already exist
func CreateClaimTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CLAIM)
	return err
}

// Finds Claims for a specific campaign
func FindClaimsByCampaignId(
	db Queryable,
	campaignId int64,
) ([]*Claim, error) {
	claims := make([]*Claim, 0)
	// Submit the query
	rows, err := db.Query(SQL_SELECT_CLAIM_BY_CAMPAIGN_ID, campaignId)
	if err != nil {
		return nil, err
	}
	// Read the rows
	defer rows.Close()
	for rows.Next() {
		var currentClaim Claim
		var currentClaimer User
		// Read row data
		err = rows.Scan(
			&currentClaim.Id, &currentClaim.Description, &currentClaim.ClaimerId, &currentClaim.CampaignId, &currentClaim.Active, &currentClaim.CreatedAt, &currentClaim.UpdatedAt, &currentClaim.DeletedAt, // The contribution fields
			&currentClaimer.Id, &currentClaimer.FirstName, &currentClaimer.LastName, &currentClaimer.Email, &currentClaimer.HashedPassword, &currentClaimer.StripeId, &currentClaimer.PictureUrl, &currentClaimer.Active, &currentClaimer.CreatedAt, &currentClaimer.UpdatedAt, &currentClaimer.DeletedAt, // The contributor fields
		)
		// Exit if there was a problem
		if err != nil {
			rows.Close()
			return nil, err
		} else {
			currentClaim.Claimer = &currentClaimer
			claims = append(claims, &currentClaim)
		}
	}
	// Return the results
	return claims, nil
}
