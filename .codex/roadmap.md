# Engineering Roadmap (Suggested)

## 1) Manifest and Permissions Model Hardening
- Split AEM vs APM responsibilities cleanly.
- Add richer permission surface (secrets, OAuth scopes, high-risk action gates).
- Enforce permission expansion review on updates.

## 2) Signing and Provenance Productionization
- Integrate Sigstore keyless/cosign flow.
- Strengthen provenance format and validation checks.
- Document trust roots and verification policy.

## 3) Policy Engine Expansion
- Add configurable deny/allow policies for risk classes.
- Improve install-time decisions with explicit findings and remediation guidance.
- Add tests for policy regressions and edge cases.

## 4) Update and Distribution Security
- Define secure update metadata strategy (TUF-style roadmap).
- Add revocation/quarantine model and CLI support.

## 5) Developer Experience and Quality
- Add linter configuration (for example `golangci-lint`).
- Expand end-to-end tests for package/sign/verify/install flows.
- Add fixture artifacts for deterministic regression testing.
