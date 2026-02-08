# Getting Started

Get from zero to a working extension in under a minute.

## 1) Install

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
```

Or build from source:

```bash
make build
```

## 2) Scaffold a new extension

```bash
agentsec init ./my-skill --id com.example.my-skill --type skill
```

This creates `my-skill/` with:
- `aem.json` — manifest with least-privilege defaults
- `devkey.json` — Ed25519 dev signing keypair
- `policy.json` — warn-mode policy for development

## 3) Declare permissions

Edit `my-skill/aem.json` to declare what your extension needs:

```json
{
  "permissions": {
    "process": { "allow_shell": true },
    "network": { "domains": ["api.example.com"] }
  }
}
```

## 4) Package and install (dev mode)

```bash
agentsec package ./my-skill --out my-skill.aext
agentsec install my-skill.aext --dev --aem my-skill/aem.json --dest ./installed
```

## Next steps

- Add signing and policy enforcement: [Full quickstart](docs/quickstart.md)
- CLI reference: [docs/cli-reference.md](docs/cli-reference.md)
- Troubleshooting: [docs/troubleshooting.md](docs/troubleshooting.md)
- How agentsec compares to alternatives: [docs/comparison.md](docs/comparison.md)
- Security model: [docs/threat-model.md](docs/threat-model.md)
- Production readiness: [docs/production-readiness.md](docs/production-readiness.md)
