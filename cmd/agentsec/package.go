package main

import (
	"fmt"
	"os"

	"github.com/pjordan/agent-extension-security/internal/util"
)

func runPackage(args []string) {
	fs := newFlagSet("package")
	out := fs.String("out", "", "output artifact path (.aext)")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec package <dir> --out <artifact.aext>

Package a directory into a .aext artifact (zip archive). Symlinks and
non-regular files are rejected. Junk files (.git, .DS_Store, __MACOSX) are
excluded automatically.

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Example:
  agentsec package ./my-skill --out my-skill.aext
`)
	}
	dieIf(parseInterspersed(fs, args))
	if fs.NArg() < 1 {
		dieIf(fmt.Errorf("usage: agentsec package <dir> --out <artifact.aext>"))
	}
	src := fs.Arg(0)
	if *out == "" {
		dieIf(fmt.Errorf("--out is required"))
	}
	dieIf(util.ZipDir(src, *out))
	fmt.Println("packaged:", *out)
}
