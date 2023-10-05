include .env
export

APP_PATH=./cmd/app

.PHONY: install-dependencies
install-dependencies:
	go install github.com/pressly/goose/v3/cmd/goose@latest;
	go install go.uber.org/mock/mockgen@latest;
	go install github.com/swaggo/swag/cmd/swag@latest;

.PHONY: up
up:
	docker-compose up --build

.PHONY: install
install:
	@cp .env.example .env
	@make install-dependencies
	@make up

.PHONY: filebeat
filebeat:
	docker-compose --profile filebeat up

.PHONY: down
down:
	docker-compose down

.PHONY: test
test:
	go test -v ./... -race -cover -count=1

.PHONY: integration-test
integration-test:
	go test -v ./test/... -tags=integration -race -cover -count=1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...

.PHONY: swag
swag:
	swag init -g $(APP_PATH)/main.go

.PHONY: swag-fmt
swag-fmt:
	swag fmt

.PHONY: clickhouse-migrate
clickhouse-migrate:
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} up

.PHONY: clickhouse-migrate-rollback
clickhouse-migrate-rollback:
	cd ./migrations/clickhouse && goose clickhouse ${CLICKHOUSE_ADDRESS} reset