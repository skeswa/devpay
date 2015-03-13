package main

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// The Contribution model represents an amount paid by a user to a Campaign
type Contribution struct {
	Id       int64   `json:"id"`       // The identifier of the contribution
	Amount   float64 `json:"amount"`   // The amount of the contribution
	StripeId string  `json:"stripeId"` // The stripe id of this transaction

	Contributor   *User     `json:"contributor,omitempty"` // The person who made this contribution
	ContributorId int64     `json:"-"`                     // The id of the contributor; Foreign key for User (belongs to)
	Campaign      *Campaign `json:"campaign,omitempty"`    // The campaign this contribution was made to
	CampaignId    int64     `json:"-"`                     // The id of the campaign; Foreign key for the Campaign (belongs to)

	Active    bool        `json:"active"`    // True if this entity has not been soft deleted
	CreatedAt time.Time   `json:"createdAt"` // The time when this contribution was created
	UpdatedAt time.Time   `json:"updatedAt"` // The time when this contribution was last updated
	DeletedAt pq.NullTime `json:"-"`         // The time when this user was soft deleted
}

const (
	TABLE_NAME_CONTRIBUTION = "contributions"

	FIELD_CONTRIBUTION_CAMPAIGN_ID    = "campaign_id"
	FIELD_CONTRIBUTION_CONTRIBUTOR_ID = "contributor_id"

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
	SQL_SELECT_CONTRIBUTION_BY_CAMPAIGN_ID = `
		SELECT * FROM ` + TABLE_NAME_CONTRIBUTION + `
			LEFT JOIN ` + TABLE_NAME_USER + ` as contributors ON ` + TABLE_NAME_CONTRIBUTION + `.` + FIELD_CONTRIBUTION_CONTRIBUTOR_ID + `=contributors.id
		WHERE (` + FIELD_CONTRIBUTION_CAMPAIGN_ID + ` = $1);
	`
)

// Creates the Contribution table if it doesn't already exist
func CreateContributionTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CONTRIBUTION)
	return err
}

// Finds Contributions to a specific campaign
func FindContributionsByCampaignId(
	db Queryable,
	campaignId int64,
) ([]*Contribution, error) {
	contributions := make([]*Contribution, 0)
	// Submit the query
	rows, err := db.Query(SQL_SELECT_CONTRIBUTION_BY_CAMPAIGN_ID, campaignId)
	if err != nil {
		return nil, err
	}
	// Read the rows
	defer rows.Close()
	for rows.Next() {
		var currentContribution Contribution
		var currentContributor User
		// Read row data
		err = rows.Scan(
			&currentContribution.Id, &currentContribution.Amount, &currentContribution.StripeId, &currentContribution.ContributorId, &currentContribution.CampaignId, &currentContribution.Active, &currentContribution.CreatedAt, &currentContribution.UpdatedAt, &currentContribution.DeletedAt, // The contribution fields
			&currentContributor.Id, &currentContributor.FirstName, &currentContributor.LastName, &currentContributor.Email, &currentContributor.HashedPassword, &currentContributor.StripeId, &currentContributor.PictureUrl, &currentContributor.Active, &currentContributor.CreatedAt, &currentContributor.UpdatedAt, &currentContributor.DeletedAt, // The contributor fields
		)
		// Exit if there was a problem
		if err != nil {
			rows.Close()
			return nil, err
		} else {
			currentContribution.Contributor = &currentContributor
			contributions = append(contributions, &currentContribution)
		}
	}
	// Return the results
	return contributions, nil
}
