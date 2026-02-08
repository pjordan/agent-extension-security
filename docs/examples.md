# Examples & Policy Templates

The `examples/` directory contains sample extensions and policy files that
demonstrate the `agentsec` security model at different permission levels.

## Example extensions

### Skills

| Example | Permissions | Key concept |
|---------|-------------|-------------|
| [hello-world](https://github.com/pjordan/agent-extension-security/tree/main/examples/skills/hello-world) | None | Minimal packaging — no special permissions needed |
| [file-reader](https://github.com/pjordan/agent-extension-security/tree/main/examples/skills/file-reader) | `files.read`, `process.allow_shell` | File access + shell — blocked by strict policy |
| [web-fetcher](https://github.com/pjordan/agent-extension-security/tree/main/examples/skills/web-fetcher) | `network.domains`, `process.allow_shell`, `process.allow_subprocess` | Network access + scanner findings |

### MCP servers

| Example | Permissions | Key concept |
|---------|-------------|-------------|
| [echo-server](https://github.com/pjordan/agent-extension-security/tree/main/examples/mcp/echo-server) | `process.allow_subprocess` | MCP server type — subprocess without shell |

## Permission gradient

The examples form a progression from zero permissions to increasingly
sensitive capabilities:

1. **hello-world** — No permissions. Passes every policy.
2. **echo-server** — Subprocess only. Passes strict and permissive policies.
3. **file-reader** — File read + shell. Blocked by strict policy (`allow_shell` denied).
4. **web-fetcher** — Network + shell + subprocess. Blocked by strict and enterprise policies. Intentionally triggers scanner findings.

## Policy templates

Three policy files in `examples/policies/` demonstrate different security postures:

| Policy | Mode | Denies | Use case |
|--------|------|--------|----------|
| [`permissive.json`](https://github.com/pjordan/agent-extension-security/tree/main/examples/policies/permissive.json) | `warn` | Nothing | Local dev, experimentation |
| [`strict.json`](https://github.com/pjordan/agent-extension-security/tree/main/examples/policies/strict.json) | `enforce` | `allow_shell`, `allow_ip_literals` | Staging, security-conscious teams |
| [`enterprise.json`](https://github.com/pjordan/agent-extension-security/tree/main/examples/policies/enterprise.json) | `enforce` | `allow_shell`, `allow_ip_literals` | Production, regulated environments |

The quickstart also uses [`docs/policy.example.json`](https://github.com/pjordan/agent-extension-security/blob/main/docs/policy.example.json),
an enforce-mode policy that denies shell and IP literal access — functionally
equivalent to `strict.json`.

## Try it: policy interaction

```bash
# Package file-reader
agentsec package examples/skills/file-reader --out file-reader.aext
agentsec keygen --out devkey.json
agentsec sign file-reader.aext --key devkey.json --out file-reader.sig.json

# Passes with permissive policy (warn mode)
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/permissive.json \
  --dest ./test-install

# Fails with strict policy (file-reader uses allow_shell)
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/strict.json \
  --dest ./test-install
# => error: install: policy denied install:
#     - denied: process.allow_shell=true
```

See each example's README for full details on permissions and usage.
