# security-review

Use this skill to review or implement security-sensitive changes.

## Review Areas

- Signature generation and verification (`internal/crypto/*`).
- Artifact integrity and hashing (`internal/util/hash.go`, packaging/install paths).
- Policy enforcement correctness (`internal/policy/*`).
- Scan coverage and false-negative risks (`cmd/agentsec/scan.go`).
- Install-time trust assumptions (`cmd/agentsec/install.go`).
- `--dev` mode install behavior — ensure it is clearly advisory-only and never bypasses policy silently.

## Workflow

1. Identify threat class affected (tampering, exfiltration, privilege expansion, update compromise).
2. Trace control points in code and confirm enforcement is real, not advisory.
3. Check default behavior is secure-by-default.
4. Validate error handling fails closed for verification/policy checks.
5. Run:
   ```bash
   make build && make test && make cover
   ```
   Coverage must stay >= 80% on `./internal/...`.
6. Review against `docs/security-hardening.md` and `docs/threat-model.md` for consistency.
7. Summarize residual risks and follow-up hardening tasks.

## Output Expectations

- Findings first (severity ordered).
- Exact file references for each finding.
- Explicit note when behavior is scaffold-level only (non-production).
- Coverage delta if tests were added or removed.
