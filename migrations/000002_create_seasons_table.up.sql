CREATE TABLE IF NOT EXISTS seasons (
    id bigserial PRIMARY KEY,
    year integer NOT NULL,
    funds integer,
    version integer NOT NULL DEFAULT 1
);