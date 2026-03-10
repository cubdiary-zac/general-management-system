COMPOSE=docker compose --env-file .env -f deploy/docker-compose.yml

.PHONY: dev up down build test test-unit test-blackbox test-regression test-smoke test-smoke-live ci-gate fmt

GO_BIN := $(shell command -v go 2>/dev/null)
ifeq ($(GO_BIN),)
GO_CMD = docker run --rm -v $(CURDIR)/apps/backend:/src -w /src golang:1.23 /usr/local/go/bin/go
GOFMT_CMD = docker run --rm -v $(CURDIR)/apps/backend:/src -w /src golang:1.23 /usr/local/go/bin/gofmt
else
GO_CMD = go
GOFMT_CMD = gofmt
endif

dev:
	@command -v go >/dev/null 2>&1 || (echo "make dev requires local Go installed" && exit 1)
	@set -a; [ -f .env ] && . ./.env; set +a; \
	( cd apps/backend && go run ./cmd/server ) & \
	( cd apps/web && npm run dev -- --host ) & \
	wait

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

build:
	cd apps/backend && $(GO_CMD) build ./...
	cd apps/web && npm run build

test: test-unit test-blackbox test-regression

test-unit:
	cd apps/backend && $(GO_CMD) test ./internal/auth ./internal/models

test-blackbox:
	cd apps/backend && $(GO_CMD) test ./internal/handlers -run Blackbox

test-regression:
	cd apps/backend && $(GO_CMD) test ./internal/handlers -run Regression

test-smoke:
	cd apps/backend && $(GO_CMD) test ./internal/handlers -run Smoke

test-smoke-live:
	./scripts/smoke_api.sh

ci-gate:
	./scripts/ci_gate.sh

fmt:
	cd apps/backend && $(GOFMT_CMD) -w ./cmd ./internal
