# Insurance Policy Core API - OpenAPI Documentation

## 📚 Dokumentasi API Lengkap

API komprehensif untuk sistem manajemen asuransi dengan 50+ endpoint yang mencakup autentikasi, produk, aplikasi, polis, klaim, billing, dokumen, notifikasi, fraud detection, OCR, dan analytics.

## 🚀 Akses Swagger UI

**Swagger UI tersedia di:** `http://localhost:8080/api/docs/index.html`

Setelah server berjalan, buka browser dan akses URL di atas untuk dokumentasi interaktif lengkap dengan:
- ✅ Deskripsi setiap endpoint
- ✅ Request/response schema
- ✅ Try-it-out functionality
- ✅ Authentication requirements
- ✅ RBAC role documentation
- ✅ Example values

## 📖 File Dokumentasi

1. **`docs/swagger.yaml`** - OpenAPI 3.0 specification (425 lines)
2. **`docs/swagger.json`** - JSON format (600 lines)
3. **`docs/docs.go`** - Go embedded docs

## 🔐 Authentication

API menggunakan dua metode autentikasi:

### 1. Bearer Token (Recommended)
```bash
Authorization: Bearer <jwt-token>
```

### 2. Cookie Auth
```bash
Cookie: auth_token=<jwt-token>
```

