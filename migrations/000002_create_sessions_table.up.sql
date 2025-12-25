CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY
    , user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
    , created_at TIMESTAMP NOT NULL
    , expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
