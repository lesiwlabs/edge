package main

import (
	"os"
	"sync"

	"labs.lesiw.io/ops/git"
	"labs.lesiw.io/ops/github"
	"lesiw.io/cmdio"
	"lesiw.io/cmdio/ctr"
	"lesiw.io/cmdio/sys"
	"lesiw.io/ops"
)

type Ops struct{ github.Ops }

var rnr *cmdio.Runner
var secrets = map[string]string{
	"CLOUDFLARE_API_TOKEN": "cloudflare/wrangler/api",
}

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "build")
	}
	github.Repo = "lesiwlabs/edge"
	github.Secrets = secrets
	if os.Args[1] == "secrets" {
		rnr = sys.Runner()
	} else {
		rnr = mustv(ctr.New(bldctr()))
	}
	git.CopyWorktree(rnr, sys.Runner())
	ops.Defer(func() { _ = rnr.Close() })

	ops.Handle(Ops{})
}

func (Ops) Build() {
	rnr.MustRun(
		"go", "run",
		"github.com/syumai/workers/cmd/workers-assets-gen@v0.18.0",
	)
	rnr.MustRun(
		"tinygo", "build",
		"-o", "./build/app.wasm",
		"-target", "wasm",
		"-no-debug",
		"-panic=trap",
		"-gc=leaking",
		"-opt=2",
		"./...",
	)
}

func (op Ops) Deploy() {
	op.Build()
	rnr := rnr.WithEnv(fetchSecrets(secrets))
	rnr.MustRun("npx", "wrangler", "deploy")
}

func (op Ops) Dev() {
	op.Build()
	rnr.MustRun("npx", "wrangler", "dev")
}

func mustv[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func fetchSecrets(s map[string]string) map[string]string {
	rnr := sys.Runner()
	ret := make(map[string]string)
	spkez := sync.OnceValue(func() string {
		if _, err := rnr.Get("which", "spkez"); err != nil {
			rnr.MustRun("go", "install", "lesiw.io/spkez@latest")
		}
		return rnr.MustGet("which", "spkez").Out
	})
	for k, v := range s {
		if value := rnr.Env(k); value != "" {
			ret[k] = value
		} else {
			ret[k] = rnr.MustGet(spkez(), "get", v).Out
		}
	}
	return ret
}
