# RBAC Role Requirements per Endpoint

Dokumentasi lengkap requirement RBAC role untuk setiap endpoint di Insurance Policy Core API.

## 🔓 Public Endpoints (No Authentication Required)

- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login
- `POST /webhooks/payment` - Midtrans payment webhook

## 🔐 Authenticated Endpoints (Any Logged-in User)

### Customer Access
- `GET /auth/me` - Get current user info
- `POST /auth/logout` - Logout
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password

### Products (Public or Customer)
- `GET /products` - List products
- `GET /products/{id}` - Get product detail
- `POST /products/{id}/calculate-premium` - Calculate premium
- `POST /products/search` - Semantic search

### Applications
- `POST /applications` - Create application (customer)
- `GET /applications/{id}` - Get application detail (owner or admin)

### Policies (Owner Access)
- `GET /policies` - List user's policies
- `GET /policies/{id}` - Get policy detail
- `POST /policies/{id}/endorse` - Endorse policy
- `POST /policies/{id}/renew` - Renew policy

### Claims (Owner Access)
- `POST /claims` - Create claim
- `GET /claims/{id}` - Get claim detail
- `PUT /claims/{id}/documents` - Upload claim documents
- `GET /claims/{id}/timeline` - Get claim timeline

### Billing (Owner Access)
- `GET /billing/invoices` - Get user's invoices
- `POST /billing/pay` - Pay invoice
- `GET /billing/history` - Get payment history

### Documents
- `GET /documents` - List documents
- `POST /documents/upload` - Upload document
- `GET /documents/{id}/download` - Download document
- `DELETE /documents/{id}` - Delete document

### Chat
- `POST /chat` - Send message to AI chatbot
- `GET /chat/history` - Get chat history

### Notifications
- `GET /notifications/ws` - WebSocket connection
- `GET /notifications` - List notifications
- `GET /notifications/unread-count` - Get unread count
- `PUT /notifications/{id}/read` - Mark as read
- `PUT /notifications/read-all` - Mark all as read
- `DELETE /notifications/{id}` - Delete notification
- `GET /notification-preferences` - Get preferences
- `PUT /notification-preferences` - Update preferences

### OCR
- `POST /ocr/extract` - Extract ID card data

## 👔 Admin Endpoints

### Super Admin + Underwriter
**Roles Required:** `super_admin`, `underwriter`

- `GET /admin/applications` - List all applications
- `PUT /admin/applications/{id}/status` - Update application status
- `POST /admin/applications/bulk-status` - Bulk update application status
- `POST /admin/email/application-approved/{id}` - Send approval email
- `POST /admin/ai-review/{id}` - AI review application
- `POST /admin/applications/{id}/analyze-risk` - Fraud risk analysis
- `GET /admin/fraud/high-risk` - Get high-risk applications
- `GET /admin/audit-logs/actions` - Get available audit actions
- `GET /admin/audit-logs/entity-types` - Get entity types
- `GET /admin/audit-logs/entity/{type}/{id}` - Get entity audit logs
- `GET /reports/analytics/excel` - Export analytics summary (+ finance)

### Super Admin + Claims Officer
**Roles Required:** `super_admin`, `claims_officer`

- `GET /admin/claims` - List all claims
- `PUT /admin/claims/{id}/status` - Update claim status
- `PUT /admin/claims/{id}/approve` - Approve claim
- `POST /admin/claims/bulk-status` - Bulk update claim status
- `POST /admin/email/claim-status/{id}` - Send claim status email
- `GET /reports/claims/excel` - Export claims report
- `GET /admin/audit-logs/actions` - Get available audit actions (shared)
- `GET /admin/audit-logs/entity-types` - Get entity types (shared)
- `GET /admin/audit-logs/entity/{type}/{id}` - Get entity audit logs (shared)

### Super Admin + Finance
**Roles Required:** `super_admin`, `finance`

- `GET /admin/billing/invoices` - List all invoices
- `POST /admin/billing/invoices` - Create invoice
- `PUT /admin/billing/invoices/{id}/status` - Update invoice status
- `GET /reports/billing/{id}/pdf` - Export billing statement PDF
- `GET /reports/analytics/excel` - Export analytics summary (+ underwriter)

### Super Admin Only
**Roles Required:** `super_admin`

- `GET /admin/users` - List all users
- `POST /admin/users/{id}/roles` - Assign role to user
- `DELETE /admin/users/{id}/roles` - Remove role from user
- `GET /admin/roles` - Get all roles
- `POST /admin/email/bulk-send` - Send bulk email
- `GET /admin/audit-logs` - List audit logs
- `GET /admin/audit-logs/{id}` - Get audit log detail
- `GET /reports/customers/excel` - Export customer list

### Multi-Role Analytics (Admin Roles)
**Roles Required:** `super_admin`, `underwriter`, `claims_officer`, `finance`

