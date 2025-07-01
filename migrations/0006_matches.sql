-- +goose up
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    club_id UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mode SMALLINT DEFAULT 0,
    ranked BOOLEAN NOT NULL DEFAULT FALSE,
    sets TEXT[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id, mode) REFERENCES game_modes(game_id, mode)
);

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    club_id UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS team_members (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    UNIQUE (team_id, member_id)
);

CREATE TABLE IF NOT EXISTS match_teams (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    team_number BIGINT NOT NULL,
    UNIQUE (match_id, team_id)
);

CREATE INDEX IF NOT EXISTS idx_match_teams_match_id ON match_teams(match_id);
CREATE INDEX IF NOT EXISTS idx_match_teams_team_id ON match_teams(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_member_id ON team_members(member_id);
CREATE INDEX IF NOT EXISTS idx_matches_club_id ON matches(club_id);
CREATE INDEX IF NOT EXISTS idx_matches_game_id ON matches(game_id);
CREATE INDEX IF NOT EXISTS idx_matches_created_at ON matches(created_at);

-- +goose down
DROP INDEX IF EXISTS idx_matches_created_at;
DROP INDEX IF EXISTS idx_matches_game_id;
DROP INDEX IF EXISTS idx_matches_club_id;
DROP INDEX IF EXISTS idx_team_members_member_id;
DROP INDEX IF EXISTS idx_match_teams_team_id;
DROP INDEX IF EXISTS idx_match_teams_match_id;

DROP TABLE IF EXISTS match_teams;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS matches;
