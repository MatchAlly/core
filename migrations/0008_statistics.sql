-- +goose up
CREATE TABLE IF NOT EXISTS member_statistics (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    wins INT NOT NULL DEFAULT 0,
    losses INT NOT NULL DEFAULT 0,
    draws INT NOT NULL DEFAULT 0,
    streak INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (member_id, game_id)
);

CREATE INDEX IF NOT EXISTS idx_member_statistics_member_id ON member_statistics(member_id);
CREATE INDEX IF NOT EXISTS idx_member_statistics_game_id ON member_statistics(game_id);

-- +goose down
DROP INDEX IF EXISTS idx_member_statistics_game_id;
DROP INDEX IF EXISTS idx_member_statistics_member_id;

DROP TABLE IF EXISTS member_statistics;