package main

import "net/http"

type url string

func (u url) handleRequest(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, string(u), http.StatusTemporaryRedirect)
}

// moved is a url that has permanently relocated.
type moved string

func (m moved) handleRequest(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, string(m), http.StatusMovedPermanently)
}
