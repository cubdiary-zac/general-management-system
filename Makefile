COMPOSE=docker compose --env-file .env -f deploy/docker-compose.yml

.PHONY: dev up down build test fmt

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

test:
	cd apps/backend && $(GO_CMD) test ./...

fmt:
	cd apps/backend && $(GOFMT_CMD) -w ./cmd ./internal
