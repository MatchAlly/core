-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    managed_organization_ids BIGINT[], 
    total_managed_users BIGINT DEFAULT 0,
    tier SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_id ON subscriptions(id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_subscriptions_id;
DROP INDEX IF EXISTS idx_subscriptions_user_id;

DROP TABLE IF EXISTS subscriptions;