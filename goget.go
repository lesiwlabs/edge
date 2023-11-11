package main

import (
	"fmt"
	"net/http"
)

var gogetTpl = `<html>
<head>
<meta name="go-import" content="%[1]s git %[2]s">
<meta name="go-source" content="%[1]s _ %[2]s{/dir} %[2]s{/dir}/{file}#L{line}">
</head>
<body>
go get %[1]s
</body>
</html>`

func goget(w http.ResponseWriter, pattern, redirect string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(gogetTpl, pattern, redirect)))
}
