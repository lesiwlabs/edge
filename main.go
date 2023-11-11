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
	key := stripSubdomain(r.URL.Host) + r.URL.Path
	key = strings.TrimSuffix(key, "/")
	target, ok := redirects[key]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.URL.Query().Get("go-get") == "1" {
		goget(w, key, target)
		return
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
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
