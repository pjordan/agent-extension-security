# Agent Extension Security (`agentsec`)

[![CI](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/ci.yml)
[![CodeQL](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/codeql.yml)
[![Docs](https://github.com/pjordan/agent-extension-security/actions/workflows/docs.yml/badge.svg)](https://github.com/pjordan/agent-extension-security/actions/workflows/docs.yml)
[![Release](https://img.shields.io/github/v/release/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/releases)
[![License](https://img.shields.io/github/license/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/pjordan/agent-extension-security)](https://github.com/pjordan/agent-extension-security/blob/main/go.mod)

Open-source reference implementation for supply chain security in agent ecosystems (skills, MCP servers, plugins, connectors).

## Why agentsec?

- **Package → sign → verify → install** in one CLI with no external dependencies
- **Install-time policy enforcement** — deny specific permissions before an extension runs
- **Least-privilege defaults** — scaffolded manifests start with zero permissions
- **Hardened archives** — symlink, path traversal, and resource-limit protections built in

See [how agentsec compares](https://pjordan.github.io/agent-extension-security/comparison/) to Sigstore/Cosign, npm/PyPI security, and rolling your own.

## Get started in 60 seconds

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
agentsec init ./my-skill --id com.example.my-skill --type skill
agentsec package ./my-skill --out my-skill.aext
agentsec install my-skill.aext --dev --aem my-skill/aem.json --dest ./installed
```

Ready for signing and policy enforcement? See the [full quickstart](https://pjordan.github.io/agent-extension-security/quickstart/).

## Documentation

Full docs at **<https://pjordan.github.io/agent-extension-security/>**

| Topic | Link |
|-------|------|
| Quickstart | [quickstart](https://pjordan.github.io/agent-extension-security/quickstart/) |
| CLI Reference | [cli-reference](https://pjordan.github.io/agent-extension-security/cli-reference/) |
| Install & Release Verification | [install](https://pjordan.github.io/agent-extension-security/install/) |
| Guides | [Claude Code](https://pjordan.github.io/agent-extension-security/guides/claude-code/), [OpenClaw](https://pjordan.github.io/agent-extension-security/guides/openclaw/), [Codex](https://pjordan.github.io/agent-extension-security/guides/codex/), [Pipeline](https://pjordan.github.io/agent-extension-security/guides/pipeline/) |
| Examples & Policies | [examples](https://pjordan.github.io/agent-extension-security/examples/) |
| Security Model | [threat-model](https://pjordan.github.io/agent-extension-security/threat-model/), [security-hardening](https://pjordan.github.io/agent-extension-security/security-hardening/), [permissions](https://pjordan.github.io/agent-extension-security/permissions/) |
| Production Readiness | [production-readiness](https://pjordan.github.io/agent-extension-security/production-readiness/) |
| Troubleshooting | [troubleshooting](https://pjordan.github.io/agent-extension-security/troubleshooting/) |
| Comparison | [comparison](https://pjordan.github.io/agent-extension-security/comparison/) |

## Current status

This is a reference scaffold — hardened where it matters, but not yet a full production supply-chain platform. See [Production Readiness](https://pjordan.github.io/agent-extension-security/production-readiness/) for the capability matrix and roadmap.

## Contributing and security

- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [Code of conduct](CODE_OF_CONDUCT.md)

## License

Apache-2.0 ([LICENSE](LICENSE)).
