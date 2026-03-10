# 通用管理系统 (General Management System)

First runnable scaffold for a private, local-deployable modular-monolith.

## Stack
- Backend: Go, Gin, GORM, PostgreSQL, JWT
- Frontend: React, Vite, TypeScript, TanStack Query
- Infra (local): Docker Compose with Postgres + Redis + MinIO + Backend + Web

## Monorepo
- `apps/backend`
- `apps/web`
- `deploy/docker-compose.yml`
- `docs/architecture.md`

## Quick Start (Local Dev)
1. Copy env file:
   ```bash
   cp .env.example .env
   ```
2. Start infrastructure:
   ```bash
   docker compose --env-file .env -f deploy/docker-compose.yml up -d postgres redis minio
   ```
3. Backend:
   ```bash
   cd apps/backend
   go mod tidy
   go run ./cmd/server
   ```
4. Frontend (new terminal):
   ```bash
   cd apps/web
   npm install
   npm run dev
   ```
5. Open `http://localhost:5173`

Default seed login:
- Email: `admin@gms.local`
- Password: `admin123`

## Quick Start (Full Docker)
```bash
cp .env.example .env
make up
```

Endpoints:
- Web: `http://localhost:5174`
- API health: `http://localhost:8080/api/health`
- MinIO API: `http://localhost:9000`
- MinIO Console: `http://localhost:9001`

Shutdown:
```bash
make down
```

## API Scaffold
- `GET /api/health`
- `POST /api/auth/login`
- `GET /api/auth/me`
- `GET /api/pm/projects`
- `POST /api/pm/projects`
- `GET /api/pm/tasks?projectId=`
- `POST /api/pm/tasks`
- `PATCH /api/pm/tasks/:id/status`

Task status flow:
- `todo -> in_progress -> in_review -> done`

## Useful Commands
- `make dev` run backend + web locally
- `make build` build backend and frontend
- `make test` run backend tests
