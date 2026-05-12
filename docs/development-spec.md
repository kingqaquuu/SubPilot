# SubPilot Phase Plan

## Purpose

This document is the concrete implementation plan for SubPilot. Project principles, stack choices, repository layout, API rules, testing commands, and workflow rules are defined in `AGENTS.md` and are not repeated here.

Each phase below defines:

- Objective
- Prerequisites
- Deliverables
- Implementation tasks
- Acceptance criteria
- Exit gate

## Phase 0: Project Baseline

### Objective

Create the minimum repository foundation needed for all later work.

### Prerequisites

- Git repository exists on `master`.
- `AGENTS.md` exists and is treated as the source of project rules.

### Deliverables

- `README.md`
- `.gitignore`
- `.env.example`
- `docs/development-spec.md`
- Base directories:
  - `apps/server`
  - `apps/web`
  - `apps/mobile`
  - `deploy/docker`
  - `deploy/nginx`
  - `deploy/scripts`
  - `docs`
  - `screenshots`

### Implementation Tasks

1. Create the base directory structure.
2. Add `.gitignore` for Go, Node, Flutter, Docker, environment files, logs, and IDE artifacts.
3. Add `.env.example` with placeholder values for backend, database, Redis, JWT, and frontend API base URL.
4. Add initial `README.md` with project description, stack summary, local development placeholder, and roadmap placeholder.
5. Confirm no secret values are committed.

### Acceptance Criteria

- Repository structure matches the expected project skeleton.
- `.env.example` documents required configuration names without real secrets.
- README can orient a new contributor.

### Exit Gate

- Commit baseline files before starting backend or frontend implementation.

## Phase 1: Local Infrastructure And Backend Skeleton

### Objective

Start the backend service and required local dependencies through Docker.

### Prerequisites

- Phase 0 complete.
- Docker and Docker Compose available locally.

### Deliverables

- `docker-compose.yml`
- `apps/server/go.mod`
- `apps/server/cmd/server/main.go`
- Backend config loader
- Backend logger bootstrap
- Health check endpoint
- Base router setup
- Initial Swagger setup

### Implementation Tasks

1. Create PostgreSQL and Redis services in Docker Compose.
2. Create the Go backend module.
3. Add configuration loading from environment and config file.
4. Add structured logger initialization.
5. Add Gin router with `/api/v1/health`.
6. Add unified response helper.
7. Add Swagger scaffolding for future API docs.
8. Add backend Dockerfile.

### Acceptance Criteria

- `docker compose up --build` starts PostgreSQL, Redis, and backend.
- `GET /api/v1/health` returns a successful unified response.
- Backend starts without hardcoded secrets.

### Exit Gate

- Backend skeleton runs locally through Docker.

## Phase 2: Database Foundation

### Objective

Establish database connection, migration flow, and core data models needed by MVP.

### Prerequisites

- Phase 1 complete.
- PostgreSQL service available through Docker.

### Deliverables

- Database connection package
- Migration mechanism
- Base GORM models
- Models for:
  - users
  - categories
  - subscriptions
  - reminders
  - payment_records
- Seed or development bootstrap notes if needed

### Implementation Tasks

1. Add database configuration and connection initialization.
2. Add migration execution path for local development.
3. Define model fields and relationships.
4. Add indexes for user ownership and common list queries.
5. Add repository interfaces and initial repository implementations.
6. Add minimal database connection tests where practical.

### Acceptance Criteria

- Backend can connect to PostgreSQL from Docker.
- Migrations create all MVP tables.
- Models include timestamps and user isolation fields where required.
- Common query paths have indexes.

### Exit Gate

- Database schema is ready for authentication and subscription APIs.

## Phase 3: Authentication

### Objective

Implement secure user registration, login, JWT issuance, and authenticated API access.

### Prerequisites

- Phase 2 complete.
- `users` table exists.

### Deliverables

- Register endpoint
- Login endpoint
- Password hashing
- JWT creation and validation
- Auth middleware
- Current user endpoint
- Auth DTOs
- Auth service and repository methods

### Implementation Tasks

1. Implement register request validation.
2. Store users with hashed passwords only.
3. Implement login credential verification.
4. Issue JWT access tokens.
5. Add middleware that rejects invalid or missing tokens.
6. Add current user lookup.
7. Add Swagger docs for auth endpoints.
8. Add focused backend tests for auth logic.

### Acceptance Criteria

- New users can register.
- Registered users can log in.
- Invalid credentials are rejected.
- Protected routes reject missing or invalid JWTs.
- Password hashes are not reversible plaintext.

### Exit Gate

- Authentication is usable by later user-scoped APIs.

## Phase 4: Category Management

### Objective

Provide user-owned category management for organizing subscriptions.

### Prerequisites

- Phase 3 complete.
- Authenticated user ID is available in request context.

### Deliverables

