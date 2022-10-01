include .env
export

APP_PATH=./cmd/app
BIN_PATH=./cmd/bin
PROVIDERS_PATH=./internal/app/providers

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
	docker-compose up

.PHONY: compose-build
compose-build:
	docker-compose up --build

.PHONY: down
down:
	docker-compose down

.PHONY: test
test:
	go test -v ./... -cover -count=1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: wire
wire:
	(cd $(PROVIDERS_PATH) && go run github.com/google/wire/cmd/wire)
