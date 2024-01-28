CREATE SCHEMA IF NOT EXISTS metrics;

CREATE TABLE metrics.metrics_counter
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    value      BIGINT       NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE TABLE metrics.metrics_gauge
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)     NOT NULL,
    value      DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp
);

CREATE UNIQUE INDEX metrics_counter_name_idx ON metrics.metrics_counter (name);
CREATE UNIQUE INDEX metrics_gauge_name_idx ON metrics.metrics_gauge (name);