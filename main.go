package main

import (
	"net/http"
	"strings"

	"github.com/syumai/workers"
)

type AnySubdomain struct {
	handlers map[string]http.HandlerFunc
}

func (h *AnySubdomain) HandleFunc(pattern string, handler http.HandlerFunc) {
	h.handlers[pattern] = handler
}

func (h *AnySubdomain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for pattern, handler := range h.handlers {
		host, path, _ := strings.Cut(pattern, "/")
		if stripSubdomain(r.URL.Host) == host && r.URL.Path == "/"+path {
			handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func stripSubdomain(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) > 2 {
		return strings.Join(parts[1:], ".")
	}
	return url
}

func main() {
	hdl := &AnySubdomain{handlers: make(map[string]http.HandlerFunc)}
	hdl.HandleFunc("chrislesiw.com", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.linkedin.com/in/christopher-lesiw/",
			http.StatusTemporaryRedirect)
	})
	hdl.HandleFunc("chrislesiw.com/github", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/lesiw", http.StatusTemporaryRedirect)
	})
	workers.Serve(hdl)
}
