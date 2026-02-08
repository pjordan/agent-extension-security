# Quickstart

This walkthrough takes you from source checkout to a signed, verified, policy-checked install.

## Prerequisites

- Go 1.22+
- Make (optional)

## 1) Build

```bash
make build
./bin/agentsec version
```

## 2) Package a skill

```bash
mkdir -p ./_demo
./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext
```

## 3) Create and validate an AEM manifest

```bash
./bin/agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json

./bin/agentsec manifest validate ./_demo/aem.json
```

## 4) Generate attestations

```bash
./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json

./bin/agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json

./bin/agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json
```

## 5) Sign and verify

```bash
./bin/agentsec keygen --out ./_demo/devkey.json
./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --pub ./_demo/devkey.json
```

`--allow-embedded-key` exists for insecure/dev-only compatibility.

## 6) Install with policy enforcement

```bash
./bin/agentsec install ./_demo/hello-world.aext \
  --sig ./_demo/hello-world.sig.json \
  --pub ./_demo/devkey.json \
  --aem ./_demo/aem.json \
  --policy ./docs/policy.example.json \
  --dest ./_demo/install
```

## 7) Validate artifacts

```bash
ls -lah ./_demo
cat ./_demo/scan.json
```

## One-command demo

```bash
bash scripts/demo.sh
```
