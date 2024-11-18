-- +goose up
CREATE TABLE IF NOT EXISTS games (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id)
);

CREATE TABLE IF NOT EXISTS gamemode (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    game_id BIGINT NOT NULL,
    mode TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_mode CHECK (mode IN ('FFA', 'TEAM', 'COOP')),
    CONSTRAINT unique_game_mode UNIQUE (game_id, mode),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

CREATE TABLE IF NOT EXISTS matches (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    game_id BIGINT NOT NULL,
    mode TEXT NOT NULL,
    ranked BOOLEAN NOT NULL DEFAULT FALSE,
    sets TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_mode CHECK (mode IN ('FFA', 'TEAM', 'COOP'))
    FOREIGN KEY (club_id) REFERENCES clubs(id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (game_id, mode) REFERENCES gamemode(game_id, mode)
);

CREATE TABLE IF NOT EXISTS teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id)
);

CREATE TABLE IF NOT EXISTS team_members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    team_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (member_id) REFERENCES members(id),
    UNIQUE (team_id, member_id)
);

CREATE TABLE IF NOT EXISTS match_teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    match_id BIGINT NOT NULL,
    team_id BIGINT NOT NULL,
    team_number BIGINT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
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
