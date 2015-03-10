package main

import (
	"database/sql"
	"github.com/go-martini/martini"
)

func SetupMiddleware(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Authentication & session management
	m.Use(Sessionize)
	// Bundle the responder in with req. handlers
	m.Use(Responderize)
}
