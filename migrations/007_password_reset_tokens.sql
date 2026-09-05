-- Password reset tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_token (token),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
);

COMMENT ON TABLE password_reset_tokens IS 'Stores password reset tokens with 1-hour expiry';
COMMENT ON COLUMN password_reset_tokens.token IS 'Unique token sent via email';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Token expires 1 hour after creation';
COMMENT ON COLUMN password_reset_tokens.used IS 'Marks if token has been used';
