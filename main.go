package main

import (
	"net/http"
	"strings"

	"github.com/syumai/workers"
)

var redirects = map[string]string{
	"chrislesiw.com":        "https://www.linkedin.com/in/christopher-lesiw/",
	"chrislesiw.com/github": "https://github.com/lesiw",

	"lesiw.io/inter": "https://github.com/lesiw/inter",
}

type RedirectHandler struct{}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for pattern, redirect := range redirects {
		host, path, _ := strings.Cut(pattern, "/")
		if stripSubdomain(r.URL.Host) == host && r.URL.Path == "/"+path {
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
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
	workers.Serve(&RedirectHandler{})
}
