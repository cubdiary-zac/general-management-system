# M5 Milestone: Configurable Industry Project Templates

- Date: 2026-03-11
- Repository: `cubdiary-zac/general-management-system`
- Theme: template engine foundation

## Goal
Deliver a backend-first template engine that lets teams define reusable industry/project pipelines (industry -> project template -> stage -> form -> field) with ordering and publish-state controls.

## Milestones

### M5.1 Foundation (this delivery)
- Deliverables:
  - Template models: `IndustryTemplate`, `ProjectTemplate`, `StageTemplate`, `FormTemplate`, `FormFieldTemplate`
  - Shared template status support (`draft` / `published`) and field widget enum support (`input`, `textarea`, `attachment`, `select`, `date`)
  - Migration wiring for runtime and handler blackbox tests
  - Authenticated `/api/tmpl` module with list/create endpoints and RBAC guards
  - Minimal blackbox create/list chain test across all template levels
- Success criteria:
  - All new template endpoints are reachable and enforce baseline validation
  - `make test` passes with template flow coverage

### M5.2 Publish + Read Models
- Deliverables:
  - Publish/unpublish endpoints and validation rules for hierarchical consistency
  - Read/query improvements (filter by parent IDs, status, version)
  - Audit metadata for publish operations
- Success criteria:
  - Published templates are queryable by consumers without draft leakage

### M5.3 Instantiation Engine
- Deliverables:
  - Service to instantiate runtime project/stage/form/field records from templates
  - Transactional creation pipeline with rollback safety
  - Integration tests for full instantiation path
- Success criteria:
  - One API call can materialize a runtime project skeleton from a published template

### M5.4 Governance + UI Integration
- Deliverables:
  - Template permissions hardening and conflict-safe edit flow
  - Frontend management UI for template lifecycle
  - Version-diff visibility and migration guidance for existing projects
- Success criteria:
  - Operations can safely manage template versions with clear rollout controls

## Next Steps
1. Add parent/status/version query filters to each `/api/tmpl/*` list endpoint.
2. Add publish/unpublish endpoints with hierarchy guards (no child published under draft parent).
3. Define the runtime instantiation contract (input/output schema and transaction boundaries).
4. Extend blackbox tests for invalid payloads and RBAC denial paths.
