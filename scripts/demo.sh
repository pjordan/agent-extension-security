#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="${ROOT}/bin/agentsec"

mkdir -p "${ROOT}/_demo"

echo "[1/8] Build"
(cd "${ROOT}" && make build)

echo "[2/8] Package"
"${BIN}" package "${ROOT}/examples/skills/hello-world" --out "${ROOT}/_demo/hello-world.aext"

echo "[3/8] Manifest"
"${BIN}" manifest init "${ROOT}/examples/skills/hello-world" --id com.example.hello-world --type skill --version 0.1.0 --out "${ROOT}/_demo/aem.json"
"${BIN}" manifest validate "${ROOT}/_demo/aem.json"

echo "[4/8] SBOM"
"${BIN}" sbom "${ROOT}/_demo/hello-world.aext" --out "${ROOT}/_demo/sbom.spdx.json"

echo "[5/8] Provenance"
"${BIN}" provenance "${ROOT}/_demo/hello-world.aext" --source-repo "local-demo" --source-rev "HEAD" --out "${ROOT}/_demo/provenance.json"

echo "[6/8] Scan"
"${BIN}" scan "${ROOT}/_demo/hello-world.aext" --out "${ROOT}/_demo/scan.json"

echo "[7/8] Sign + verify"
"${BIN}" keygen --out "${ROOT}/_demo/devkey.json"
"${BIN}" sign "${ROOT}/_demo/hello-world.aext" --key "${ROOT}/_demo/devkey.json" --out "${ROOT}/_demo/hello-world.sig.json"
"${BIN}" verify "${ROOT}/_demo/hello-world.aext" --sig "${ROOT}/_demo/hello-world.sig.json"

echo "[8/8] Install"
"${BIN}" install "${ROOT}/_demo/hello-world.aext" --sig "${ROOT}/_demo/hello-world.sig.json" --dest "${ROOT}/_demo/install"

echo "Demo complete. See ${ROOT}/_demo/"
