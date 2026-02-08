# Project Context

## Goal
`agentsec` is a reference CLI for agent extension supply-chain security. It packages, signs, verifies, scans, and installs extension artifacts, with manifest schemas and policy checks.

## Architecture Map
- `cmd/agentsec/*`: CLI command handlers and argument parsing.
- `internal/manifest/*`: Agent Extension Manifest (AEM) structures and validation.
- `internal/policy/*`: install-time policy loading and evaluation.
- `internal/crypto/*`: dev key generation and signature verification.
- `internal/util/*`: archive/hash helper utilities.
- `spec/aem/v0/*`: AEM schema.
- `spec/apm/v0/*`: APM schema scaffold.
- `docs/*`: threat model, permissions model, Sigstore roadmap.

## Current Scope
Implemented now:
- local ed25519 dev key signing/verification,
- packaging and install flow,
- minimal SBOM/provenance/scan output,
- policy gates for selected permissions.

Planned next:
- production Sigstore flows,
- stronger provenance and update security,
- split and enforce richer APM model.

## Invariants
- Keep the CLI behavior deterministic and script-friendly.
- Preserve backward compatibility for existing command flags unless a migration note is added.
- Align code and schema changes with docs (`README.md`, `GETTING_STARTED.md`, `docs/*`).
