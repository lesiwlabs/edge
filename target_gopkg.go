package main

import (
	"fmt"
	"net/http"
	"strings"

	_ "embed"
)

type gopkg struct {
	app string
	pkg string
	src string
}

//go:embed goinstall.sh
var installScript string

var gogetTpl = `<html>
<head>
<meta name="go-import" content="%[1]s git %[2]s">
<meta name="go-source" content="%[1]s _ %[2]s{/dir} %[2]s{/dir}/{file}#L{line}">
</head>
<body>
go get %[1]s
</body>
</html>
`

func (p *gopkg) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("go-get") == "1" {
		p.get(w)
		return
	}
	ua := strings.ToLower(r.UserAgent())
	if p.app != "" && (strings.Contains(ua, "curl") || strings.Contains(ua, "wget")) {
		p.install(w)
		return
	}
	http.Redirect(w, r, p.src, http.StatusTemporaryRedirect)
}

func (p *gopkg) get(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(gogetTpl, p.pkg, p.src)))
}

func (p *gopkg) install(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	script := expandenv(installScript, map[string]string{
		"APP":        p.app,
		"BIN_URL":    p.src + "/releases/latest/download",
		"TICKET_URL": p.src + "/issues",
	})
	w.Write([]byte(script))
}

func expandenv(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "$"+k, v)
	}
	return s
}
