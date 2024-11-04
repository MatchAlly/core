-- +goose up
CREATE DATABASE core;

-- +goose down
DROP DATABASE IF EXISTS core;