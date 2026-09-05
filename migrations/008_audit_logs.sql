-- Audit logs table for tracking admin actions
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- approve_application, reject_application, status_change, email_sent, etc.
    entity_type VARCHAR(50) NOT NULL, -- application, claim, policy, product, user
    entity_id VARCHAR(255) NOT NULL,
    changes_json JSONB, -- stores before/after values
    ip_address VARCHAR(45), -- IPv4 or IPv6
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for efficient querying
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- Composite index for common filter combinations
CREATE INDEX idx_audit_logs_user_entity ON audit_logs(user_id, entity_type, created_at DESC);
CREATE INDEX idx_audit_logs_action_date ON audit_logs(action, created_at DESC);
