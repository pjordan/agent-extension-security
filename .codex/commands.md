# Command Recipes

Use these commands as the default execution playbook for this repository.

## Build and Test
- Build CLI:
  - `make build`
- Run test suite:
  - `go test ./...`
- Format Go files:
  - `make fmt`

## Local Demo Flow
- Build + run version:
  - `make build && ./bin/agentsec version`
- Package example skill:
  - `./bin/agentsec package ./examples/skills/hello-world --out ./_demo/hello-world.aext`
- Generate manifest:
  - `./bin/agentsec manifest init ./examples/skills/hello-world --id com.example.hello-world --type skill --version 0.1.0 --out ./_demo/aem.json`
- Validate manifest:
  - `./bin/agentsec manifest validate ./_demo/aem.json`
- Generate SBOM:
  - `./bin/agentsec sbom ./_demo/hello-world.aext --out ./_demo/sbom.spdx.json`
- Generate provenance:
  - `./bin/agentsec provenance ./_demo/hello-world.aext --source-repo https://github.com/pjordan/agent-extension-security --source-rev "$(git rev-parse HEAD)" --out ./_demo/provenance.json`
- Scan artifact:
  - `./bin/agentsec scan ./_demo/hello-world.aext --out ./_demo/scan.json`
- Generate key + sign + verify:
  - `./bin/agentsec keygen --out ./_demo/devkey.json`
  - `./bin/agentsec sign ./_demo/hello-world.aext --key ./_demo/devkey.json --out ./_demo/hello-world.sig.json`
  - `./bin/agentsec verify ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json`
- Install (simulated):
  - `./bin/agentsec install ./_demo/hello-world.aext --sig ./_demo/hello-world.sig.json --dest ./_demo/install`

## Change Validation Matrix
- CLI-only change:
  - `make build && go test ./...`
- Manifest/schema/policy change:
  - `make build && go test ./...`
  - Re-run Local Demo Flow manifest + scan + install steps.
- Crypto/signature change:
  - `make build && go test ./...`
  - Re-run keygen/sign/verify/install sequence.

## Cleanup
- Remove generated artifacts:
  - `make clean`

## Pre-PR Check
- Recommended final command set:
  - `make build`
  - `go test ./...`
  - `git status --short`
