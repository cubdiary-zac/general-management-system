# Architecture Overview - 通用管理系统

## Goals
- Private and local-deployable system for small teams (~100 users)
- Fast UI interactions and simple operations
- Modular-monolith to keep complexity low while preserving clear module boundaries

## Monorepo Layout
- `apps/backend`: Go API server (Gin + GORM)
- `apps/web`: React + Vite frontend
- `deploy/`: Docker Compose and local deployment assets
- `docs/`: Architecture and design notes

## Modular-Monolith Pattern
Single deployable backend process, but code is organized by modules under explicit boundaries:
- `core`: auth, users, RBAC primitives, shared middleware
- `pm`: project/task APIs and business rules
- future modules: CRM / HR / FIN as isolated packages

Module rules:
- Keep module domain models and handlers scoped to module packages
- Cross-module calls happen through service interfaces, not direct DB coupling
- Shared infrastructure lives in `internal/` packages (config, DB, middleware)

## Plugin-Style Expansion Path (PM/CRM/HR/FIN)
Current runtime is static, but module code is shaped like plugin slices:
- each module exposes route registration function
- each module owns migrations/seeds for its entities
- feature flags can conditionally mount module routes

Later options:
- Enable/disable modules from config
- Separate module schemas/tables while retaining one runtime
- Extract a module to independent service only if scale/ownership demands it

## Data Layer
- PostgreSQL via GORM
- `AutoMigrate` for scaffold stage
- seed owner account on startup (idempotent)

## Auth & Access Control
- Local JWT auth
- role enum: `owner`, `admin`, `member`, `viewer`
- middleware-based RBAC gate as baseline placeholder

## Runtime
- Backend REST API under `/api`
- Frontend SPA with TanStack Query for selective data refresh
- Docker Compose runtime includes: postgres, redis, minio, backend, web
