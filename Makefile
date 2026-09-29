APP_NAME := app
CMD_PATH := ./cmd/api
SEEDER_PATH := ./cmd/seeder
WORKER_PATH := ./cmd/worker
BUILD_DIR := ./bin
MIGRATIONS_DIR := ./migrations

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Database configuration for migrations (reads from .env / environment with fallback defaults)
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= nextstore
DB_SSLMODE ?= disable

DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
GOOSE := $(shell which goose 2>/dev/null || which $(shell go env GOPATH)/bin/goose 2>/dev/null || echo "go run github.com/pressly/goose/v3/cmd/goose@latest")


.PHONY: all run build start test test-cover fmt vet lint tidy clean docker-up docker-down docker-logs \
	migrate-up migrate-down migrate-status migrate-reset migrate-create seed worker db-info help

all: help

## Run application
run:
	go run $(CMD_PATH)

## Run background worker
worker:
	go run $(WORKER_PATH)

## Build application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

## Run built binary
start: build
	$(BUILD_DIR)/$(APP_NAME)

## Run tests
test:
	go test -v ./...

## Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

## Format code
fmt:
	go fmt ./...

## Run go vet
vet:
	go vet ./...

## Run linter
lint:
	golangci-lint run

## Update dependencies
tidy:
	go mod tidy

## Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

## Development mode
dev:
	go run $(CMD_PATH)

## Start Docker services (PostgreSQL)
docker-up:
	docker compose up -d

## Stop Docker services
docker-down:
	docker compose down

## View Docker service logs
docker-logs:
	docker compose logs -f

## Run all pending Goose migrations
migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

## Rollback the last Goose migration
migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down

## Check current Goose migration status
migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) status

## Rollback all Goose migrations
migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) reset

## Create a new Goose migration: make migrate-create name=add_column
migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: 'name' is required. Usage: make migrate-create name=<migration_name>"; exit 1; fi
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

## Run database seeders
seed:
	go run $(SEEDER_PATH)

## Show database configuration parsed from .env
db-info:
	@echo "Database configuration from .env:"
	@echo "  Host:     $(DB_HOST)"
	@echo "  Port:     $(DB_PORT)"
	@echo "  User:     $(DB_USER)"
	@echo "  Database: $(DB_NAME)"
	@echo "  SSL Mode: $(DB_SSLMODE)"
	@echo "  DSN:      postgres://$(DB_USER):****@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

## Show available commands
help:
	@echo "Available commands:"
	@echo "  make run            - Run application"
	@echo "  make build          - Build application binary"
	@echo "  make start          - Build and run application"
	@echo "  make test           - Run unit and integration tests"
	@echo "  make test-cover     - Run tests with HTML coverage report"
	@echo "  make fmt            - Format code with go fmt"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make tidy           - Run go mod tidy"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make docker-up      - Start PostgreSQL container"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - Follow Docker container logs"
	@echo "  make migrate-up     - Run pending Goose migrations"
	@echo "  make migrate-down   - Rollback the last Goose migration"
	@echo "  make migrate-status - Check Goose migration status"
	@echo "  make migrate-reset  - Rollback all Goose migrations"
	@echo "  make migrate-create - Create a new migration file (e.g. make migrate-create name=add_items)"
	@echo "  make seed           - Run database seeders"
	@echo "  make worker         - Run background worker (Asynq)"
	@echo "  make db-info        - Display database configuration loaded from .env"