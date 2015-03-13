package main

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// The Campaign model represents a funding effort with a clear goal and a deadline
type Campaign struct {
	Id                  int64     `json:"id"`                  // The identifier of the campaign
	Title               string    `json:"title"`               // The title of the campaign
	Description         string    `json:"description"`         // The description of the campaign
	CoverPictureUrl     string    `json:"coverPictureUrl"`     // The URL of this campaign's cover picture
	ThumbnailPictureUrl string    `json:"thumbnailPictureUrl"` // The URL of this campaign's thumbnail picture
	Amount              float64   `json:"amount"`              // The current amount that this campaign has raised
	Deadline            time.Time `json:"deadline"`            // When this campaign expires
	Finished            bool      `json:"finished"`            // True if the campaign is over

	Creator       *User           `json:"creator,omitempty"` // The person who started this campaign; One-To-Many relationship (has one)
	CreatorId     int64           `json:"-"`                 // The id of the creator; Foreign key for User (belongs to)
	Claimer       *User           `json:"claimer,omitempty"` // The person who successfully claimed the Campaign; One-To-Many relationship (has one)
	ClaimerId     sql.NullInt64   `json:"-"`                 // The id of the person who successfully claimed the Campaign; Foreign key for User (belongs to)
	Contributions []*Contribution `json:"contributions"`     // All the contributions to this campaign; One-To-Many relationship (has many)
	Claims        []*Claim        `json:"claims"`            // All the claims for this campaign; One-To-Many relationship (has many)

	Active    bool        `json:"active"`    // True if this entity has not been soft deleted
	CreatedAt time.Time   `json:"createdAt"` // The time when this campaign was created
	UpdatedAt time.Time   `json:"updatedAt"` // The time when this campaign was last updated
	DeletedAt pq.NullTime `json:"-"`         // The time when this user was soft deleted
}

const (
	TABLE_NAME_CAMPAIGN = "campaigns"

	FIELD_CAMPAIGN_CREATED_AT  = "created_at"
	FIELD_CAMPAIGN_CREATOR_ID  = "creator_id"
	FIELD_CAMPAIGN_CLAIMER_ID  = "claimer_id"
	FIELD_CAMPAIGN_TITLE       = "title"
	FIELD_CAMPAIGN_DESCRIPTION = "description"

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
	SQL_CREATE_NEW_CAMPAIGN = `
		INSERT INTO ` + TABLE_NAME_CAMPAIGN + `
		(title, description, cover_picture_url, thumbnail_picture_url, amount, deadline, finished, creator_id, active, created_at, updated_at) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
	`
	SQL_SELECT_CAMPAIGN_BY_ID = `
		SELECT * FROM ` + TABLE_NAME_CAMPAIGN + ` WHERE (id = $1);
	`
	SQL_SELECT_FULL_CAMPAIGN_BY_ID = `
		SELECT * FROM ` + TABLE_NAME_CAMPAIGN + `
			LEFT JOIN ` + TABLE_NAME_USER + ` as creators ON ` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATOR_ID + `=creators.id
			LEFT JOIN ` + TABLE_NAME_USER + ` as claimers ON ` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CLAIMER_ID + `=claimers.id
		WHERE (` + TABLE_NAME_CAMPAIGN + `.id = $1);
	`
	SQL_SELECT_CAMPAIGNS = `
		SELECT * FROM ` + TABLE_NAME_CAMPAIGN + `
			LEFT JOIN ` + TABLE_NAME_USER + ` as creators ON ` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATOR_ID + `=creators.id
		ORDER BY ` + FIELD_CAMPAIGN_CREATED_AT + ` DESC
		OFFSET $1 LIMIT $2;
	`
	SQL_SELECT_AND_FILTER_CAMPAIGNS = `
		SELECT * FROM ` + TABLE_NAME_CAMPAIGN + `
			LEFT JOIN ` + TABLE_NAME_USER + ` as creators ON ` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATOR_ID + `=creators.id
		WHERE ((` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATED_AT + ` LIKE $1) OR (` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATED_AT + ` LIKE $1))
		ORDER BY ` + TABLE_NAME_CAMPAIGN + `.` + FIELD_CAMPAIGN_CREATED_AT + ` DESC
		OFFSET $2 LIMIT $3;
	`
)

// Creates the Campaign table if it doesn't already exist
func CreateCampaignTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CAMPAIGN)
	return err
}

// Creates a new Campaign in the database; returns the id of the new campaign
func CreateNewCampaign(
	db Queryable, // The database
	Title string, // The first name of the user
	Description string, // The last name of the user
	CoverPictureUrl string, // The email address of the user (indexed)
	ThumbnailPictureUrl string, // The bcrypted password of the user
	Amount float64, // The id of the user with Stripe's API
	Deadline time.Time, // The URL to user's picture
	CreatorId int64,
) (int64, error) {
	var (
		id  int64
		now = time.Now()
	)
	err := db.QueryRow(SQL_CREATE_NEW_CAMPAIGN, Title, Description, CoverPictureUrl, ThumbnailPictureUrl, Amount, Deadline, false, CreatorId, true, now, now).Scan(&id)
	if err != nil {
		return -1, err
	} else {
		return id, nil
	}
}

