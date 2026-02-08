package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/pjordan/agent-extension-security/internal/util"
)

// This is a minimal reference SBOM format. It's *not* a full SPDX implementation.
// Replace this with Syft or a real SPDX/CycloneDX generator as a next step.
type RefSBOM struct {
    Format    string    `json:"format"`
    CreatedAt time.Time `json:"created_at"`
    Artifact  struct {
        Path   string `json:"path"`
        Digest string `json:"digest"`
    } `json:"artifact"`
    Notes string `json:"notes"`
}

func runSBOM(args []string) {
    fs := newFlagSet("sbom")
    out := fs.String("out", "", "output sbom json path")
    dieIf(fs.Parse(args))
    if fs.NArg() < 1 || *out == "" {
        dieIf(fmt.Errorf("usage: agentsec sbom <artifact.aext> --out <sbom.json>"))
    }
    art := fs.Arg(0)
    sha, err := util.Sha256File(art)
    dieIf(err)

    var sb RefSBOM
    sb.Format = "SPDX-JSON (reference)"
    sb.CreatedAt = time.Now().UTC()
    sb.Artifact.Path = art
    sb.Artifact.Digest = "sha256:" + sha
    sb.Notes = "This is a placeholder SBOM. Use Syft or another generator for production."

    b, err := json.MarshalIndent(&sb, "", "  ")
    dieIf(err)
    dieIf(os.WriteFile(*out, append(b, '\n'), 0o644))
    fmt.Println("wrote sbom:", *out)
}
