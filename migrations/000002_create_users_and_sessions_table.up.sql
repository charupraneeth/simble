CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    github_id   BIGINT UNIQUE NOT NULL,
    email       TEXT,
    username    TEXT NOT NULL,
    avatar_url  TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    token       TEXT PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

ALTER TABLE sites ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;