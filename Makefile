include .env
export

APP_PATH=./cmd/app
BIN_PATH=./cmd/bin
PROVIDERS_PATH=./internal/container

.PHONY: build
build:
	go build -o $(BIN_PATH)/scraper $(APP_PATH)

.PHONY: run
run:
	docker-compose up -d redis;
	(cd $(BIN_PATH) && ./scraper);

.PHONY: clean
clean:
	if [ -f $(BIN_PATH)/scraper ]; then rm $(BIN_PATH)/scraper; fi

.PHONY: all
all : clean build run

.PHONY: up
up:
	docker-compose up --build

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
	go test -v ./tests/... -tags=integration -race -cover -count=1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...

.PHONY: wire
wire:
	(cd $(PROVIDERS_PATH) && go run github.com/google/wire/cmd/wire@latest)

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