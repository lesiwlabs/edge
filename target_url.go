package main

import "net/http"

type url struct {
	target string
}

func (u *url) handleRequest(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, u.target, http.StatusTemporaryRedirect)
}
