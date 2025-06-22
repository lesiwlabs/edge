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
	"chrislesiw.com":        url("https://www.linkedin.com/in/christopher-lesiw/"),
	"chrislesiw.com/github": url("https://github.com/lesiw"),

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
	"lesiw.io/chrono": &gopkg{
		pkg: "lesiw.io/chrono",
		src: "https://github.com/lesiw/chrono",
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
	"lesiw.io/fill": &gopkg{
		pkg: "lesiw.io/fill",
		src: "https://github.com/lesiw/fill",
	},
	"lesiw.io/flag": &gopkg{
		pkg: "lesiw.io/flag",
		src: "https://github.com/lesiw/flag",
	},
	"lesiw.io/http2https": &gopkg{
		pkg: "lesiw.io/http2https",
		src: "https://github.com/lesiw/http2https",
	},
	"lesiw.io/hue": &gopkg{
		app: "hue",
		pkg: "lesiw.io/hue",
		src: "https://github.com/lesiw/hue",
	},
	"lesiw.io/igo": &gopkg{
		app: "igo",
		pkg: "lesiw.io/igo",
		src: "https://github.com/lesiw/igo",
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
	"lesiw.io/notes": &gopkg{
		app: "notes",
		pkg: "lesiw.io/notes",
		src: "https://github.com/lesiw/notes",
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
	"lesiw.io/testdetect": &gopkg{
		pkg: "lesiw.io/testdetect",
		src: "https://github.com/lesiw/testdetect",
	},
	"lesiw.io/zync": &gopkg{
		pkg: "lesiw.io/zync",
		src: "https://github.com/lesiw/zync",
	},

	"lesiw.io/datastax": url("https://github.com/lesiw/datastax"),
	"lesiw.io/talks":    url("https://github.com/lesiw/talks"),

	"lesiw.chat": url("https://discord.gg/EYWxqssV99"),

	"labs.lesiw.io/ctr": &gopkg{
		pkg: "labs.lesiw.io/ctr",
		src: "https://github.com/lesiwlabs/ctr",
	},
	"labs.lesiw.io/discord": url("https://github.com/lesiwlabs/discord"),
	"labs.lesiw.io/echo":    url("https://github.com/lesiwlabs/echo"),
	"labs.lesiw.io/edge": &gopkg{
		pkg: "labs.lesiw.io/edge",
		src: "https://github.com/lesiwlabs/edge",
	},
	"labs.lesiw.io/feed": url("https://github.com/lesiwlabs/feed"),
	"labs.lesiw.io/k8s":  url("https://github.com/lesiwlabs/k8s"),
	"labs.lesiw.io/ops": &gopkg{
		pkg: "labs.lesiw.io/ops",
		src: "https://github.com/lesiwlabs/ops",
	},
	"labs.lesiw.io/pass": url("https://github.com/lesiwlabs/pass"),

	"origin.lesiw.dev":         url("https://github.com/lesiw"),
	"origin.lesiw.dev/discord": url("https://discord.gg/EYWxqssV99"),
	"origin.lesiw.dev/twitch":  url("https://twitch.tv/lesiwlabs"),
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Host, "www.") {
		http.Redirect(w, r,
			"https://"+strings.TrimPrefix(r.URL.Host, "www.")+r.URL.Path,
			http.StatusMovedPermanently,
		)
		return
	}
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
