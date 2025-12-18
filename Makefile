.PHONY: run dev build install docs help migrate-create migrate-up migrate-down migrate-status migrate-reset

GO ?= go
SWAG ?= $(shell go env GOPATH)/bin/swag
AIR ?= $(shell go env GOPATH)/bin/air
PKG := ./cmd/api/main.go
BIN := go-tasks

# Separate migrations directories by database dialect
MIGRATION_DIR := ./migrations/postgres
MIGRATION_SQLITE_DIR := ./migrations/sqlite

# Auto-load environment variables from .env if present
ENV_FILE := .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
# Export all keys from .env to child processes
export $(shell sed -n 's/^\([A-Za-z0-9_][A-Za-z0-9_]*\)=.*/\1/p' $(ENV_FILE))
endif

help:
	@echo "Available targets:"
	@echo "  dev            - run with hot reloading using Air"
	@echo "  debug          - run with hot reloading and debugging (Delve on :2345)"
	@echo "  debug-headless - run headless debugger without hot reload"
	@echo "  run            - go run ./cmd/api/main.go (basic run without hot reload)"
	@echo "  build          - go build -o go-tasks ./cmd/api/main.go"
	@echo "  install        - go install ./cmd/api/main.go"
	@echo "  docs           - swag init -g ./cmd/api/main.go -o ./docs"
	@echo "  test           - run unit tests"
	@echo "  test-race      - run tests with race detector"
	@echo "  test-cover     - run tests with coverage summary"
	@echo "  test-cover-html- run tests and open HTML coverage report"
	@echo "  secret         - generate a strong JWT secret (Go CLI)"
	@echo "  migrate-create         - create new Postgres migration: make migrate-create name=add_table"
	@echo "  migrate-create-sqlite  - create new SQLite migration: make migrate-create-sqlite name=add_table"
	@echo "  migrate-up     - apply all pending migrations"
	@echo "  migrate-down   - rollback a single migration"
	@echo "  migrate-status - show migration status"
	@echo "  migrate-reset  - rollback all migrations"
	@echo "  migrate-up-sqlite     - apply all migrations using sqlite3 driver"
	@echo "  migrate-down-sqlite   - rollback a single migration using sqlite3 driver"
	@echo "  migrate-status-sqlite - show migration status using sqlite3 driver"
	@echo "  migrate-reset-sqlite  - rollback all migrations using sqlite3 driver"

dev:
	@echo "Starting development server with hot reloading..."
	@command -v air >/dev/null 2>&1 || { echo "Installing Air..."; $(GO) install github.com/air-verse/air@latest; }
	$(AIR)

debug:
	@echo "Starting development server with hot reloading and debugging..."
	@echo "Debugger will be available on :2345"
	@echo "Connect your IDE debugger to localhost:2345"
	@command -v air >/dev/null 2>&1 || { echo "Installing Air..."; $(GO) install github.com/air-verse/air@latest; }
	@command -v dlv >/dev/null 2>&1 || { echo "Installing Delve..."; $(GO) install github.com/go-delve/delve/cmd/dlv@latest; }
	@mkdir -p tmp
	air -c .air-debug.toml

debug-headless:
	@echo "Starting headless debugger without hot reload..."
	@echo "Debugger will be available on :2345"
	@command -v dlv >/dev/null 2>&1 || { echo "Installing Delve..."; $(GO) install github.com/go-delve/delve/cmd/dlv@latest; }
	dlv debug ./cmd/api/main.go --headless --listen=:2345 --api-version=2 --accept-multiclient

debug-stop:
	@echo "Stopping debug processes..."
	@pkill -f "dlv exec" || true
	@pkill -f "main-debug" || true

run:
	$(GO) run $(PKG)

build:
	$(GO) build -o $(BIN) $(PKG)

install:
	$(GO) install $(PKG)

docs:
	$(SWAG) init -g $(PKG) -o ./docs --parseDependency --parseInternal

# Testing
test:
	$(GO) test ./...

test-race:
	$(GO) test ./... -race

test-cover:
	$(GO) test ./... -cover

test-cover-html:
	$(GO) test ./... -coverprofile=coverage.out && $(GO) tool cover -html=coverage.out

# Goose-based SQL migrations using DATABASE_URL (postgres DSN)
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=your_migration_name"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_DIR) create $(name) sql

migrate-create-sqlite:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create-sqlite name=your_migration_name"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_SQLITE_DIR) create $(name) sql

migrate-up:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_DIR) postgres "$$DATABASE_URL" up

migrate-down:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_DIR) postgres "$$DATABASE_URL" down

migrate-status:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_DIR) postgres "$$DATABASE_URL" status

migrate-reset:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_DIR) postgres "$$DATABASE_URL" reset

# SQLite3 variants (require DATABASE_URL set to sqlite file path, e.g. ./db/go-tasks.db)
migrate-up-sqlite:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_SQLITE_DIR) sqlite3 "$$DATABASE_URL" up

migrate-down-sqlite:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_SQLITE_DIR) sqlite3 "$$DATABASE_URL" down

migrate-status-sqlite:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_SQLITE_DIR) sqlite3 "$$DATABASE_URL" status

migrate-reset-sqlite:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL env var is required"; exit 2; fi
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATION_SQLITE_DIR) sqlite3 "$$DATABASE_URL" reset

# Generate a strong JWT secret using Go CLI (defaults: 48 bytes, base64url)
secret:
	$(GO) run ./cmd/secret/main.go

# Examples:
#  make secret                  # base64url, 48 bytes
#  make secret ARGS="-bytes=64" # increase entropy
#  make secret ARGS="-format=hex" # hex output
.PHONY: secret