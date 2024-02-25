-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    date       timestamp,
    url        String,
    ip         String,
    user_agent String
) ENGINE = MergeTree()
      PRIMARY KEY (date, url)
      PARTITION BY toYYYYMM(date)
      ORDER BY (date, url);

-- +goose Down
DROP TABLE events;