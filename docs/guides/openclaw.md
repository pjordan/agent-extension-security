# Securing OpenClaw Skills with agentsec

This guide covers applying `agentsec` supply-chain security to OpenClaw skill packages.

## What is an OpenClaw skill?

OpenClaw skills typically use `.mdc` (Markdown Components) files to define skill content, rules, and templates. The packaging and security workflow is the same as any other skill format — only the content format and install location differ.

## Prerequisites

- Go 1.23+
- `agentsec` installed: `go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest`

## 1) Scaffold the skill project

```bash
agentsec init ./my-openclaw-skill --id com.example.my-openclaw-skill --type skill
```

## 2) Add your `.mdc` content

Place your `.mdc` files and any supporting resources in the scaffolded directory:

```
my-openclaw-skill/
├── aem.json          # manifest (from init)
├── devkey.json       # dev signing key (from init)
├── policy.json       # dev policy (from init)
├── main.mdc          # your skill definition
└── templates/
    └── response.mdc  # supporting templates
```

## 3) Declare permissions

Edit `my-openclaw-skill/aem.json` to match your skill's actual needs. Start with the least-privilege defaults from `agentsec init` and add only what's required.

## 4) Package, sign, and install

```bash
# Package
agentsec package ./my-openclaw-skill --out my-openclaw-skill.aext

# Sign
agentsec sign my-openclaw-skill.aext \
  --key my-openclaw-skill/devkey.json \
  --out my-openclaw-skill.sig.json

# Verify and install with policy enforcement
agentsec install my-openclaw-skill.aext \
  --sig my-openclaw-skill.sig.json \
  --pub my-openclaw-skill/devkey.json \
  --aem my-openclaw-skill/aem.json \
  --policy my-openclaw-skill/policy.json \
  --dest ./installed-skills/my-openclaw-skill
```

Or use dev mode for quick iteration:

```bash
agentsec package ./my-openclaw-skill --out my-openclaw-skill.aext
agentsec install my-openclaw-skill.aext --dev --aem my-openclaw-skill/aem.json --dest ./installed-skills/my-openclaw-skill
```

## Scanner limitations

The current `agentsec scan` heuristics inspect `SKILL.md`, `.sh`, and `.ps1` files. If your skill corpus is primarily `.mdc`, scanner findings may be limited. Consider extending scanner rules for `.mdc` content — see [Contributing](../../CONTRIBUTING.md) for how to add scan rules.

## Next steps

- [Generic Pipeline](pipeline.md) — full attestation and CI/CD details
- [Examples & Policies](../examples.md) — permission gradients and policy templates
- [CLI Reference](../cli-reference.md) — every command and flag
