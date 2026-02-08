# Threat model (scaffold)

This document describes the initial threat model for agent extensions (skills, MCP servers, plugins).

## Assets to protect
- User secrets (API keys, OAuth tokens)
- Confidential documents and emails
- Account integrity (email/calendar/ticketing)
- Host machine integrity (files, processes)

## Threat actors
- Malicious publishers (typosquats, impersonators)
- Compromised maintainers / CI
- Registry compromise
- Prompt/social engineering leading to unsafe installs

## Attack classes
- Exfiltration via over-privileged tool servers
- "Instruction malware" inside SKILL.md (copy/paste commands)
- Malicious dependencies (npm/pip)
- Update channel compromise (serve malicious "latest")

## Mitigations (project roadmap)
- Signed artifacts (Sigstore/cosign + transparency)
- Provenance attestations (SLSA / in-toto)
- SBOM + vuln scanning (OSV, etc.)
- Least-privilege permission manifests enforced at runtime
- Secure update metadata (TUF)
- Revocation/quarantine mechanism

## Implemented hardening in this scaffold
- Verification requires a trusted public key by default (`--pub`).
- Install path enforces policy in fail-closed mode (`mode=enforce`).
- Manifest parsing and policy parsing reject unknown fields.
- Packaging/extraction rejects symlink entries.
- Extraction applies limits for entry count, entry size, total expanded bytes, and suspicious compression ratios.
