# Comparison with Alternatives

How does `agentsec` compare to other approaches for securing extensions and plugins?

## agentsec vs. Sigstore / Cosign

[Sigstore](https://sigstore.dev/) and [Cosign](https://github.com/sigstore/cosign) provide keyless signing and verification tied to identity (via OIDC). They are excellent tools for artifact signing in CI/CD pipelines.

| Capability | agentsec | Cosign |
|-----------|----------|--------|
| Artifact signing | Ed25519 dev keys (Sigstore planned) | Keyless OIDC, KMS, or local keys |
| Permission manifest | AEM schema with typed permissions | Not included |
| Install-time policy enforcement | Built-in (enforce/warn modes) | Not included |
| Hardened archive handling | ZipSlip protection, size limits, ratio checks | Not included (signing only) |
| Heuristic scanning | Built-in (SKILL.md, shell scripts) | Not included |
| Packaging | Built-in `.aext` format | Not included |

**When to use both:** `agentsec` is designed to adopt Sigstore/Cosign for its signing layer (see the [Sigstore roadmap](contributing.md#sigstore-integration)). Use `agentsec` for the full extension security lifecycle; use Cosign when you only need artifact signing.

**Key difference:** Cosign handles *"was this artifact signed by a trusted identity?"* — `agentsec` handles *"what does this extension want to do, and should we allow it?"*

## agentsec vs. npm / PyPI security

Language package managers (npm, PyPI, cargo, etc.) provide supply-chain security features for their respective ecosystems.

| Capability | agentsec | npm / PyPI |
|-----------|----------|------------|
| Runtime scope | Cross-runtime (any language/framework) | Single language ecosystem |
| Target | Agent extensions (skills, MCP servers, plugins) | Language packages/libraries |
| Permission model | Explicit manifest (files, network, process) | Not included (npm) / limited (PyPI trusted publishers) |
| Install policy | Enforce or warn against permission declarations | Advisory-based (npm audit, pip-audit) |
| Signing | Ed25519 per-artifact | npm provenance, PyPI attestations |
| Sandbox enforcement | Manifest-declared (runtime enforcement planned) | Not included |

**When to use agentsec:** When you are building or consuming agent extensions that may span multiple languages (a Python MCP server, a shell-based skill, a Node.js plugin) and need a consistent security model across all of them.

**When to use npm/PyPI security:** When you are managing language-specific dependencies within a single ecosystem. These tools and `agentsec` are complementary — `agentsec` secures the *extension* layer, while npm/PyPI secure the *dependency* layer within an extension.

## agentsec vs. rolling your own

Many teams consider building custom signing, verification, and policy enforcement for their extension systems.

| Concern | agentsec | Custom implementation |
|---------|----------|----------------------|
| Time to first artifact | Minutes (scaffold with `agentsec init`) | Weeks to months |
| Archive hardening | Tested: ZipSlip, symlinks, size limits, compression bombs | Must implement and test each attack vector |
| Adversarial test coverage | Included (fuzzing, malformed inputs, path traversal) | Must write from scratch |
| Manifest schema | Defined (AEM v0), extensible | Must design and version |
| Policy engine | Built-in with enforce/warn modes | Must design and implement |
| Maintenance | Community-maintained, open source | Team-maintained |
| Audit surface | Small, stdlib-only Go codebase | Varies |

**Key risk of rolling your own:** Supply-chain security has many subtle attack vectors (ZipSlip, symlink traversal, decompression bombs, key confusion). Each must be independently discovered, implemented, and tested. `agentsec` provides these primitives with existing test coverage.

## When to use agentsec

`agentsec` is a good fit when:

- [x] You are building an agent or platform that loads extensions (skills, MCP servers, plugins)
- [x] Extensions come from third parties or untrusted sources
- [x] You need to enforce permission boundaries at install time
- [x] You want a cross-runtime security model (not tied to one language)
- [x] You want to start with a working scaffold and iterate toward production

`agentsec` may not be the right fit when:

- [ ] You only need container image signing (use Cosign directly)
- [ ] You are managing language-specific package dependencies (use npm audit, pip-audit)
- [ ] You need runtime sandboxing today (agentsec declares permissions; runtime enforcement is planned)

## Current limitations

`agentsec` is intentionally a scaffold. Current limitations include:

- Ed25519 local keys only (Sigstore keyless signing is planned)
- Reference-quality SBOM and provenance (not yet SLSA-compliant)
- Heuristic scanning only (no static analysis or sandbox execution)
- No runtime enforcement (manifest declares permissions; enforcement is at install time)

See [Production Readiness](reference/production-readiness.md) for the full capability matrix and roadmap.
