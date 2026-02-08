package main

import (
    "fmt"
    "os"

    "github.com/pjordan/agent-extension-security/internal/manifest"
)

func runManifest(args []string) {
    if len(args) < 1 {
        dieIf(fmt.Errorf("usage: agentsec manifest <init|validate> ..."))
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
    dieIf(fs.Parse(args))
    if fs.NArg() < 1 {
        dieIf(fmt.Errorf("usage: agentsec manifest init <dir> --id ... --type ... --version ... --out ..."))
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
                Read:  []string{"~/Documents/**"},
                Write: []string{"~/Documents/**"},
            },
            Network: manifest.NetPerms{
                Domains:         []string{},
                AllowIPLiterals: false,
            },
            Process: manifest.ProcessPerm{
                AllowShell:      false,
                AllowSubprocess: true,
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
    dieIf(fs.Parse(args))
    if fs.NArg() < 1 {
        dieIf(fmt.Errorf("usage: agentsec manifest validate <aem.json>"))
    }
    p := fs.Arg(0)
    m, err := manifest.LoadAEM(p)
    dieIf(err)
    dieIf(m.Validate())
    fmt.Println("manifest ok:", p)
}
