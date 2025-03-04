-- +goose up
CREATE TYPE role AS ENUM ("none", "observer", "member", "manager", "admin", "owner");

CREATE TABLE IF NOT EXISTS club_members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT,
    role role DEFAULT "none",
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_club_members_id ON members(id);
CREATE INDEX IF NOT EXISTS idx_club_members_club_id ON members(club_id);

-- +goose down
DROP INDEX IF EXISTS idx_club_members_club_id;
DROP INDEX IF EXISTS idx_club_members_id;

DROP TABLE IF EXISTS club_members;
