# Declaring Permissions

Every extension declares its permission requirements in an AEM manifest. This page explains how to choose the right values and why least-privilege defaults matter.

## Permission fields

| Field | Type | Default | What it controls |
|-------|------|---------|------------------|
| `permissions.files.read` | `string[]` | `[]` | Path globs the extension may read |
| `permissions.files.write` | `string[]` | `[]` | Path globs the extension may write |
| `permissions.network.domains` | `string[]` | `[]` | DNS names the extension may contact |
| `permissions.network.allow_ip_literals` | `bool` | `false` | Whether raw IP addresses are allowed |
| `permissions.process.allow_shell` | `bool` | `false` | Whether shell execution is allowed |
| `permissions.process.allow_subprocess` | `bool` | `false` | Whether subprocess spawning is allowed |

## Secure defaults

Both `agentsec init` and `agentsec manifest init` generate manifests with **zero permissions**:

```json
{
  "permissions": {
    "files": { "read": [], "write": [] },
    "network": { "domains": [], "allow_ip_literals": false },
    "process": { "allow_shell": false, "allow_subprocess": false }
  }
}
```

Start here and add only what your extension genuinely needs. Every permission you add is a permission that consumers must evaluate and accept.

## Choosing the right values

### Read-only skill (documentation, templates)

No special permissions needed. The default manifest works:

```json
{
  "permissions": {
    "files": { "read": ["*.md", "*.txt"], "write": [] },
    "network": { "domains": [], "allow_ip_literals": false },
    "process": { "allow_shell": false, "allow_subprocess": false }
  }
}
```

### API-connected skill (fetches external data)

Declare specific domains. Never use `allow_ip_literals` unless you have a concrete reason:

```json
{
  "permissions": {
    "files": { "read": ["*.json"], "write": [] },
    "network": { "domains": ["api.example.com"], "allow_ip_literals": false },
    "process": { "allow_shell": false, "allow_subprocess": false }
  }
}
```

### Build tool skill (runs commands)

If your skill needs to execute build commands, declare shell and/or subprocess access:

```json
{
  "permissions": {
    "files": { "read": ["*"], "write": ["build/**", "dist/**"] },
    "network": { "domains": [], "allow_ip_literals": false },
    "process": { "allow_shell": true, "allow_subprocess": true }
  }
}
```

!!! warning
    `allow_shell=true` is the single most scrutinized permission. Consumers running strict policies will deny it. Only declare it if your extension genuinely needs to execute shell commands.

## How permissions interact with policy

At install time, the consumer's policy file defines which permissions are denied. If your manifest declares a permission that the consumer's policy denies, the install fails (in `enforce` mode) or warns (in `warn` mode).

Example: if a consumer's policy denies `allow_shell`, then any extension declaring `allow_shell=true` will be rejected.

See [Writing Policy](../consuming/policy.md) for the consumer's perspective and [Examples & Policies](../examples.md) for worked interaction examples.

## Next steps

- [AEM Manifest Schema](../reference/aem-schema.md) — full schema reference
- [Examples & Policies](../examples.md) — permission gradient and policy interaction demos
- [CI/CD Pipeline](ci-cd.md) — automate manifest validation in CI
