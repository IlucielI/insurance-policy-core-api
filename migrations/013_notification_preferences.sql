-- Notification preferences table
CREATE TABLE notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Email notification preferences
    promotional_emails BOOLEAN DEFAULT TRUE,
    policy_update_emails BOOLEAN DEFAULT TRUE,
    claim_notification_emails BOOLEAN DEFAULT TRUE,
    newsletter_emails BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_user_preferences UNIQUE(user_id)
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);

-- Insert default preferences for existing users
INSERT INTO notification_preferences (user_id, promotional_emails, policy_update_emails, claim_notification_emails, newsletter_emails)
SELECT id, TRUE, TRUE, TRUE, TRUE FROM users
ON CONFLICT (user_id) DO NOTHING;
