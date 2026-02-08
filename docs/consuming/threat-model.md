# Threat Model

This document describes the threat model for agent extensions (skills, MCP servers, plugins, connectors).

## Assets to protect

| Asset | Risk |
|-------|------|
| User secrets (API keys, OAuth tokens) | Exfiltration via over-privileged extensions |
| Confidential documents and emails | Unauthorized read access |
| Account integrity (email, calendar, ticketing) | Unauthorized actions via tool access |
| Host machine integrity (files, processes) | Arbitrary code execution, persistence |

## Threat actors

| Actor | Attack vector |
|-------|---------------|
| Malicious publishers | Typosquats, impersonation of legitimate extensions |
| Compromised maintainers / CI | Supply-chain injection through trusted update channels |
| Registry compromise | Serving malicious artifacts from a trusted source |
| Social engineering | Prompting users to install unverified extensions |

## Attack classes

| Attack | Description |
|--------|-------------|
| Exfiltration via tool servers | Over-privileged MCP servers or skills leak data to external endpoints |
| Instruction malware | Malicious commands embedded in `SKILL.md` content |
| Dependency attacks | Malicious npm/pip packages bundled inside extensions |
| Update channel compromise | Serving a malicious "latest" version through a legitimate update path |
| Archive attacks | ZipSlip, symlink traversal, decompression bombs in `.aext` files |

## Mitigations

| Mitigation | Status |
|------------|--------|
| Signature verification with trusted keys (`--pub`) | Implemented |
| Install-time policy enforcement (fail closed) | Implemented |
| Least-privilege manifest defaults | Implemented |
| Strict JSON parsing (unknown fields rejected) | Implemented |
| Archive hardening (symlink blocking, size limits, ratio checks) | Implemented |
| Heuristic scanning of skill content and scripts | Implemented |
| Sigstore/Cosign keyless signing with identity binding | Planned |
| SLSA/in-toto provenance with verification | Planned |
| Real SBOM generation and vulnerability scanning | Planned |
| Runtime permission enforcement (not just install-time) | Planned |
| Secure update metadata (TUF) | Planned |
| Revocation/quarantine mechanism | Planned |

## What this scaffold does not protect against

- **Runtime enforcement**: Permissions are checked at install time only. A skill that declares `allow_shell=false` is not sandboxed at runtime — the manifest is a declaration, not an enforcement boundary.
- **Identity verification**: Current Ed25519 dev keys prove integrity (the artifact wasn't modified) but not authenticity (who signed it). Sigstore integration will add identity binding.
- **Dependency analysis**: The scanner checks skill content and shell scripts, not transitive dependencies (npm, pip, etc.).

## Next steps

- [Security Guarantees](security.md) — what `agentsec` enforces today
- [Writing Policy](policy.md) — customize policy for your threat posture
- [Production Readiness](../reference/production-readiness.md) — capability matrix and roadmap
