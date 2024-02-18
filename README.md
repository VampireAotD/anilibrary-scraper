# Anilibrary Scraper

:warning: **Currently this project is under development and is not considered stable**

Microservice for scraping anime data.

[![tests](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/test.yml)
[![golangci-lint](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/VampireAotD/anilibrary-scraper/actions/workflows/golangci-lint.yml)

---

## Configuration

Firstly, to start working with scraper you need to fill up environmental variables with your values in **.env** file.
The recommended way to gain it is to make it from **.env.example** like this:

```shell
cp .env.example .env
```

because some variables already has a default value, or you can create in manually.

List of all variables and their description:

| Variable                        | Default value                | Description                                                                                                                                        |
|---------------------------------|------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| **APP_ENV**                     | local                        | Specify current project environment. Also used in tracing.                                                                                         | 
| **APP_NAME**                    | anilibrary-scraper           | Custom name for service.                                                                                                                           | 
| **HTTP_ADDR**                   | 0.0.0.0                      | Host to run service on.                                                                                                                            | 
| **HTTP_MAIN_PORT**              | 8080                         | Port to run main server with endpoints.                                                                                                            | 
| **HTTP_MONITORING_PORT**        | 8081                         | Port to run secondary server with metrics.                                                                                                         | 
| **JWT_SECRET**                  |                              | Token to communicate with other Anilibrary services.                                                                                               | 
| **TIMEZONE**                    | Europe/Kiev                  | Specify current timezone in container.                                                                                                             | 
| **REDIS_VERSION**               | latest                       | Specify desired Redis version.                                                                                                                     | 
| **REDIS_ADDRESS**               | redis:6379                   | Specify Redis connection address.                                                                                                                  | 
| **REDIS_PASSWORD**              |                              | Specify Redis password.                                                                                                                            |
| **REDIS_POOL_TIMEOUT**          | 5s                           | Max wait time for a connection from the pool, preventing hangs when all connections are busy.                                                      |
| **OTEL_EXPORTER_OTLP_ENDPOINT** | http://localhost:4318        | Specify endpoint on where to send traces. By default traces are sent to [monitoring service](https://github.com/VampireAotD/anilibrary-monitoring) |
| **FILEBEAT_VERSION**            | 8.12.0                       | Specify Filebeat version.                                                                                                                          |
| **FILEBEAT_OUTPUT**             |                              | Specify url on where to send logs to. Logs can be visualized by Kibana in [ELK service](https://github.com/VampireAotD/anilibrary-elk).            |
| **FILEBEAT_USER**               |                              | Specify login for Filebeat user.                                                                                                                   |
| **FILEBEAT_PASSWORD**           |                              | Specify password for Filebeat user.                                                                                                                |
| **KAFKA_VERSION**               | latest                       | Specify desired Kafka image version.                                                                                                               |
| **KAFKA_ADDRESS**               | kafka:9092                   | Specify Kafka connection address.                                                                                                                  |
| **KAFKA_CLIENT_USERS**          | example                      | Specify Kafka username for client connection.                                                                                                      |
| **KAFKA_CLIENT_PASSWORDS**      |                              | Specify Kafka password for client connection.                                                                                                      |
| **KAFKA_INTER_BROKER_USER**     | kafka                        | Specify Kafka inter broker username for broker communication.                                                                                      |
| **KAFKA_INTER_BROKER_PASSWORD** |                              | Specify Kafka inter broker password for broker communication.                                                                                      |
| **KAFKA_TOPIC**                 |                              | Specify Kafka topic.                                                                                                                               |
| **KAFKA_PARTITION**             |                              | Specify Kafka partition.                                                                                                                           |
| **CLICKHOUSE_VERSION**          | latest                       | Specify desired Clickhouse version.                                                                                                                |
| **CLICKHOUSE_USER**             | example                      | Specify username for Clickhouse connection.                                                                                                        |
| **CLICKHOUSE_PASSWORD**         |                              | Specify password for Clickhouse connection.                                                                                                        |
| **CLICKHOUSE_KAFKA_USER**       | KAFKA_CLIENT_USERS value     | Specify Kafka client name for Kafka table engine.                                                                                                  |
| **CLICKHOUSE_KAFKA_PASSWORD**   | KAFKA_CLIENT_PASSWORDS value | Specify Kafka client password for Kafka table engine.                                                                                              |
| **CLICKHOUSE_ADDRESS**          |                              | Specify Clickhouse connection address.                                                                                                             |

---

## Deployment

After configuration you can start your work with scraper.
Firstly, make sure to install dependencies like **golangci-lint**, **swaggo**, **mockgen**, etc. You can do it manually
or using:

```shell
make install-dependencies
```

After that you can build and run the service using:

```shell
make up
```

or you can do both of this commands using:

```shell
make install
```

To run unit tests use:

```shell
make test
```

and to run integration test use:

```shell
make integration-test
```

To find more useful commands make sure to check Makefile.

List of used services and their ports:

| Service     | Port       |
|-------------|------------|
| Application | 8080       |
| Monitoring  | 8081       |
| Redis       | 6379       |
| Clickhouse  | 8123, 9005 |
| Kafka       | 9092       |