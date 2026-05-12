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
