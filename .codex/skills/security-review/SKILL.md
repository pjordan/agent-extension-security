# security-review

Use this skill to review or implement security-sensitive changes.

## Review Areas
- Signature generation and verification (`internal/crypto/*`).
- Artifact integrity and hashing (`internal/util/hash.go`, packaging/install paths).
- Policy enforcement correctness (`internal/policy/*`).
- Scan coverage and false-negative risks (`cmd/agentsec/scan.go`).
- Install-time trust assumptions (`cmd/agentsec/install.go`).

## Workflow
1. Identify threat class affected (tampering, exfiltration, privilege expansion, update compromise).
2. Trace control points in code and confirm enforcement is real, not advisory.
3. Check default behavior is secure-by-default.
4. Validate error handling fails closed for verification/policy checks.
5. Run:
   - `make build`
   - `go test ./...`
6. Summarize residual risks and follow-up hardening tasks.

## Output Expectations
- Findings first (severity ordered).
- Exact file references for each finding.
- Explicit note when behavior is scaffold-level only (non-production).
