# Securing Claude Code Skills with agentsec

This guide walks through securing a Claude Code skill end-to-end: from scaffolding to signed, policy-checked installation.

## What is a Claude Code skill?

A Claude Code skill is a folder containing a `SKILL.md` file (and optional resources) that lives under `.claude/skills/` in your project or user directory. Skills teach Claude how to perform specific tasks — code generation patterns, API integrations, workflow automation.

Because skills execute with Claude's full tool access, verifying their integrity before loading matters.

## Prerequisites

- Go 1.23+
- `agentsec` installed: `go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest`

## 1) Scaffold the skill project

```bash
agentsec init ./my-claude-skill --id com.example.my-claude-skill --type skill
```

This creates:

- `my-claude-skill/aem.json` — manifest with least-privilege defaults
- `my-claude-skill/devkey.json` — Ed25519 dev signing keypair
- `my-claude-skill/policy.json` — warn-mode policy for development

## 2) Add your skill content

Create your `SKILL.md` inside the scaffolded directory:

```bash
cat > my-claude-skill/SKILL.md << 'EOF'
# My Skill

Instructions for Claude go here.
EOF
```

Add any supporting files (scripts, templates, config) alongside it.

## 3) Declare permissions

Edit `my-claude-skill/aem.json` to declare what your skill actually needs:

```json
{
  "permissions": {
    "files": {
      "read": ["*.md", "*.txt"],
      "write": []
    },
    "network": {
      "domains": [],
      "allow_ip_literals": false
    },
    "process": {
      "allow_shell": false,
      "allow_subprocess": false
    }
  }
}
```

Start with the least permissions possible. The [Permissions & Policy](../permissions.md) reference documents every field.

## 4) Package and sign

```bash
agentsec package ./my-claude-skill --out my-claude-skill.aext
agentsec sign my-claude-skill.aext --key my-claude-skill/devkey.json --out my-claude-skill.sig.json
```

## 5) Verify and install

```bash
# Verify signature
agentsec verify my-claude-skill.aext \
  --sig my-claude-skill.sig.json \
  --pub my-claude-skill/devkey.json

# Install with policy enforcement
agentsec install my-claude-skill.aext \
  --sig my-claude-skill.sig.json \
  --pub my-claude-skill/devkey.json \
  --aem my-claude-skill/aem.json \
  --policy my-claude-skill/policy.json \
  --dest .claude/skills/my-claude-skill
```

The verified skill is now installed at `.claude/skills/my-claude-skill/`.

## Dev mode (quick iteration)

During development, skip signing and use the permissive policy:

```bash
agentsec package ./my-claude-skill --out my-claude-skill.aext
agentsec install my-claude-skill.aext --dev --aem my-claude-skill/aem.json --dest .claude/skills/my-claude-skill
```

When you're ready for production, switch to the full flow above.

## Optional: generate attestations

For supply-chain evidence, generate SBOM, provenance, and scan reports before signing:

```bash
agentsec sbom my-claude-skill.aext --out sbom.spdx.json
agentsec provenance my-claude-skill.aext \
  --source-repo https://github.com/your-org/your-skill-repo \
  --source-rev "$(git rev-parse HEAD)" \
  --out provenance.json
agentsec scan my-claude-skill.aext --out scan.json
```

See the [Generic Pipeline](pipeline.md) for full attestation details.

## Next steps

- [Add a pre-load verification hook](claude-code-hook.md) to gate all skill directory writes
- [Examples & Policies](../examples.md) — see permission gradients and policy templates
- [CLI Reference](../cli-reference.md) — every command and flag
