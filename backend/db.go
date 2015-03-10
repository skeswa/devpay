package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	// String templates
	TEMPLATE_PG_CONN_STRING = "user=%s dbname=%s sslmode=disable"
)

func SetupDatabase(env *Environment) (*sql.DB, error) {
	// Connect using the parameters above
	connString := fmt.Sprintf(TEMPLATE_PG_CONN_STRING, env.dbUser, env.dbName)
	if env.dbPass != "" {
		connString += " password=" + env.dbPass
	}
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, errors.New("Failed to connect to the database with connection string \"" + connString + "\": " + err.Error())
	}

	// Create tables for all the models
	err = CreateUserTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_USER, err.Error()))
	}
	err = CreateCampaignTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_CAMPAIGN, err.Error()))
	}
	err = CreateContributionTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_CONTRIBUTION, err.Error()))
	}
	err = CreateClaimTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_CLAIM, err.Error()))
	}
	err = CreateClaimEvidenceTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_CLAIM_EVIDENCE, err.Error()))
	}
	err = CreateClaimVoteTable(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_TABLE_CREATION_FAILED, TABLE_NAME_CLAIM_VOTE, err.Error()))
	}

	return db, nil
}
