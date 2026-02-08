# Permissions (APM) - scaffold

This repo includes a small permission model in the AEM (JSON).

In a production system, permissions should:
- be explicit
- be minimal by default
- be enforced (not advisory)
- be diffed between versions, and expansions should require approval

## Current install-time behavior

- `agentsec install` requires an AEM file (`--aem`) and a policy file (`--policy`).
- Policy `mode` defaults to `enforce`; denied findings block install.
- Policy `mode=warn` reports findings but allows install.
- Unknown JSON fields in manifest/policy files are rejected.

## Current fields (AEM (JSON))
- files.read: list of path globs
- files.write: list of path globs
- network.domains: allowlist
- network.allow_ip_literals: boolean
- process.allow_shell: boolean
- process.allow_subprocess: boolean

## Secure defaults

`agentsec manifest init` now emits least-privilege defaults:
- `files.read=[]`
- `files.write=[]`
- `network.domains=[]`
- `network.allow_ip_literals=false`
- `process.allow_shell=false`
- `process.allow_subprocess=false`

## Next steps
- Split AEM (JSON) vs APM schemas
- Add secret requirements and OAuth scopes
- Add human-approval gates for high-risk actions
