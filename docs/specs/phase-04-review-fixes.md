# Phase 04: Category Review Fixes Spec

## Status

- Approved for implementation

## Context

Phase 4 added category soft deletion. Review found that soft deleting a category does not trigger the database-level `ON DELETE SET NULL` constraint for subscriptions because the category row remains present with `deleted_at` set.

The chosen product behavior is: deleting a category turns subscriptions in that category into uncategorized subscriptions.

## Requirements

- REQ-001: Deleting a category must set `category_id` to `NULL` for subscriptions owned by the same user and referencing that category.
- REQ-002: Clearing subscription category references and soft deleting the category must be atomic.
- REQ-003: Deleting one user's category must not clear another user's subscription references.
- REQ-004: Tests must cover the new delete side effect.

## Non-Goals

- NOGOAL-001: Do not implement subscription CRUD in this fix.
- NOGOAL-002: Do not hard delete categories.
- NOGOAL-003: Do not change the category API response shape.

## Design

- Keep the public category service and API unchanged.
- Update the GORM category repository delete operation to run in a transaction.
- In the transaction, update `subscriptions` rows where `user_id` and `category_id` match the deleted category, setting `category_id` to `NULL`.
- Then soft delete the category.

## Data Model

- No schema changes.

## API Contract

- No API changes.
- `DELETE /api/v1/categories/{id}` keeps returning the deleted category object.

## Task Breakdown

- TASK-001: Update category repository delete behavior to clear same-user subscription references in a transaction.
- TASK-002: Add focused service test coverage for delete clearing behavior.
- TASK-003: Run Go tests, build, and Compose config verification.

## Acceptance Criteria

- AC-001: Deleting a category clears same-user subscription references to that category.
- AC-002: Deleting a category does not clear other-user subscription references.
- AC-003: Existing category delete API behavior still passes tests.

## Verification Plan

- VERIFY-001: Run `go test ./...` in `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` in `apps/server`.
- VERIFY-003: Run `JWT_SECRET=<non-default> docker compose config`.

## Traceability

```text
REQ-001 -> TASK-001,TASK-002 -> AC-001 -> VERIFY-001
REQ-002 -> TASK-001 -> AC-003 -> VERIFY-001
REQ-003 -> TASK-001,TASK-002 -> AC-002 -> VERIFY-001
REQ-004 -> TASK-002 -> AC-001,AC-002 -> VERIFY-001
```

## Open Questions

- None.
