CREATE TABLE IF NOT EXISTS bids (
    id bigserial PRIMARY KEY,
    ts timestamp NOT NULL,
    player bigserial NOT NULL,
    bidder bigserial NOT NULL,
    amount integer,
    version integer NOT NULL DEFAULT 1
);