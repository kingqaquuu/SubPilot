# Phase 04: Category Management Spec

## Status

- Approved for implementation

## Context

Phase 3 established authentication and request-scoped user identity. Phase 4 adds user-owned category management so later subscription APIs can attach subscriptions to organized categories.

Categories are private to each user. A category name must be unique only within that user, not globally.

## Requirements

- REQ-001: Authenticated users must be able to create categories with a name and optional color.
- REQ-002: Category names must be required, trimmed, length-limited, and unique per user.
- REQ-003: Category colors must be optional and default to the model default when omitted.
- REQ-004: Authenticated users must be able to list only their own categories.
- REQ-005: Authenticated users must be able to update only their own categories.
- REQ-006: Authenticated users must be able to delete only their own categories.
- REQ-007: Category endpoints must reject missing or invalid authentication.
- REQ-008: Swagger docs must describe the category endpoints.
- REQ-009: Tests must cover service behavior, handler behavior, duplicate names, not-found cases, and user ownership boundaries.

## Non-Goals

- NOGOAL-001: Do not implement subscription-category assignment in this phase.
- NOGOAL-002: Do not add category ordering, icons, budgets, or analytics.
- NOGOAL-003: Do not implement hard deletes for soft-deleted models.

## Design

- Add category DTOs under `internal/dto`.
- Extend `CategoryRepository` with update, delete, and per-user name lookup helpers.
- Keep duplicate-name handling in the service layer and map database unique constraint errors to the same domain error.
- Add `CategoryService` for validation, normalization, ownership-scoped operations, and response mapping.
- Add `CategoryHandler` for HTTP binding, auth context lookup, response mapping, and error status mapping.
- Mount category routes under authenticated `/categories`.

## Data Model

- No schema changes.
- Existing `categories` table fields are used:
  - `id`
  - `user_id`
  - `name`
  - `color`
  - timestamps and soft delete
- Existing unique index `idx_categories_user_name` enforces per-user unique names.

## API Contract

All endpoints require bearer authentication.

- `POST /api/v1/categories`
  - Request: `{ "name": "Streaming", "color": "#0ea5e9" }`
  - Response: category object
  - Errors: `400 invalid_request`, `401`, `409 category_name_exists`
- `GET /api/v1/categories`
  - Response: category object array
  - Errors: `401`
- `PUT /api/v1/categories/{id}`
  - Request: `{ "name": "Music", "color": "#a855f7" }`
  - Response: category object
  - Errors: `400 invalid_request`, `401`, `404 category_not_found`, `409 category_name_exists`
- `DELETE /api/v1/categories/{id}`
  - Response: deleted category object
  - Errors: `400 invalid_category_id`, `401`, `404 category_not_found`

Category response fields:

- `id`
- `name`
- `color`
- `created_at`
- `updated_at`

## Task Breakdown

- TASK-001: Add category DTOs.
- TASK-002: Extend category repository methods for update, delete, and name lookup.
- TASK-003: Add category service validation and user-scoped operations.
- TASK-004: Add category handler and authenticated routes.
- TASK-005: Update Swagger docs.
- TASK-006: Add focused service and handler tests.
- TASK-007: Run Go tests, build, Compose config, and API smoke verification.

## Acceptance Criteria

- AC-001: Users can create, list, update, and delete their own categories.
- AC-002: Users cannot access another user's categories.
- AC-003: Duplicate category names for the same user return a conflict error.
- AC-004: The same category name may be used by different users.
- AC-005: Invalid input returns structured errors.
- AC-006: Missing or invalid auth rejects category routes.

## Verification Plan

- VERIFY-001: Run `go test ./...` in `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` in `apps/server`.
- VERIFY-003: Run `JWT_SECRET=<non-default> docker compose config`.
- VERIFY-004: Run Docker startup and smoke-test auth plus category create/list/update/delete.
- VERIFY-005: Run `docker compose down` after verification.

## Traceability

```text
REQ-001 -> TASK-001,TASK-002,TASK-003,TASK-004 -> AC-001 -> VERIFY-001,VERIFY-004
REQ-002 -> TASK-003,TASK-006 -> AC-003,AC-005 -> VERIFY-001,VERIFY-004
REQ-003 -> TASK-001,TASK-003,TASK-006 -> AC-001 -> VERIFY-001
REQ-004 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002 -> VERIFY-001,VERIFY-004
REQ-005 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002 -> VERIFY-001,VERIFY-004
REQ-006 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002 -> VERIFY-001,VERIFY-004
REQ-007 -> TASK-004,TASK-006 -> AC-006 -> VERIFY-001,VERIFY-004
REQ-008 -> TASK-005 -> AC-001 -> VERIFY-004
REQ-009 -> TASK-006 -> AC-001,AC-002,AC-003,AC-005,AC-006 -> VERIFY-001
```

## Open Questions

- None.
