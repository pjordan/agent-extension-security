# Publishing Checklist

Use this checklist before publishing an extension to ensure it meets security and quality standards.

## Pre-publish verification

- [ ] **Manifest validates** — `agentsec manifest validate aem.json` passes
- [ ] **Permissions are minimal** — only declare what the extension genuinely needs
- [ ] **Scan passes** — `agentsec scan` output reviewed, no unexpected findings
- [ ] **Signature verifies** — `agentsec verify` succeeds with your signing key
- [ ] **Policy install works** — `agentsec install` succeeds with a representative policy
- [ ] **Attestations generated** — SBOM, provenance, and scan reports are current

## Bundle file list

Publish all of these together:

| File | Required | Purpose |
|------|----------|---------|
| `<name>.aext` | Yes | Packaged extension artifact |
| `aem.json` | Yes | Permission manifest |
| `<name>.sig.json` | Yes | Ed25519 signature |
| `scan.json` | Recommended | Heuristic scan results |
| `sbom.spdx.json` | Recommended | Software bill of materials |
| `provenance.json` | Recommended | Build provenance record |

## Maturity caveats

Current scaffold limitations to communicate to consumers:

- Signatures use local Ed25519 dev keys (no identity binding)
- SBOM and provenance are reference-quality, not compliance-grade
- Scanner coverage is limited to `SKILL.md`, `.sh`, and `.ps1` heuristics

See [Production Readiness](../reference/production-readiness.md) for the full capability matrix.

## Consumer experience test

Before publishing, test the full consumer flow:

```bash
# Simulate what a consumer will do
agentsec install your-extension.aext \
  --sig your-extension.sig.json \
  --pub your-signing-key.json \
  --aem aem.json \
  --policy examples/policies/strict.json \
  --dest ./test-consumer-install
```

If the install fails due to policy denial, decide whether to:

1. Reduce your permissions (preferred)
2. Document why the permission is needed so consumers can adjust their policy

## Next steps

- [Examples & Policies](../examples.md) — test against different policy templates
- [CI/CD Pipeline](ci-cd.md) — automate this checklist in CI
