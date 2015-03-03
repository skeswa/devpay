package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()
	// Middleware
	m.Use(render.Renderer(render.Options{
		Directory:  "public/html",
		Extensions: []string{".html"},
	}))
	// Page routes
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})
	// Start the server
	m.Run()
}
