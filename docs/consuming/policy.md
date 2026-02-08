# Writing Policy

Policy files control which permissions are allowed or denied at install time. This page explains how to write and customize policy files for your environment.

## Policy file structure

A policy file is a JSON object with two fields:

```json
{
  "mode": "enforce",
  "deny": {
    "permissions": {
      "process": {
        "allow_shell": true
      },
      "network": {
        "allow_ip_literals": true
      }
    }
  }
}
```

### `mode`

| Value | Behavior |
|-------|----------|
| `enforce` | Denied findings block installation (fail closed) |
| `warn` | Denied findings are printed to stderr but installation continues |

### `deny`

The `deny` object lists permission values that should be rejected. At install time, the extension's AEM manifest is evaluated against these deny rules. Any matching permission triggers a finding.

## Common policy patterns

### Block shell access (recommended default)

```json
{
  "mode": "enforce",
  "deny": {
    "permissions": {
      "process": {
        "allow_shell": true
      }
    }
  }
}
```

Blocks any extension that declares `allow_shell=true`. This is the single most effective policy rule.

### Block shell and direct IP access

```json
{
  "mode": "enforce",
  "deny": {
    "permissions": {
      "process": {
        "allow_shell": true
      },
      "network": {
        "allow_ip_literals": true
      }
    }
  }
}
```

This is the `strict.json` policy template. It blocks shell access and prevents extensions from connecting to raw IP addresses (which bypass domain-based allowlists).

### Warn-only (development)

```json
{
  "mode": "warn",
  "deny": {
    "permissions": {
      "process": {
        "allow_shell": true
      }
    }
  }
}
```

Same deny rules, but violations are reported without blocking installation. Use during development and testing.

### Permissive (no restrictions)

```json
{
  "mode": "warn"
}
```

No deny rules. All permissions are allowed. The `permissive.json` template uses this.

## How policy interacts with AEM permissions

At install time, `agentsec install` does the following:

1. Loads the extension's AEM manifest
2. Loads your policy file
3. For each deny rule, checks whether the AEM declares a matching permission
4. If a match is found:
    - `mode=enforce` → install fails with a list of denied findings
    - `mode=warn` → findings are printed, install continues

**Example interaction:**

Extension AEM declares:
```json
{ "permissions": { "process": { "allow_shell": true } } }
```

Policy denies `allow_shell`:
```json
{ "deny": { "permissions": { "process": { "allow_shell": true } } } }
```

Result (enforce mode):
```
install: policy denied install:
 - denied: process.allow_shell=true
```

## Policy templates

Three templates are provided in `examples/policies/`:

| Template | Mode | Denies | Best for |
|----------|------|--------|----------|
| `permissive.json` | `warn` | Nothing | Local dev, experimentation |
| `strict.json` | `enforce` | `allow_shell`, `allow_ip_literals` | Staging, security-conscious teams |
| `enterprise.json` | `enforce` | `allow_shell`, `allow_ip_literals` | Production, regulated environments |

See [Examples & Policies](../examples.md) for interactive demos of each template.

## Strict parsing

Policy files are parsed with `DisallowUnknownFields`. Any typo or unknown field in the JSON will be rejected immediately:

```
load policy policy.json: json: unknown field "mdoe"
```

See [Troubleshooting](../troubleshooting.md) for common parsing errors.

## Next steps

- [Policy File Schema](../reference/policy-schema.md) — full schema reference
- [Examples & Policies](../examples.md) — see policy interaction with different extension types
- [Threat Model](threat-model.md) — understand what policies protect against
