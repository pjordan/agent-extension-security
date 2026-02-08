# CI/CD Pipeline

This guide covers automating the `agentsec` supply-chain pipeline in CI/CD. The same flow works for any extension format — skills, MCP servers, plugins.

## Pipeline overview

Every extension release should go through these stages:

| Stage | Command | Purpose |
|-------|---------|---------|
| Scaffold | `agentsec init` | Create project with manifest, key, policy |
| Package | `agentsec package` | Zip directory into `.aext` artifact |
| Validate | `agentsec manifest validate` | Verify manifest schema and constraints |
| Attest | `agentsec sbom`, `provenance`, `scan` | Generate supply-chain evidence |
| Sign | `agentsec sign` | Sign artifact digest |
| Verify | `agentsec verify` | Self-verify before publishing |
| Policy test | `agentsec install` | Test install with target policy |

## Example: GitHub Actions workflow

```yaml
name: Release Extension
on:
  push:
    tags: ["v*"]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"

      - name: Install agentsec
        run: go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest

      - name: Package
        run: agentsec package ./my-skill --out my-skill.aext

      - name: Validate manifest
        run: agentsec manifest validate ./my-skill/aem.json

      - name: Generate attestations
        run: |
          agentsec sbom my-skill.aext --out sbom.spdx.json
          agentsec provenance my-skill.aext \
            --source-repo "${{ github.server_url }}/${{ github.repository }}" \
            --source-rev "${{ github.sha }}" \
            --out provenance.json
          agentsec scan my-skill.aext --out scan.json

      - name: Sign
        run: agentsec sign my-skill.aext --key ./my-skill/devkey.json --out my-skill.sig.json

      - name: Verify
        run: agentsec verify my-skill.aext --sig my-skill.sig.json --pub ./my-skill/devkey.json

      - name: Policy install test
        run: |
          agentsec install my-skill.aext \
            --sig my-skill.sig.json \
            --pub ./my-skill/devkey.json \
            --aem ./my-skill/aem.json \
            --policy ./my-skill/policy.json \
            --dest ./test-install

      - name: Upload release artifacts
        uses: softprops/action-gh-release@v2
        with:
          files: |
            my-skill.aext
            my-skill/aem.json
            my-skill.sig.json
            sbom.spdx.json
            provenance.json
            scan.json
```

## Key management

!!! warning
    `agentsec keygen` generates local Ed25519 dev keys. These are suitable for development and testing but do not provide identity binding.

For CI pipelines:

- Store the signing key as a **repository secret** and write it to a temp file during the job
- Publish the corresponding public key in your repository or key server
- Consumers verify with `--pub` pointing to your published key

For production, the planned path is Sigstore/Cosign keyless signing — see [Contributing](../contributing.md) for the roadmap.

## Generic pipeline reference

For the full step-by-step pipeline with input variables and all commands, see the [CLI Reference](../reference/cli.md).

## Next steps

- [Publishing Checklist](publishing-checklist.md) — final checks before release
- [Examples & Policies](../examples.md) — policy templates and interaction demos
- [Consuming Extensions](../consuming/quickstart.md) — what consumers do with your artifacts
