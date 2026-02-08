package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pjordan/agent-extension-security/internal/crypto"
	"github.com/pjordan/agent-extension-security/internal/manifest"
	"github.com/pjordan/agent-extension-security/internal/policy"
)

func runInit(args []string) {
	fs := newFlagSet("init")
	id := fs.String("id", "", "extension id (e.g., com.example.hello)")
	typ := fs.String("type", "", "skill|mcp-server|plugin")
	ver := fs.String("version", "0.1.0", "version (e.g., 0.1.0)")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec init <dir> --id <id> --type <type> [--version <ver>]

Scaffold a new extension project. Creates the following files inside <dir>:

  aem.json      AEM manifest with secure, least-privilege defaults
  devkey.json   Ed25519 signing keypair (keep private!)
  policy.json   Permissive warn-only install policy

Does not overwrite existing files.

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Example:
  agentsec init ./my-skill --id com.example.myskill --type skill
`)
	}
	dieIf(parseInterspersed(fs, args))

	if fs.NArg() < 1 || *id == "" || *typ == "" {
		dieIf(fmt.Errorf("usage: agentsec init <dir> --id <id> --type <type> [--version <ver>]"))
	}
	dir := fs.Arg(0)

	// Validate type and version early
	m := &manifest.AEM{
		Schema:  "aessf.dev/aem/v0",
		ID:      *id,
		Type:    *typ,
		Version: *ver,
		Permissions: manifest.Permissions{
			Files: manifest.FilePerms{
				Read:  []string{},
				Write: []string{},
			},
			Network: manifest.NetPerms{
				Domains:         []string{},
				AllowIPLiterals: false,
			},
			Process: manifest.ProcessPerm{
				AllowShell:      false,
				AllowSubprocess: false,
			},
		},
	}
	dieIf(m.Validate())

	// Create directory if needed
	if err := os.MkdirAll(dir, 0o755); err != nil {
		dieIf(fmt.Errorf("create directory %s: %w", dir, err))
	}

	// Check for existing files before creating any
	aemPath := filepath.Join(dir, "aem.json")
	keyPath := filepath.Join(dir, "devkey.json")
	policyPath := filepath.Join(dir, "policy.json")

	for _, p := range []string{aemPath, keyPath, policyPath} {
		if _, err := os.Stat(p); err == nil {
			dieIf(fmt.Errorf("%s already exists (will not overwrite)", p))
		}
	}

	// Write AEM manifest
	aemBytes, err := m.ToJSON()
	dieIf(err)
	dieIf(os.WriteFile(aemPath, append(aemBytes, '\n'), 0o644))

	// Generate and write dev keypair
	key, err := crypto.GenerateDevKeypair()
	dieIf(err)
	dieIf(crypto.SaveKey(keyPath, key))

	// Write permissive policy
	p := policy.DefaultPermissivePolicy()
	policyBytes, err := json.MarshalIndent(p, "", "  ")
	dieIf(err)
	dieIf(os.WriteFile(policyPath, append(policyBytes, '\n'), 0o644))

	// Print summary
	fmt.Printf(`Created:
  %s    (manifest — edit permissions for your use case)
  %s   (signing key — keep private, share public key)
  %s   (install policy — customize deny rules)

Next steps:
  agentsec package %s --out %s.aext
  agentsec sign %s.aext --key %s --out %s.sig.json
  agentsec install %s.aext --sig %s.sig.json --pub %s \
    --aem %s --policy %s --dest ./installed

Or use dev mode for quick local testing:
  agentsec package %s --out %s.aext
  agentsec install %s.aext --dev --aem %s --dest ./installed
`,
		aemPath, keyPath, policyPath,
		dir, filepath.Base(dir),
		filepath.Base(dir), keyPath, filepath.Base(dir),
		filepath.Base(dir), filepath.Base(dir), keyPath,
		aemPath, policyPath,
		dir, filepath.Base(dir),
		filepath.Base(dir), aemPath,
	)
}
