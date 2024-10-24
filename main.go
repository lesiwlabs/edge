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

	"lesiw.io/bump": &gopkg{
		app: "bump",
		pkg: "lesiw.io/bump",
		src: "https://github.com/lesiw/bump",
	},
	"lesiw.io/buzzybox": &gopkg{
		app: "buzzybox",
		pkg: "lesiw.io/buzzybox",
		src: "https://github.com/lesiw/buzzybox",
	},
	"lesiw.io/clerk": &gopkg{
		pkg: "lesiw.io/clerk",
		src: "https://github.com/lesiw/clerk",
	},
	"lesiw.io/cmdio": &gopkg{
		pkg: "lesiw.io/cmdio",
		src: "https://github.com/lesiw/cmdio",
	},
	"lesiw.io/cmdio/x/busybox": &gopkg{
		pkg: "lesiw.io/cmdio/x/busybox",
		src: "https://github.com/lesiw/cmdio-busybox",
	},
	"lesiw.io/ctrctl": &gopkg{
		pkg: "lesiw.io/ctrctl",
		src: "https://github.com/lesiw/ctrctl",
	},
	"lesiw.io/dataer": &gopkg{
		app: "dataer",
		pkg: "lesiw.io/dataer",
		src: "https://github.com/lesiw/dataer",
	},
	"lesiw.io/defers": &gopkg{
		pkg: "lesiw.io/defers",
		src: "https://github.com/lesiw/defers",
	},
	"lesiw.io/flag": &gopkg{
		pkg: "lesiw.io/flag",
		src: "https://github.com/lesiw/flag",
	},
	"lesiw.io/hue": &gopkg{
		app: "hue",
		pkg: "lesiw.io/hue",
		src: "https://github.com/lesiw/hue",
	},
	"lesiw.io/inter": &gopkg{
		app: "inter",
		pkg: "lesiw.io/inter",
		src: "https://github.com/lesiw/inter",
	},
	"lesiw.io/moxie": &gopkg{
		app: "moxie",
		pkg: "lesiw.io/moxie",
		src: "https://github.com/lesiw/moxie",
	},
	"lesiw.io/op": &gopkg{
		app: "op",
		pkg: "lesiw.io/op",
		src: "https://github.com/lesiw/op",
	},
	"lesiw.io/ops": &gopkg{
		app: "op",
		bin: "https://github.com/lesiw/op",
		pkg: "lesiw.io/ops",
		src: "https://github.com/lesiw/ops",
	},
	"lesiw.io/plain": &gopkg{
		pkg: "lesiw.io/plain",
		src: "https://github.com/lesiw/plain",
	},
	"lesiw.io/prefix": &gopkg{
		pkg: "lesiw.io/prefix",
		src: "https://github.com/lesiw/prefix",
	},
	"lesiw.io/repo": &gopkg{
		app: "repo",
		pkg: "lesiw.io/repo",
		src: "https://github.com/lesiw/repo",
	},
	"lesiw.io/run": &gopkg{
		app: "run",
		pkg: "lesiw.io/run",
		src: "https://github.com/lesiw/run",
	},
	"lesiw.io/smol": &gopkg{
		pkg: "lesiw.io/smol",
		src: "https://github.com/lesiw/smol",
	},
	"lesiw.io/spkez": &gopkg{
		app: "spkez",
		pkg: "lesiw.io/spkez",
		src: "https://github.com/lesiw/spkez",
	},

	"lesiw.io/talks": &url{"https://github.com/lesiw/talks"},

	"lesiw.dev":         &url{"https://github.com/lesiw"},
	"lesiw.dev/discord": &url{"https://discord.gg/EYWxqssV99"},
	"lesiw.dev/twitch":  &url{"https://twitch.tv/lesiwlabs"},

	"lesiw.chat": &url{"https://discord.gg/EYWxqssV99"},

	"labs.lesiw.io/ops": &gopkg{
		pkg: "labs.lesiw.io/ops",
		src: "https://github.com/lesiwlabs/ops",
	},
	"labs.lesiw.io/edge": &gopkg{
		pkg: "labs.lesiw.io/edge",
		src: "https://github.com/lesiwlabs/edge",
	},
	"labs.lesiw.io/pass": &url{"https://github.com/lesiwlabs/pass"},
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Host, "www.") + r.URL.Path
	key = strings.TrimSuffix(key, "/")
	target, ok := targets[key]
	if !ok {
		http.NotFound(w, r)
		return
	}
	target.handleRequest(w, r)
}

func main() {
	workers.Serve(&handler{})
}
