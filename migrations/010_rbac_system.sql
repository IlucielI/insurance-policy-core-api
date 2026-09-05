-- RBAC System Migration
-- Create roles table and user_roles junction table

-- Roles table
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User roles junction table (many-to-many)
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Insert predefined roles
INSERT INTO roles (name, display_name, description, permissions) VALUES
(
    'super_admin',
    'Super Admin',
    'Full system access - can manage all modules and users',
    '["users.manage", "products.manage", "applications.manage", "policies.manage", "claims.manage", "billing.manage", "reports.view", "settings.manage"]'::jsonb
),
(
    'underwriter',
    'Underwriter',
    'Can review and approve applications, manage policies',
    '["applications.view", "applications.review", "applications.approve", "policies.view", "policies.manage", "products.view"]'::jsonb
),
(
    'claims_officer',
    'Claims Officer',
    'Can review and process insurance claims only',
    '["claims.view", "claims.review", "claims.approve", "claims.process", "policies.view"]'::jsonb
),
(
    'finance',
    'Finance Officer',
    'Can manage billing, invoices, and payment records',
    '["billing.view", "billing.manage", "invoices.view", "invoices.create", "payments.view", "reports.finance"]'::jsonb
),
(
    'customer',
    'Customer',
    'Regular customer with access to own policies and claims',
    '["policies.view_own", "claims.create", "claims.view_own", "billing.view_own"]'::jsonb
);

-- Migrate existing user roles to new system
-- Insert role assignments for existing users based on their current role column
INSERT INTO user_roles (user_id, role_id)
SELECT 
    u.id,
    r.id
FROM users u
CROSS JOIN roles r
WHERE u.role = r.name;

-- Add indexes for performance
CREATE INDEX idx_roles_name ON roles(name);

-- Add updated_at trigger for roles table
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Keep the old role column for backward compatibility (will be deprecated later)
-- Add comment to indicate it's deprecated
COMMENT ON COLUMN users.role IS 'DEPRECATED: Use user_roles table instead. Kept for backward compatibility.';
