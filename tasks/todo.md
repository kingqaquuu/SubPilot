# Repository Initialization Todo

- [x] Read `AGENTS.md` and follow project workflow rules.
- [x] Inspect the current directory and existing Git metadata.
- [x] Initialize or repair the local Git repository.
- [x] Configure local Git identity.
- [x] Create an initial commit on `master`.
- [x] Create the GitHub repository. Remote provided: `https://github.com/kingqaquuu/SubPilot`.
- [x] Connect local `master` to the GitHub remote and push.
- [x] Verify local and remote repository state.

## Review

- Initialized the local Git repository on `master`.
- Configured local Git identity as `kingqaquuu <kingqaquuu@outlook.com>`.
- Added `origin` remote: `https://github.com/kingqaquuu/SubPilot`.
- Pushed `master` to GitHub and set upstream tracking.

---

# Project Development Specification Todo

- [x] Read `AGENTS.md` and existing task notes.
- [x] Confirm the repository has no conflicting local changes.
- [x] Create a phased project development specification under `docs/`.
- [x] Verify the document structure and content.
- [x] Record review notes.

## Review

- Added `docs/development-spec.md`, then revised it after user feedback to avoid duplicating `AGENTS.md`.
- The document now focuses on concrete phases, prerequisites, deliverables, implementation tasks, acceptance criteria, exit gates, dependencies, and release milestones.
- Added `tasks/lessons.md` to record that phase-plan documents should reference `AGENTS.md` instead of repeating its rules.

---

# Phase 0 Project Baseline Todo

- [x] Review Phase 0 scope in `docs/development-spec.md`.
- [x] Confirm the repository has no conflicting local changes.
- [x] Create base directory structure.
- [x] Add `.gitignore`.
- [x] Add `.env.example`.
- [x] Add initial `README.md`.
- [x] Verify no real secrets are committed.
- [x] Commit and push Phase 0 baseline.

## Review

- Created the Phase 0 base directory structure with `.gitkeep` placeholders for empty directories.
- Added `.gitignore`, `.env.example`, and initial `README.md`.
- Verified `.env.example` uses placeholder values and no real secrets were introduced.

---

# Phase 1 Local Infrastructure And Backend Skeleton Todo

- [x] Review Phase 1 scope in `docs/development-spec.md`.
- [x] Confirm the repository has no conflicting local changes.
- [x] Create PostgreSQL, Redis, and backend services in Docker Compose.
- [x] Create Go backend module.
- [x] Add configuration loading from environment and optional config file.
- [x] Add structured logger initialization.
- [x] Add Gin router with `/api/v1/health`.
- [x] Add unified response helper.
- [x] Add Swagger scaffolding.
- [x] Add backend Dockerfile.
- [x] Verify backend tests/build and Docker startup path.
- [x] Commit and push Phase 1 skeleton.

## Review

- Added Docker Compose services for PostgreSQL, Redis, and the backend server.
- Added the Go backend module with config loading, Zap logger setup, Gin router, unified response helper, health endpoint, and Swagger scaffold.
- Added the backend Dockerfile and example YAML config.
- Verified `go test ./...`, `go build -o /tmp/subpilot-server ./cmd/server`, `docker compose config`, `docker compose up --build -d`, and `GET /api/v1/health`.

---

# Phase 1 Review Fixes Todo

- [x] Review reported Phase 1 issues.
- [x] Serve Swagger spec at the path advertised by the placeholder endpoint.
- [x] Prevent production startup with default development secrets.
- [x] Add config loading tests for environment variables and secret validation.
- [x] Run Go tests, build, Compose config, and health endpoint verification.
- [x] Commit and push Phase 1 review fixes.

## Review

- Served `docs/swagger.yaml` through `/docs/swagger.yaml` and added route coverage.
- Added explicit environment variable bindings and config tests for env loading and production secret validation.
- Added a Docker entrypoint guard that rejects default development secrets when `APP_ENV=production`.
- Verified `go test ./...`, `go build -o /tmp/subpilot-server ./cmd/server`, `docker compose config`, `docker compose up --build -d`, `GET /api/v1/health`, `GET /docs/swagger.yaml`, and production default secret rejection.

