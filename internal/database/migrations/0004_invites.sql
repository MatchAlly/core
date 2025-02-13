-- +goose Up
CREATE TABLE IF NOT EXISTS invites (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    initiator SMALLINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
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