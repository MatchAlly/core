-- +goose up
CREATE TABLE games (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id)
);

CREATE INDEX idx_games_club_id ON games(club_id);

-- +goose down
DROP INDEX IF EXISTS idx_games_club_id;

DROP TABLE IF EXISTS games;