---

# Phase 2 Database Foundation Todo

- [x] Create and review `docs/specs/phase-02-database-foundation.md`.
- [x] Extend config with PostgreSQL SSL mode and DSN support.
- [x] Add `internal/database` connection package.
- [x] Add `internal/model` base model and MVP models.
- [x] Add `internal/migration` runner.
- [x] Wire database connection and migration into backend startup.
- [x] Add repository interfaces and initial GORM implementations.
- [x] Add unit tests for DSN and migration model coverage.
- [x] Verify Go tests, build, Compose config, Docker startup, health endpoint, and MVP tables.
- [x] Commit and push Phase 2 database foundation.

## Review

- Added `docs/specs/phase-02-database-foundation.md` before implementation according to SDD.
- Added PostgreSQL SSL mode config, DSN construction, GORM database connection, and startup migration flow.
- Added MVP models for users, categories, subscriptions, reminders, and payment records with timestamps, soft delete support, user isolation fields, and query indexes.
- Added repository interfaces and initial GORM-backed implementations for all MVP tables.
- Verified `go test ./...`, `go build -o /tmp/subpilot-server ./cmd/server`, `docker compose config`, `docker compose up --build -d`, `GET /api/v1/health`, database migration logs, MVP table creation, and key indexes.

---

# Phase 3 Authentication Todo

- [x] Create and review `docs/specs/phase-03-authentication.md`.
- [x] Add auth DTOs for register, login, auth response, and current user response.
- [x] Add JWT token manager with signing, expiry, parsing, and validation.
- [x] Add auth service for register, login, password hashing, and current user lookup.
- [x] Add auth middleware and context helpers.
- [x] Add auth handler endpoints.
- [x] Wire auth routes into router.
- [x] Update Swagger docs for auth endpoints and bearer auth.
- [x] Add tests for token manager, service, middleware, and handler behavior.
- [x] Verify Go tests, build, Compose config, Docker startup, auth endpoints, and stored password hash.
- [x] Commit and push Phase 3 authentication.

## Review

- Added `docs/specs/phase-03-authentication.md` before implementation according to SDD.
- Added auth DTOs, JWT token manager, auth context helpers, auth middleware, auth service, and auth handler endpoints.
- Wired `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, and protected `GET /api/v1/auth/me` into the router.
- Updated Swagger docs with auth endpoints and bearer auth.
- Added focused tests for JWT behavior, auth service behavior, middleware authorization, and handler auth flow.
- Verified `go test ./...`, `go build -o /tmp/subpilot-server ./cmd/server`, `docker compose config`, `docker compose up --build -d`, register/login/me flows, missing token `401`, invalid credentials `401`, and stored bcrypt password hash.

---

# Phase 3 Authentication Review Fixes Todo

- [x] Create and review `docs/specs/phase-03-review-fixes.md`.
- [x] Map database duplicate email conflicts to the existing auth conflict error.
- [x] Require an explicit non-default JWT secret for Docker Compose startup.
- [x] Add JWT expiration regression coverage.
- [x] Verify Go tests, build, Compose config failure/success, and entrypoint secret rejection.
- [x] Commit and push Phase 3 review fixes.

## Review

- Added `docs/specs/phase-03-review-fixes.md` before implementation according to SDD.
- Mapped Postgres unique constraint errors from user creation to the existing duplicate email auth error.
- Required Docker Compose callers to provide `JWT_SECRET` explicitly and rejected the committed default secret at entrypoint startup.
- Added JWT expiration regression coverage and non-production default JWT secret config coverage.
- Verified `go test ./...`, `go build -o /tmp/subpilot-server ./cmd/server`, `docker compose config` missing-secret failure, `JWT_SECRET=local-review-fix-secret docker compose config`, default-secret entrypoint rejection, and non-default entrypoint success.
