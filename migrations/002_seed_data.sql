-- Seed sample products
INSERT INTO products (id, name, type, description, base_premium, coverage_amount, min_age, max_age, features, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Asuransi Jiwa Premium', 'life', 'Perlindungan finansial untuk keluarga tercinta dengan santunan meninggal dunia dan nilai tunai', 500000, 1000000000, 18, 65, '["Santunan meninggal dunia", "Premi fleksibel", "Nilai tunai", "Perlindungan hingga 65 tahun"]', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'Asuransi Kesehatan Plus', 'health', 'Biaya perawatan medis dan rawat inap dengan coverage lengkap', 300000, 500000000, 18, 70, '["Rawat inap", "Rawat jalan", "Obat-obatan", "Checkup rutin"]', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'Asuransi Kendaraan Comprehensive', 'vehicle', 'Perlindungan mobil dan motor dari risiko all risk dan TLO', 2000000, 300000000, 21, 70, '["All risk", "TLO", "Tanggung jawab pihak ketiga", "Banjir & gempa"]', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'Asuransi Jiwa Basic', 'life', 'Perlindungan dasar dengan premi terjangkau', 200000, 500000000, 18, 60, '["Santunan meninggal", "Premi murah", "Proses cepat"]', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'Asuransi Kesehatan Basic', 'health', 'Coverage rawat inap dasar', 150000, 200000000, 18, 65, '["Rawat inap", "ICU", "Operasi"]', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed sample admin user
INSERT INTO users (id, email, password_hash, full_name, role, created_at, updated_at) VALUES
('650e8400-e29b-41d4-a716-446655440001', 'admin@insurance.com', '$2a$10$YourHashedPasswordHere', 'Admin User', 'admin', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Seed sample applications
INSERT INTO applications (id, product_id, applicant_name, applicant_email, applicant_phone, applicant_dob, applicant_address, applicant_occupation, sum_assured, payment_term, premium_amount, status, health_declaration, beneficiary_name, beneficiary_relation, created_at, updated_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'Budi Santoso', 'budi@example.com', '081234567890', '1985-05-15', 'Jl. Sudirman No. 123, Jakarta Selatan', 'Software Engineer', 500000000, 10, 5500000, 'submitted', '{"smoking": false, "chronic_diseases": [], "recent_surgery": false}', 'Siti Santoso', 'wife', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'Ani Wijaya', 'ani@example.com', '081234567891', '1990-08-20', 'Jl. Gatot Subroto No. 45, Jakarta Pusat', 'Marketing Manager', 300000000, 5, 1600000, 'under_review', '{"smoking": false, "chronic_diseases": ["diabetes"], "recent_surgery": false}', 'Budi Wijaya', 'husband', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 'Doni Prakoso', 'doni@example.com', '081234567892', '1988-12-10', 'Jl. Thamrin No. 67, Jakarta Pusat', 'Business Owner', 200000000, 3, 6200000, 'draft', '{"smoking": true, "chronic_diseases": [], "recent_surgery": false}', 'Rina Prakoso', 'wife', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
