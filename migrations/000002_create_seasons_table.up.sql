CREATE TABLE IF NOT EXISTS seasons (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    creator bigserial NOT NULL,
    year integer NOT NULL,
    funds integer,
    version integer NOT NULL DEFAULT 1
);