# Anilibrary-scraper

Microservice from scraping anime data

[![tests](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml)
[![golangci-lint](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml)

---

## Config

Copy `.env.example` to `.env` and fill the values.

---

## Build
Compile into binary:

```shell
make all
```

Deploy locally with docker:

```shell
make up # docker-compose up --build
```

### Ports

* `8080` - **HTTP**
* `6379` - **Redis**
* `9090` - **Prometheus** 
* `3000` - **Grafana**
* `16686` - **Jaeger**

---

## Tests

```shell
make test # go test -v ./...
```