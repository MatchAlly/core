-- +goose up
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE clubs (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE members (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    club_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    display_name VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (club_id) REFERENCES clubs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_users_id ON users(id);
CREATE INDEX idx_clubs_id ON clubs(id);
CREATE INDEX idx_members_id ON members(id);
CREATE INDEX idx_members_club_id ON members(club_id);

-- +goose down
DROP INDEX IF EXISTS idx_members_club_id;
DROP INDEX IF EXISTS idx_members_id;
DROP INDEX IF EXISTS idx_clubs_id;
DROP INDEX IF EXISTS idx_users_id;

DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS clubs;
DROP TABLE IF EXISTS users;