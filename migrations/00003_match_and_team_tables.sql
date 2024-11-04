-- +goose up
CREATE TABLE matches (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    game_id BIGINT NOT NULL,
    result CHAR(1),
    sets TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

CREATE TABLE teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id)
);

CREATE TABLE team_members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    team_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (member_id) REFERENCES members(id),
    UNIQUE (team_id, member_id)
);

CREATE TABLE match_teams (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    match_id BIGINT NOT NULL,
    team_id BIGINT NOT NULL,
    team_number BIGINT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    UNIQUE (match_id, team_id)
);

CREATE INDEX idx_match_teams_match_id ON match_teams(match_id);
CREATE INDEX idx_match_teams_team_id ON match_teams(team_id);
CREATE INDEX idx_team_members_member_id ON team_members(member_id);

-- +goose down
DROP INDEX IF EXISTS idx_team_members_member_id;
DROP INDEX IF EXISTS idx_match_teams_team_id;
DROP INDEX IF EXISTS idx_match_teams_match_id;

DROP TABLE match_teams;
DROP TABLE team_members;
DROP TABLE teams;
DROP TABLE matches;