- Category create API
- Category list API
- Category update API
- Category delete API
- Category service and repository methods
- Swagger docs

### Implementation Tasks

1. Define category DTOs.
2. Implement create with duplicate-name handling per user.
3. Implement list scoped by authenticated user.
4. Implement update scoped by authenticated user.
5. Implement delete scoped by authenticated user.
6. Prevent one user from reading or mutating another user's categories.
7. Add API tests or service tests for ownership behavior.

### Acceptance Criteria

- Users can manage only their own categories.
- Category list returns frontend-ready data.
- Duplicate or invalid category input returns structured errors.

### Exit Gate

- Categories can be used by subscription APIs.

## Phase 5: Subscription Management

### Objective

Implement the core subscription CRUD workflow.

### Prerequisites

- Phase 4 complete.
- Category APIs are available.

### Deliverables

- Subscription create API
- Subscription list API
- Subscription detail API
- Subscription update API
- Subscription delete API
- Filtering and pagination
- Subscription service and repository methods
- Swagger docs

### Implementation Tasks

1. Define subscription DTOs.
2. Implement create with category ownership validation.
3. Implement paginated list by authenticated user.
4. Add filters for status, category, and billing cycle.
5. Implement detail retrieval scoped by authenticated user.
6. Implement update with ownership checks.
7. Implement delete or soft delete.
8. Add tests for ownership, pagination, and validation.

### Acceptance Criteria

- Users can complete subscription CRUD through APIs.
- Users cannot access another user's subscriptions.
- List endpoint supports pagination.
- Invalid billing data returns structured errors.

### Exit Gate

- Backend supports the primary MVP domain workflow.

## Phase 6: Reminders And Payment Records

### Objective

Track upcoming payments and historical payment records.

### Prerequisites

- Phase 5 complete.
- Subscription ownership rules are implemented.

### Deliverables

- Reminder configuration APIs
- Payment record create/list APIs
- Payment record association with subscriptions
- Upcoming payment query support
- Swagger docs

### Implementation Tasks

1. Implement reminder create/update behavior for subscriptions.
2. Implement payment record creation scoped by subscription owner.
3. Implement payment record listing with pagination.
4. Add upcoming payment query by authenticated user.
5. Add validation for dates, amounts, and currency.
6. Add tests for ownership and date-based queries.

### Acceptance Criteria

- Users can configure reminders for their subscriptions.
- Users can record payments for owned subscriptions.
- Upcoming payment data is queryable for dashboard use.

### Exit Gate

- Backend has enough data for MVP dashboard statistics.

## Phase 7: Dashboard APIs

### Objective

Expose aggregated subscription data for the web dashboard.

### Prerequisites

- Phase 6 complete.
- Subscription, category, reminder, and payment data exist.

### Deliverables

- Dashboard summary endpoint
- Monthly spending endpoint or response section
- Category distribution endpoint or response section
- Upcoming payments endpoint or response section
- Active subscription count
- Swagger docs

### Implementation Tasks

1. Define dashboard response DTOs.
2. Implement monthly spending aggregation.
3. Implement category distribution aggregation.
4. Implement upcoming payment summary.
5. Implement active subscription count.
6. Review query performance and add indexes if needed.
7. Add tests for user-scoped aggregation.

### Acceptance Criteria

- Dashboard APIs return user-scoped data only.
- Aggregation responses are directly usable by the frontend.
- Query performance is acceptable for MVP data volumes.

### Exit Gate

- Backend MVP API surface is complete.

## Phase 8: Web App Foundation

### Objective

Create the Vue web application foundation and connect it to backend authentication.

### Prerequisites

- Phase 3 complete for auth integration.
- Backend base URL is configurable.

### Deliverables

- Vue app scaffold
- Router setup
- Pinia setup
- Axios client
- Auth store
- Login page
- Register page
- Authenticated layout
- i18n baseline
- Dark mode baseline

### Implementation Tasks

1. Scaffold the web app under `apps/web`.
2. Add routes for login, register, dashboard, subscriptions, and categories.
3. Add API client with token handling.
4. Implement auth store.
5. Build login and register pages.
6. Add route guards for authenticated pages.
7. Add basic app layout and navigation.
8. Add English translation keys for visible text.

### Acceptance Criteria

- Users can register and log in through the web UI.
- Authenticated routes are protected.
- Frontend builds successfully.
- API base URL is configurable.

### Exit Gate

- Web app can authenticate against the backend.

## Phase 9: Web Subscription Workflows

### Objective

Implement web UI for category and subscription management.

### Prerequisites

- Phase 8 complete.
- Phase 5 backend APIs available.

### Deliverables

- Category list and form UI
- Subscription list UI
- Subscription create/edit form
- Subscription detail view
- Delete confirmation flow
- Loading and error states

### Implementation Tasks

