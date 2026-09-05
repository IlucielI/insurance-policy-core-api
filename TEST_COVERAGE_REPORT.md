# Test Coverage Report

**Tanggal**: 2026-09-05
**Target**: 80%+ coverage untuk backend Go API
**Status**: Coverage saat ini **4.9%**

## Test Files yang Dibuat

1. **internal/middleware/auth_test.go** (225 baris) ✅
   - TestGenerateToken
   - TestAuthRequired (4 scenarios: valid, missing, invalid, expired token)
   - TestRequireRole (3 scenarios)
   - TestRequirePermission
   - TestGetUserFromContext
   - **Coverage: 46.3%**

2. **internal/middleware/audit_middleware_test.go** (63 baris) ✅
   - TestExtractAuditContext (3 scenarios)
   - **All tests passing**

3. **internal/repository/fraud_repository_test.go** (97 baris) ✅
   - TestFraudRepository_UpdateRiskScore (2 scenarios)
   - TestFraudRepository_GetHighRiskApplications (3 scenarios)
   - **All tests passing**

## Test yang Dihapus (Compile Errors)

Test yang dihapus karena mismatch dengan implementasi SQL asli:
- `internal/repository/billing_repository_test.go` - SQL mock tidak cocok dengan implementasi
- `internal/repository/user_repository_test.go` - Masalah NULL handling
- `internal/usecase/*_test.go` - Interface mocking terlalu kompleks
- `internal/delivery/http/*_test.go` - Handler test memerlukan refactor

## Coverage Breakdown

```
middleware:     46.3% ✅
repository:      1.7% ❌
cache:           0.0% ❌
Total:           4.9% ❌
```

## Masalah yang Dihadapi

1. **SQL Mock Mismatch**: sqlmock expectations tidak cocok dengan SQL queries di repository yang sebenarnya
2. **NULL Handling**: Field seperti `payment_method`, `payment_reference` di Invoice perlu sql.NullString
3. **Interface Complexity**: Banyak interface dengan method yang kompleks sulit di-mock
4. **Build Errors**: Beberapa package memiliki build errors (`covdata` tool missing)

## Rekomendasi untuk Mencapai 80%+

### Prioritas Tinggi (Critical Paths)
1. **Auth Flow**:
   - ✅ Middleware auth sudah 46.3%
   - ❌ Perlu: auth_usecase tests, auth_handler tests
   
2. **Payment Flow**:
   - ❌ billing_repository tests
   - ❌ billing_usecase tests
   - ❌ billing_handler tests
   
3. **Fraud Detection**:
   - ✅ fraud_repository tests sudah ada
   - ❌ fraud_usecase tests
   - ❌ fraud_handler tests

### Solusi Teknis

1. **Ganti SQL Mock dengan Testcontainers**: Gunakan PostgreSQL real container untuk integration tests
2. **Simplify Mocks**: Buat mock generators (mockery) untuk interface yang kompleks
3. **Refactor Nullable Fields**: Ubah field yang nullable jadi `sql.NullString`, `sql.NullInt64`, dll
4. **Fix Build Errors**: Update Go version atau fix `covdata` issue

## Files yang Sudah Dibuat

```bash
internal/middleware/auth_test.go                     ✅ 225 lines
internal/middleware/audit_middleware_test.go         ✅  63 lines
internal/repository/fraud_repository_test.go         ✅  97 lines
internal/infrastructure/cache/redis_client_test.go   ✅  56 lines (skipped)
```

**Total test code**: 441 baris

## Langkah Selanjutnya

Untuk mencapai 80%+ coverage, perlu:
1. Refactor domain models untuk NULL handling yang proper
2. Setup testcontainers untuk integration tests
3. Generate mocks dengan mockery untuk semua interfaces
4. Tulis ~2000+ baris test code tambahan untuk:
   - Semua repositories (user, billing, policy, claim, dll)
   - Semua usecases (auth, billing, fraud, policy, claim, dll)
   - Semua handlers (HTTP endpoints)
   
Estimasi: **2-3 hari kerja** untuk coverage 80%+ yang solid dengan proper integration tests.
