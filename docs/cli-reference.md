# CLI Reference

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

## Commands

### `version`

Print the CLI version string.

### `init`

Scaffold a new extension project with a manifest, dev signing key, and policy file.

```bash
./bin/agentsec init ./my-skill --id com.example.my-skill --type skill
```

Creates the target directory (if needed) containing:

- `aem.json` — AEM manifest with least-privilege defaults
- `devkey.json` — Ed25519 dev signing keypair (mode 0600)
- `policy.json` — warn-mode policy for development

Flags:

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--id` | Yes | — | Extension identifier (e.g., `com.example.hello`) |
| `--type` | Yes | — | Extension type: `skill`, `mcp-server`, or `plugin` |
| `--version` | No | `0.1.0` | Semantic version |

Will not overwrite existing files. Edit the generated `aem.json` to declare the permissions your extension actually needs.

### `keygen`

Generate an Ed25519 dev keypair JSON.

```bash
./bin/agentsec keygen --out ./_demo/devkey.json
```

### `package`

Zip a directory into a `.aext` artifact.

```bash
./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext
```

### `manifest init`

Create an AEM JSON manifest with least-privilege defaults.

```bash
./bin/agentsec manifest init ./examples/skills/hello-world \
  --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json
```

### `manifest validate`

Validate manifest schema/version/type/semver checks.

```bash
./bin/agentsec manifest validate ./_demo/aem.json
```

### `sbom`

Emit a reference SBOM JSON with artifact digest metadata.

```bash
./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json
```

### `provenance`

Emit reference provenance JSON with source metadata and digest.

```bash
./bin/agentsec provenance ./_demo/hello-world.aext \
  --source-repo https://github.com/pjordan/agent-extension-security \
  --source-rev "$(git rev-parse HEAD)" \
  --out ./_demo/provenance.json
```

### `scan`

Run heuristic scanning on `SKILL.md`, `.sh`, and `.ps1` files in the artifact.

```bash
./bin/agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json
```

### `sign`

Sign an artifact digest with a local Ed25519 dev key.

```bash
./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json
```

### `verify`

Verify signature using a trusted public key (`--pub`) or insecure embedded-key mode.

```bash
./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --pub ./_demo/devkey.json
```

### `install`

Verify signature, evaluate policy against AEM, then extract artifact.

```bash
./bin/agentsec install ./_demo/hello-world.aext \
  --sig ./_demo/hello-world.sig.json \
  --pub ./_demo/devkey.json \
  --aem ./_demo/aem.json \
  --policy ./docs/policy.example.json \
  --dest ./_demo/install
```

**Dev mode:** Skip signature verification and use a permissive warn-only policy (for local development only):

```bash
./bin/agentsec install ./_demo/hello-world.aext \
  --dev \
  --aem ./_demo/aem.json \
  --dest ./_demo/install
```

!!! warning
    `--dev` mode skips signature verification entirely. Do not use in production.
