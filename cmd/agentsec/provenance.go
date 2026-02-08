package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/pjordan/agent-extension-security/internal/util"
)

type RefProvenance struct {
    Schema    string    `json:"schema"`
    CreatedAt time.Time `json:"created_at"`
    Source    struct {
        Repo string `json:"repo"`
        Rev  string `json:"rev"`
    } `json:"source"`
    Artifact struct {
        Digest string `json:"digest"`
    } `json:"artifact"`
    Notes string `json:"notes"`
}

func runProvenance(args []string) {
    fs := newFlagSet("provenance")
    out := fs.String("out", "", "output provenance json path")
    sourceRepo := fs.String("source-repo", "", "source repository url")
    sourceRev := fs.String("source-rev", "", "source revision (commit/tag)")
    dieIf(fs.Parse(args))
    if fs.NArg() < 1 || *out == "" || *sourceRepo == "" || *sourceRev == "" {
        dieIf(fmt.Errorf("usage: agentsec provenance <artifact.aext> --source-repo <url> --source-rev <rev> --out <prov.json>"))
    }
    art := fs.Arg(0)
    sha, err := util.Sha256File(art)
    dieIf(err)

    var p RefProvenance
    p.Schema = "aessf.dev/attestation/provenance/v0"
    p.CreatedAt = time.Now().UTC()
    p.Source.Repo = *sourceRepo
    p.Source.Rev = *sourceRev
    p.Artifact.Digest = "sha256:" + sha
    p.Notes = "Reference provenance only. For production, emit SLSA/in-toto compatible provenance and sign it."

    b, err := json.MarshalIndent(&p, "", "  ")
    dieIf(err)
    dieIf(os.WriteFile(*out, append(b, '\n'), 0o644))
    fmt.Println("wrote provenance:", *out)
}
