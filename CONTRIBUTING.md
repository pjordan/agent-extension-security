# Contributing

Thanks for contributing! This repo is intentionally small and easy to extend.

## Development setup

- Go 1.22+
- `make build` builds `./bin/agentsec`
- `make test` runs unit tests

## Suggested first contributions

- Add Cosign/Sigstore-backed signing + verification (see `docs/sigstore.md`)
- Expand `scan` rules (skill markdown + packaged scripts)
- Add dynamic sandbox runner (container/VM) and behavior attestations
- Add more schemas + strict validation for AEM/APM

## Style

- Keep dependencies minimal
- Prefer small, composable packages under `internal/`
- Add tests for anything that parses manifests or touches security boundaries

