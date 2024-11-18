-- +goose up
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    hash TEXT NOT NULL,
    created_at TIMESTAMPZ WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clubs (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    display_name TEXT,
    role TEXT NOT NULL DEFAULT 'MEMBER',
    created_at TIMESTAMPZ WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_role CHECK (role IN ('ADMIN', 'MANAGER', 'MEMBER')),
    FOREIGN KEY (club_id) REFERENCES clubs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
CREATE INDEX IF NOT EXISTS idx_clubs_id ON clubs(id);
CREATE INDEX IF NOT EXISTS idx_members_id ON members(id);
CREATE INDEX IF NOT EXISTS idx_members_club_id ON members(club_id);

-- +goose down
DROP INDEX IF EXISTS idx_members_club_id;
DROP INDEX IF EXISTS idx_members_id;
DROP INDEX IF EXISTS idx_clubs_id;
DROP INDEX IF EXISTS idx_users_id;

DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS clubs;
DROP TABLE IF EXISTS users;
