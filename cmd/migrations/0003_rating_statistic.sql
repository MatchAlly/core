-- +goose up
CREATE TABLE IF NOT EXISTS member_statistics (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    member_id BIGINT NOT NULL UNIQUE,
    game_id BIGINT NOT NULL UNIQUE,
    wins INT NOT NULL DEFAULT 0,
    losses INT NOT NULL DEFAULT 0,
    draws INT NOT NULL DEFAULT 0,
    streak INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

CREATE TABLE IF NOT EXISTS member_ratings (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    member_id BIGINT NOT NULL UNIQUE,
    game_id BIGINT NOT NULL UNIQUE,
    value NUMERIC(5, 2) NOT NULL DEFAULT 0,
    deviation NUMERIC(5, 2) NOT NULL DEFAULT 0,
    volatility NUMERIC(5, 2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

CREATE INDEX IF NOT EXISTS idx_member_ratings_member_id ON member_ratings(member_id);
CREATE INDEX IF NOT EXISTS idx_member_ratings_game_id ON member_ratings(game_id);

CREATE INDEX IF NOT EXISTS idx_member_statistics_member_id ON member_statistics(member_id);
CREATE INDEX IF NOT EXISTS idx_member_statistics_game_id ON member_statistics(game_id);

-- +goose down
DROP INDEX IF EXISTS idx_member_ratings_member_id;
DROP INDEX IF EXISTS idx_member_ratings_game_id;

DROP INDEX IF EXISTS idx_member_statistics_game_id;
DROP INDEX IF EXISTS idx_member_statistics_game_id;

DROP TABLE IF EXISTS member_statistics;
DROP TABLE IF EXISTS member_ratings;