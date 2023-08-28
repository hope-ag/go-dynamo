package handlers

import "net/http"

type Interface interface {
	Get(e http.ResponseWriter, r *http.Request)
	Post(e http.ResponseWriter, r *http.Request)
	Put(e http.ResponseWriter, r *http.Request)
	Patch(e http.ResponseWriter, r *http.Request)
	Delete(e http.ResponseWriter, r *http.Request)
	Options(e http.ResponseWriter, r *http.Request)
}