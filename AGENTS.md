## Skills
A skill is a set of local instructions to follow stored in a `SKILL.md` file. Below is the list of skills available for this repo.

### Available skills
- skill-creator: Guide for creating effective skills. Use when creating or updating a skill that extends Codex capabilities. (file: /Users/pjordan/.codex/skills/.system/skill-creator/SKILL.md)
- skill-installer: Install Codex skills into `$CODEX_HOME/skills` from a curated list or GitHub repo path. (file: /Users/pjordan/.codex/skills/.system/skill-installer/SKILL.md)
- agentsec-cli-change: Workflow for adding or modifying `agentsec` CLI commands and keeping docs/tests in sync. (file: /Users/pjordan/git/agent-extension-security/.codex/skills/agentsec-cli-change/SKILL.md)
- spec-evolution: Workflow for evolving AEM/APM schemas and corresponding manifest/policy behavior. (file: /Users/pjordan/git/agent-extension-security/.codex/skills/spec-evolution/SKILL.md)
- security-review: Security-focused review and implementation checklist for signature, policy, scan, and install paths. (file: /Users/pjordan/git/agent-extension-security/.codex/skills/security-review/SKILL.md)

### How to use skills
- Trigger rules: If a user names a listed skill (with `$SkillName` or plain text), use it for that turn.
- Multiple skills: Use the minimal set needed and state the sequence.
- Missing/blocked: If a skill path cannot be read, say so briefly and continue with the best fallback.
- Progressive disclosure: Open `SKILL.md` and load only files needed to complete the current request.
