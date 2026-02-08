# Example Policies

These policy files demonstrate different security postures for `agentsec install`.

## Policy files

### `permissive.json` — Warn-only, no deny rules

- **Mode:** `warn`
- **Use case:** Local development, experimentation, CI testing
- Allows all permissions but prints warnings for anything a stricter policy would deny.

### `strict.json` — Denies shell + IP literals, enforce mode

- **Mode:** `enforce`
- **Denies:** `process.allow_shell: true`, `network.allow_ip_literals: true`
- **Use case:** Staging environments, security-conscious teams
- Extensions that require shell access or IP literal connections will be **blocked**.

### `enterprise.json` — Full lockdown, enforce mode

- **Mode:** `enforce`
- **Denies:** `process.allow_shell: true`, `network.allow_ip_literals: true`
- **Use case:** Production environments, regulated industries
- The strictest posture — blocks shell access and IP literals.

## How policies work

At install time, `agentsec install` compares the extension's AEM manifest
permissions against the policy's deny rules:

1. If a permission is listed in the policy's `deny` section and the extension
   requests that permission, it's a **finding**.
2. In `enforce` mode, any findings cause installation to **fail**.
3. In `warn` mode, findings are printed as warnings but installation proceeds.

## Example: testing policy interaction

```bash
# Package the file-reader example
agentsec package examples/skills/file-reader --out file-reader.aext
agentsec sign file-reader.aext --key devkey.json --out file-reader.sig.json

# Try with permissive policy (succeeds with no warnings)
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/permissive.json \
  --dest ./test-install

# Try with strict policy (fails — file-reader uses allow_shell)
agentsec install file-reader.aext \
  --sig file-reader.sig.json --pub devkey.json \
  --aem examples/skills/file-reader/aem.json \
  --policy examples/policies/strict.json \
  --dest ./test-install
# => error: policy denied install: denied: process.allow_shell=true
```
