package main

import (
	"database/sql"
	"time"
)

const (
	TABLE_NAME_CLAIM_EVIDENCE = "claim_evidence"

	SQL_CREATE_TABLE_CLAIM_EVIDENCE = `
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME_CLAIM_EVIDENCE + `(
			id			BIGSERIAL		PRIMARY KEY,
			type		INTEGER			NOT NULL,
			url			VARCHAR(511)	NOT NULL,

			claim_id	BIGINT REFERENCES ` + TABLE_NAME_CLAIM + `(id)	NOT NULL,

			active			BOOLEAN				NOT NULL,
			created_at		TIMESTAMPTZ			NOT NULL,
			updated_at		TIMESTAMPTZ			NOT NULL,
			deleted_at		TIMESTAMPTZ
		);
	`
)

// The ClaimEvidence model represents proof that supports a Claim
type ClaimEvidence struct {
	Id   int64  // The identifier of the contribution
	Type int    // The type of the claim
	Url  string // The url of the evidence

	ClaimId sql.NullInt64 // The id of the claim; Foreign key for Claim (belongs to)

	Active    bool      // True if this entity has not been soft deleted
	CreatedAt time.Time // The time when this claim evidence was created
	UpdatedAt time.Time // The time when this claim evidence was last updated
	DeletedAt time.Time // The time when this claim evidence was soft deleted
}

// Creates the ClaimEvidence table if it doesn't already exist
func CreateClaimEvidenceTable(db *sql.DB) error {
	_, err := db.Exec(SQL_CREATE_TABLE_CLAIM_EVIDENCE)
	return err
}
