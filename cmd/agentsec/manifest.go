package main

import (
	"fmt"
	"os"

	"github.com/pjordan/agent-extension-security/internal/manifest"
)

func runManifest(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, `Usage: agentsec manifest <init|validate> <args>

Subcommands:
  init       Create an AEM manifest with least-privilege defaults
  validate   Validate an existing AEM manifest

Run 'agentsec manifest <subcommand> -h' for subcommand help.
`)
		os.Exit(2)
	}
	sub := args[0]
	subArgs := args[1:]
	switch sub {
	case "init":
		runManifestInit(subArgs)
	case "validate":
		runManifestValidate(subArgs)
	default:
		dieIf(fmt.Errorf("unknown manifest subcommand: %s", sub))
	}
}

func runManifestInit(args []string) {
	fs := newFlagSet("manifest init")
	id := fs.String("id", "", "extension id (e.g., com.example.hello)")
	typ := fs.String("type", "", "skill|mcp-server|plugin")
	ver := fs.String("version", "", "version (e.g., 0.1.0)")
	out := fs.String("out", "", "output manifest path (json)")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec manifest init <dir> --id <id> --type <type> --version <ver> --out <path>

Create an AEM (Agent Extension Manifest) with secure, least-privilege defaults.
All permissions start empty/false. Edit the generated file to add only the
permissions your extension requires.

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Example:
  agentsec manifest init ./my-skill --id com.example.hello --type skill --version 0.1.0 --out aem.json
`)
	}
	dieIf(parseInterspersed(fs, args))
	if fs.NArg() < 1 {
		dieIf(fmt.Errorf("usage: agentsec manifest init <dir> --id <id> --type <type> --version <ver> --out <path>"))
	}
	if *id == "" || *typ == "" || *ver == "" || *out == "" {
		dieIf(fmt.Errorf("--id, --type, --version, --out are required"))
	}

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

	if err := m.Validate(); err != nil {
		dieIf(err)
	}
	b, err := m.ToJSON()
	dieIf(err)
	dieIf(os.WriteFile(*out, append(b, '\n'), 0o644))
	fmt.Println("wrote manifest:", *out)
}

func runManifestValidate(args []string) {
	fs := newFlagSet("manifest validate")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec manifest validate <aem.json>

Validate an existing AEM manifest. Checks schema version, required fields,
type enum, version format, and rejects unknown fields.

Example:
  agentsec manifest validate aem.json
`)
	}
	dieIf(parseInterspersed(fs, args))
	if fs.NArg() < 1 {
		dieIf(fmt.Errorf("usage: agentsec manifest validate <aem.json>"))
	}
	p := fs.Arg(0)
	m, err := manifest.LoadAEM(p)
	dieIf(err)
	dieIf(m.Validate())
	fmt.Println("manifest ok:", p)
}
