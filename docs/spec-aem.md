# AEM Schema

Current schema file:

- `spec/aem/v0/aem.schema.json`

## Runtime validation contract

In addition to JSON decoding, runtime validation enforces:

- `schema == "aessf.dev/aem/v0"`
- `id` is required
- `type` is one of `skill | mcp-server | plugin`
- `version` must be semantic version format (`1.2.3`, optional prerelease/build suffix)

## Core fields

- `schema`
- `id`
- `type`
- `version`
- `source_repo` (optional)
- `source_rev` (optional)
- `permissions`

## Permission fields

- `permissions.files.read[]`
- `permissions.files.write[]`
- `permissions.network.domains[]`
- `permissions.network.allow_ip_literals`
- `permissions.process.allow_shell`
- `permissions.process.allow_subprocess`

See [Permissions & Policy](permissions.md) for install-time behavior.
