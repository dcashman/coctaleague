CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    owner bigserial NOT NULL,
    season bigserial NOT NULL, -- Should match seasons table id
    spreadsheet_position integer, -- relative position in draft spreadsheet for the season
    version integer NOT NULL DEFAULT 1
);