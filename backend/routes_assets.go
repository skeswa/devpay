package main

import (
	"database/sql"
	"github.com/go-martini/martini"
)

const (
	PATH_INDEX_HTML = "public/html/index.html"
	PATH_404_HTML   = "public/html/404.html"
)

func SetupAssetRoutes(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Default page route
	m.Get("/", func(responder *Responder) {
		responder.Page(PATH_INDEX_HTML)
	})
	// Setup the 404
	m.NotFound(func(responder *Responder) {
		responder.Page(PATH_404_HTML)
	})
}
