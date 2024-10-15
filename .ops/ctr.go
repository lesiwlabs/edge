package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sub"
	"lesiw.io/cmdio/sys"
)

const goversion = "1.23.2"
const tgoversion = "0.33.0"
const ctrimg = "lledgectr" // TODO: derive from uuid
const ctrfile = ".ops/ctr.go"

var unametr = map[string]string{
	"aarch64": "arm64",
	"x86_64":  "amd64",
}

func bldctr() string {
	rnr := sys.Runner()
	defer rnr.Close()
	if !shouldbuild(rnr, ctrfile, ctrimg) {
		return ctrimg
	}

	uname := rnr.MustGet("uname", "-m").Out

	ctrid := rnr.MustGet("docker", "run", "-d", "-i", "node:22", "cat").Out
	defer rnr.Run("docker", "rm", "-f", ctrid)

	bld := sub.WithRunner(rnr, "docker", "exec", "-i", ctrid)

	bld.MustRun("npx", "--yes", "wrangler", "--version")
	bld.MustRun("wget", "-O", "/tmp/go.tar.gz",
		"https://go.dev/dl/go"+goversion+".linux-"+unametr[uname]+".tar.gz")
	bld.MustRun("tar", "-C", "/usr/local", "-xzf", "/tmp/go.tar.gz")
	bld.MustRun("wget", "-O", "/tmp/tinygo.deb",
		"https://github.com/tinygo-org/tinygo/releases/download/v"+
			tgoversion+"/tinygo_"+tgoversion+"_"+unametr[uname]+".deb")
	bld.MustRun("dpkg", "-i", "/tmp/tinygo.deb")
	cmdio.MustPipe(
		bld.Command("curl", "lesiw.io/op"),
		bld.Command("sh"),
	)

	rnr.MustRun("docker", "container", "commit",
		"-c", "WORKDIR /work",
		"-c", "ENV PATH="+bld.Env("PATH")+":/usr/local/go/bin",
		ctrid, ctrimg)

	return ctrimg
}

func getMtime(path string) (mtime int64, err error) {
	var info fs.FileInfo
	info, err = os.Lstat(path)
	if err != nil {
		return
	}
	mtime = info.ModTime().Unix()
	return
}

func shouldbuild(rnr *cmdio.Runner, file, img string) bool {
	insp, err := rnr.Get("docker", "image", "inspect",
		"--format", "{{.Created}}", img)
	mtime := mustv(getMtime(file))
	if err != nil {
		return true // Image does not exist.
	}
	ctime, err := time.Parse(time.RFC3339, insp.Out)
	if err != nil {
		panic(fmt.Sprintf(
			"failed to parse container created timestamp '%s': %s",
			insp.Out, err,
		))
	}
	return ctime.Unix() < mtime
}
