# Phase 02: Database Foundation Spec

## Status

- Approved for implementation

## Context

Phase 1 created the Docker-first backend skeleton. Phase 2 establishes the database foundation that later authentication, category, subscription, reminder, payment record, and dashboard phases will build on.

The goal is to make the backend connect to PostgreSQL, run MVP schema migrations, define GORM models, and provide initial repository boundaries without implementing business APIs yet.

## Requirements

- REQ-001: The backend must initialize a PostgreSQL connection from environment/config values.
- REQ-002: The backend must run a local development migration path during startup.
- REQ-003: The schema must include MVP tables: `users`, `categories`, `subscriptions`, `reminders`, and `payment_records`.
- REQ-004: User-owned tables must include `user_id` isolation fields.
- REQ-005: MVP models must include `created_at` and `updated_at`.
- REQ-006: Tables that need soft deletion must include deleted-at support.
- REQ-007: Common user-scoped list/query paths must have indexes.
- REQ-008: Repository interfaces and initial GORM implementations must exist for MVP tables.
- REQ-009: Backend startup through Docker must connect to PostgreSQL and complete migrations.

## Non-Goals

- NOGOAL-001: Do not implement authentication APIs.
- NOGOAL-002: Do not implement subscription CRUD APIs.
- NOGOAL-003: Do not implement category CRUD APIs.
- NOGOAL-004: Do not implement dashboard aggregation APIs.
- NOGOAL-005: Do not introduce a separate external migration CLI unless required.

## Design

Phase 2 adds these backend modules:

- `internal/database`: PostgreSQL DSN construction, GORM connection initialization, and connection close helper.
- `internal/migration`: schema migration runner for local development startup.
- `internal/model`: MVP GORM models and enum-like constants.
- `internal/repository`: repository interfaces plus initial GORM-backed constructors and implementations.

Startup flow:

```text
config.Load()
logger.New()
database.Open()
migration.Run()
router.New()
http.ListenAndServe()
```

Migration approach:

- Use GORM `AutoMigrate` for the early MVP foundation.
- Enable PostgreSQL `pgcrypto` so UUID defaults can use `gen_random_uuid()`.
- Keep the migration runner isolated behind `internal/migration` so it can be replaced by versioned SQL migrations later if needed.

Repository boundary:

- Repositories own database access.
- Handlers and services must not use `*gorm.DB` directly.
- Phase 2 creates minimal repository methods that later phases can expand.

## Data Model

### users

- `id` UUID primary key
- `email` unique, indexed
- `password_hash`
- `name`
- `created_at`
- `updated_at`
- `deleted_at`

### categories

- `id` UUID primary key
- `user_id` indexed
- `name`
- `color`
- `created_at`
- `updated_at`
- `deleted_at`
- Unique per user: `user_id`, `name`

### subscriptions

- `id` UUID primary key
- `user_id` indexed
- `category_id` nullable, indexed
- `name`
- `amount`
- `currency`
- `billing_cycle`
- `next_payment_date` indexed with user scope
- `status`
- `notes`
- `created_at`
- `updated_at`
- `deleted_at`

### reminders

- `id` UUID primary key
- `user_id` indexed
- `subscription_id` indexed
- `remind_before_days`
- `enabled`
- `created_at`
- `updated_at`
- `deleted_at`

### payment_records

- `id` UUID primary key
- `user_id` indexed
- `subscription_id` indexed
- `amount`
- `currency`
- `paid_at` indexed with user scope
- `note`
- `created_at`
- `updated_at`
- `deleted_at`

## API Contract

- No API changes.

## Task Breakdown

- TASK-001: Extend config with PostgreSQL SSL mode and DSN support.
- TASK-002: Add `internal/database` connection package.
- TASK-003: Add `internal/model` base model and MVP models.
- TASK-004: Add `internal/migration` runner with `pgcrypto` enablement and GORM migrations.
- TASK-005: Wire database connection and migration into backend startup.
- TASK-006: Add repository interfaces and initial GORM implementations.
- TASK-007: Add unit tests for DSN construction and migration model coverage where practical.
- TASK-008: Verify Docker startup creates all MVP tables.

## Acceptance Criteria

- AC-001: `go test ./...` passes.
- AC-002: `go build -o /tmp/subpilot-server ./cmd/server` passes.
- AC-003: `docker compose up --build -d` starts PostgreSQL, Redis, and backend.
- AC-004: Backend logs show the database connected and migrations completed.
- AC-005: PostgreSQL contains `users`, `categories`, `subscriptions`, `reminders`, and `payment_records`.
- AC-006: Models include timestamps and user isolation fields where required.
- AC-007: Repository interfaces exist for all MVP tables.

## Verification Plan

- VERIFY-001: Run `go test ./...` under `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` under `apps/server`.
- VERIFY-003: Run `docker compose config`.
- VERIFY-004: Run `docker compose up --build -d`.
- VERIFY-005: Request `GET http://127.0.0.1:18080/api/v1/health`.
- VERIFY-006: Query PostgreSQL through Docker to confirm MVP tables exist.
- VERIFY-007: Run `docker compose down` after verification.

## Traceability

```text
REQ-001 -> TASK-001,TASK-002,TASK-005 -> AC-003,AC-004 -> VERIFY-003,VERIFY-004
REQ-002 -> TASK-004,TASK-005 -> AC-004 -> VERIFY-004
REQ-003 -> TASK-003,TASK-004 -> AC-005 -> VERIFY-006
REQ-004 -> TASK-003 -> AC-006 -> VERIFY-001
REQ-005 -> TASK-003 -> AC-006 -> VERIFY-001
REQ-006 -> TASK-003 -> AC-006 -> VERIFY-001
REQ-007 -> TASK-003 -> AC-006 -> VERIFY-006
REQ-008 -> TASK-006 -> AC-007 -> VERIFY-001
REQ-009 -> TASK-005,TASK-008 -> AC-003,AC-004 -> VERIFY-004,VERIFY-005
```

## Open Questions

- None.
