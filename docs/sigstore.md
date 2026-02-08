# Sigstore integration plan (scaffold)

The current `agentsec sign/verify` commands use **local dev keys** (ed25519) to keep the scaffold self-contained.

For real publishing, the intended path is Sigstore + Cosign:

- publish extensions as OCI artifacts (or blobs) in an OCI registry
- sign with `cosign sign` (keyless in CI)
- verify with identity constraints:
  - `--certificate-identity` (e.g., workflow identity URI)
  - `--certificate-oidc-issuer` (e.g., GitHub Actions token issuer)
- require transparency log inclusion

Planned CLI additions:
- `agentsec sign --keyless` (shell out to cosign)
- `agentsec verify --sigstore-policy <policy.yaml>`

See:
- https://docs.sigstore.dev/
- https://docs.sigstore.dev/cosign/verifying/verify/

