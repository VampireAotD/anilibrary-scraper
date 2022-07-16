build_path=./app/cmd/app

.PHONY: build
build:
	go build -o $(build_path)/app $(build_path)

.PHONY: up
up:
	docker-compose up

.PHONY: compose-build
	docker-compuse up --build

.PHONY: down
down:
	docker-compose down