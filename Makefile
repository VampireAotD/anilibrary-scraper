include .env
export

APP_PATH=./app/cmd/app
BIN_PATH=./app/cmd/bin

.PHONY: build
build:
	go build -o $(BIN_PATH)/scraper $(APP_PATH)

.PHONY: run
run:
	(cd $(BIN_PATH) && ./scraper)

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
	go test -v ./... -tags test -count=1

.PHONY: lint
lint:
	golangci-lint run