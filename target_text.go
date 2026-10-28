package main

import (
	_ "embed"
	"net/http"
)

type text string

//go:embed review.txt
var review text

func (t text) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(t))
}
