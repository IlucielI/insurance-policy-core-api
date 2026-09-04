-- Seed sample products (matching actual schema)
INSERT INTO products (id, name, slug, category, description, coverage_details, min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, base_premium_rate, age_factor, is_active, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Asuransi Jiwa Premium', 'asuransi-jiwa-premium', 'life', 'Perlindungan finansial untuk keluarga tercinta dengan santunan meninggal dunia dan nilai tunai', '{"features": ["Santunan meninggal dunia", "Premi fleksibel", "Nilai tunai", "Perlindungan hingga 65 tahun"]}', 100000000, 2000000000, 12, 240, 0.0055, '{"18-30": 1.0, "31-40": 1.2, "41-50": 1.5, "51-65": 2.0}', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'Asuransi Kesehatan Plus', 'asuransi-kesehatan-plus', 'health', 'Biaya perawatan medis dan rawat inap dengan coverage lengkap', '{"features": ["Rawat inap", "Rawat jalan", "Obat-obatan", "Checkup rutin"]}', 50000000, 1000000000, 12, 120, 0.003, '{"18-30": 1.0, "31-40": 1.1, "41-50": 1.3, "51-70": 1.8}', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'Asuransi Kendaraan Comprehensive', 'asuransi-kendaraan-comprehensive', 'vehicle', 'Perlindungan mobil dan motor dari risiko all risk dan TLO', '{"features": ["All risk", "TLO", "Tanggung jawab pihak ketiga", "Banjir & gempa"]}', 50000000, 500000000, 12, 36, 0.062, '{"21-30": 1.0, "31-40": 0.9, "41-70": 1.1}', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed sample admin user  
INSERT INTO users (id, email, password_hash, full_name, role, is_verified, created_at, updated_at) VALUES
('650e8400-e29b-41d4-a716-446655440001', 'admin@insurance.com', '$2a$10$example', 'Admin User', 'admin', true, NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440002', 'user1@example.com', '$2a$10$example', 'Budi Santoso', 'customer', true, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Seed sample applications
INSERT INTO applications (id, user_id, product_id, applicant_data, sum_assured, payment_term, premium_amount, status, health_questions, submitted_at, created_at, updated_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', '{"name": "Budi Santoso", "email": "budi@example.com", "phone": "081234567890", "dob": "1985-05-15", "address": "Jl. Sudirman No. 123, Jakarta", "occupation": "Software Engineer"}', 500000000, 120, 33000000, 'submitted', '{"smoking": false, "chronic_diseases": []}', NOW(), NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', '{"name": "Ani Wijaya", "email": "ani@example.com", "phone": "081234567891", "dob": "1990-08-20", "address": "Jl. Gatot Subroto No. 45, Jakarta", "occupation": "Manager"}', 300000000, 60, 9900000, 'under_review', '{"smoking": false, "chronic_diseases": ["diabetes"]}', NOW(), NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', '{"name": "Doni Prakoso", "email": "doni@example.com", "phone": "081234567892", "dob": "1988-12-10", "address": "Jl. Thamrin No. 67, Jakarta", "occupation": "Business Owner"}', 200000000, 36, 44640000, 'draft', '{"smoking": true}', NULL, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
