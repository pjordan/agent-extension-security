# Quickstart

This walkthrough takes you from install to a signed, verified, policy-checked extension.

## Prerequisites

- Go 1.23+

## Install

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
agentsec version
```

## Express path (3 commands)

Scaffold a new extension, package it, and install it in dev mode:

```bash
# Scaffold a new extension project
agentsec init ./my-skill --id com.example.my-skill --type skill

# Edit my-skill/aem.json to declare the permissions your extension needs, then:
agentsec package ./my-skill --out my-skill.aext

# Install in dev mode (skips signature verification, uses permissive policy)
agentsec install my-skill.aext --dev --aem my-skill/aem.json --dest ./installed
```

When you are ready for production, follow the full flow below to add signing and policy enforcement.

## Full flow

### 1) Package a skill

```bash
mkdir -p ./_demo
agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext
```

### 2) Create and validate an AEM manifest

```bash
agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json

agentsec manifest validate ./_demo/aem.json
```

### 3) Generate attestations

```bash
agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json

agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json

agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json
```

### 4) Sign and verify

```bash
agentsec keygen --out ./_demo/devkey.json
agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --pub ./_demo/devkey.json
```

### 5) Install with policy enforcement

```bash
agentsec install ./_demo/hello-world.aext \
  --sig ./_demo/hello-world.sig.json \
  --pub ./_demo/devkey.json \
  --aem ./_demo/aem.json \
  --policy ./docs/policy.example.json \
  --dest ./_demo/install
```

### 6) Validate artifacts

```bash
ls -lah ./_demo
cat ./_demo/scan.json
```

## One-command demo (from source)

If you have cloned the repo and want to run the full flow automatically:

```bash
make build
bash scripts/demo.sh
```
