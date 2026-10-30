package main

import (
	"embed"
	"net/http"
)

//go:embed *.txt
var docs embed.FS

// doc is a plain-text document, named by the file it was embedded from.
type doc string

func (d doc) handleRequest(w http.ResponseWriter, r *http.Request) {
	b, err := docs.ReadFile(string(d))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(b)
}
