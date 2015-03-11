package main

import (
	"database/sql"
	"github.com/go-martini/martini"
)

func SetupMiddleware(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Add environment vars
	m.Use(func(c martini.Context) {
		c.Map(env)
	})
	// Authentication & session management
	m.Use(Sessionize)
	// Bundle the responder in with req. handlers
	m.Use(Responderize)
}