1. Build category management screens.
2. Build subscription list with pagination and filters.
3. Build subscription create and edit forms.
4. Add subscription detail view.
5. Add delete confirmation behavior.
6. Add loading, empty, and error states.
7. Verify responsive behavior on desktop and mobile widths.

### Acceptance Criteria

- Users can manage categories through the web UI.
- Users can complete subscription CRUD through the web UI.
- UI handles loading, empty, and error states.
- Frontend build passes.

### Exit Gate

- Primary MVP workflow is usable in the browser.

## Phase 10: Web Dashboard

### Objective

Build the MVP dashboard using backend summary APIs.

### Prerequisites

- Phase 7 complete.
- Phase 8 web foundation complete.

### Deliverables

- Dashboard summary cards
- Monthly spending chart
- Category distribution chart
- Upcoming payments panel
- Active subscriptions metric

### Implementation Tasks

1. Connect dashboard API client methods.
2. Build dashboard summary components.
3. Add ECharts visualizations.
4. Add loading and empty states.
5. Verify dark mode rendering.
6. Verify responsive dashboard layout.

### Acceptance Criteria

- Dashboard renders real user-scoped backend data.
- Charts remain readable in light and dark modes.
- Dashboard is responsive.

### Exit Gate

- Web MVP feature set is complete.

## Phase 11: Deployment Documentation And Release Prep

### Objective

Make the MVP reproducible for developers and self-hosters.

### Prerequisites

- Backend MVP complete.
- Web MVP complete.
- Docker Compose setup working.

### Deliverables

- Complete README quick start
- Docker deployment guide
- Environment variable reference
- Swagger API documentation
- Nginx example
- Initial screenshots
- MVP release checklist

### Implementation Tasks

1. Update README with local startup steps.
2. Document Docker deployment.
3. Document environment variables.
4. Confirm Swagger docs match backend routes.
5. Add Nginx example config.
6. Capture screenshots for README.
7. Create release checklist.

### Acceptance Criteria

- A new developer can start the project from the README.
- Docker startup path is documented and verified.
- API docs match implemented APIs.

### Exit Gate

- Web MVP is ready for a first tagged release.

## Phase 12: Mobile MVP

### Objective

Create a lightweight Flutter client that reuses the existing backend APIs.

### Prerequisites

- Backend MVP APIs are stable.
- Web MVP has validated core workflows.

### Deliverables

- Flutter app scaffold
- Auth flow
- Subscription list
- Subscription detail
- Subscription create/edit flow
- Dashboard summary screen
- Dio API client
- Riverpod state management

### Implementation Tasks

1. Scaffold mobile app under `apps/mobile`.
2. Add API client and environment configuration.
3. Implement login and register screens.
4. Implement authenticated navigation.
5. Implement subscription list and detail screens.
6. Implement create/edit flow.
7. Implement dashboard summary screen.
8. Test against local backend.

### Acceptance Criteria

- Mobile app authenticates against the backend.
- Mobile app can manage subscriptions through existing APIs.
- Mobile UI works on common phone screen sizes.

### Exit Gate

- Mobile MVP is usable against the shared backend.

## Phase 13: Hardening

### Objective

Reduce release risk before broader use.

### Prerequisites

- MVP web flow complete.
- Optional mobile MVP complete if included in the release.

### Deliverables

- Security review notes
- Error handling review notes
- Logging review notes
- Performance review notes
- Backup and restore notes
- Known limitations list

### Implementation Tasks

1. Review authentication and authorization boundaries.
2. Review user-owned data isolation.
3. Review sensitive log output.
4. Review API error responses.
5. Review common database query plans.
6. Document backup and restore expectations.
7. Document known MVP limitations.

### Acceptance Criteria

- No known critical auth or data isolation issues remain.
- Runtime errors are logged with useful context.
- User-facing errors do not expose internals.
- Known limitations are explicit.

### Exit Gate

- MVP is ready for public release or private beta.

## Phase Dependency Map

```text
Phase 0
  -> Phase 1
  -> Phase 2
  -> Phase 3
  -> Phase 4
  -> Phase 5
  -> Phase 6
  -> Phase 7
  -> Phase 8
  -> Phase 9
  -> Phase 10
  -> Phase 11
  -> Phase 13

Phase 12 depends on stable backend APIs from Phase 7 and validated UX from Phase 10.
```

## Release Milestones

### Milestone A: Backend MVP

Includes:

- Phase 1 through Phase 7

Result:

- API surface is ready for frontend development and external testing.

### Milestone B: Web MVP

Includes:

- Phase 8 through Phase 11

Result:

- Browser-based MVP is ready for first release.

### Milestone C: Mobile MVP

Includes:

- Phase 12

Result:

- Mobile client is ready for shared backend validation.

### Milestone D: Release Candidate

Includes:

- Phase 13

Result:

- MVP is hardened enough for public release or private beta.
