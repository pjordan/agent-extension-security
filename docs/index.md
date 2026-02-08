# Agent Extension Security

`agentsec` is a reference CLI and policy engine for securing agent extensions —
skills, MCP servers, plugins, and connectors.

!!! tip "New here?"
    Jump straight to the [Quickstart](quickstart.md) to scaffold, sign, and
    install your first extension in under a minute.

## The pipeline

Every extension goes through a five-stage supply-chain pipeline:

| Stage | Command | What it does |
|-------|---------|--------------|
| **Scaffold** | `agentsec init` | Creates a project with a manifest, dev key, and policy |
| **Package** | `agentsec package` | Zips the directory into a `.aext` artifact |
| **Attest** | `agentsec sbom`, `provenance`, `scan` | Generates SBOM, provenance, and scan reports |
| **Sign** | `agentsec sign` | Signs the artifact digest with an Ed25519 key |
| **Install** | `agentsec install` | Verifies signature, enforces policy, then extracts |

See the [CLI Reference](cli-reference.md) for full command details.

## Choose your path

**Getting started**

- [Quickstart](quickstart.md) — end-to-end walkthrough in 3 minutes
- [Install & Release Verification](install.md) — prebuilt binaries and checksum verification
- [Examples & Policies](examples.md) — sample extensions and policy templates

**Integrating**

- [Guides](guides/index.md) — platform-specific walkthroughs
    - [Claude Code](guides/claude-code.md) — secure and install skills
    - [Claude Code Hook](guides/claude-code-hook.md) — pre-load verification
    - [OpenClaw](guides/openclaw.md) — `.mdc` skill security
    - [Codex](guides/codex.md) — Codex skill security
    - [Generic Pipeline](guides/pipeline.md) — the full pipeline for any format
- [Examples & Policies](examples.md) — permission gradient and policy interaction

**Security model**

- [Threat Model](threat-model.md) — assumptions and trust boundaries
- [Security Hardening](security-hardening.md) — archive, manifest, and policy hardening
- [Permissions & Policy](permissions.md) — AEM permissions and deny rules
- [Sigstore Plan](sigstore.md) — planned keyless signing roadmap

**Specs & reference**

- [AEM Schema](spec-aem.md) — Agent Extension Manifest specification
- [APM Schema](spec-apm.md) — Agent Policy Manifest specification
- [CLI Reference](cli-reference.md) — every command and flag
- [Comparison](comparison.md) — how agentsec compares to alternatives
- [Troubleshooting](troubleshooting.md) — common errors and fixes

## Project status

This project is intentionally a scaffold: hardened where it matters for a
reference implementation, but not yet a full production supply-chain platform.

See [Production Readiness](production-readiness.md) for an explicit capability
matrix and next-step roadmap.
