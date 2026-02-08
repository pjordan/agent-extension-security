# Production Readiness

This repo is a reference implementation. Use this matrix to understand what is scaffold-level vs production-ready.

| Capability | Current State | For Production |
|---|---|---|
| Packaging / extraction hardening | Implemented (symlink blocking, unzip limits, ZipSlip defense) | Keep current controls, add additional archive-format fuzzing |
| Signature flow | Local Ed25519 dev keys | Sigstore keyless signing with identity constraints + transparency log verification |
| Policy enforcement | Install-time deny/warn on selected permissions | Expand policy model, add version-diff approvals and richer deny/allow semantics |
| Manifest validation | Strict JSON decoding + runtime schema/type/semver checks | Publish fully versioned schemas and compatibility policy |
| SBOM generation | Reference JSON placeholder | Real SPDX/CycloneDX generator (for example Syft) + signing/attestation chain |
| Provenance | Reference JSON placeholder | SLSA/in-toto provenance with signature and verification policy |
| Scanner | Lightweight heuristics over instruction/script files | Multi-engine static/dynamic analysis and CVE/dependency checks |
| Release process | Multi-platform GitHub release artifacts + checksums | Signed release artifacts, provenance for binaries, policy-gated release promotion |

## Recommended adoption posture

- Use current scaffold for prototyping policy and workflow shape.
- Treat `sbom`, `provenance`, and `scan` outputs as references, not compliance-grade outputs.
- Keep `--allow-embedded-key` strictly to local/dev workflows.
