# Agent Extension Security (agentsec)

[![CI](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml)
[![CodeQL](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml)
[![Release](https://img.shields.io/github/v/release/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/releases)
[![License](https://img.shields.io/github/license/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/go.mod)

**Agent Extension Security** is an open-source reference implementation for **supply chain security** in agent ecosystems:
Agent Skills, MCP servers, plugins, and connectors.

This repo provides:

- `agentsec` CLI: package, sign, verify, scan, and install agent extensions
- A starter **spec**: Agent Extension Manifest (AEM (JSON) + Agent Permission Manifest (APM) (planned)
- A simple, embeddable **policy engine**
- Example extensions (Skill + MCP server skeleton)

## Why this exists

Knowledge-work agents are becoming *action systems* (email, calendar, docs, tickets).
Extensions are powerful—and therefore a high-value attack surface. This project makes it easier to adopt:

- **Signed artifacts**
- **Verifiable provenance**
- **SBOMs**
- **Automated scanning**
- **Least-privilege permissions**
- **Safe install/update policies**

## Quickstart

```bash
make build
./bin/agentsec version

# Package an example skill
./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext

# Generate & validate a manifest
./bin/agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json
./bin/agentsec manifest validate ./_demo/aem.json

# Generate SBOM & provenance (minimal reference formats)
./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json
./bin/agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json

# Generate a dev keypair, sign, verify, then install
./bin/agentsec keygen --out ./_demo/devkey.json
./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json
./bin/agentsec install ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --dest ./_demo/install
```

For a more complete walkthrough, see **[GETTING_STARTED.md](GETTING_STARTED.md)**.

## Status

This is an initial scaffold intended to be easy to extend.
The signing flow currently supports **local dev keys (ed25519)** out of the box.
Keyless Sigstore flows are represented as placeholders under `docs/sigstore.md`.

## Contributing and security

- Contribution guide: [CONTRIBUTING.md](CONTRIBUTING.md)
- Security policy: [SECURITY.md](SECURITY.md)
- Code of Conduct: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## License

Apache-2.0. See [LICENSE](LICENSE).
