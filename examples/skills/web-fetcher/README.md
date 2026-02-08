# web-fetcher example

A skill that fetches data from approved domains via `curl`. Demonstrates:

- **Network permissions** (`network.domains: ["api.example.com"]`)
- **Shell and subprocess access** (`process.allow_shell: true`, `process.allow_subprocess: true`)
- **Scanner interaction** — the `fetch.sh` script will be flagged by `agentsec scan`

## Package, sign, and install

```bash
# 1. Package
agentsec package examples/skills/web-fetcher --out web-fetcher.aext

# 2. Scan (see the heuristic scanner in action)
agentsec scan web-fetcher.aext --out scan-report.json
cat scan-report.json
# You'll see findings about shell/subprocess usage

# 3. Sign
agentsec sign web-fetcher.aext --key devkey.json --out web-fetcher.sig.json

# 4. Install
agentsec install web-fetcher.aext \
  --sig web-fetcher.sig.json --pub devkey.json \
  --aem examples/skills/web-fetcher/aem.json \
  --policy examples/policies/permissive.json \
  --dest ./installed
```

## Policy interaction

- A **strict policy** that denies `allow_shell: true` will **block** this
  extension.
- An **enterprise policy** that also denies subprocess access will block it
  for two reasons.
- A **permissive policy** will allow it with warnings.

## Scanner output

Running `agentsec scan` on this artifact will produce findings because
`fetch.sh` contains shell commands. This is intentional — the scanner's
purpose is to surface risk for human review, not to block legitimate tools.
