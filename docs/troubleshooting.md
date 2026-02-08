# Troubleshooting

Common errors you may encounter when using `agentsec`, and how to resolve them.

## "policy denied install"

**Error:**
```
install: policy denied install:
 - denied: process.allow_shell=true
```

**What happened:** The extension's AEM manifest declares permissions that your policy file blocks. The install was rejected because the policy mode is `enforce`.

**How to fix:**

1. **Inspect which permission triggered the denial.** The error lists every denied finding. Common ones:
   - `denied: process.allow_shell=true` — the extension requests shell access
   - `denied: process.allow_subprocess=true` — the extension requests subprocess spawning
   - `denied: network.allow_ip_literals=true` — the extension requests direct IP network access

2. **If the permission is expected,** update your policy file to allow it. For example, to permit shell access, remove `allow_shell` from the deny list in your policy JSON.

3. **If you want to see warnings instead of hard failures,** switch the policy mode to `warn`:
   ```json
   {
     "mode": "warn",
     "deny": { ... }
   }
   ```
   In `warn` mode, policy violations are printed to stderr but do not block installation.

4. **For development/testing,** use `--dev` mode which applies a permissive warn-only policy by default:
   ```bash
   agentsec install hello-world.aext --dev --aem aem.json --dest ./install
   ```

## "signature verification failed"

**Error:**
```
verify: signature verification failed
```

**What happened:** The Ed25519 signature does not match the artifact digest using the provided public key.

**Common causes:**

| Cause | How to check |
|-------|-------------|
| **Key mismatch** | You are verifying with a different key than the one used to sign. Compare the `public_key` field in the `.sig.json` file with your `--pub` key file. |
| **File was modified** | The artifact file changed after signing (re-download or re-package). Re-sign the artifact with `agentsec sign`. |
| **Wrong file** | You passed the wrong artifact or signature file. Double-check your file paths. |

**How to re-sign:**
```bash
agentsec sign my-extension.aext --key devkey.json --out my-extension.sig.json
agentsec verify my-extension.aext --sig my-extension.sig.json --pub devkey.json
```

## "refusing to package symlink"

**Error:**
```
refusing to package symlink: config/settings.yaml
```

**What happened:** The directory you are packaging contains a symbolic link. `agentsec` rejects symlinks during both packaging and extraction to prevent [ZipSlip](https://security.snyk.io/research/zip-slip-vulnerability) and path-traversal attacks.

**How to fix:**

1. Replace the symlink with the actual file (copy instead of link).
2. Restructure your extension directory so symlinks are not needed.
3. If the symlink points outside the extension directory, this is the exact attack vector `agentsec` is designed to prevent.

**Related error:** `refusing to package non-regular file: ...` — same treatment applies to device files, sockets, and other non-regular files.

## Archive size and entry limits

**Errors:**
```
archive contains too many entries: 15000
archive entry too large: data/model.bin
archive exceeds total uncompressed size limit
suspicious compression ratio: data/payload.dat
```

**What happened:** The artifact exceeds one of the hardened archive limits:

| Limit | Value | Purpose |
|-------|-------|---------|
| Max entries | 10,000 | Prevents zip bomb with many small files |
| Max total uncompressed size | 512 MiB | Prevents disk exhaustion |
| Max single entry size | 64 MiB | Prevents individual file bombs |
| Max compression ratio | 200:1 | Detects zip bombs and decompression attacks |

**How to fix:**

- **Too many entries:** Reduce the number of files in your extension. Consider whether all files are necessary, or consolidate small files.
- **Entry/archive too large:** Remove large binary assets from the extension. Extensions should contain code and configuration, not large data files.
- **Suspicious compression ratio:** This usually indicates a crafted archive designed to expand to a very large size. If you have legitimately highly compressible data, consider restructuring the extension.

## "unknown field" errors

**Error:**
```
load manifest example/aem.json: json: unknown field "permisions"
load policy policy.json: json: unknown field "mdoe"
```

**What happened:** Both manifest and policy JSON files are parsed with strict mode (`DisallowUnknownFields`). Any field name not in the schema is rejected immediately.

**How to fix:**

1. Check for **typos** in field names. Common mistakes:
   - `permisions` → `permissions`
   - `mdoe` → `mode`
   - `allow_ip_literal` → `allow_ip_literals`
   - `allow_subprocesses` → `allow_subprocess`

2. Check for **extra fields** you may have added. Only documented fields are accepted:
   - Manifest (AEM): `schema`, `id`, `type`, `version`, `source_repo`, `source_rev`, `permissions`
   - Policy: `mode`, `deny`

3. Validate a known-good file to confirm your tooling works:
   ```bash
   agentsec manifest validate examples/skills/file-reader/aem.json
   ```

## "no trusted key provided"

**Error:**
```
verify: no trusted key provided: pass --pub <pubkey.json> (or --allow-embedded-key for insecure/dev mode)
```

**What happened:** You ran `verify` or `install` without specifying how to obtain the trusted public key.

**Options:**

| Flag | Security | Use case |
|------|----------|----------|
| `--pub <pubkey.json>` | **Recommended** | You have the signer's public key file. The signature is verified against this trusted key. |
| `--allow-embedded-key` | **Insecure** | Uses the public key embedded in the signature file itself. This only proves the artifact was signed by *someone* — it does not verify *who*. |

**Best practice:** Always use `--pub` with a key you obtained out-of-band (e.g., from the extension author's repository or a key server). Reserve `--allow-embedded-key` for local development only.

**For development workflows,** `--dev` mode on `install` skips signature verification entirely:
```bash
agentsec install my-ext.aext --dev --aem aem.json --dest ./install
```

## Common gotchas

### Forgetting to edit `aem.json` after scaffolding

Both `agentsec init` and `agentsec manifest init` generate manifests with **least-privilege defaults** (no file, network, or process permissions). If your extension needs shell access, network access, or file access, you must edit the `permissions` section before packaging:

```json
{
  "permissions": {
    "process": {
      "allow_shell": true
    },
    "network": {
      "domains": ["api.example.com"]
    }
  }
}
```

### Using dev keys in production

The `keygen` command generates local Ed25519 keypairs for development. These keys:

- Have no identity binding (no certificate, no email, no org)
- Are stored as plain JSON files
- Provide integrity verification but not authenticity

For production use, plan to migrate to Sigstore/Cosign keyless signing (see the [Sigstore roadmap](contributing.md#sigstore-integration)).

### Manifest says "single JSON object"

**Error:**
```
load manifest aem.json: manifest must contain a single JSON object
```

The manifest file contains trailing data after the JSON object (e.g., a second JSON object, trailing text, or extra whitespace with content). Ensure the file contains exactly one JSON object and nothing else.

### Policy mode must be `enforce` or `warn`

**Error:**
```
load policy policy.json: mode must be enforce|warn
```

The `mode` field in your policy file must be exactly `"enforce"` or `"warn"`. Check for typos, capitalization, or missing quotes.
