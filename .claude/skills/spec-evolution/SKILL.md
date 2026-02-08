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
   - `docs/spec-aem.md` — AEM spec reference
   - `docs/spec-apm.md` — APM spec reference
   - `docs/permissions.md` — if permission model changed
   - `docs/threat-model.md` — if threat assumptions changed
   - `docs/examples.md` — example manifests and policies use real schemas
6. Run:
   ```bash
   make build && make test && make examples-test
   ```
   `make examples-test` validates example manifests against current schema behavior.

## Compatibility Rules

- Prefer additive changes in `v0` unless explicitly breaking.
- If breaking, call it out clearly and propose migration path.
- New optional fields are always safe; removing or renaming fields requires a migration note.
