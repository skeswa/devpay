package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
)

const (
	ERR_COULDNT_START = "Couldn't start the the server: "
)

func main() {
	// Read environment variables
	env, err := NewEnvironment()
	if err != nil {
		log.Fatalln(ERR_COULDNT_START + err.Error())
	}
	// Setup the database
	_, err = SetupDatabase(env)
	if err != nil {
		log.Fatalln(ERR_COULDNT_START + err.Error())
	}

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
