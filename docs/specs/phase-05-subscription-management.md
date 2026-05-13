# Phase 05: Subscription Management Spec

## Status

- Approved for implementation

## Context

Phase 4 added user-owned categories. Phase 5 implements the core MVP workflow: authenticated users can manage their own subscriptions and optionally attach them to categories they own.

## Requirements

- REQ-001: Authenticated users must be able to create subscriptions.
- REQ-002: Create and update must validate name, amount, currency, billing cycle, status, next payment date, and optional category ownership.
- REQ-003: Authenticated users must be able to list only their own subscriptions with pagination.
- REQ-004: List must support filters for status, category, and billing cycle.
- REQ-005: Authenticated users must be able to retrieve only their own subscription detail.
- REQ-006: Authenticated users must be able to update only their own subscriptions.
- REQ-007: Authenticated users must be able to soft delete only their own subscriptions.
- REQ-008: Subscription endpoints must reject missing or invalid authentication.
- REQ-009: Swagger docs must describe subscription endpoints.
- REQ-010: Tests must cover ownership boundaries, pagination, filters, validation, and category ownership.

## Non-Goals

- NOGOAL-001: Do not implement reminders or payment records in this phase.
- NOGOAL-002: Do not implement dashboard analytics.
- NOGOAL-003: Do not implement currency conversion.
- NOGOAL-004: Do not implement recurring payment date advancement.

## Design

- Add subscription DTOs under `internal/dto`.
- Extend `SubscriptionRepository` with filtered pagination, count, update, and delete operations.
- Add `SubscriptionService` for validation, category ownership checks, user-scoped operations, and response mapping.
- Add `SubscriptionHandler` for HTTP binding, auth context lookup, query parsing, response mapping, and error status mapping.
- Mount subscription routes under authenticated `/subscriptions`.

## Data Model

- No schema changes.
- Existing `subscriptions` table is used.

## API Contract

All endpoints require bearer authentication.

- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions`
- `GET /api/v1/subscriptions/{id}`
- `PUT /api/v1/subscriptions/{id}`
- `DELETE /api/v1/subscriptions/{id}`

Subscription request fields:

- `name`: required, trimmed, max 160
- `amount`: required, positive
- `currency`: required three-letter code, normalized uppercase
- `billing_cycle`: `monthly`, `yearly`, or `custom`
- `next_payment_date`: `YYYY-MM-DD` or RFC3339 timestamp
- `status`: optional, defaults to `active`, allowed `active`, `paused`, `canceled`
- `category_id`: optional UUID; must belong to the authenticated user
- `notes`: optional, max 2000

List query fields:

- `page`: optional, default 1
- `page_size`: optional, default 20, max 100
- `status`: optional
- `category_id`: optional UUID
- `billing_cycle`: optional

## Task Breakdown

- TASK-001: Add subscription DTOs and pagination DTO.
- TASK-002: Extend subscription repository with filters, count, update, and delete.
- TASK-003: Add subscription service validation and category ownership checks.
- TASK-004: Add subscription handler and authenticated routes.
- TASK-005: Update Swagger docs.
- TASK-006: Add focused service and handler tests.
- TASK-007: Run Go tests, build, Compose config, and API smoke verification.

## Acceptance Criteria

- AC-001: Users can create, list, retrieve, update, and delete their own subscriptions.
- AC-002: Users cannot access another user's subscriptions.
- AC-003: Users cannot attach subscriptions to another user's categories.
- AC-004: List endpoint supports pagination and filters.
- AC-005: Invalid billing data returns structured errors.
- AC-006: Missing or invalid auth rejects subscription routes.

## Verification Plan

- VERIFY-001: Run `go test ./...` in `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` in `apps/server`.
- VERIFY-003: Run `JWT_SECRET=<non-default> docker compose config`.
- VERIFY-004: Run Docker startup and smoke-test auth plus subscription create/list/detail/update/delete.
- VERIFY-005: Run `docker compose down` after verification.

## Traceability

```text
REQ-001 -> TASK-001,TASK-002,TASK-003,TASK-004 -> AC-001 -> VERIFY-001,VERIFY-004
REQ-002 -> TASK-003,TASK-006 -> AC-003,AC-005 -> VERIFY-001,VERIFY-004
REQ-003 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002,AC-004 -> VERIFY-001,VERIFY-004
REQ-004 -> TASK-002,TASK-004,TASK-006 -> AC-004 -> VERIFY-001,VERIFY-004
REQ-005 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002 -> VERIFY-001,VERIFY-004
REQ-006 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002,AC-003 -> VERIFY-001,VERIFY-004
REQ-007 -> TASK-002,TASK-003,TASK-004,TASK-006 -> AC-001,AC-002 -> VERIFY-001,VERIFY-004
REQ-008 -> TASK-004,TASK-006 -> AC-006 -> VERIFY-001,VERIFY-004
REQ-009 -> TASK-005 -> AC-001 -> VERIFY-004
REQ-010 -> TASK-006 -> AC-001,AC-002,AC-003,AC-004,AC-005,AC-006 -> VERIFY-001
```

## Open Questions

- None.
