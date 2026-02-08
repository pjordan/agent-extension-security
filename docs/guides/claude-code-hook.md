# Adding a Pre-Load Verification Hook in Claude Code

This guide shows how to use a Claude Code `PreToolUse` hook to enforce that all skill directory writes go through `agentsec` verification first.

## Why a hook?

Claude Code doesn't expose a dedicated "skill load" event. The practical equivalent is a `PreToolUse` hook that intercepts tool calls (like `Bash`) before they can modify skill directories. Any write to `.claude/skills/` that doesn't go through your verified installer is denied.

## Prerequisites

- `agentsec` installed: `go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest`
- A trusted signing key (from `agentsec keygen` or your team's key)
- `jq` available on `PATH`

## 1) Add the hook to `.claude/settings.json`

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

This intercepts every `Bash` tool call and runs your verification script before allowing it.

## 2) Create the hook script

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

## 3) Create the verified install wrapper

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

Make it executable:

```bash
chmod +x ./scripts/install-verified-skill.sh
```

## How it works

1. Claude (or a user) runs a `Bash` command that touches `.claude/skills/`
2. The `PreToolUse` hook fires and checks the command
3. If the command uses `install-verified-skill.sh`, it's allowed through
4. Otherwise, the hook denies the command with an explanation
5. The install wrapper runs `agentsec verify` → `manifest validate` → `install --policy` → promote

This gives you a real, enforceable pre-load gate in Claude workflows.

## Example usage

```bash
# Package and sign a skill
agentsec package ./my-skill --out my-skill.aext
agentsec sign my-skill.aext --key devkey.json --out my-skill.sig.json

# Install through the verified wrapper
./scripts/install-verified-skill.sh \
  my-skill.aext \
  my-skill.sig.json \
  devkey.json \
  my-skill/aem.json \
  examples/policies/strict.json \
  my-skill
```

## Next steps

- [Securing Claude Code Skills](claude-code.md) — full skill creation walkthrough
- [Examples & Policies](../examples.md) — policy templates (permissive, strict, enterprise)
- [Troubleshooting](../troubleshooting.md) — common errors and fixes
