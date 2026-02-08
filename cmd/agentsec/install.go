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
	dev := fs.Bool("dev", false, "dev mode: skip signature verification, use permissive policy if --policy omitted")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec install <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key) --aem <aem.json> --policy <policy.json> --dest <dir>
       agentsec install <artifact.aext> --dev --aem <aem.json> --dest <dir>

Complete secure installation pipeline: verify the artifact's signature,
evaluate the AEM manifest against the install policy, and extract the artifact
to the destination directory.

In --dev mode, signature verification is skipped and a permissive warn-only
policy is used if --policy is not provided. The --aem and --dest flags are
still required.

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  # Full secure install
  agentsec install my-skill.aext \
    --sig sig.json --pub devkey.json \
    --aem aem.json --policy policy.json \
    --dest ./installed

  # Dev mode (skip signature, permissive policy)
  agentsec install my-skill.aext --dev --aem aem.json --dest ./installed
`)
	}
	dieIf(parseInterspersed(fs, args))

	if *dev {
		// Dev mode: relaxed requirements
		if fs.NArg() < 1 || *aemPath == "" || *dest == "" {
			dieIf(fmt.Errorf("usage: agentsec install <artifact.aext> --dev --aem <aem.json> --dest <dir>"))
		}
		fmt.Fprintln(os.Stderr, "⚠ WARNING: --dev mode skips signature verification. Do not use in production.")
	} else {
		// Normal mode: all flags required
		if fs.NArg() < 1 || *sigPath == "" || *dest == "" || *aemPath == "" || *policyPath == "" {
			dieIf(fmt.Errorf("usage: agentsec install <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key) --aem <aem.json> --policy <policy.json> --dest <dir>"))
		}
	}
	art := fs.Arg(0)

	// Signature verification (skipped in dev mode)
	if !*dev {
		dieIf(verifyArtifact(art, *sigPath, *pubPath, *allowEmbeddedKey))
	}

	// Policy evaluation
	p, findings, err := evaluateInstallPolicy(*aemPath, *policyPath, *dev)
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

func evaluateInstallPolicy(aemPath, policyPath string, devMode bool) (*policy.Policy, []string, error) {
	aem, err := manifest.LoadAEM(aemPath)
	if err != nil {
		return nil, nil, fmt.Errorf("install: %w", err)
	}
	if err := aem.Validate(); err != nil {
		return nil, nil, fmt.Errorf("install: validate manifest: %w", err)
	}

	var p *policy.Policy
	if policyPath != "" {
		p, err = policy.Load(policyPath)
		if err != nil {
			return nil, nil, fmt.Errorf("install: %w", err)
		}
	} else if devMode {
		p = policy.DefaultPermissivePolicy()
	} else {
		return nil, nil, fmt.Errorf("install: --policy is required (use --dev for permissive default)")
	}

	findings := policy.Evaluate(p, aem)
	if len(findings) > 0 && p.Mode == "enforce" {
		return nil, nil, fmt.Errorf("install: policy denied install:\n - %s", strings.Join(findings, "\n - "))
	}
	return p, findings, nil
}
