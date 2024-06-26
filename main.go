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
	"lesiw.io/buzzybox": &ghGoPkg{
		app: "buzzybox",
		pkg: "lesiw.io/buzzybox",
		src: "https://github.com/lesiw/buzzybox",
	},
	"lesiw.io/ctrctl": &ghGoPkg{
		app: "ctrctl",
		pkg: "lesiw.io/ctrctl",
		src: "https://github.com/lesiw/ctrctl",
	},
	"lesiw.io/defers": &ghGoPkg{
		app: "defers",
		pkg: "lesiw.io/defers",
		src: "https://github.com/lesiw/defers",
	},
	"lesiw.io/flag": &ghGoPkg{
		app: "flag",
		pkg: "lesiw.io/flag",
		src: "https://github.com/lesiw/flag",
	},
	"lesiw.io/hue": &ghGoPkg{
		app: "hue",
		pkg: "lesiw.io/hue",
		src: "https://github.com/lesiw/hue",
	},
	"lesiw.io/inter": &ghGoPkg{
		app: "inter",
		pkg: "lesiw.io/inter",
		src: "https://github.com/lesiw/inter",
	},
	"lesiw.io/plain": &ghGoPkg{
		app: "repo",
		pkg: "lesiw.io/plain",
		src: "https://github.com/lesiw/plain",
	},
	"lesiw.io/repo": &ghGoPkg{
		app: "repo",
		pkg: "lesiw.io/repo",
		src: "https://github.com/lesiw/repo",
	},
	"lesiw.io/run": &ghGoPkg{
		app: "run",
		pkg: "lesiw.io/run",
		src: "https://github.com/lesiw/run",
	},
	"lesiw.io/smol": &ghGoPkg{
		app: "smol",
		pkg: "lesiw.io/smol",
		src: "https://github.com/lesiw/smol",
	},

	"lesiw.dev":         &url{"https://github.com/lesiw"},
	"lesiw.dev/discord": &url{"https://discord.gg/EYWxqssV99"},
	"lesiw.dev/twitch":  &url{"https://twitch.tv/lesiwlabs"},

	"lesiw.chat": &url{"https://discord.gg/EYWxqssV99"},
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
