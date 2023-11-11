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
	key := domain(r.URL.Host) + r.URL.Path
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

func domain(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) < 1 {
		return ""
	} else if len(parts) < 2 {
		return parts[0]
	}
	return strings.Join(parts[len(parts)-2:], ".")
}

func main() {
	workers.Serve(&RedirectHandler{})
}