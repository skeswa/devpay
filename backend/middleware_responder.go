package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"net/http"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentJSON    = "application/json"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	ContentXML     = "text/xml"
	defaultCharset = "UTF-8"
)

type Responder struct {
	response http.ResponseWriter
	request  *http.Request
}

func (r *Responder) Page(path string) {
	// Set content type to json
	r.response.Header().Set(ContentType, ContentHTML)
	http.ServeFile(r.response, r.request, path)
}

func (r *Responder) PageWithStatus(status int, path string) {
	// Set content type to json
	r.response.WriteHeader(status)
	r.response.Header().Set(ContentType, ContentHTML)
	http.ServeFile(r.response, r.request, path)
}

func (r *Responder) Text(text string) {
	r.response.Write([]byte(text))
}

func (r *Responder) Json(v interface{}) {
	// Set content type to json
	r.response.Header().Set(ContentType, ContentJSON)
	// Perform json marshalling
	result, err := json.Marshal(v)
	if err != nil {
		Debug("Could not render the json result: ", err)
		r.response.WriteHeader(http.StatusInternalServerError)
		r.response.Write(PUBERR_INTERNAL_SERVER_ERROR.Json)
		return
	}
	// Ship the json
	r.response.WriteHeader(http.StatusOK)
	r.response.Write(result)
}

func (r *Responder) Error(v interface{}) {
	// Set content type to json
	r.response.Header().Set(ContentType, ContentJSON)
	if pubErr, ok := v.(*PublicError); ok {
		r.response.WriteHeader(pubErr.Status)
		r.response.Write(pubErr.Json)
	} else {
		Debug("Private error was squashed: ", v)
		r.response.WriteHeader(http.StatusInternalServerError)
		r.response.Write(PUBERR_INTERNAL_SERVER_ERROR.Json)
	}
}

func Responderize(res http.ResponseWriter, req *http.Request, c martini.Context) {
	c.Map(&Responder{
		request:  req,
		response: res,
	})
}
