# CLI Reference

## Command overview

| Command | Persona | Description |
|---------|---------|-------------|
| `version` | Both | Print CLI version |
| `init` | Creator | Scaffold a new extension project |
| `package` | Creator | Zip a directory into a `.aext` artifact |
| `manifest init` | Creator | Create an AEM manifest with least-privilege defaults |
| `manifest validate` | Both | Validate manifest schema and constraints |
| `sbom` | Creator | Generate reference SBOM |
| `provenance` | Creator | Generate reference provenance record |
| `scan` | Both | Run heuristic scan on artifact contents |
| `keygen` | Creator | Generate an Ed25519 dev keypair |
| `sign` | Creator | Sign an artifact digest |
| `verify` | Consumer | Verify signature with a trusted key |
| `install` | Consumer | Verify, enforce policy, and extract |

## Global usage

```text
agentsec version
agentsec init <dir> --id <id> --type <skill|mcp-server|plugin> [--version <ver>]
agentsec keygen --out <file>
agentsec package <dir> --out <artifact.aext>
agentsec manifest init <dir> --id <id> --type <skill|mcp-server|plugin> --version <ver> --out <aem.json>
agentsec manifest validate <aem.json>
agentsec sbom <artifact.aext> --out <sbom.spdx.json>
agentsec provenance <artifact.aext> --source-repo <url> --source-rev <rev> --out <provenance.json>
agentsec scan <artifact.aext> --out <scan.json>
agentsec sign <artifact.aext> --key <devkey.json> --out <sig.json>
agentsec verify <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key)
agentsec install <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key) --aem <aem.json> --policy <policy.json> --dest <dir>
agentsec install <artifact.aext> --dev --aem <aem.json> --dest <dir>
```

Flags may appear before or after positional arguments.

---

## Creator commands

### `init`

Scaffold a new extension project with a manifest, dev signing key, and policy file.

```bash
agentsec init ./my-skill --id com.example.my-skill --type skill
```

Creates the target directory (if needed) containing:

- `aem.json` — AEM manifest with least-privilege defaults
- `devkey.json` — Ed25519 dev signing keypair (mode 0600)
- `policy.json` — warn-mode policy for development

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--id` | Yes | — | Extension identifier (e.g., `com.example.hello`) |
| `--type` | Yes | — | Extension type: `skill`, `mcp-server`, or `plugin` |
| `--version` | No | `0.1.0` | Semantic version |

Will not overwrite existing files. Edit the generated `aem.json` to declare the permissions your extension actually needs.

### `keygen`

Generate an Ed25519 dev keypair JSON.

```bash
agentsec keygen --out devkey.json
```

### `package`

Zip a directory into a `.aext` artifact. Rejects symlinks and non-regular files.

```bash
agentsec package ./my-skill --out my-skill.aext
```

### `manifest init`

Create an AEM JSON manifest with least-privilege defaults.

```bash
agentsec manifest init ./my-skill \
  --id com.example.my-skill --type skill --version 0.1.0 --out aem.json
```

### `manifest validate`

Validate manifest schema, type, and version constraints.

```bash
agentsec manifest validate aem.json
```

### `sbom`

Emit a reference SBOM JSON with artifact digest metadata.

```bash
agentsec sbom my-skill.aext --out sbom.spdx.json
```

### `provenance`

Emit reference provenance JSON with source metadata and digest.

```bash
agentsec provenance my-skill.aext \
  --source-repo https://github.com/your-org/your-repo \
  --source-rev "$(git rev-parse HEAD)" \
  --out provenance.json
```

### `scan`

Run heuristic scanning on `SKILL.md`, `.sh`, and `.ps1` files in the artifact.

```bash
agentsec scan my-skill.aext --out scan.json
```

### `sign`

Sign an artifact digest with a local Ed25519 dev key.

```bash
agentsec sign my-skill.aext --key devkey.json --out my-skill.sig.json
```

---

## Consumer commands

### `verify`

Verify signature using a trusted public key (`--pub`) or insecure embedded-key mode.

```bash
agentsec verify my-skill.aext --sig my-skill.sig.json --pub publisher-key.json
```

| Flag | Description |
|------|-------------|
| `--pub <file>` | Trusted public key file (recommended) |
| `--allow-embedded-key` | Use key from signature file (insecure, dev-only) |

### `install`

Verify signature, evaluate policy against AEM, then extract artifact.

```bash
agentsec install my-skill.aext \
  --sig my-skill.sig.json \
  --pub publisher-key.json \
  --aem aem.json \
  --policy policy.json \
  --dest ./installed/my-skill
```

| Flag | Required | Description |
|------|----------|-------------|
| `--sig` | Yes (unless `--dev`) | Signature file |
| `--pub` | Yes (unless `--dev` or `--allow-embedded-key`) | Trusted public key |
| `--aem` | Yes | AEM manifest file |
| `--policy` | Yes (unless `--dev`) | Policy file |
| `--dest` | Yes | Installation directory |
| `--dev` | No | Dev mode: skip signature, permissive policy |

**Dev mode:** Skip signature verification and use a permissive warn-only policy (for local development only):

```bash
agentsec install my-skill.aext --dev --aem aem.json --dest ./installed/my-skill
```

!!! warning
    `--dev` mode skips signature verification entirely. Do not use in production.

---

## Common commands

### `version`

Print the CLI version string.

```bash
agentsec version
```
