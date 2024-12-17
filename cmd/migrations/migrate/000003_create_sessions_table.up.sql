CREATE TABLE IF NOT EXISTS sessions(
    id              UUID PRIMARY KEY,
    refresh_token   VARCHAR(255) NOT NULL UNIQUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at      TIMESTAMP NOT NULL,
    user_id         UUID REFERENCES users_auth(id)
)