-- +goose up
CREATE DATABASE IF NOT EXISTS core;

-- +goose down
DROP DATABASE IF EXISTS core;