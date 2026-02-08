package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pjordan/agent-extension-security/internal/manifest"
	"github.com/pjordan/agent-extension-security/internal/policy"
	"github.com/pjordan/agent-extension-security/internal/util"
)

func runInstall(args []string) {
	fs := newFlagSet("install")
	sigPath := fs.String("sig", "", "signature file (json)")
	pubPath := fs.String("pub", "", "trusted pubkey file (json)")
	allowEmbeddedKey := fs.Bool("allow-embedded-key", false, "allow trusting signature public_key (insecure/dev-only)")
	aemPath := fs.String("aem", "", "agent extension manifest (json)")
	policyPath := fs.String("policy", "", "install policy file (json)")
	dest := fs.String("dest", "", "destination directory")
	dieIf(fs.Parse(args))
	if fs.NArg() < 1 || *sigPath == "" || *dest == "" || *aemPath == "" || *policyPath == "" {
		dieIf(fmt.Errorf("usage: agentsec install <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key) --aem <aem.json> --policy <policy.json> --dest <dir>"))
	}
	art := fs.Arg(0)

	dieIf(verifyArtifact(art, *sigPath, *pubPath, *allowEmbeddedKey))
	p, findings, err := evaluateInstallPolicy(*aemPath, *policyPath)
	dieIf(err)
	if len(findings) > 0 && p.Mode == "warn" {
		fmt.Fprintln(os.Stderr, "policy warnings:")
		for _, f := range findings {
			fmt.Fprintln(os.Stderr, " -", f)
		}
	}

	base := filepath.Base(art)
	outDir := filepath.Join(*dest, base)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		dieIf(err)
	}
	dieIf(util.UnzipFile(art, outDir))
	fmt.Println("installed to:", outDir)
}

func evaluateInstallPolicy(aemPath, policyPath string) (*policy.Policy, []string, error) {
	aem, err := manifest.LoadAEM(aemPath)
	if err != nil {
		return nil, nil, err
	}
	if err := aem.Validate(); err != nil {
		return nil, nil, err
	}

	p, err := policy.Load(policyPath)
	if err != nil {
		return nil, nil, err
	}
	findings := policy.Evaluate(p, aem)
	if len(findings) > 0 && p.Mode == "enforce" {
		return nil, nil, fmt.Errorf("policy denied install:\n - %s", strings.Join(findings, "\n - "))
	}
	return p, findings, nil
}
