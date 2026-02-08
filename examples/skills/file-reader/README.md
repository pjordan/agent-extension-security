# file-reader example

A skill that reads configuration files and summarizes them. Demonstrates:

- **File read permissions** (`files.read: ["~/.config/*"]`)
- **Shell access** (`process.allow_shell: true`)
- A helper shell script (`read-config.sh`)

## Package, sign, and install

```bash
# 1. Generate a signing key (if you don't have one)
agentsec keygen --out devkey.json

# 2. Package
agentsec package examples/skills/file-reader --out file-reader.aext

# 3. Sign
agentsec sign file-reader.aext --key devkey.json --out file-reader.sig.json

# 4. Verify
agentsec verify file-reader.aext --sig file-reader.sig.json --pub devkey.json

# 5. Install (full secure flow)
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/permissive.json \
  --dest ./installed

# Or install in dev mode (skip signature verification)
agentsec install file-reader.aext --dev \
  --aem examples/skills/file-reader/aem.json \
  --dest ./installed
```

## Policy interaction

- A **strict policy** that denies `allow_shell: true` will **block** this
  extension at install time.
- A **permissive policy** (warn mode) will allow it with a warning.

Try it:

```bash
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/strict.json \
  --dest ./installed
# => error: policy denied install
```
