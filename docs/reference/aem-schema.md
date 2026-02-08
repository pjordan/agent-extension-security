# AEM Manifest Schema

The Agent Extension Manifest (AEM) declares an extension's identity, type, version, and permission requirements.

## Schema file

The canonical JSON schema is at `spec/aem/v0/aem.schema.json` in the repository.

## Core fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schema` | `string` | Yes | Must be `"aessf.dev/aem/v0"` |
| `id` | `string` | Yes | Extension identifier (e.g., `com.example.hello-world`) |
| `type` | `string` | Yes | Extension type: `skill`, `mcp-server`, or `plugin` |
| `version` | `string` | Yes | Semantic version (e.g., `1.0.0`) |
| `source_repo` | `string` | No | Source repository URL |
| `source_rev` | `string` | No | Source revision (commit hash) |
| `permissions` | `object` | Yes | Permission declarations (see below) |

## Permission fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `permissions.files.read` | `string[]` | `[]` | Path globs the extension may read |
| `permissions.files.write` | `string[]` | `[]` | Path globs the extension may write |
| `permissions.network.domains` | `string[]` | `[]` | DNS names the extension may contact |
| `permissions.network.allow_ip_literals` | `bool` | `false` | Whether raw IP addresses are allowed |
| `permissions.process.allow_shell` | `bool` | `false` | Whether shell execution is allowed |
| `permissions.process.allow_subprocess` | `bool` | `false` | Whether subprocess spawning is allowed |

## Runtime validation

In addition to JSON decoding, `agentsec manifest validate` enforces:

- `schema` must equal `"aessf.dev/aem/v0"`
- `id` is required and non-empty
- `type` must be one of `skill`, `mcp-server`, `plugin`
- `version` must be valid semver format (e.g., `1.2.3`, with optional prerelease/build suffix)
- Unknown JSON fields are rejected (`DisallowUnknownFields`)
- Trailing JSON values after the root object are rejected

## Example manifest

```json
{
  "schema": "aessf.dev/aem/v0",
  "id": "com.example.hello-world",
  "type": "skill",
  "version": "0.1.0",
  "permissions": {
    "files": {
      "read": [],
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

??? note "Full JSON schema"
    The JSON schema is in `spec/aem/v0/aem.schema.json`. Key constraints:

    - `schema` is `const: "aessf.dev/aem/v0"`
    - `type` is `enum: ["skill", "mcp-server", "plugin"]`
    - `version` is validated as semver at runtime
    - `additionalProperties: false` at all levels

## Future: APM split

The current AEM includes both identity and permission fields. A planned evolution is to split the permission model into a standalone APM (Agent Permission Manifest):

- Separate permission declarations from extension identity
- Add policy-relevant dimensions (OAuth scopes, secret requirements)
- Add compatibility and migration guidance between schema versions

The APM schema placeholder is at `spec/apm/v0/apm.schema.json`.

## Next steps

- [Declaring Permissions](../creating/permissions.md) — guidance on choosing values
- [Policy File Schema](policy-schema.md) — how policies evaluate AEM permissions
- [CLI Reference](cli.md) — `manifest init` and `manifest validate` commands
