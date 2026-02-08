# Policy File Schema

Policy files control which extension permissions are allowed or denied at install time.

## Structure

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

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mode` | `string` | Yes | `"enforce"` or `"warn"` |
| `deny` | `object` | No | Deny rules (if omitted, nothing is denied) |
| `deny.permissions` | `object` | No | Permission values to deny |

### `mode`

| Value | Behavior |
|-------|----------|
| `enforce` | Denied findings block installation. Install fails with a list of denied permissions. |
| `warn` | Denied findings are printed to stderr. Installation continues. |

Any other value is rejected.

### `deny.permissions`

The deny permissions object mirrors the AEM permissions structure. Only include the fields you want to deny:

| Deny field | Denies extensions that declare |
|-----------|-------------------------------|
| `deny.permissions.process.allow_shell: true` | `permissions.process.allow_shell: true` |
| `deny.permissions.process.allow_subprocess: true` | `permissions.process.allow_subprocess: true` |
| `deny.permissions.network.allow_ip_literals: true` | `permissions.network.allow_ip_literals: true` |

## Parsing behavior

Policy files are parsed with strict JSON decoding:

- **Unknown fields are rejected** — any field not in the schema causes an immediate error
- **Trailing content is rejected** — only a single JSON object is accepted
- **Mode validation** — `mode` must be exactly `"enforce"` or `"warn"`

## Provided templates

| Template | Mode | Denies | File |
|----------|------|--------|------|
| Permissive | `warn` | Nothing | `examples/policies/permissive.json` |
| Strict | `enforce` | `allow_shell`, `allow_ip_literals` | `examples/policies/strict.json` |
| Enterprise | `enforce` | `allow_shell`, `allow_ip_literals` | `examples/policies/enterprise.json` |
| Docs example | `enforce` | `allow_shell`, `allow_ip_literals` | `docs/policy.example.json` |

## Evaluation logic

At install time (`agentsec install --policy`):

1. Load the policy file (strict JSON parsing)
2. Load the extension's AEM manifest
3. For each deny rule, check the corresponding AEM permission value
4. Collect all findings (permission field + declared value)
5. If `mode=enforce` and findings exist → reject install
6. If `mode=warn` and findings exist → print findings, continue

## Next steps

- [Writing Policy](../consuming/policy.md) — practical guide to writing policies
- [AEM Manifest Schema](aem-schema.md) — the manifest fields that policy evaluates
- [Examples & Policies](../examples.md) — interactive policy demos
