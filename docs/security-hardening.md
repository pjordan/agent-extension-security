# Security Hardening Changes

This document tracks the hardening changes made to the scaffold implementation.

## 1) Trusted-key verification by default

### What changed
- `agentsec verify` now requires `--pub <pubkey.json>` by default.
- `--allow-embedded-key` is available for insecure/dev-only compatibility.

### Why
- Prevents trust-root confusion where an attacker can self-sign and embed their own key.

### Migration note
- Existing scripts that relied on embedded signature keys must now pass either:
  - `--pub <pubkey.json>` (recommended), or
  - `--allow-embedded-key` (legacy/dev-only).

## 2) Install-time policy enforcement

### What changed
- `agentsec install` now requires:
  - `--pub <pubkey.json>` or `--allow-embedded-key`
  - `--aem <aem.json>`
  - `--policy <policy.json>`
- Install validates the AEM, evaluates policy findings, and fails closed in `mode=enforce`.
- In `mode=warn`, findings are printed and install continues.

### Why
- Moves policy controls from advisory to enforced in the installation path.

### Migration note
- Existing install scripts must provide `--aem` and `--policy`.

## 3) Least-privilege manifest defaults

### What changed
- `agentsec manifest init` now defaults to:
  - `files.read=[]`
  - `files.write=[]`
  - `network.domains=[]`
  - `network.allow_ip_literals=false`
  - `process.allow_shell=false`
  - `process.allow_subprocess=false`

### Why
- Reduces accidental over-privilege in newly generated manifests.

## 4) Packaging/extraction symlink hardening

### What changed
- Packaging refuses symlink and non-regular file entries.
- Extraction refuses symlink entries.

### Why
- Prevents symlink-based exfiltration and path redirection attacks.

## 5) Archive extraction limits

### What changed
- Extraction now enforces:
  - maximum file count
  - maximum per-entry uncompressed size
  - maximum total uncompressed size
  - suspicious compression ratio detection

### Why
- Mitigates zip bomb and resource-exhaustion risks.

## 6) Strict JSON decoding and schema/version checks

### What changed
- Manifest and policy loaders now:
  - reject unknown fields
  - reject trailing JSON values
- AEM validation now enforces:
  - `schema == "aessf.dev/aem/v0"`
  - semantic version format (e.g., `1.2.3`)
- Policy validation now enforces `mode` in `enforce|warn`.

### Why
- Prevents silent config drift and malformed-policy bypasses.

## 7) Adversarial tests

### What changed
- Added tests for:
  - trusted-key enforcement vs embedded-key fallback
  - install policy enforcement (`enforce` vs `warn`)
  - symlink packaging rejection
  - suspicious compression ratio rejection
  - strict parser behavior (unknown fields, trailing JSON, invalid modes)

### Why
- Keeps hardening controls regression-resistant.
