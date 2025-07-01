-- +goose up
CREATE TABLE IF NOT EXISTS statistics (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    wins INT NOT NULL DEFAULT 0,
    losses INT NOT NULL DEFAULT 0,
    draws INT NOT NULL DEFAULT 0,
    streak INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (member_id, game_id)
);

CREATE INDEX IF NOT EXISTS idx_statistics_member_id ON statistics(member_id);
CREATE INDEX IF NOT EXISTS idx_statistics_game_id ON statistics(game_id);

CREATE TRIGGER update_statistics_updated_at
    BEFORE UPDATE ON statistics
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- +goose down
DROP TRIGGER IF EXISTS update_statistics_updated_at ON statistics;
DROP INDEX IF EXISTS idx_statistics_game_id;
DROP INDEX IF EXISTS idx_statistics_member_id;

DROP TABLE IF EXISTS statistics;