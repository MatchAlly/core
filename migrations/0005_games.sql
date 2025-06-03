-- +goose up
CREATE TABLE IF NOT EXISTS games (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS game_modes (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mode SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (game_id, mode)
);

CREATE INDEX IF NOT EXISTS idx_games_club_id ON games(club_id);
CREATE INDEX IF NOT EXISTS idx_game_modes_game_id ON game_modes(game_id);

-- +goose down
DROP INDEX IF EXISTS idx_games_club_id;
DROP INDEX IF EXISTS idx_game_modes_game_id;

DROP TABLE game_modes;
DROP TABLE games;
