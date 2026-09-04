-- Policies table (issued policies from approved applications)
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_number VARCHAR(50) UNIQUE NOT NULL,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    
    -- Policy details
    sum_assured BIGINT NOT NULL,
    premium_amount BIGINT NOT NULL,
    payment_frequency VARCHAR(20) NOT NULL, -- monthly, quarterly, annually
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, lapsed, surrendered, expired
    
    -- Dates
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    last_premium_paid_date DATE,
    next_premium_due_date DATE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_policies_user_id ON policies(user_id);
CREATE INDEX idx_policies_policy_number ON policies(policy_number);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_application_id ON policies(application_id);

-- Claims table
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_number VARCHAR(50) UNIQUE NOT NULL,
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Claim details
    claim_type VARCHAR(50) NOT NULL, -- death, maturity, health, accident, vehicle_damage
    claim_amount BIGINT NOT NULL,
    incident_date DATE NOT NULL,
    incident_description TEXT NOT NULL,
    
    -- Status workflow
    status VARCHAR(50) NOT NULL DEFAULT 'submitted', -- submitted, under_review, approved, rejected, paid
    reviewer_id UUID REFERENCES users(id),
    reviewer_notes TEXT,
    rejection_reason TEXT,
    approved_amount BIGINT,
    
    -- Timestamps
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    approved_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_user_id ON claims(user_id);
CREATE INDEX idx_claims_claim_number ON claims(claim_number);
CREATE INDEX idx_claims_status ON claims(status);

-- Claim documents table
CREATE TABLE claim_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    
    document_type VARCHAR(50) NOT NULL, -- medical_report, police_report, invoice, photo, other
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_claim_documents_claim_id ON claim_documents(claim_id);

-- Claim timeline/activity table
CREATE TABLE claim_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    
    action VARCHAR(100) NOT NULL, -- submitted, document_uploaded, status_changed, comment_added
    description TEXT NOT NULL,
    actor_id UUID REFERENCES users(id),
    actor_name VARCHAR(255),
    metadata JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_claim_timeline_claim_id ON claim_timeline(claim_id);
CREATE INDEX idx_claim_timeline_created_at ON claim_timeline(created_at);

-- Invoices table (for premium billing)
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Invoice details
    amount BIGINT NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, paid, overdue, cancelled
    
    -- Payment details
    paid_amount BIGINT DEFAULT 0,
    paid_at TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(50), -- credit_card, bank_transfer, gopay, etc.
    payment_reference VARCHAR(255),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_invoices_policy_id ON invoices(policy_id);
CREATE INDEX idx_invoices_user_id ON invoices(user_id);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);

-- Documents table (general policy documents)
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID REFERENCES policies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    document_type VARCHAR(50) NOT NULL, -- policy_certificate, endorsement, receipt, notice, other
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_documents_policy_id ON documents(policy_id);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_document_type ON documents(document_type);

-- Policy endorsements table (policy changes/amendments)
CREATE TABLE policy_endorsements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endorsement_number VARCHAR(50) UNIQUE NOT NULL,
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    
    endorsement_type VARCHAR(50) NOT NULL, -- coverage_change, beneficiary_change, premium_adjustment, other
    description TEXT NOT NULL,
    effective_date DATE NOT NULL,
    
    old_values JSONB,
    new_values JSONB,
    
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_policy_endorsements_policy_id ON policy_endorsements(policy_id);
CREATE INDEX idx_policy_endorsements_status ON policy_endorsements(status);
