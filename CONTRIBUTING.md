# Contributing

Thanks for contributing! This repo is intentionally small and easy to extend.

## Development setup

- Go 1.23+
- `make build` builds `./bin/agentsec`
- `make test` runs unit tests
- `make examples-test` validates example extensions (manifests, packaging, scanning)
- `bash scripts/docs-smoke.sh` runs the documented quickstart flow
- `make fmt-check` verifies formatting
- `make lint` runs static analysis (`golangci-lint`)
- `make hooks` enables local pre-commit checks

## Documentation workflow

- Docs source lives under `docs/` and is built with MkDocs Material.
- Install docs tooling with `pip install -r docs/requirements.txt`.
- Run `mkdocs build --strict` locally before larger docs PRs.
- Keep command examples copy/paste-valid. If command behavior changes, update:
  - `docs/quickstart.md` (authoritative full walkthrough)
  - `docs/cli-reference.md`
  - `scripts/docs-smoke.sh`
  - The README contains only the 3-command express path — update it only if `init`, `package`, or `install --dev` change.

## Suggested first contributions

- Add Cosign/Sigstore-backed signing + verification (see `docs/sigstore.md`)
- Expand `scan` rules (skill markdown + packaged scripts)
- Add dynamic sandbox runner (container/VM) and behavior attestations
- Add more schemas + strict validation for AEM/APM

## Style

- Keep dependencies minimal
- Prefer small, composable packages under `internal/`
- Add tests for anything that parses manifests or touches security boundaries

## Pull request expectations

- Keep pull requests scoped to one logical change.
- Update docs when commands, behavior, or security assumptions change.
- Include tests for new behavior and edge cases.
- Use clear commit messages that explain intent and security implications when relevant.
