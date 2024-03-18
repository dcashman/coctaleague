CREATE TABLE IF NOT EXISTS seasons (
    id bigserial PRIMARY KEY,
    spreadsheet_id text,  --identifier for the particular spreadsheet on Google Sheets.
    name text NOT NULL,
    creator bigserial NOT NULL,
    year integer NOT NULL,
    funds integer,
    version integer NOT NULL DEFAULT 1
);