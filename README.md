# SubPilot

SubPilot is a modern open-source subscription management platform for tracking recurring payments, categories, reminders, and spending insights.

The project is planned as a self-hostable product with web, mobile, Docker deployment, multi-user synchronization, and internationalization support.

## Status

SubPilot is in early development. The current repository state is the project baseline for Phase 0.

## Planned Stack

Backend:

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- JWT
- Viper
- Zap
- Swagger

Frontend:

- Vue 3
- Vite
- TypeScript
- Pinia
- Element Plus
- Axios
- ECharts

Mobile:

- Flutter
- Riverpod
- Dio
- GoRouter

DevOps:

- Docker
- Docker Compose
- Nginx

## Repository Layout

```text
.
├── apps
│   ├── server
│   ├── web
│   └── mobile
├── deploy
│   ├── docker
│   ├── nginx
│   └── scripts
├── docs
├── screenshots
├── tasks
├── .env.example
├── AGENTS.md
└── README.md
```

## Local Development

Local development will be Docker-first.

The target startup command is:

```bash
docker compose up --build
```

This command will become available after the Docker and backend skeleton are implemented in Phase 1.

## Configuration

Copy `.env.example` to `.env` for local development:

```bash
cp .env.example .env
```

Update placeholder values before running services. Do not commit `.env`.

## Roadmap

- Phase 0: Project baseline
- Phase 1: Local infrastructure and backend skeleton
- Phase 2: Database foundation
- Phase 3: Authentication
- Phase 4: Category management
- Phase 5: Subscription management
- Phase 6: Reminders and payment records
- Phase 7: Dashboard APIs
- Phase 8: Web app foundation
- Phase 9: Web subscription workflows
- Phase 10: Web dashboard
- Phase 11: Deployment documentation and release prep
- Phase 12: Mobile MVP
- Phase 13: Hardening

See [docs/development-spec.md](docs/development-spec.md) for the detailed phase plan.

## Development Rules

Project rules are defined in [AGENTS.md](AGENTS.md).

## License

MIT
