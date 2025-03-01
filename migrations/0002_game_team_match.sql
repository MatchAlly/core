-- +goose up
CREATE TABLE IF NOT EXISTS games (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gamemode (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mode SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (game_id, mode)
);

CREATE TABLE IF NOT EXISTS matches (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    game_id BIGINT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    mode SMALLINT DEFAULT 0,
    ranked BOOLEAN NOT NULL DEFAULT FALSE,
    sets TEXT[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id, mode) REFERENCES gamemode(game_id, mode)
);

CREATE TABLE IF NOT EXISTS teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS team_members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    UNIQUE (team_id, member_id)
);

CREATE TABLE IF NOT EXISTS match_teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    team_number BIGINT NOT NULL,
    UNIQUE (match_id, team_id)
);

CREATE INDEX IF NOT EXISTS idx_games_club_id ON games(club_id);
CREATE INDEX IF NOT EXISTS idx_games_gamemode_id ON gamemode(game_id);
CREATE INDEX IF NOT EXISTS idx_match_teams_match_id ON match_teams(match_id);
CREATE INDEX IF NOT EXISTS idx_match_teams_team_id ON match_teams(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_member_id ON team_members(member_id);

-- +goose down
DROP INDEX IF EXISTS idx_team_members_member_id;
DROP INDEX IF EXISTS idx_match_teams_team_id;
DROP INDEX IF EXISTS idx_match_teams_match_id;
DROP INDEX IF EXISTS idx_games_club_id;

DROP TABLE games;
DROP TABLE gamemode;
DROP TABLE match_teams;
DROP TABLE team_members;
DROP TABLE teams;
DROP TABLE matches;
