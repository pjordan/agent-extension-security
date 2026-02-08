# Definition of Done

A change is complete only if all apply:

- Build passes: `make build`.
- Tests pass: `go test ./...`.
- CLI UX is coherent: usage, errors, and outputs are updated.
- Security impact is reviewed and noted in the summary.
- Docs/specs are in sync with code behavior.
- No unrelated files are modified.

## PR Summary Template
- What changed
- Why it changed
- Security impact
- Validation run (commands + pass/fail)
- Follow-ups (if any)
