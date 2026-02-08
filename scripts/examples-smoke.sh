#!/bin/sh
# examples-smoke.sh — Smoke test all examples: package, validate manifests, scan.
# Exit on first failure.

set -e

AGENTSEC="${AGENTSEC:-./agentsec}"
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "=== Building agentsec ==="
go build -o "$AGENTSEC" ./cmd/agentsec/

echo ""
echo "=== Validating example manifests ==="
for manifest in examples/skills/file-reader/aem.json \
                examples/skills/web-fetcher/aem.json \
                examples/mcp/echo-server/aem.json; do
  echo "  validate: $manifest"
  "$AGENTSEC" manifest validate "$manifest"
done

echo ""
echo "=== Packaging examples ==="
for dir in examples/skills/hello-world \
           examples/skills/file-reader \
           examples/skills/web-fetcher \
           examples/mcp/echo-server; do
  name=$(basename "$dir")
  out="$TMPDIR/${name}.aext"
  echo "  package: $dir -> $out"
  "$AGENTSEC" package "$dir" --out "$out"
done

echo ""
echo "=== Scanning examples ==="
for dir in examples/skills/hello-world \
           examples/skills/file-reader \
           examples/skills/web-fetcher \
           examples/mcp/echo-server; do
  name=$(basename "$dir")
  art="$TMPDIR/${name}.aext"
  report="$TMPDIR/${name}-scan.json"
  echo "  scan: $art -> $report"
  "$AGENTSEC" scan "$art" --out "$report"
done

echo ""
echo "=== All example smoke tests passed ==="
