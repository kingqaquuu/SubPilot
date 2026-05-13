# Phase 03: Authentication Spec

## Status

- Approved for implementation

## Context

Phase 2 created the PostgreSQL schema and repository foundation. Phase 3 adds the first user-facing backend workflow: account registration, login, JWT issuance, JWT validation, and access to the current authenticated user.

This phase must establish the authentication boundary used by later user-scoped APIs.

## Requirements

- REQ-001: Users must be able to register with email, password, and name.
- REQ-002: Registration must validate required fields, email format, password length, and duplicate email conflicts.
- REQ-003: Passwords must be stored only as bcrypt hashes.
- REQ-004: Users must be able to log in with email and password.
- REQ-005: Invalid credentials must be rejected without exposing whether the email or password was wrong.
- REQ-006: Successful login and registration must return a signed JWT access token.
- REQ-007: JWT validation must check token signature and expiry.
- REQ-008: Auth middleware must reject missing, malformed, invalid, or expired Bearer tokens.
- REQ-009: Authenticated routes must be able to access the current user ID.
- REQ-010: `GET /api/v1/auth/me` must return the current authenticated user.
- REQ-011: Swagger docs must include authentication endpoints.
- REQ-012: Tests must cover password hashing, login validation, JWT issuance/validation, middleware behavior, and auth route behavior.

## Non-Goals

- NOGOAL-001: Do not implement OAuth login.
- NOGOAL-002: Do not implement refresh tokens.
- NOGOAL-003: Do not implement password reset.
- NOGOAL-004: Do not implement email verification.
- NOGOAL-005: Do not implement rate limiting yet.

## Design

New modules:

- `internal/dto`: request and response DTOs for auth endpoints.
- `internal/service`: auth service for registration, login, password hashing, and current user lookup.
- `internal/middleware`: JWT authentication middleware.
- `internal/auth`: JWT token manager and context helpers.

Layering:

```text
handler -> service -> repository -> database
middleware -> auth token manager -> repository lookup when needed
```

Security behavior:

- Password hashing uses bcrypt.
- JWT claims include `sub` as the user ID and `exp` as expiry.
- Login errors use a generic `invalid credentials` message.
- API responses use the existing unified response envelope.

Route structure:

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/me
```

`/auth/me` must be protected by auth middleware.

## Data Model

- No schema changes beyond Phase 2 `users` table.

## API Contract

### POST /api/v1/auth/register

Request:

```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "User"
}
```

Success:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "jwt",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "User"
    }
  }
}
```

Errors:

- `400` invalid request
- `409` email already exists

### POST /api/v1/auth/login

Request:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Success: same response as register.

Errors:

- `400` invalid request
- `401` invalid credentials

### GET /api/v1/auth/me

Headers:

```text
Authorization: Bearer <token>
```

Success:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User"
  }
}
```

Errors:

- `401` missing or invalid token
- `404` user not found

## Task Breakdown

- TASK-001: Add auth DTOs for register, login, auth response, and current user response.
- TASK-002: Add JWT token manager with signing, expiry, parsing, and validation.
- TASK-003: Add auth service for register, login, and current user lookup.
- TASK-004: Add auth middleware and context helpers.
- TASK-005: Add auth handler endpoints.
- TASK-006: Wire auth routes into router with public and protected groups.
- TASK-007: Update Swagger docs for auth endpoints and bearer auth.
- TASK-008: Add tests for token manager, auth service, middleware, and auth handler behavior.
- TASK-009: Verify Docker startup and auth endpoints against PostgreSQL.

## Acceptance Criteria

- AC-001: New users can register and receive a JWT.
- AC-002: Registered users can log in and receive a JWT.
- AC-003: Stored passwords are bcrypt hashes, not plaintext.
- AC-004: Invalid credentials return `401`.
- AC-005: Protected routes reject missing or invalid tokens.
- AC-006: Valid tokens allow access to `/api/v1/auth/me`.
- AC-007: Auth implementation follows the existing layered architecture.
- AC-008: Swagger documents auth endpoints and Bearer token auth.

## Verification Plan

- VERIFY-001: Run `go test ./...` under `apps/server`.
- VERIFY-002: Run `go build -o /tmp/subpilot-server ./cmd/server` under `apps/server`.
- VERIFY-003: Run `docker compose config`.
- VERIFY-004: Run `docker compose up --build -d`.
- VERIFY-005: Register a test user through `POST /api/v1/auth/register`.
- VERIFY-006: Log in through `POST /api/v1/auth/login`.
- VERIFY-007: Request `GET /api/v1/auth/me` with the returned token.
- VERIFY-008: Confirm `GET /api/v1/auth/me` without a token returns `401`.
- VERIFY-009: Query PostgreSQL to confirm password hash does not equal plaintext.
- VERIFY-010: Run `docker compose down` after verification.

## Traceability

```text
REQ-001 -> TASK-001,TASK-003,TASK-005,TASK-006 -> AC-001 -> VERIFY-005
REQ-002 -> TASK-001,TASK-003,TASK-005 -> AC-001 -> VERIFY-001,VERIFY-005
REQ-003 -> TASK-003 -> AC-003 -> VERIFY-001,VERIFY-009
REQ-004 -> TASK-001,TASK-003,TASK-005,TASK-006 -> AC-002 -> VERIFY-006
REQ-005 -> TASK-003,TASK-005 -> AC-004 -> VERIFY-001,VERIFY-006
REQ-006 -> TASK-002,TASK-003 -> AC-001,AC-002 -> VERIFY-005,VERIFY-006
REQ-007 -> TASK-002,TASK-004 -> AC-005,AC-006 -> VERIFY-001,VERIFY-007
REQ-008 -> TASK-004,TASK-006 -> AC-005 -> VERIFY-001,VERIFY-008
REQ-009 -> TASK-004,TASK-006 -> AC-006 -> VERIFY-001,VERIFY-007
REQ-010 -> TASK-003,TASK-005,TASK-006 -> AC-006 -> VERIFY-007
REQ-011 -> TASK-007 -> AC-008 -> VERIFY-001
REQ-012 -> TASK-008 -> AC-001,AC-002,AC-003,AC-004,AC-005,AC-006 -> VERIFY-001
```

## Open Questions

- None.
