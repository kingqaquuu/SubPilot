# AGENTS.md

## Project

SubPilot

A modern open-source subscription management platform.

Supports:

- Web
- Mobile
- Docker self-hosting
- Multi-user synchronization
- Internationalization (i18n)

License:

- MIT

---

## Tech Stack

### Backend

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- JWT
- Viper
- Zap
- Swagger

### Frontend

- Vue3
- Vite
- TypeScript
- Pinia
- Element Plus
- Axios
- ECharts

### Mobile

- Flutter
- Riverpod
- Dio
- GoRouter

### DevOps

- Docker
- Docker Compose
- Nginx

---

## Development Principles

- Keep architecture simple and maintainable.
- Prefer readability over abstraction.
- Avoid over-engineering.
- Prefer explicit logic.
- Keep modules loosely coupled.
- Follow Docker-first development.
- All services must run with docker compose.
- All APIs must be RESTful.
- All code must be production-oriented.
- Avoid generating unnecessary files.
- Avoid giant functions and giant packages.

---

## Workflow Orchestration

### Plan Node Default

- Any non-trivial task, including tasks with 3+ steps or architecture decisions, must enter planning mode first.
- If implementation drifts from the plan, stop immediately and re-plan before continuing.
- Verification work must also be planned, not only build or implementation work.
- Write detailed specifications early to reduce ambiguity.

### Subagent Strategy

- Use subagents heavily to keep the main context window clean.
- Delegate research, exploration, and parallel analysis to subagents.
- For complex problems, use subagents to invest additional reasoning capacity.
- Each subagent must focus on one direction only.

### Self-Improvement Loop

- After any user correction, immediately update `tasks/lessons.md` and record the pattern.
- Write rules that prevent repeated mistakes.
- Ruthlessly iterate lessons until the error rate goes down.
- At the start of each session, review project-relevant lessons first.

### Verification Before Done

- Never mark work complete before proving that it works.
- Compare behavior against the main branch when necessary.
- Ask: "Would a senior engineer approve this?"
- Run tests, inspect logs, and demonstrate correctness.

### Demand Elegance (Balanced)

- For non-trivial changes, pause and ask: "Is there a more elegant way?"
- If a fix feels hacky, implement the elegant solution based on the current full context.
- Do not over-engineer simple problems.
- Challenge the work before presenting it.

### Autonomous Bug Fixing

- When receiving a bug report, fix it directly without requiring the user to hand-hold the process.
- Point to logs, errors, or failing tests, then solve the problem.
- Avoid forcing the user into unnecessary context switching.
- Automatically fix failing CI tests when possible.

---

## Task Management

1. Plan first: write the plan to `tasks/todo.md` as a checkable checklist.
2. Validate the plan: check in before implementation.
3. Track progress: mark each item as soon as it is completed.
4. Explain changes: provide a high-level summary for each step.
5. Record results: append a review section to the end of `tasks/todo.md`.
6. Capture lessons: after corrections, update `tasks/lessons.md`.

---

## Core Principles

- Simplicity first: keep every change as simple as possible and affect the smallest necessary code surface.
- Never be lazy: find the root cause, avoid temporary fixes, and hold senior-engineer standards.
- Minimal impact: change only what is necessary and avoid introducing new bugs.

---

## Repository Structure

```text
subpilot
├── apps
│   ├── server
│   ├── web
│   └── mobile
│
├── deploy
│   ├── docker
│   ├── nginx
│   └── scripts
│
├── docs
├── screenshots
├── .env
├── .env.example
├── docker-compose.yml
└── AGENTS.md
```

---

## Backend Rules

### Stack

- Go 1.25+
- Gin
- GORM
- PostgreSQL
- Redis
- JWT
- Viper
- Zap

### Architecture

Use layered architecture:

```text
handler -> service -> repository -> database
```

Never access database directly inside handlers.

### Backend Structure

```text
apps/server
├── cmd
│   └── server
│       └── main.go
│
├── internal
│   ├── config
│   ├── router
│   ├── middleware
│   ├── handler
│   ├── service
│   ├── repository
│   ├── model
│   ├── dto
│   ├── task
│   ├── response
│   └── utils
│
├── configs
│   ├── config.yaml
│   └── config.example.yaml
│
├── docs
├── migrations
├── go.mod
└── Dockerfile
```

---

## API Rules

- Use RESTful API design.
- Use `/api/v1` prefix.
- Return unified response structure.
- Use JWT authentication.
- Validate all request parameters.
- Use pagination for list APIs.
- Never expose internal database fields.
- Keep APIs frontend-friendly.
- Use semantic HTTP status codes.