// Gets a Campaign from the database by id
func GetCampaign(
	db Queryable,
	id int64,
) (*Campaign, error) {
	rows, err := db.Query(SQL_SELECT_CAMPAIGN_BY_ID, id)
	if err != nil {
		// TODO standardize all database error returns
		return nil, PUBERR_ENTITY_NOT_FOUND
	}
	// Read the rows
	defer rows.Close()
	var campaign Campaign
	for rows.Next() {
		err = rows.Scan(&campaign.Id, &campaign.Title, &campaign.Description, &campaign.CoverPictureUrl, &campaign.ThumbnailPictureUrl, &campaign.Amount, &campaign.Deadline, &campaign.Finished, &campaign.CreatorId, &campaign.ClaimerId, &campaign.Active, &campaign.CreatedAt, &campaign.UpdatedAt, &campaign.DeletedAt)
		if err != nil {
			return nil, err
		} else {
			return &campaign, nil
		}
	}
	// We didn't find any users
	return nil, PUBERR_ENTITY_NOT_FOUND
}

// Gets a Campaign from the database by id; has all its relationships fulfilled
func GetFullCampaign(
	db Queryable,
	id int64,
) (*Campaign, error) {
	var (
		foundResults = false
		creator      User
		campaign     Campaign
		// Claimer fields are null b/c claimer is optional
		claimerId         sql.NullInt64
		claimerFirstName  sql.NullString
		claimerLastName   sql.NullString
		claimerEmail      sql.NullString
		claimerPictureUrl sql.NullString
		ignoredField      interface{}
	)
	// Query the db
	rows, err := db.Query(SQL_SELECT_FULL_CAMPAIGN_BY_ID, id)
	if err != nil {
		return nil, err
	}
	// Read the rows
	for rows.Next() {
		foundResults = true
		// Scan the results
		err = rows.Scan(
			&campaign.Id, &campaign.Title, &campaign.Description, &campaign.CoverPictureUrl, &campaign.ThumbnailPictureUrl, &campaign.Amount, &campaign.Deadline, &campaign.Finished, &campaign.CreatorId, &campaign.ClaimerId, &campaign.Active, &campaign.CreatedAt, &campaign.UpdatedAt, &campaign.DeletedAt, // The campaign fields
			&creator.Id, &creator.FirstName, &creator.LastName, &creator.Email, &creator.HashedPassword, &creator.StripeId, &creator.PictureUrl, &creator.Active, &creator.CreatedAt, &creator.UpdatedAt, &creator.DeletedAt, // The creator fields
			&claimerId, &claimerFirstName, &claimerLastName, &claimerEmail, &ignoredField, &ignoredField, &claimerPictureUrl, &ignoredField, &ignoredField, &ignoredField, &ignoredField, // The claimer fields
		)
		// Exit if there was a problem
		if err != nil {
			rows.Close()
			return nil, err
		} else {
			// Nest the entities
			campaign.Creator = &creator
			if campaign.ClaimerId.Valid {
				campaign.Claimer = &User{}
				claimerId.Scan(&campaign.Claimer.Id)
				claimerFirstName.Scan(&campaign.Claimer.FirstName)
				claimerLastName.Scan(&campaign.Claimer.LastName)
				claimerEmail.Scan(&campaign.Claimer.Email)
				claimerPictureUrl.Scan(&campaign.Claimer.PictureUrl)
			}
			// Close rows and break
			rows.Close()
			break
		}
	}
	// Exit if there were no results
	if !foundResults {
		return nil, PUBERR_ENTITY_NOT_FOUND
	}
	// Grab the contributions
	contributions, err := FindContributionsByCampaignId(db, campaign.Id)
	if err != nil {
		return nil, err
	} else {
		campaign.Contributions = contributions
	}
	// Grab the claims
	claims, err := FindClaimsByCampaignId(db, campaign.Id)
	if err != nil {
		return nil, err
	} else {
		campaign.Claims = claims
	}
	// We didn't find any users
	return &campaign, nil
}

// Get all campaigns
func GetCampaigns(
	db Queryable,
	offset int,
	limit int,
) ([]*Campaign, error) {
	rows, err := db.Query(SQL_SELECT_CAMPAIGNS, offset, limit)
	if err != nil {
		return nil, err
	}
	// Read the rows
	defer rows.Close()
	campaigns := make([]*Campaign, 0, limit)
	for rows.Next() {
		var (
			campaign Campaign
			creator  User
		)
		err = rows.Scan(
			&campaign.Id, &campaign.Title, &campaign.Description, &campaign.CoverPictureUrl, &campaign.ThumbnailPictureUrl, &campaign.Amount, &campaign.Deadline, &campaign.Deadline, &campaign.CreatorId, &campaign.ClaimerId, &campaign.Active, &campaign.CreatedAt, &campaign.UpdatedAt, &campaign.DeletedAt, // The campaign fields
			&creator.Id, &creator.FirstName, &creator.LastName, &creator.Email, &creator.HashedPassword, &creator.StripeId, &creator.PictureUrl, &creator.Active, &creator.CreatedAt, &creator.UpdatedAt, &creator.DeletedAt, // The creator fields
		)
		if err != nil {
			return nil, err
		} else {
			campaign.Creator = &creator
			campaigns = append(campaigns, &campaign)
		}
	}
	// We didn't find any users
	return campaigns, nil
}

// Filter all campaigns
func FilterCampaigns(
	db Queryable,
	offset int,
	limit int,
) ([]*Campaign, error) {
	return nil, nil
}