### Mendapatkan Token

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...},
  "roles": ["customer"]
}
```

## 🎯 RBAC Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| **`customer`** | User biasa | Akses personal data, aplikasi, polis, klaim, billing |
| **`underwriter`** | Underwriter | Approve aplikasi, risk analysis, fraud detection |
| **`claims_officer`** | Claims Officer | Manajemen klaim, approve/reject claims |
| **`finance`** | Finance Officer | Billing, invoice management, payment tracking |
| **`super_admin`** | Super Admin | Full access ke semua endpoint |

## 📋 Endpoint Summary

### Authentication & Users (7 endpoints)
- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login dan dapatkan JWT token
- `POST /auth/logout` - Logout user
- `GET /auth/me` - Get current user info
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password dengan token
- `GET /admin/users` - List all users (admin) 🔒

### Products (4 endpoints)
- `GET /products` - List produk asuransi
- `GET /products/{id}` - Detail produk
- `POST /products/{id}/calculate-premium` - Hitung premi
- `POST /products/search` - Semantic search produk

### Applications (5 endpoints)
- `POST /applications` - Buat aplikasi baru
- `GET /applications/{id}` - Detail aplikasi
- `GET /admin/applications` - List aplikasi (admin) 🔒
- `PUT /admin/applications/{id}/status` - Update status 🔒
- `POST /admin/applications/bulk-status` - Bulk update 🔒

### Policies (4 endpoints)
- `GET /policies` - List polis user 🔒
- `GET /policies/{id}` - Detail polis
- `POST /policies/{id}/endorse` - Endorse polis
- `POST /policies/{id}/renew` - Renew polis

### Claims (8 endpoints)
- `POST /claims` - Buat klaim baru
- `GET /claims/{id}` - Detail klaim
- `PUT /claims/{id}/documents` - Upload dokumen
- `GET /claims/{id}/timeline` - Timeline klaim
- `GET /admin/claims` - List klaim (admin) 🔒
- `PUT /admin/claims/{id}/status` - Update status 🔒
- `PUT /admin/claims/{id}/approve` - Approve klaim 🔒
- `POST /admin/claims/bulk-status` - Bulk update 🔒

### Billing & Payment (8 endpoints)
- `GET /billing/invoices` - List invoice user 🔒
- `POST /billing/pay` - Bayar invoice (Midtrans)
- `GET /billing/history` - Riwayat pembayaran 🔒
- `POST /webhooks/payment` - Midtrans webhook (public)
- `GET /admin/billing/invoices` - List all invoices 🔒
- `POST /admin/billing/invoices` - Buat invoice 🔒
- `PUT /admin/billing/invoices/{id}/status` - Update status 🔒

### Documents (4 endpoints)
- `GET /documents` - List dokumen
- `POST /documents/upload` - Upload dokumen
- `GET /documents/{id}/download` - Download dokumen
- `DELETE /documents/{id}` - Hapus dokumen

### Notifications (7 endpoints)
- `GET /notifications/ws` - WebSocket real-time
- `GET /notifications` - List notifikasi
- `GET /notifications/unread-count` - Jumlah unread
- `PUT /notifications/{id}/read` - Mark as read
- `PUT /notifications/read-all` - Mark all as read
- `DELETE /notifications/{id}` - Hapus notifikasi
- `GET /notification-preferences` - Get preferences 🔒
- `PUT /notification-preferences` - Update preferences 🔒

### OCR (1 endpoint)
- `POST /ocr/extract` - Ekstraksi data KTP 🔒

### Fraud Detection (2 endpoints)
- `POST /admin/applications/{id}/analyze-risk` - Analisis fraud 🔒
- `GET /admin/fraud/high-risk` - List high-risk apps 🔒

### AI & Chat (3 endpoints)
- `POST /chat` - Kirim pesan ke AI chatbot
- `GET /chat/history` - Riwayat chat
- `POST /admin/ai-review/{id}` - AI review aplikasi 🔒

### Analytics & Reports (8 endpoints)
- `GET /admin/analytics/dashboard` - Dashboard analytics 🔒
- `GET /admin/analytics/revenue` - Revenue data 🔒
- `GET /admin/analytics/claims` - Claims distribution 🔒
- `GET /admin/analytics/products` - Top products 🔒
- `GET /reports/billing/{id}/pdf` - Export billing PDF 🔒
- `GET /reports/claims/excel` - Export claims Excel 🔒
- `GET /reports/customers/excel` - Export customer list 🔒
- `GET /reports/analytics/excel` - Export analytics 🔒

### Audit Logs (5 endpoints)
- `GET /admin/audit-logs` - List audit logs 🔒
- `GET /admin/audit-logs/{id}` - Detail audit log 🔒
- `GET /admin/audit-logs/actions` - Available actions 🔒
- `GET /admin/audit-logs/entity-types` - Entity types 🔒
- `GET /admin/audit-logs/entity/{type}/{id}` - Entity logs 🔒

### Admin Email (4 endpoints)
- `POST /admin/email/application-approved/{id}` - Send approval email 🔒
- `POST /admin/email/claim-status/{id}` - Send claim email 🔒
- `POST /admin/email/bulk-send` - Bulk email 🔒
- `GET /admin/email/preview/{type}` - Preview template 🔒

### Health & Metrics (3 endpoints)
- `GET /health` - Liveness probe
- `GET /health/ready` - Readiness probe
- `GET /metrics` - Prometheus metrics

🔒 = Requires authentication

---

## 📝 Example Requests

### 1. Register & Login

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123",
    "full_name": "John Doe",
    "phone": "+628123456789"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### 2. Search Products

**Semantic Search:**
```bash
curl -X POST http://localhost:8080/api/v1/products/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "asuransi kesehatan untuk keluarga dengan biaya terjangkau",
    "limit": 5
  }'
```

### 3. Calculate Premium

```bash
curl -X POST http://localhost:8080/api/v1/products/{product-id}/calculate-premium \
  -H "Content-Type: application/json" \
  -d '{
    "age": 30,
    "sum_assured": 100000000,
    "payment_term": 12
  }'
```

### 4. Create Application

```bash
curl -X POST http://localhost:8080/api/v1/applications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "product_id": "prod-uuid",
    "user_id": "user-uuid",
    "sum_assured": 100000000,
    "form_data": {
      "occupation": "Software Engineer",
      "income": "15000000"
    }
  }'
```

### 5. Create Claim

```bash
curl -X POST http://localhost:8080/api/v1/claims \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "policy_id": "policy-uuid",
    "claim_type": "health",
    "claim_amount": 5000000,
    "description": "Medical expenses for hospitalization",
    "incident_date": "2026-09-01"
  }'