### Response Format

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### Error Handling

- Never panic in business logic.
- Return structured errors.
- Log internal errors with Zap.
- Do not expose internal server errors to clients.

---

## Authentication Rules

- Use JWT authentication.
- Use bcrypt password hashing.
- Never store plaintext passwords.
- JWT secret must come from environment variables.
- Support multi-user isolation.

---

## Configuration Rules

Use:

- `.env`
- `config.yaml`

Rules:

- Sensitive configuration must come from `.env`.
- Support development and production environments.
- Never hardcode secrets.
- Provide `.env.example`.
- Provide `config.example.yaml`.

---

## Database Rules

### PostgreSQL

Use UUID primary keys when possible.

### Required Tables

- users
- subscriptions
- categories
- reminders
- payment_records

### Rules

- Add `created_at`.
- Add `updated_at`.
- Use soft delete when appropriate.
- Add indexes for high-frequency queries.
- Add `user_id` isolation for all user-owned resources.

---

## Subscription Rules

Each subscription should support:

- name
- amount
- currency
- billing_cycle
- next_payment_date
- category
- reminder
- status
- notes

Supported billing cycles:

- monthly
- yearly
- custom

---

## Frontend Rules

### Stack

- Vue3
- Vite
- TypeScript
- Pinia
- Element Plus

### UI Principles

- Modern SaaS style.
- Responsive design.
- Dark mode support.
- Card-based layout.
- Mobile-friendly.
- Avoid admin-template style UI.
- Focus on usability and clean spacing.

### Frontend Structure

```text
apps/web/src
├── api
├── assets
├── components
├── composables
├── layouts
├── pages
├── router
├── stores
├── styles
├── types
└── utils
```

---

## Mobile Rules

### Stack

- Flutter
- Riverpod
- Dio
- GoRouter

### Principles

- Reuse backend APIs.
- Keep mobile UI lightweight.
- Follow platform-native interaction patterns.
- Use responsive layouts.

---

## Internationalization Rules

- Support i18n architecture from the beginning.
- Default language: English.
- Secondary language: Chinese.
- Avoid hardcoded UI text.
- Frontend text must support translation keys.

---

## Docker Rules

- Every service must support Docker.
- Use docker compose for local development.
- PostgreSQL and Redis must run in containers.
- Do not hardcode ports.
- Use environment variables.
- Support one-command startup.

Required command:

```bash
docker compose up --build
```

---

## Logging Rules

Use Zap logger.

Rules:

- Structured logs.
- Log errors with context.
- Avoid excessive logging.
- Do not log sensitive data.

---

## Security Rules

- Never trust frontend input.
- Validate all request parameters.
- Sanitize user-generated content.
- Protect authenticated routes.
- Use rate limiting when appropriate.
- Do not expose internal stack traces.

---

## Git Rules

### Main Branch

- master

### Feature Branches

- feature/*
- fix/*
- refactor/*

### Commit Convention

```text
feat:
fix:
docs:
refactor:
chore:
style:
test:
```

### Commit Message Detail

- Every commit message must include as many key details as possible.
- Use a concise subject plus a body that explains what changed, why it changed, and how it was verified.
- Historical commits should be rewritten or supplemented when their messages lack enough detail.

---

## Documentation Rules

Every major feature must include:

- API documentation
- Database changes
- Docker updates
- README updates if needed

Swagger documentation must stay synchronized with APIs.

---

## README Requirements

README must include:

- Project introduction
- Screenshots
- Features
- Tech stack
- Quick start
- Docker deployment
- Environment variables
- Roadmap
- Contribution guide
- License

---

## Testing Rules

### Backend

```bash
go test ./...
```

### Frontend

```bash
npm run build
```

### Docker

```bash
docker compose up --build
```

All major features must pass build checks before commit.

---

## Performance Rules

- Avoid N+1 queries.
- Add indexes for frequent queries.
- Use pagination.
- Cache high-frequency statistics when appropriate.

---

## First Milestone: MVP

MVP goals:

- User register/login
- JWT authentication
- Subscription CRUD
- Category management
- Dashboard statistics
- Docker deployment
- Swagger API docs
- Responsive web UI
- Dark mode
- English i18n support

Do not implement advanced features before MVP is complete.

---

## Non-Goals For MVP

Do not implement yet:

- OAuth login
- AI recommendations
- Auto bank sync
- Payment gateway integration
- Real-time websocket sync
- Microservice architecture
- Kubernetes deployment

Keep MVP simple and maintainable.
