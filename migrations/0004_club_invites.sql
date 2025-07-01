-- +goose up
CREATE TABLE IF NOT EXISTS club_invites (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    club_id UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    initiator SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (club_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_club_invites_id ON club_invites(id);
CREATE INDEX IF NOT EXISTS idx_club_invites_club_id ON club_invites(club_id);
CREATE INDEX IF NOT EXISTS idx_club_invites_user_id ON club_invites(user_id);

-- +goose down
DROP INDEX IF EXISTS idx_club_invites_user_id;
DROP INDEX IF EXISTS idx_club_invites_club_id;
DROP INDEX IF EXISTS idx_club_invites_id;

DROP TABLE IF EXISTS club_invites;