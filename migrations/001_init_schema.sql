-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Users table (for both customers and admins)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    role VARCHAR(20) NOT NULL DEFAULT 'customer', -- customer, admin, underwriter
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    category VARCHAR(50) NOT NULL, -- life, health, vehicle
    description TEXT,
    coverage_details JSONB,
    min_sum_assured BIGINT NOT NULL,
    max_sum_assured BIGINT NOT NULL,
    min_payment_term INT NOT NULL, -- in months
    max_payment_term INT NOT NULL,
    base_premium_rate DECIMAL(10, 6) NOT NULL, -- percentage
    age_factor JSONB, -- age-based multipliers
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active);

-- Applications table
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    
    -- Applicant data
    applicant_data JSONB NOT NULL, -- name, dob, address, ktp, etc.
    
    -- Policy details
    sum_assured BIGINT NOT NULL,
    payment_term INT NOT NULL, -- in months
    premium_amount BIGINT NOT NULL, -- in cents/rupiah smallest unit
    
    -- Health questions (optional, depends on product)
    health_questions JSONB,
    
    -- Status workflow
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, submitted, under_review, approved, rejected
    underwriter_id UUID REFERENCES users(id),
    underwriter_notes TEXT,
    rejection_reason TEXT,
    
    -- Timestamps
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_product_id ON applications(product_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_underwriter_id ON applications(underwriter_id);
CREATE INDEX idx_applications_submitted_at ON applications(submitted_at);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    
    -- Midtrans integration
    order_id VARCHAR(255) UNIQUE NOT NULL,
    midtrans_transaction_id VARCHAR(255),
    payment_type VARCHAR(50), -- credit_card, bank_transfer, gopay, etc.
    
    -- Amount
    gross_amount BIGINT NOT NULL,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, success, failed, expired
    
    -- Timestamps
    paid_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_application_id ON payments(application_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);

-- Chat sessions table (for RAG chatbot)
CREATE TABLE chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- nullable for anonymous users
    session_id VARCHAR(255) UNIQUE NOT NULL, -- client-generated or server-generated
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chat_sessions_user_id ON chat_sessions(user_id);
CREATE INDEX idx_chat_sessions_session_id ON chat_sessions(session_id);

-- Chat messages table
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    
    role VARCHAR(20) NOT NULL, -- user, assistant
    content TEXT NOT NULL,
    
    -- RAG context (documents retrieved for this message)
    context_docs JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_chat_session_id ON chat_messages(chat_session_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);

-- Product embeddings table (for RAG semantic search)
CREATE TABLE product_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    -- Chunk metadata
    chunk_type VARCHAR(50) NOT NULL, -- description, benefits, exclusions, faq
    chunk_text TEXT NOT NULL,
    
    -- Vector embedding (using pgvector)
    embedding vector(1024), -- dimension depends on embedding model (bge-m3 = 1024)
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_product_embeddings_product_id ON product_embeddings(product_id);
CREATE INDEX idx_product_embeddings_embedding ON product_embeddings USING ivfflat (embedding vector_cosine_ops);

-- Activity log table (audit trail for admin actions)
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    action VARCHAR(100) NOT NULL, -- application_status_changed, product_created, etc.
    entity_type VARCHAR(50) NOT NULL, -- application, product, user
    entity_id UUID NOT NULL,
    
    metadata JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_entity_type ON activity_logs(entity_type);
CREATE INDEX idx_activity_logs_entity_id ON activity_logs(entity_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
