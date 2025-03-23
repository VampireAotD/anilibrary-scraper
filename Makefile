include .env
export

compose := $(shell command -v docker-compose || echo docker compose)

.PHONY: help install-dependencies install up filebeat down test integration-test lint swag \
        swag-fmt clickhouse-migrate clickhouse-migrate-rollback

help:
	@printf "Usage: make <command>\n"
	@grep -hE '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-dependencies: ## Install dependencies to create mocks, OpenAPI specs, manage database migrations, etc.
	go install github.com/pressly/goose/v3/cmd/goose@v3.24.1;
	go install go.uber.org/mock/mockgen@v0.5.0;
	go install github.com/swaggo/swag/cmd/swag@v1.16.4;
	go install github.com/evilmartians/lefthook@v1.11.3;

install: ## Install dependencies and build application.
	@make install-dependencies
	@make up

up: ## Start application.
	$(compose) -f docker/compose.yml up --build

filebeat: ## Start application with Filebeat.
	$(compose) -f docker/compose.yml --profile filebeat up --build

down: ## Stop application.
	$(compose) -f docker/compose.yml down --remove-orphans

generate: ## Generate mocks.
	go generate ./...

test: ## Run unit tests.
	go test -v ./... -race -cover -count=1

integration-test: ## Run integration tests.
	go test -v ./test/... -tags=integration -race -cover -count=1

lint: ## Run linters.
	golangci-lint run

swag: ## Generate OpenAPI specification.
	swag init -g cmd/app/main.go

swag-fmt: ## Format OpenAPI specification.
	swag fmt

clickhouse-migrate: ## Run ClickHouse migrations.
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} up

clickhouse-migrate-rollback: ## Rollback ClickHouse migrations.
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} reset