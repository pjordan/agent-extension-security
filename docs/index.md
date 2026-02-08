# Agent Extension Security

`agentsec` is a reference CLI and policy engine for securing agent extensions —
skills, MCP servers, plugins, and connectors.

## The pipeline

Every extension goes through a five-stage supply-chain pipeline:

| Stage | Command | What it does |
|-------|---------|--------------|
| **Scaffold** | `agentsec init` | Creates a project with a manifest, dev key, and policy |
| **Package** | `agentsec package` | Zips the directory into a `.aext` artifact |
| **Attest** | `agentsec sbom`, `provenance`, `scan` | Generates SBOM, provenance, and scan reports |
| **Sign** | `agentsec sign` | Signs the artifact digest with an Ed25519 key |
| **Install** | `agentsec install` | Verifies signature, enforces policy, then extracts |

## Choose your path

<div class="grid cards" markdown>

-   **Building an extension?**

    ---

    Package, sign, and publish skills, MCP servers, or plugins.

    [:octicons-arrow-right-24: Creator Quickstart](creating/quickstart.md)

-   **Installing an extension?**

    ---

    Verify, enforce policy, and safely install extensions from others.

    [:octicons-arrow-right-24: Consumer Quickstart](consuming/quickstart.md)

</div>

## What you can do today

- **Scaffold** a new extension project with `agentsec init`
- **Package** any skill directory into a signed `.aext` artifact
- **Declare** least-privilege permissions in an AEM manifest
- **Enforce** install-time policy with deny rules (fail closed or warn)
- **Verify** signatures against trusted public keys
- **Scan** skill content and scripts for common risk patterns
- **Integrate** with [Claude Code](creating/platforms/claude-code.md), [OpenClaw](creating/platforms/openclaw.md), and [Codex](creating/platforms/codex.md)

## Quick links

| Topic | Link |
|-------|------|
| Install the CLI | [Install](install.md) |
| Examples & policy templates | [Examples & Policies](examples.md) |
| CLI command reference | [CLI Reference](reference/cli.md) |
| How agentsec compares | [Comparison](comparison.md) |
| Common errors and fixes | [Troubleshooting](troubleshooting.md) |
| Security model | [Threat Model](consuming/threat-model.md), [Security Guarantees](consuming/security.md) |
| Production readiness | [Production Readiness](reference/production-readiness.md) |

## Project status

This project is intentionally a scaffold: hardened where it matters for a
reference implementation, but not yet a full production supply-chain platform.

See [Production Readiness](reference/production-readiness.md) for an explicit capability
matrix and next-step roadmap.
