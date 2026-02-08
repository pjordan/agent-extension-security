# Quickstart: Verify & Install

This walkthrough covers the consumer side: receiving an extension bundle from a publisher, verifying its integrity, and installing with policy enforcement.

## What you receive

A published extension bundle typically contains:

| File | Purpose |
|------|---------|
| `<name>.aext` | Packaged extension artifact |
| `aem.json` | Permission manifest (what the extension requests) |
| `<name>.sig.json` | Ed25519 signature |
| `scan.json` | Heuristic scan results |
| `sbom.spdx.json` | Software bill of materials |
| `provenance.json` | Build provenance record |

## Step-by-step

### 1) Obtain the publisher's public key

You need the publisher's public key to verify the signature. Obtain it out-of-band — from the publisher's repository, key server, or documentation.

```bash
# Example: download from the publisher's repo
curl -o publisher-key.json https://example.com/keys/publisher-key.json
```

!!! warning
    Never use `--allow-embedded-key` for third-party extensions. It only proves the artifact was signed by *someone*, not *who*.

### 2) Verify the signature

```bash
agentsec verify extension.aext --sig extension.sig.json --pub publisher-key.json
```

This checks that:

- The artifact digest matches the signature
- The signature was produced by the trusted public key

### 3) Validate the manifest

```bash
agentsec manifest validate aem.json
```

This checks:

- Schema string matches `aessf.dev/aem/v0`
- Extension type is valid (`skill`, `mcp-server`, `plugin`)
- Version is valid semver
- No unknown fields in the JSON

### 4) Review the scan results

```bash
cat scan.json
```

Look for findings that indicate risky behavior (shell commands in SKILL.md, suspicious script patterns). The scanner is heuristic — it catches common patterns, not everything.

### 5) Review the manifest permissions

```bash
cat aem.json | python3 -m json.tool
```

Key questions to ask:

- Does the extension need `allow_shell`? Why?
- Does it declare network access? To which domains?
- Does it need file write access? To which paths?

### 6) Write or choose a policy

Create a policy file that matches your security posture. Three templates are provided:

| Policy | Mode | Denies | Use case |
|--------|------|--------|----------|
| `examples/policies/permissive.json` | `warn` | Nothing | Local dev |
| `examples/policies/strict.json` | `enforce` | `allow_shell`, `allow_ip_literals` | Staging |
| `examples/policies/enterprise.json` | `enforce` | `allow_shell`, `allow_ip_literals` | Production |

See [Writing Policy](policy.md) for how to customize policy files.

### 7) Install with policy enforcement

```bash
agentsec install extension.aext \
  --sig extension.sig.json \
  --pub publisher-key.json \
  --aem aem.json \
  --policy examples/policies/strict.json \
  --dest ./installed-extensions/extension-name
```

What happens under the hood:

1. Signature is re-verified against the trusted public key
2. AEM manifest is validated
3. Manifest permissions are evaluated against policy deny rules
4. If `mode=enforce` and any denied permission is declared, install fails
5. If `mode=warn`, denied findings are printed but install continues
6. Artifact is extracted with hardened archive handling (symlink, size, ratio checks)

### Dev mode (for testing only)

For local evaluation, `--dev` skips signature verification and uses a permissive policy:

```bash
agentsec install extension.aext --dev --aem aem.json --dest ./test-install
```

!!! warning
    Never use `--dev` for third-party extensions in production.

## Next steps

- [Writing Policy](policy.md) — customize policy for your environment
- [Threat Model](threat-model.md) — understand the threats extensions can pose
- [Examples & Policies](../examples.md) — see permission gradient and policy interaction
- [Troubleshooting](../troubleshooting.md) — common errors and fixes
