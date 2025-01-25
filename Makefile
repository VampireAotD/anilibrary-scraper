include .env
export

compose := $(shell command -v docker-compose || echo docker compose)

.PHONY:help
help:
	@printf "\nUsage: make <command>\n"
	@grep -hE '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install-dependencies
install-dependencies: ## Install dependencies to create mocks, OpenAPI specs, manage database migrations, etc.
	go install github.com/pressly/goose/v3/cmd/goose@latest;
	go install go.uber.org/mock/mockgen@latest;
	go install github.com/swaggo/swag/cmd/swag@latest;
	go install github.com/evilmartians/lefthook@latest;

.PHONY: install
install: ## Install dependencies and build application.
	@make install-dependencies
	@make up

.PHONY: up
up: ## Start application.
	$(compose) -f docker/compose.yml up --build

.PHONY: filebeat
filebeat: ## Start application with Filebeat.
	$(compose) -f docker/compose.yml --profile filebeat up --build

.PHONY: down
down: ## Stop application.
	$(compose) -f docker/compose.yml down

.PHONY: test
test: ## Run unit tests.
	go test -v ./... -race -cover -count=1

.PHONY: integration-test
integration-test: ## Run integration tests.
	go test -v ./test/... -tags=integration -race -cover -count=1

.PHONY: lint
lint: ## Run linters.
	golangci-lint run

.PHONY: swag
swag: ## Generate OpenAPI specification.
	swag init -g cmd/app/main.go

.PHONY: swag-fmt
swag-fmt: ## Format OpenAPI specification.
	swag fmt

.PHONY: clickhouse-migrate
clickhouse-migrate: ## Run ClickHouse migrations.
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} up

.PHONY: clickhouse-migrate-rollback
clickhouse-migrate-rollback: ## Rollback ClickHouse migrations.
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} reset