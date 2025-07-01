-- +goose up
CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mu NUMERIC(5, 2) NOT NULL DEFAULT 25.00,
    sigma NUMERIC(5, 2) NOT NULL DEFAULT 3.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (member_id, game_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_member_id ON ratings(member_id);
CREATE INDEX IF NOT EXISTS idx_ratings_game_id ON ratings(game_id);

CREATE TRIGGER update_ratings_updated_at
    BEFORE UPDATE ON ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- +goose down
DROP TRIGGER IF EXISTS update_ratings_updated_at ON ratings;
DROP INDEX IF EXISTS idx_ratings_member_id;
DROP INDEX IF EXISTS idx_ratings_game_id;

DROP TABLE IF EXISTS ratings;