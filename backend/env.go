package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	ENV_VAR_DB_NAME    = "DB_NAME"    // Name of the database name environment variable
	ENV_VAR_DB_USER    = "DB_USER"    // Name of the database user environment variable
	ENV_VAR_DB_PASS    = "DB_PASS"    // Name of the database password environment variable
	ENV_VAR_JWT_SECRET = "JWT_SECRET" // Name of JWT secret environment variable
	ENV_VAR_PORT       = "PORT"       // Name of the HTTP poer environment variable

	ERR_ENV_VAR_MISSING = "The environment variable \"%s\" was either missing or invalid"
)

type Environment struct {
	dbName    string
	dbUser    string
	dbPass    string
	jwtSecret string
	port      int
}

func NewEnvironment() (*Environment, error) {
	// Required variables
	dbName := os.Getenv(ENV_VAR_DB_NAME)
	if dbName == "" {
		return nil, errors.New(fmt.Sprintf(ERR_ENV_VAR_MISSING, ENV_VAR_DB_NAME))
	}
	dbUser := os.Getenv(ENV_VAR_DB_USER)
	if dbUser == "" {
		return nil, errors.New(fmt.Sprintf(ERR_ENV_VAR_MISSING, ENV_VAR_DB_USER))
	}
	jwtSecret := os.Getenv(ENV_VAR_JWT_SECRET)
	if jwtSecret == "" {
		return nil, errors.New(fmt.Sprintf(ERR_ENV_VAR_MISSING, ENV_VAR_JWT_SECRET))
	}
	portStr := os.Getenv(ENV_VAR_PORT)
	if portStr == "" {
		return nil, errors.New(fmt.Sprintf(ERR_ENV_VAR_MISSING, ENV_VAR_PORT))
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(ERR_ENV_VAR_MISSING, ENV_VAR_PORT))
	}
	// Optional variables
	dbPass := os.Getenv(ENV_VAR_DB_PASS)

	return &Environment{
		dbName:    dbName,
		dbUser:    dbUser,
		dbPass:    dbPass,
		jwtSecret: jwtSecret,
		port:      port,
	}, nil
}
