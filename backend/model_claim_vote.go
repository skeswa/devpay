package main

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

// The ClaimVote model represents a vote in favor of, or against a Claim
type ClaimVote struct {
	Id          int64 // The identifier of the contribution
	Affirmative bool  // True if in favor of the Claim

	VoterId sql.NullInt64 // The id of the voter; Foreign key for User (belongs to)
	ClaimId sql.NullInt64 // The id of the claim; Foreign key for Claim (belongs to)

	Active    bool        // True if this entity has not been soft deleted
	CreatedAt time.Time   // The time when this claim evidence was created
	UpdatedAt time.Time   // The time when this claim evidence was last updated
	DeletedAt pq.NullTime // The time when this user was soft deleted
}

const (
	TABLE_NAME_CLAIM_VOTE = "claim_vote"

	SQL_CREATE_TABLE_CLAIM_VOTE = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_CLAIM_VOTE + `(
			id			BIGSERIAL		PRIMARY KEY,
			affirmative	BOOLEAN			NOT NULL,

			voter_id	BIGINT REFERENCES ` + TABLE_NAME_USER + `(id)	NOT NULL,
			claim_id	BIGINT REFERENCES ` + TABLE_NAME_CLAIM + `(id)	NOT NULL,

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// Creates the ClaimVote table if it doesn't already exist
func CreateClaimVoteTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CLAIM_VOTE)
	return err
}
