package main

import "net/http"

type url string

func (u url) handleRequest(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, string(u), http.StatusTemporaryRedirect)
}
