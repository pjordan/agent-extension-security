# Getting Started

This guide gets you from zero to a **signed, verified, policy-checked** agent extension artifact.

## Prerequisites

- Go 1.22+
- Make (optional)

## Build the CLI

```bash
make build
./bin/agentsec version
```

Optional but recommended for contributors:

```bash
make hooks
```

## Quick demo (skill)

### 1) Package an example skill

```bash
./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext
```

### 2) Generate and validate an Agent Extension Manifest (AEM)

```bash
./bin/agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json

./bin/agentsec manifest validate ./_demo/aem.json
```

### 3) Produce an SBOM (reference format)

```bash
./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json
```

### 4) Produce provenance (reference format)

```bash
./bin/agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json
```

### 5) Scan for obvious risk patterns

```bash
./bin/agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json
```

### 6) Sign and verify (trusted key required)

```bash
./bin/agentsec keygen --out ./_demo/devkey.json
./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --pub ./_demo/devkey.json
```

`--allow-embedded-key` is available for insecure/dev-only compatibility.

### 7) Install (simulated, policy enforced)

```bash
./bin/agentsec install ./_demo/hello-world.aext \
  --sig ./_demo/hello-world.sig.json \
  --pub ./_demo/devkey.json \
  --aem ./_demo/aem.json \
  --policy ./docs/policy.example.json \
  --dest ./_demo/install
```

## Next steps

- Read `docs/threat-model.md`
- Read `docs/permissions.md`
- Read `docs/security-hardening.md`
- Review `spec/aem/v0/aem.schema.json` and `spec/apm/v0/apm.schema.json`
- See `docs/sigstore.md` for how Sigstore/Cosign would plug in