```

### 6. Pay Invoice (Midtrans)

```bash
curl -X POST http://localhost:8080/api/v1/billing/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "invoice_id": "invoice-uuid"
  }'
```

**Response:**
```json
{
  "message": "payment transaction created",
  "snap_token": "abc123...",
  "redirect_url": "https://app.midtrans.com/snap/v2/vtweb/..."
}
```

### 7. OCR Extract KTP

```bash
curl -X POST http://localhost:8080/api/v1/ocr/extract \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "image_base64": "data:image/jpeg;base64,/9j/4AAQSkZJRg..."
  }'
```

### 8. WebSocket Notifications

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/notifications/ws?user_id=user-uuid');

ws.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  console.log('New notification:', notification);
};
```

### 9. Admin - Approve Application

```bash
curl -X PUT http://localhost:8080/api/v1/admin/applications/{app-id}/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin-token>" \
  -d '{
    "status": "approved",
    "notes": "Application approved after verification"
  }'
```

### 10. Fraud Analysis

```bash
curl -X POST http://localhost:8080/api/v1/admin/applications/{app-id}/analyze-risk \
  -H "Authorization: Bearer <admin-token>"
```

---

## 🔧 Common Response Formats

### Success Response
```json
{
  "message": "Operation successful",
  "data": {...}
}
```

### Paginated Response
```json
{
  "data": [...],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

### Error Response
```json
{
  "error": "Error message describing what went wrong"
}
```

---

## 🚦 HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## 📦 Export OpenAPI File

**YAML Format:**
```bash
curl http://localhost:8080/api/docs/swagger.yaml > openapi.yaml
```

**JSON Format:**
```bash
curl http://localhost:8080/api/docs/swagger.json > openapi.json
```

---

## 🛠️ Development

### Generate/Update Documentation

Setelah menambah atau mengubah endpoint:

```bash
# Install swag tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
cd /home/bayu/Project/insurance-policy-core-api
~/go/bin/swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### Annotation Format

```go
// HandlerName godoc
// @Summary Short description
// @Description Detailed description
// @Tags Tag Name
// @Accept json
// @Produce json
// @Param id path string true "Parameter description"
// @Param request body RequestType true "Request body"
// @Success 200 {object} ResponseType "Success description"
// @Failure 400 {object} ErrorType "Error description"
// @Security BearerAuth
// @Router /endpoint [method]
func (h *Handler) HandlerName(c *fiber.Ctx) error {
    // handler implementation
}
```

---

## 📚 Additional Resources

- **Swagger UI**: http://localhost:8080/api/docs/index.html
- **OpenAPI 3.0 Spec**: https://swagger.io/specification/
- **Fiber Framework**: https://docs.gofiber.io/
- **Swaggo**: https://github.com/swaggo/swag

---

## 📊 Statistics

- **Total Endpoints**: 50+
- **Swagger YAML Lines**: 425
- **Swagger JSON Lines**: 600
- **Authentication Methods**: 2 (Bearer Token, Cookie)
- **RBAC Roles**: 5
- **API Tags**: 15+

---

## 🔍 Testing with Swagger UI

1. Start the server: `go run cmd/api/main.go`
2. Open browser: `http://localhost:8080/api/docs/index.html`
3. Click "Authorize" button (🔓)
4. Enter your JWT token: `Bearer <your-token>`
5. Try any endpoint with "Try it out" button
6. View request/response in real-time

---

## 📝 Notes

- Semua endpoint yang memerlukan autentikasi ditandai dengan 🔒
- Endpoint admin memerlukan role tertentu (super_admin, underwriter, claims_officer, finance)
- WebSocket endpoint menggunakan protokol `ws://` (production: `wss://`)
- File upload menggunakan `multipart/form-data`
- Tanggal menggunakan format `YYYY-MM-DD`
- Timestamp menggunakan format ISO 8601 `YYYY-MM-DDTHH:MM:SSZ`

---

**Generated**: 2026-09-05  
**Version**: 1.0  
**Author**: Insurance Policy Core API Team
