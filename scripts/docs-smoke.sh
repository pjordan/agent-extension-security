#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="${ROOT}/bin/agentsec"
WORKDIR="$(mktemp -d)"
trap 'rm -rf "${WORKDIR}"' EXIT

mkdir -p "${WORKDIR}/demo"

"${BIN}" version >/dev/null

"${BIN}" package "${ROOT}/examples/skills/hello-world" --out "${WORKDIR}/demo/hello-world.aext"

"${BIN}" manifest init "${ROOT}/examples/skills/hello-world" \
  --id com.example.hello-world --type skill --version 0.1.0 --out "${WORKDIR}/demo/aem.json"
"${BIN}" manifest validate "${WORKDIR}/demo/aem.json"

"${BIN}" sbom "${WORKDIR}/demo/hello-world.aext" --out "${WORKDIR}/demo/sbom.spdx.json"
"${BIN}" provenance "${WORKDIR}/demo/hello-world.aext" \
  --source-repo "https://github.com/pjordan/agent-extension-security" \
  --source-rev "HEAD" \
  --out "${WORKDIR}/demo/provenance.json"
"${BIN}" scan "${WORKDIR}/demo/hello-world.aext" --out "${WORKDIR}/demo/scan.json"

"${BIN}" keygen --out "${WORKDIR}/demo/devkey.json"
"${BIN}" sign "${WORKDIR}/demo/hello-world.aext" --key "${WORKDIR}/demo/devkey.json" --out "${WORKDIR}/demo/hello-world.sig.json"
"${BIN}" verify "${WORKDIR}/demo/hello-world.aext" --sig "${WORKDIR}/demo/hello-world.sig.json" --pub "${WORKDIR}/demo/devkey.json"

"${BIN}" install "${WORKDIR}/demo/hello-world.aext" \
  --sig "${WORKDIR}/demo/hello-world.sig.json" \
  --pub "${WORKDIR}/demo/devkey.json" \
  --aem "${WORKDIR}/demo/aem.json" \
  --policy "${ROOT}/docs/policy.example.json" \
  --dest "${WORKDIR}/demo/install"

test -f "${WORKDIR}/demo/install/hello-world.aext/SKILL.md"
test -f "${WORKDIR}/demo/sbom.spdx.json"
test -f "${WORKDIR}/demo/provenance.json"
test -f "${WORKDIR}/demo/scan.json"

echo "docs smoke test passed"
