# Makefile

.PHONY: help

# Project name drives the docker-compose network prefix (defaults to directory name)
COMPOSE_PROJECT_NAME ?= $(notdir $(CURDIR))
MIGRATE_NETWORK      := $(COMPOSE_PROJECT_NAME)_auth-network
MIGRATIONS_DIR       := $(shell pwd)/db/migrations/sql
DB_MIGRATE_URL       := cockroachdb://root@auth-database:26257/defaultdb?sslmode=disable

# Force BuildKit and serialize builds to avoid OOM on low-RAM machines
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-25s %s\n", $$1, $$2}'

## Development Commands

compose-build: ## Build Docker images
	docker compose build

compose-with-debug: compose-build ## Start with debug logs
	@echo "Starting in the debug mode for container"
	@docker compose up

compose-without-app: compose-build ## Start without app container
	@echo "Starting in the debug mode for container"
	@docker compose up --scale app=0 -d

compose-up: compose-build ## Start all containers in background
	@docker compose up -d

compose-stop: ## Stop all containers
	@echo "stopping docker compose in background"
	@docker compose down

compose-clean: compose-stop ## Stop and remove containers
	docker compose rm -f

compose-build-no-cache: ## Build without cache
	docker compose build --no-cache

compose-migrate-debug: compose-build ## Build, run SQL migrations against DB, then start all services with debug logs
	@echo "Tearing down any existing containers and volumes for a clean state..."
	@docker compose down -v 2>/dev/null || true
	@echo "Starting infrastructure (db, redis, rabbitmq, mailpit)..."
	@docker compose up -d db redis rabbitmq mailpit
	@echo "Waiting for database to be healthy..."
	@until [ "$$(docker inspect --format='{{.State.Health.Status}}' auth-database 2>/dev/null)" = "healthy" ]; do \
		printf '.'; sleep 2; \
	done
	@echo " Database is healthy."
	@echo "Running SQL migrations from $(MIGRATIONS_DIR)..."
	docker run --rm \
		--network $(MIGRATE_NETWORK) \
		-v $(MIGRATIONS_DIR):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_MIGRATE_URL)" \
		up
	@echo "Starting all services in debug mode (logs visible)..."
	docker compose up

## Testing Commands

test-setup: ## Setup test environment (first time only)
	@cd test-suite && chmod +x setup.sh && ./setup.sh

test-isolated: ## Run tests with isolated database environment
	@cd test-suite && chmod +x run_tests.sh && ./run_tests.sh

test-quick: ## Run tests using existing containers (faster)
	@cd test-suite && go test -v -timeout 3m ./...

test-start: ## Start test containers without running tests
	@cd test-suite && docker compose -f test-suite/docker-compose.test.yml up -d

test-stop: ## Stop test containers
	@cd test-suite && docker compose -f test-suite/docker-compose.test.yml down

test-clean: ## Remove test containers and volumes
	@cd test-suite && docker compose -f test-suite/docker-compose.test.yml down -v

test-logs: ## Show test container logs
	@cd test-suite && docker compose -f test-suite/docker-compose.test.yml logs -f

test-status: ## Show test container status
	@cd test-suite && docker compose -f test-suite/docker-compose.test.yml ps

.DEFAULT_GOAL := help