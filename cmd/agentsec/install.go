package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pjordan/agent-extension-security/internal/util"
)

func runInstall(args []string) {
	fs := newFlagSet("install")
	sigPath := fs.String("sig", "", "signature file (json)")
	dest := fs.String("dest", "", "destination directory")
	dieIf(fs.Parse(args))
	if fs.NArg() < 1 || *sigPath == "" || *dest == "" {
		dieIf(fmt.Errorf("usage: agentsec install <artifact.aext> --sig <sig.json> --dest <dir>"))
	}
	art := fs.Arg(0)

	// Verify first (reuse verify logic by calling runVerify? keep simple here)
	runVerify([]string{art, "--sig", *sigPath})

	base := filepath.Base(art)
	outDir := filepath.Join(*dest, base)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		dieIf(err)
	}
	dieIf(util.UnzipFile(art, outDir))
	fmt.Println("installed to:", outDir)
}
