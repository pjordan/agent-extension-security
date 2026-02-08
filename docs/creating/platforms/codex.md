# Securing Codex Skills

This guide covers applying `agentsec` supply-chain security to Codex skill packages.

## What is a Codex skill?

A Codex skill is a folder containing a `SKILL.md` file and supporting resources. Skills are typically stored under `$CODEX_HOME/skills` (or a workspace-specific skill path) and loaded by the Codex runtime.

## Prerequisites

- Go 1.23+
- `agentsec` installed: `go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest`

## 1) Scaffold the skill project

```bash
agentsec init ./my-codex-skill --id com.example.my-codex-skill --type skill
```

## 2) Add your skill content

Place your `SKILL.md` and supporting files in the scaffolded directory:

```
my-codex-skill/
├── aem.json          # manifest (from init)
├── devkey.json       # dev signing key (from init)
├── policy.json       # dev policy (from init)
├── SKILL.md          # your skill definition
└── resources/
    └── config.json   # supporting resources
```

## 3) Declare permissions

Edit `my-codex-skill/aem.json` to match your skill's actual needs. Start with the least-privilege defaults from `agentsec init` and add only what's required.

See [Declaring Permissions](../permissions.md) for guidance on choosing values.

## 4) Package, sign, and install

```bash
# Package
agentsec package ./my-codex-skill --out my-codex-skill.aext

# Sign
agentsec sign my-codex-skill.aext \
  --key my-codex-skill/devkey.json \
  --out my-codex-skill.sig.json

# Verify and install with policy enforcement
agentsec install my-codex-skill.aext \
  --sig my-codex-skill.sig.json \
  --pub my-codex-skill/devkey.json \
  --aem my-codex-skill/aem.json \
  --policy my-codex-skill/policy.json \
  --dest "${CODEX_HOME}/skills/my-codex-skill"
```

Or use dev mode for quick iteration:

```bash
agentsec package ./my-codex-skill --out my-codex-skill.aext
agentsec install my-codex-skill.aext --dev --aem my-codex-skill/aem.json --dest "${CODEX_HOME}/skills/my-codex-skill"
```

## Consumer-side verification

When receiving skills from others, verify before promoting:

1. Verify signature with a trusted key
2. Evaluate policy against the skill's AEM manifest
3. Extract into a staging directory
4. Promote staged files into `$CODEX_HOME/skills/` only after all checks pass

```bash
agentsec install received-skill.aext \
  --sig received-skill.sig.json \
  --pub trusted-publisher-key.json \
  --aem received-skill-aem.json \
  --policy examples/policies/strict.json \
  --dest "${CODEX_HOME}/skills/received-skill"
```

See [Consuming Extensions](../../consuming/quickstart.md) for the full consumer workflow.

## Next steps

- [CI/CD Pipeline](../ci-cd.md) — full attestation and CI/CD details
- [Examples & Policies](../../examples.md) — permission gradients and policy templates
- [CLI Reference](../../reference/cli.md) — every command and flag
