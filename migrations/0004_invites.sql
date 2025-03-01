-- +goose Up
CREATE TABLE IF NOT EXISTS invites (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    initiator SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (club_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_invites_id ON invites(id);
CREATE INDEX IF NOT EXISTS idx_invites_club_id ON invites(club_id);
CREATE INDEX IF NOT EXISTS idx_invites_user_id ON invites(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_invites_user_id;
DROP INDEX IF EXISTS idx_invites_club_id;
DROP INDEX IF EXISTS idx_invites_id;

DROP TABLE IF EXISTS invites;