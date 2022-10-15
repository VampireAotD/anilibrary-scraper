# Anilibrary-scraper

Microservice from scraping anime data

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

---

## Tests

```shell
make test # go test -v ./...
```