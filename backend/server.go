package main

import "github.com/go-martini/martini"

type Server struct {
	martini *martini.ClassicMartini
	port    int
}

func NewServer(port int) *Server {
	server := Server{
		port:    port,
		martini: martini.Classic(),
	}

	return &server
}
