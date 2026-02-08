# Securing Claude Code Skills

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

Start with the least permissions possible. See [Declaring Permissions](../permissions.md) for guidance on choosing the right values.

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

See the [CI/CD Pipeline](../ci-cd.md) for full attestation and automation details.

## Adding a pre-load verification hook

You can add a Claude Code `PreToolUse` hook to enforce that all skill directory writes go through `agentsec` verification. This provides a real, enforceable pre-load gate in Claude workflows.

### Why a hook?

Claude Code doesn't expose a dedicated "skill load" event. The practical equivalent is a `PreToolUse` hook that intercepts tool calls (like `Bash`) before they can modify skill directories. Any write to `.claude/skills/` that doesn't go through your verified installer is denied.

### Add the hook to `.claude/settings.json`

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": ".claude/hooks/enforce-skill-verify.sh"
          }
        ]
      }
    ]
  }
}
```

### Create the hook script

Create `.claude/hooks/enforce-skill-verify.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

INPUT="$(cat)"
CMD="$(printf '%s' "${INPUT}" | jq -r '.tool_input.command // ""')"

# Only gate commands that touch Claude skill directories.
if [[ "${CMD}" != *".claude/skills/"* ]]; then
  exit 0
fi

# Allow the team's verified install wrapper.
if [[ "${CMD}" == *"./scripts/install-verified-skill.sh"* ]]; then
  exit 0
fi

# Deny everything else that touches skill directories.
jq -n '{
  hookSpecificOutput: {
    hookEventName: "PreToolUse",
    permissionDecision: "deny",
    permissionDecisionReason: "Skill directory writes must use ./scripts/install-verified-skill.sh (agentsec verify + policy install)."
  }
}'
```

Make it executable:

```bash
chmod +x .claude/hooks/enforce-skill-verify.sh
```

### Create the verified install wrapper

Create `./scripts/install-verified-skill.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/install-verified-skill.sh <artifact> <sig> <pub> <aem> <policy> <skill-name>
ARTIFACT="$1"
SIG="$2"
PUB="$3"
AEM="$4"
POLICY="$5"
SKILL_NAME="$6"
DEST=".claude/skills/${SKILL_NAME}"
STAGING="$(mktemp -d)"
trap 'rm -rf "${STAGING}"' EXIT

# Verify signature with trusted key
agentsec verify "${ARTIFACT}" --sig "${SIG}" --pub "${PUB}"

# Validate manifest
agentsec manifest validate "${AEM}"

# Install with policy enforcement into staging
agentsec install "${ARTIFACT}" \
  --sig "${SIG}" \
  --pub "${PUB}" \
  --aem "${AEM}" \
  --policy "${POLICY}" \
  --dest "${STAGING}"

# Promote from staging to skill directory only on success
rm -rf "${DEST}"
mv "${STAGING}/$(basename "${ARTIFACT}" .aext).aext" "${DEST}"

echo "Skill installed to ${DEST}"
```

### How it works

1. Claude (or a user) runs a `Bash` command that touches `.claude/skills/`
2. The `PreToolUse` hook fires and checks the command
3. If the command uses `install-verified-skill.sh`, it's allowed through
4. Otherwise, the hook denies the command with an explanation
5. The install wrapper runs `agentsec verify` → `manifest validate` → `install --policy` → promote

## Next steps

- [Declaring Permissions](../permissions.md) — guidance on choosing permission values
- [Examples & Policies](../../examples.md) — permission gradients and policy templates
- [CLI Reference](../../reference/cli.md) — every command and flag
