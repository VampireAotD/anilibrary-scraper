# Anilibrary-scraper

Microservice for scraping anime data

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

With docker:

```shell
make up # docker-compose up --build
make filebeat # same as up, but with filebeat container
```

### Ports

* `8080` - **HTTP**
* `6379` - **Redis**
* `16686` - **Jaeger UI**

---

## Monitoring

Prometheus' metrics are sent to [monitoring](https://github.com/VampireAotD/anilibrary-monitoring) service

---

## Logs

Logs are written to file and can be sent to [elk](https://github.com/VampireAotD/anilibrary-elk) service

---

## Tests

```shell
make test # go test -v ./...
```