# agentsec-cli-change

Use this skill when adding or modifying `agentsec` CLI commands.

## Inputs
- Target command/subcommand.
- Desired behavior and output contract.
- Backward-compatibility constraints.

## Workflow
1. Locate command handler in `cmd/agentsec/*.go`.
2. Update flags/argument validation and error text.
3. Ensure usage text in `cmd/agentsec/main.go` is updated.
4. Add/adjust supporting logic in `internal/*` only when needed.
5. Run:
   - `make build`
   - `go test ./...`
6. Update docs:
   - `README.md`
   - `GETTING_STARTED.md` (if flow changed)

## Guardrails
- Keep outputs stable and machine-readable.
- Keep errors explicit and actionable.
- Do not remove existing flags/behavior without migration notes.
