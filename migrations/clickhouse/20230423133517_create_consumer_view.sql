-- +goose Up
CREATE MATERIALIZED VIEW events_consumer TO events
AS SELECT toDateTime(date) as date, url
FROM events_queue;

-- +goose Down
DROP VIEW events_consumer;