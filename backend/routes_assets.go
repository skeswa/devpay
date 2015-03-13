package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
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
	m.NotFound(func(req *http.Request, responder *Responder) {
		if strings.Index(req.URL.Path, API_PREFIX) == 0 {
			responder.Error(PUBERR_ENDPOINT_NOT_FOUND)
		} else {
			responder.PageWithStatus(http.StatusNotFound, PATH_404_HTML)
		}
	})
}
