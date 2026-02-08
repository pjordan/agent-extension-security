# Agent Extension Security (`agentsec`)

[![CI](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml)
[![CodeQL](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml)
[![Docs](https://github.com/pjordan/agent-extension-security/actions/workflows/docs.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/docs.yml)
[![Release](https://img.shields.io/github/v/release/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/releases)
[![License](https://img.shields.io/github/license/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/go.mod)

Open-source reference implementation for supply chain security in agent ecosystems (skills, MCP servers, plugins, connectors).

## What this repo provides

- `agentsec` CLI to package, sign, verify, scan, and install extensions
- Starter manifest and policy model (AEM now, APM planned)
- Hardened archive handling and strict manifest/policy parsing
- Example extensions (skill + MCP placeholder)

## Documentation

- Docs site: https://pjordan.github.io/agent-extension-security/
- Quickstart: `docs/quickstart.md`
- CLI reference: `docs/cli-reference.md`
- Security model: `docs/threat-model.md`, `docs/security-hardening.md`, `docs/permissions.md`
- Production readiness boundaries: `docs/production-readiness.md`

## Quickstart

```bash
make build
./bin/agentsec version

mkdir -p ./_demo
./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext

./bin/agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json
./bin/agentsec manifest validate ./_demo/aem.json

./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json
./bin/agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json
./bin/agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json

./bin/agentsec keygen --out ./_demo/devkey.json
./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --pub ./_demo/devkey.json

./bin/agentsec install ./_demo/hello-world.aext \
  --sig ./_demo/hello-world.sig.json \
  --pub ./_demo/devkey.json \
  --aem ./_demo/aem.json \
  --policy ./docs/policy.example.json \
  --dest ./_demo/install
```

Or run the scripted flow:

```bash
bash scripts/demo.sh
```

## Current status

This is an initial scaffold designed to be easy to extend.

Implemented hardening includes:

- Trusted-key verification by default (`--pub`)
- Install-time policy enforcement (`--aem` + `--policy`)
- Least-privilege defaults for generated manifests
- Symlink/path/resource hardening during archive package/extract
- Strict JSON decoding for manifest and policy files

Planned production-grade additions:

- Sigstore/Cosign keyless signing and identity verification
- SLSA/in-toto provenance
- Real SBOM generation and deeper scanning

See `docs/production-readiness.md` for details.

## Contributing and security

- Contributing: `CONTRIBUTING.md`
- Security policy: `SECURITY.md`
- Code of conduct: `CODE_OF_CONDUCT.md`

## License

Apache-2.0 (`LICENSE`).
