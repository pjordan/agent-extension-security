# spec-evolution

Use this skill for changes to `spec/aem/*`, `spec/apm/*`, or manifest/policy semantics.

## Workflow
1. Define whether change is additive or breaking.
2. Update schema files in `spec/*`.
3. Update Go manifest/policy handling in:
   - `internal/manifest/*`
   - `internal/policy/*`
4. Ensure CLI validation behavior reflects schema intent.
5. Update docs:
   - `docs/permissions.md`
   - `docs/threat-model.md` (if threat assumptions changed)
   - `README.md` for contract-level changes
6. Run:
   - `make build`
   - `go test ./...`

## Compatibility Rules
- Prefer additive changes in `v0` unless explicitly breaking.
- If breaking, call it out clearly and propose migration path.
