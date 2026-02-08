# Guides

Platform-specific walkthroughs for integrating `agentsec` into your agent runtime.

Each guide covers the end-to-end flow: scaffold a skill, declare permissions, package, sign, and install — tailored to a specific platform's skill format and directory layout.

| Guide | Platform | What you'll learn |
|-------|----------|-------------------|
| [Claude Code Skills](claude-code.md) | Claude Code | Secure and install `.claude/skills/` packages |
| [Claude Code Verification Hook](claude-code-hook.md) | Claude Code | Add a `PreToolUse` hook to gate skill directory writes |
| [OpenClaw Skills](openclaw.md) | OpenClaw | Secure `.mdc`-based skill content |
| [Codex Skills](codex.md) | Codex | Secure `SKILL.md` packages for Codex runtimes |

Looking for the generic pipeline that works with any format? See the [Generic Pipeline](pipeline.md).

## Prerequisites

All guides assume you have `agentsec` installed:

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
agentsec version
```