- `GET /admin/analytics/dashboard` - Get dashboard analytics
- `GET /admin/analytics/revenue` - Get revenue analytics
- `GET /admin/analytics/claims` - Get claims status
- `GET /admin/analytics/products` - Get top products

### Email Preview (Multiple Admin Roles)
**Roles Required:** `super_admin`, `underwriter`, `claims_officer`

- `GET /admin/email/preview/{type}` - Preview email template

## 📊 Role Permission Matrix

| Endpoint Category | customer | underwriter | claims_officer | finance | super_admin |
|-------------------|----------|-------------|----------------|---------|-------------|
| Auth & Profile | ✅ | ✅ | ✅ | ✅ | ✅ |
| Products | ✅ | ✅ | ✅ | ✅ | ✅ |
| Own Applications | ✅ | ✅ | ✅ | ✅ | ✅ |
| All Applications | ❌ | ✅ | ❌ | ❌ | ✅ |
| Application Approval | ❌ | ✅ | ❌ | ❌ | ✅ |
| Own Policies | ✅ | ✅ | ✅ | ✅ | ✅ |
| Own Claims | ✅ | ✅ | ✅ | ✅ | ✅ |
| All Claims | ❌ | ❌ | ✅ | ❌ | ✅ |
| Claim Approval | ❌ | ❌ | ✅ | ❌ | ✅ |
| Own Billing | ✅ | ✅ | ✅ | ✅ | ✅ |
| All Billing | ❌ | ❌ | ❌ | ✅ | ✅ |
| Invoice Management | ❌ | ❌ | ❌ | ✅ | ✅ |
| Documents | ✅ | ✅ | ✅ | ✅ | ✅ |
| Notifications | ✅ | ✅ | ✅ | ✅ | ✅ |
| Chat/AI | ✅ | ✅ | ✅ | ✅ | ✅ |
| OCR | ✅ | ✅ | ✅ | ✅ | ✅ |
| Fraud Detection | ❌ | ✅ | ❌ | ❌ | ✅ |
| Analytics | ❌ | ✅ | ✅ | ✅ | ✅ |
| Reports | ❌ | ✅ | ✅ | ✅ | ✅ |
| Audit Logs (View) | ❌ | ✅ | ✅ | ❌ | ✅ |
| Audit Logs (Full) | ❌ | ❌ | ❌ | ❌ | ✅ |
| User Management | ❌ | ❌ | ❌ | ❌ | ✅ |
| Role Management | ❌ | ❌ | ❌ | ❌ | ✅ |
| Bulk Email | ❌ | ❌ | ❌ | ❌ | ✅ |

## 🔑 Role Descriptions

### customer
- **Purpose:** Regular users who purchase insurance
- **Access:** Personal data only (own applications, policies, claims, billing)
- **Cannot:** Access admin features, view other users' data

### underwriter
- **Purpose:** Approve insurance applications and assess risk
- **Access:** All applications, fraud detection, AI review, analytics
- **Cannot:** Manage claims, billing, or user roles

### claims_officer
- **Purpose:** Process and approve insurance claims
- **Access:** All claims, claim approval, claim reports
- **Cannot:** Approve applications, manage billing, or user roles

### finance
- **Purpose:** Manage billing and invoices
- **Access:** All invoices, payment management, billing reports
- **Cannot:** Approve applications or claims, manage user roles

### super_admin
- **Purpose:** Full system administration
- **Access:** Everything - all endpoints, all data, all management
- **Special:** User management, role assignment, audit logs, bulk operations

## 🛡️ Security Notes

1. **JWT Token Required:** All authenticated endpoints require valid JWT token
2. **Role Hierarchy:** super_admin has access to all endpoints
3. **Owner Access:** Users can always access their own data
4. **Multiple Roles:** Users can have multiple roles (cumulative permissions)
5. **Audit Trail:** All admin actions are logged in audit_logs table

## 📝 Implementation

Role checking is implemented in middleware:

```go
// Requires authentication only
middleware.AuthRequired()

// Requires specific roles
middleware.RequireRole("super_admin", "underwriter")

// Example route
admin.Put("/applications/:id/status", 
    middleware.RequireRole("super_admin", "underwriter"), 
    applicationHandler.UpdateStatus)
```

## 🚀 Testing Role Access

### Get Admin Token
```bash
# Login as admin
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Test Role-Protected Endpoint
```bash
# Access admin endpoint
curl -X GET http://localhost:8080/api/v1/admin/applications \
  -H "Authorization: Bearer <admin-token>"
```

### Expected Responses

**Success (with correct role):**
```json
{
  "data": [...],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

**Unauthorized (no token):**
```json
{
  "error": "unauthorized"
}
```

**Forbidden (wrong role):**
```json
{
  "error": "forbidden - insufficient permissions"
}
```

---

**Last Updated:** 2026-09-05  
**API Version:** 1.0
