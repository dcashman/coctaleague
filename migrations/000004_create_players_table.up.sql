CREATE TABLE IF NOT EXISTS players (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    organization text,
    player_type text,
    season bigserial NOT NULL,
    spreadsheet_position integer,
    espn_id integer,
    espn_predicted_points integer,
    espn_actual_points integer,
    version integer NOT NULL DEFAULT 1
);