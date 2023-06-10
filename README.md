# Anilibrary Scraper

:warning: **This is an experimental branch and does not considered stable**

Microservice for scraping anime data

[![tests](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml)
[![golangci-lint](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml)

---

## Config

Firstly, you need to fill up environmental variables with your values, for that you can either :

- Manually copy `.env.example` with all filled variables to `.env` and start project using `make up`.
- Fill required variables in `.env.example` and use `make install`

---

## Deployment

Compile into binary:

```shell
make all
```

Using docker:

```shell
make up # docker-compose up --build
make filebeat # same as up, but with filebeat container
```

If you are deploying this project for the first time, better use:

```shell
make install # will copy .env.example to .env and deploy app using docker
```

| Service     | Port  |
|-------------|-------|
| Application | 8080  |
| Redis       | 6379  |
| Jaeger UI   | 16686 |
| Clickhouse  | 9000  |
| ZooKeeper   | 2181  |
| Kafka       | 9092  |

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