include .env
export

APP_PATH=./app/cmd/app
BIN_PATH=./app/cmd/bin

.PHONY: build
build:
	go build -o $(BIN_PATH) $(APP_PATH)

.PHONY: run
run:
	(cd $(BIN_PATH) && ./app)

.PHONY: clean
clean:
	if [ -f $(BIN_PATH)/app ]; then rm $(BIN_PATH)/app; fi

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