package main

import (
	"fmt"

	"github.com/pjordan/agent-extension-security/internal/util"
)

func runPackage(args []string) {
	fs := newFlagSet("package")
	out := fs.String("out", "", "output artifact path (.aext)")
	dieIf(fs.Parse(args))
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
