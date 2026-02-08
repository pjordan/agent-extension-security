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
5. Run validation:
   ```bash
   make build && make test && make examples-test && bash scripts/docs-smoke.sh
   ```
6. Update docs:
   - `docs/cli-reference.md` — every command and flag
   - `docs/quickstart.md` — if the main flow changed
   - `scripts/docs-smoke.sh` — smoke test must match real commands
   - `docs/guides/` — update platform-specific guides if the command is referenced there
   - README only has the express path; update only if `init`, `package`, or `install --dev` changed

## Guardrails

- Keep outputs stable and machine-readable.
- Keep errors explicit and actionable.
- Do not remove existing flags/behavior without migration notes.
- Prefer additive changes to maintain backward compatibility.
