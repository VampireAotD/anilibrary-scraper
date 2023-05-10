-- +goose Up
CREATE TABLE IF NOT EXISTS events_queue
(
    date timestamp,
    url  String
) ENGINE = Kafka('kafka:9092', 'scraper_topic', 'scraper', 'JSONEachRow');

-- +goose Down
DROP TABLE events_queue;