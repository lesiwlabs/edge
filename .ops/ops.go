package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"labs.lesiw.io/ops/github"
	"lesiw.io/command"
	"lesiw.io/command/ctr"
	"lesiw.io/command/sub"
	"lesiw.io/command/sys"
	"lesiw.io/ops"
)

type Ops struct{ github.Ops }

var rnr *command.Sh
var secrets = map[string]string{
	"CLOUDFLARE_API_TOKEN": "cloudflare/wrangler/api",
}

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "build")
	}
	github.Repo = "lesiwlabs/edge"
	github.Secrets = secrets

	if os.Args[1] != "secrets" {
		if err := setup(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	ops.Handle(Ops{})
}

func setup() error {
	ctx := context.Background()
	m := ctr.Machine(sys.Machine(), ".ops/Containerfile")
	rnr = command.Shell(m, "go", "tinygo", "npx")
	ops.Defer(func() {
		if sm, ok := m.(command.ShutdownMachine); ok {
			_ = sm.Shutdown(context.Background())
		}
	})

	local := command.Shell(sys.Machine(), "git")
	archive := command.NewReader(ctx, local, "git", "archive", "HEAD")
	dst, err := rnr.Create(ctx, ".")
	if err != nil {
		return fmt.Errorf("create worktree: %w", err)
	}
	if _, err := io.Copy(dst, archive); err != nil {
		dst.Close()
		return fmt.Errorf("copy worktree: %w", err)
	}
	if err := dst.Close(); err != nil {
		return fmt.Errorf("close worktree: %w", err)
	}
	return nil
}

func (Ops) Build() error {
	ctx := context.Background()
	if err := rnr.Exec(ctx,
		"go", "run",
		"github.com/syumai/workers/cmd/workers-assets-gen@v0.18.0",
	); err != nil {
		return fmt.Errorf("generate assets: %w", err)
	}
	return rnr.Exec(ctx,
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

func (op Ops) Deploy() error {
	if err := op.Build(); err != nil {
		return err
	}
	ctx := context.Background()
	env, err := fetchSecrets(ctx, secrets)
	if err != nil {
		return err
	}
	ctx = command.WithEnv(ctx, env)
	return rnr.Exec(ctx, "npx", "wrangler", "deploy")
}

func (op Ops) Dev() error {
	if err := op.Build(); err != nil {
		return err
	}
	return rnr.Exec(context.Background(), "npx", "wrangler", "dev")
}

func fetchSecrets(
	ctx context.Context, s map[string]string,
) (map[string]string, error) {
	sh := command.Shell(sys.Machine(), "go")
	sh.Handle("spkez", sh.Unshell())

	err := sh.Do(ctx, "spkez", "--version")
	if command.NotFound(err) {
		err = sh.Exec(ctx,
			"go", "install", "lesiw.io/spkez@latest")
		if err != nil {
			return nil, fmt.Errorf("install spkez: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("check spkez: %w", err)
	}

	spkez := sub.Machine(sh, "spkez")

	ret := make(map[string]string)
	for k, v := range s {
		if value := sh.Env(ctx, k); value != "" {
			ret[k] = value
		} else {
			val, err := command.Read(ctx, spkez, "get", v)
			if err != nil {
				return nil, fmt.Errorf(
					"get secret %s: %w", k, err)
			}
			ret[k] = val
		}
	}
	return ret, nil
}
