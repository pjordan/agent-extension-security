# Permissions (APM) - scaffold

This repo includes a small permission model in the AEM (JSON).

In a production system, permissions should:
- be explicit
- be minimal by default
- be enforced (not advisory)
- be diffed between versions, and expansions should require approval

## Current fields (AEM (JSON))
- files.read: list of path globs
- files.write: list of path globs
- network.domains: allowlist
- network.allow_ip_literals: boolean
- process.allow_shell: boolean
- process.allow_subprocess: boolean

## Next steps
- Split AEM (JSON) vs APM schemas
- Add secret requirements and OAuth scopes
- Add human-approval gates for high-risk actions

