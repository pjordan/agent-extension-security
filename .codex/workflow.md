# Agent Workflow

## Standard Loop
1. Read task + relevant files.
2. Implement minimal correct change.
3. Run validation:
   - `make build`
   - `go test ./...`
4. Update docs/specs if behavior or contract changed.
5. Summarize files changed, risks, and follow-up work.

## Validation Rules
- If CLI flags/usage changed, update:
  - `cmd/agentsec/main.go` usage text,
  - `README.md`,
  - `GETTING_STARTED.md` if flow changed.
- If manifest/policy semantics changed, update:
  - `spec/*` schemas,
  - `docs/permissions.md` and/or `docs/threat-model.md`.

## Safety Rules
- Never weaken verification or policy checks without explicit task scope.
- Prefer additive compatibility for schemas and command flags.
- Call out any scaffold placeholders that are still non-production (for example Sigstore).
