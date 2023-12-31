package main

import (
	"net/http"
	"strings"

	"github.com/syumai/workers"
)

type handler struct{}
type target interface {
	handleRequest(http.ResponseWriter, *http.Request)
}

var targets = map[string]target{
	"chrislesiw.com":        &url{"https://www.linkedin.com/in/christopher-lesiw/"},
	"chrislesiw.com/github": &url{"https://github.com/lesiw"},

	"lesiw.io/bump": &ghGoPkg{
		app: "bump",
		pkg: "lesiw.io/bump",
		src: "https://github.com/lesiw/bump",
	},
	"lesiw.io/inter": &ghGoPkg{
		app: "inter",
		pkg: "lesiw.io/inter",
		src: "https://github.com/lesiw/inter",
	},
	"lesiw.io/gobox": &ghGoPkg{
		app: "gobox",
		pkg: "lesiw.io/gobox",
		src: "https://github.com/lesiw/gobox",
	},
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := domain(r.URL.Host) + r.URL.Path
	key = strings.TrimSuffix(key, "/")
	target, ok := targets[key]
	if !ok {
		http.NotFound(w, r)
		return
	}
	target.handleRequest(w, r)
}

func domain(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) > 1 {
		return strings.Join(parts[len(parts)-2:], ".")
	} else if len(parts) > 0 {
		return parts[0]
	} else {
		return ""
	}
}

func main() {
	workers.Serve(&handler{})
}
