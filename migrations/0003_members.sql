-- +goose up
CREATE TYPE role AS ENUM ('none', 'observer', 'member', 'manager', 'admin', 'owner');

CREATE TABLE IF NOT EXISTS members (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    club_id UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role role DEFAULT 'none',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_members_id ON members(id);
CREATE INDEX IF NOT EXISTS idx_members_club_id ON members(club_id);

-- +goose down
DROP INDEX IF EXISTS idx_members_club_id;
DROP INDEX IF EXISTS idx_members_id;

DROP TABLE IF EXISTS members;
