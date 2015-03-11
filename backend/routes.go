package main

import (
	"database/sql"
	"github.com/go-martini/martini"
)

const (
	API_PREFIX = "/api"
	// Auth routes
	API_SESSION      = API_PREFIX + "/session"
	API_AUTHENTICATE = API_PREFIX + "/authenticate"
	// User routes
	API_REGISTER_USER = API_PREFIX + "/users"
	API_GET_ALL_USERS = API_PREFIX + "/users"
	API_GET_USER      = API_PREFIX + "/users/:id"
)

func SetupRoutes(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Routes that serve site assets
	SetupAssetRoutes(m, db, env)
	// Routes that handle authentication
	SetupAuthRoutes(m, db, env)
	// Routes to do with users
	SetupUserRoutes(m, db, env)
}
