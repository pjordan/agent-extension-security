# Security Guarantees

This page describes the security guarantees that `agentsec` provides today. Each guarantee includes what is enforced, how it works, and what its limitations are.

## 1. Signature verification

**Guarantee:** An installed extension's artifact has not been modified since it was signed.

**How it works:**

- `agentsec sign` computes a SHA-256 digest of the `.aext` artifact and signs it with an Ed25519 private key
- `agentsec verify` recomputes the digest and verifies the signature against a trusted public key (`--pub`)
- `agentsec install` re-verifies the signature before extracting

**Limitations:**

- Current signatures use local Ed25519 dev keys with no identity binding. They prove the artifact wasn't tampered with, but not *who* signed it.
- `--allow-embedded-key` uses the key from the signature file itself — this only proves self-consistency, not trust. It should not be used for third-party extensions.

## 2. Install-time policy enforcement

**Guarantee:** An extension will not be installed if its declared permissions violate the consumer's policy.

**How it works:**

- The consumer provides a policy file with `mode` (`enforce` or `warn`) and `deny` rules
- `agentsec install` loads the extension's AEM manifest and evaluates it against the deny rules
- In `enforce` mode, any denied finding blocks installation entirely
- In `warn` mode, findings are printed to stderr but installation continues

**Limitations:**

- Policy enforcement is install-time only. There is no runtime sandbox.
- An extension can declare fewer permissions than it actually uses. The manifest is a declaration, not an enforcement boundary.

## 3. Least-privilege defaults

**Guarantee:** Newly scaffolded extensions start with zero permissions.

**How it works:**

- `agentsec init` and `agentsec manifest init` generate manifests with all permissions set to their most restrictive values:
    - `files.read=[]`, `files.write=[]`
    - `network.domains=[]`, `network.allow_ip_literals=false`
    - `process.allow_shell=false`, `process.allow_subprocess=false`
- Extension authors must explicitly opt in to each permission

**Limitations:**

- Authors can set any permissions they want. Least-privilege defaults reduce accidental over-privilege, not intentional over-privilege. Policy enforcement handles the intentional case.

## 4. Archive safety

**Guarantee:** Packaging and extraction will not create or follow symlinks, and will not exhaust disk resources.

**How it works:**

- **Packaging** refuses symlinks and non-regular files
- **Extraction** refuses symlinks
- **Extraction** enforces:
    - Maximum file count (10,000 entries)
    - Maximum per-entry uncompressed size (64 MiB)
    - Maximum total uncompressed size (512 MiB)
    - Maximum compression ratio (200:1)

**Limitations:**

- These limits prevent common archive attacks (ZipSlip, decompression bombs) but are not a complete defense against all malicious archives.

## 5. Strict parsing

**Guarantee:** Manifest and policy files with unknown fields or invalid values will be rejected.

**How it works:**

- JSON parsing uses `DisallowUnknownFields` — any field not in the schema causes an error
- Trailing JSON values after the root object are rejected
- AEM schema, extension type, and version format are validated
- Policy mode must be exactly `enforce` or `warn`

**Limitations:**

- Strict parsing prevents config drift and typos. It does not validate the *semantics* of permission values (e.g., whether a domain allowlist is reasonable).

## Next steps

- [Threat Model](threat-model.md) — the threats these guarantees address
- [Writing Policy](policy.md) — customize enforcement for your environment
- [Production Readiness](../reference/production-readiness.md) — what's scaffold vs production-ready
