# Contributing

Thanks for contributing! This repo is intentionally small and easy to extend.

## Development setup

- Go 1.23+
- `make build` builds `./bin/agentsec`
- `make test` runs unit tests
- `make examples-test` validates example extensions (manifests, packaging, scanning)
- `make cover` runs tests with 80% coverage gate on `internal/` packages
- `bash scripts/docs-smoke.sh` runs the documented quickstart flow end-to-end
- `make fmt-check` verifies formatting
- `make lint` runs static analysis (`golangci-lint`)
- `make hooks` enables local pre-commit checks

## Documentation workflow

- Docs source lives under `docs/` and is built with MkDocs Material.
- Install docs tooling with `pip install -r docs/requirements.txt`.
- Run `mkdocs build --strict` locally before larger docs PRs.
- Keep command examples copy/paste-valid. If command behavior changes, update:
    - `docs/creating/quickstart.md`
    - `docs/reference/cli.md`
    - `scripts/docs-smoke.sh`
    - The README contains only the 3-command express path — update it only if `init`, `package`, or `install --dev` change.

## Style

- Keep dependencies minimal.
- Prefer small, composable packages under `internal/`.
- Add tests for anything that parses manifests or touches security boundaries.

## Pull request expectations

- Keep pull requests scoped to one logical change.
- Update docs when commands, behavior, or security assumptions change.
- Include tests for new behavior and edge cases.
- Use clear commit messages that explain intent and security implications when relevant.

## Suggested contributions

- Add Cosign/Sigstore-backed signing + verification (see roadmap below)
- Expand `scan` rules (skill markdown + packaged scripts)
- Add dynamic sandbox runner (container/VM) and behavior attestations
- Add more schemas + strict validation for AEM/APM

---

## Development roadmap

These are planned improvements to evolve the scaffold into a production-grade platform.

### Sigstore integration

The current `agentsec sign/verify` commands use local Ed25519 dev keys to keep the scaffold self-contained.

The intended production path is Sigstore + Cosign:

- Publish extensions as OCI artifacts (or blobs) in an OCI registry
- Sign with `cosign sign` (keyless in CI)
- Verify with identity constraints:
    - `--certificate-identity` (workflow identity URI)
    - `--certificate-oidc-issuer` (GitHub Actions token issuer)
- Require transparency log inclusion

Planned CLI additions:

- `agentsec sign --keyless` (shell out to cosign)
- `agentsec verify --sigstore-policy <policy.yaml>`

References:

- [Sigstore docs](https://docs.sigstore.dev/)
- [Cosign verification](https://docs.sigstore.dev/cosign/verifying/verify/)

### APM schema split

Effective permissions currently live inside the AEM. Planned evolution:

- Split standalone APM (Agent Permission Manifest) from AEM
- Add policy-relevant permission dimensions (OAuth scopes, secrets)
- Add compatibility and migration guidance between schema versions

### Production readiness gaps

See the "For Production" column in [Production Readiness](reference/production-readiness.md) for specific gaps:

- Real SPDX/CycloneDX SBOM generation (e.g., Syft)
- SLSA/in-toto provenance with signature + verification
- Multi-engine static/dynamic scanning and CVE checks
- Signed release artifacts with provenance
- Expanded policy model with version-diff approvals

### Permission model evolution

- Add secret requirements and OAuth scopes to manifests
- Add human-approval gates for high-risk actions
- Support runtime enforcement (not just install-time)

---

## Migration notes

These document breaking changes from past hardening work.

### Trusted-key verification (from embedded-key default)

`agentsec verify` now requires `--pub <pubkey.json>` by default. Existing scripts that relied on embedded signature keys must now pass either:

- `--pub <pubkey.json>` (recommended), or
- `--allow-embedded-key` (legacy/dev-only)

### Install-time policy enforcement (from no-policy default)

`agentsec install` now requires `--aem` and `--policy`. Existing install scripts must provide both flags. For development, use `--dev` mode which applies a permissive warn-only policy.

---

## Related files

- [`SECURITY.md`](https://github.com/pjordan/agent-extension-security/blob/main/SECURITY.md) — security vulnerability reporting
- [`CODE_OF_CONDUCT.md`](https://github.com/pjordan/agent-extension-security/blob/main/CODE_OF_CONDUCT.md) — community conduct expectations
- [`AGENTS.md`](https://github.com/pjordan/agent-extension-security/blob/main/AGENTS.md) — agent conventions for AI-assisted contribution
