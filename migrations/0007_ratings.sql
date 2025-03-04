-- +goose up
CREATE TABLE IF NOT EXISTS member_ratings (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    value NUMERIC(5, 2) NOT NULL DEFAULT 0,
    deviation NUMERIC(5, 2) NOT NULL DEFAULT 0,
    volatility NUMERIC(5, 2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (member_id, game_id)
);

CREATE INDEX IF NOT EXISTS idx_member_ratings_member_id ON member_ratings(member_id);
CREATE INDEX IF NOT EXISTS idx_member_ratings_game_id ON member_ratings(game_id);

-- +goose down
DROP INDEX IF EXISTS idx_member_ratings_member_id;
DROP INDEX IF EXISTS idx_member_ratings_game_id;

DROP TABLE IF EXISTS member_ratings;