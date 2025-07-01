-- +goose up
CREATE TABLE IF NOT EXISTS clubs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_clubs_id ON clubs(id);

-- +goose down
DROP INDEX IF EXISTS idx_clubs_id;

DROP TABLE IF EXISTS clubs;
