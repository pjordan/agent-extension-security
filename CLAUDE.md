# CLAUDE.md — Project Instructions

## Build & Test

```bash
make build          # Build CLI to ./bin/agentsec
make test           # go test ./...
make lint           # golangci-lint (optional, CI runs it)
make cover          # Coverage with 80% gate on ./internal/...
make fmt            # gofmt -w .
make examples-test  # Run examples smoke test
make docs-smoke     # Build CLI + run scripts/docs-smoke.sh
bash scripts/docs-smoke.sh       # Docs smoke test (needs CLI built)
bash scripts/examples-smoke.sh   # Examples smoke test
```

## Pre-PR Checklist

Run before every PR:

```bash
make build && make test && make examples-test && bash scripts/docs-smoke.sh
```

## Architecture

| Package | Purpose |
|---------|---------|
| `cmd/agentsec/` | CLI command handlers and argument parsing |
| `internal/manifest/` | AEM structures and validation |
| `internal/policy/` | Install-time policy loading and evaluation |
| `internal/crypto/` | Ed25519 dev key generation and signature verification |
| `internal/util/` | Archive/hash helpers |
| `spec/aem/v0/` | AEM JSON schema |
| `spec/apm/v0/` | APM JSON schema (scaffold) |
| `docs/` | MkDocs Material site source |

## Key Conventions

- **Zero external deps** — stdlib only. No third-party imports.
- **Strict JSON decoding** — always use `DisallowUnknownFields` for manifest/policy parsing.
- **Table-driven tests** — use `[]struct{ name string; ... }` pattern for test cases.
- **Contextual errors** — wrap with `fmt.Errorf("operation: %w", err)`.
- **80% coverage gate** — `make cover` enforces >= 80% on `./internal/...`.
- **Deterministic CLI** — outputs must be stable and machine-readable.

## Safety Rules

- Never weaken signature verification or policy enforcement.
- Prefer additive schema changes (new optional fields, not renames or removals).
- Security checks must fail closed — deny on error, not allow.
- Call out any scaffold/placeholder behavior that is not production-grade (e.g., Sigstore is planned, not implemented).

## Doc Maintenance

When CLI commands or flags change, update:

1. `docs/creating/quickstart.md` — creator walkthrough
2. `docs/reference/cli.md` — every command and flag
3. `scripts/docs-smoke.sh` — smoke test must match real commands
4. `docs/creating/platforms/` — platform-specific guides if affected

The README only has the 3-command express path. Update it only if `init`, `package`, or `install --dev` change.

When manifest/policy semantics change, also update:

- `docs/reference/aem-schema.md`
- `docs/creating/permissions.md`
- `docs/consuming/threat-model.md` (if threat assumptions changed)
- `docs/examples.md` (example manifests and policies)
- `docs/consuming/policy.md` and `docs/reference/policy-schema.md`
