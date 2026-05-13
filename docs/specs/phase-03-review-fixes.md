# Phase 03: Authentication Review Fixes Spec

## Status

- Approved for implementation

## Context

Phase 3 authentication is implemented, but code review found three hardening gaps:

- Duplicate email registration can return an internal error when the database unique constraint wins a race.
- Docker Compose can start with the committed default JWT secret.
- JWT expiration validation lacks a direct regression test.

This spec narrows the fix to those review findings.

## Requirements

- REQ-001: Registration must return the existing duplicate-email domain error when user creation hits a database unique constraint.
- REQ-002: Docker Compose must not start the server with the committed default JWT secret.
- REQ-003: Token tests must directly prove expired JWTs are rejected.

## Non-Goals

- NOGOAL-001: Do not change the public authentication API shape.
- NOGOAL-002: Do not introduce refresh tokens, sessions, or token revocation.
- NOGOAL-003: Do not redesign deployment topology.

## Design

- Keep duplicate email handling inside the user repository/service boundary.
- Convert repository duplicate email errors into `service.ErrEmailAlreadyExists`.
- Require Compose callers to provide `JWT_SECRET` explicitly and reject the previous committed default at entrypoint startup.
- Add focused token expiration test coverage without changing token claims.

## Data Model

- No data model changes.

## API Contract

- No API shape changes.
- Duplicate registration must consistently return the existing conflict response.

## Task Breakdown

- TASK-001: Add repository-level duplicate email error mapping.
- TASK-002: Map repository duplicate email errors to auth service duplicate email errors.
- TASK-003: Require explicit Docker Compose `JWT_SECRET` and reject the default secret outside production too.
- TASK-004: Add JWT expiration regression test.
- TASK-005: Run Go tests, build, and Compose config verification.

## Acceptance Criteria

- AC-001: Normal and database-level duplicate email registration paths both map to `ErrEmailAlreadyExists`.
- AC-002: `docker compose config` fails when `JWT_SECRET` is missing and succeeds when it is provided.
- AC-003: The Docker entrypoint rejects an empty or committed default `JWT_SECRET`.
- AC-004: Expired JWTs return an invalid-token error in tests.

## Verification Plan

- VERIFY-001: Run `go test ./...` in `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` in `apps/server`.
- VERIFY-003: Run `docker compose config` without `JWT_SECRET` and confirm it fails.
- VERIFY-004: Run `JWT_SECRET=<non-default> docker compose config` and confirm it passes.
- VERIFY-005: Run the Docker entrypoint with the committed default `JWT_SECRET` and confirm it fails.

## Traceability

```text
REQ-001 -> TASK-001,TASK-002 -> AC-001 -> VERIFY-001
REQ-002 -> TASK-003 -> AC-002,AC-003 -> VERIFY-003,VERIFY-004,VERIFY-005
REQ-003 -> TASK-004 -> AC-004 -> VERIFY-001
```

## Open Questions

- None.
