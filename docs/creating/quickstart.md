# Quickstart: Package & Sign

This walkthrough takes you from an empty directory to a signed, verified extension artifact.

## Prerequisites

- Go 1.23+
- `agentsec` installed: `go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest`

## Express path (3 commands)

Scaffold a new extension, package it, and install in dev mode:

```bash
agentsec init ./my-skill --id com.example.my-skill --type skill
agentsec package ./my-skill --out my-skill.aext
agentsec install my-skill.aext --dev --aem my-skill/aem.json --dest ./installed
```

When you're ready for signing and policy enforcement, follow the full flow below.

## Full flow

### 1) Scaffold the project

```bash
agentsec init ./my-skill --id com.example.my-skill --type skill
```

This creates:

- `my-skill/aem.json` — AEM manifest with least-privilege defaults
- `my-skill/devkey.json` — Ed25519 dev signing keypair
- `my-skill/policy.json` — warn-mode policy for development

### 2) Add your content

Create your `SKILL.md` and any supporting files:

```bash
cat > my-skill/SKILL.md << 'EOF'
# My Skill

Instructions for the agent go here.
EOF
```

### 3) Declare permissions

Edit `my-skill/aem.json` to declare what your extension actually needs. Start with the least-privilege defaults from `agentsec init` and add only what's required.

See [Declaring Permissions](permissions.md) for guidance on choosing the right values.

### 4) Package the artifact

```bash
agentsec package ./my-skill --out my-skill.aext
```

### 5) Create and validate a manifest

If you used `agentsec init`, you already have `aem.json`. Otherwise, create one:

```bash
agentsec manifest init ./my-skill \
  --id com.example.my-skill --type skill --version 0.1.0 --out my-skill/aem.json
agentsec manifest validate my-skill/aem.json
```

### 6) Generate attestations

Generate supply-chain evidence for the artifact:

```bash
agentsec sbom my-skill.aext --out sbom.spdx.json
agentsec provenance my-skill.aext \
  --source-repo https://github.com/your-org/your-skill-repo \
  --source-rev "$(git rev-parse HEAD)" \
  --out provenance.json
agentsec scan my-skill.aext --out scan.json
```

These are reference-quality outputs in the current scaffold. See [Production Readiness](../reference/production-readiness.md) for the roadmap.

### 7) Sign and verify

```bash
agentsec sign my-skill.aext --key my-skill/devkey.json --out my-skill.sig.json
agentsec verify my-skill.aext --sig my-skill.sig.json --pub my-skill/devkey.json
```

### What you've built

Your publishable bundle contains:

| File | Purpose |
|------|---------|
| `my-skill.aext` | Packaged extension artifact |
| `aem.json` | Permission manifest |
| `my-skill.sig.json` | Ed25519 signature |
| `sbom.spdx.json` | Software bill of materials |
| `provenance.json` | Build provenance record |
| `scan.json` | Heuristic scan results |

## Self-test with policy install

Verify the full consumer experience works before publishing:

```bash
agentsec install my-skill.aext \
  --sig my-skill.sig.json \
  --pub my-skill/devkey.json \
  --aem my-skill/aem.json \
  --policy my-skill/policy.json \
  --dest ./test-install
```

See [Consuming Extensions](../consuming/quickstart.md) for the consumer's perspective.

## Next steps

- [Declaring Permissions](permissions.md) — guidance on choosing permission values
- [Platform Integration](platforms/claude-code.md) — platform-specific guides
- [CI/CD Pipeline](ci-cd.md) — automate the pipeline in CI
- [Publishing Checklist](publishing-checklist.md) — final checks before publishing
