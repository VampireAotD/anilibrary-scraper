include .env
export

APP_PATH=./cmd/app

.PHONY: install-dependencies
install-dependencies:
	go install github.com/pressly/goose/v3/cmd/goose@latest;
	go install go.uber.org/mock/mockgen@latest;
	go install github.com/swaggo/swag/cmd/swag@latest;

.PHONY: install
install:
	@make install-dependencies
	@make up

.PHONY: up
up:
	docker-compose -f docker/docker-compose.yml up --build

.PHONY: filebeat
filebeat:
	docker-compose -f docker/docker-compose.yml --profile filebeat up --build

.PHONY: down
down:
	docker-compose -f docker/docker-compose.yml down

.PHONY: test
test:
	go test -v ./... -race -cover -count=1

.PHONY: integration-test
integration-test:
	go test -v ./test/... -tags=integration -race -cover -count=1

.PHONY: lint
lint:
	golangci-lint run

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