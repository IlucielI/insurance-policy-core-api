# Midtrans Payment Integration - Ringkasan

## Status Implementasi ✅

### Backend (Go API)
**Lokasi**: `/home/bayu/Project/insurance-policy-core-api`

#### SDK & Dependencies
- ✅ Midtrans Go SDK v1.3.7 installed (`go.mod`)
- ✅ Midtrans client initialized di `cmd/api/main.go`

#### Endpoints
1. ✅ **POST `/api/v1/billing/pay`** - Create payment transaction
   - Handler: `billing_handler.go:PayInvoice()`
   - Creates Midtrans Snap transaction
   - Returns `snap_token` & `redirect_url`

2. ✅ **POST `/api/v1/webhooks/payment`** - Payment webhook callback
   - Handler: `billing_handler.go:HandlePaymentWebhook()`
   - Processes Midtrans notifications
   - Updates payment status

3. ✅ **GET `/api/v1/billing/invoices`** - List invoices
4. ✅ **GET `/api/v1/billing/history`** - Payment history

#### Infrastructure
- ✅ `internal/infrastructure/payment/midtrans_client.go`
  - `CreateTransaction()` - Snap API integration
  - `GetStatus()` - Check transaction status
  - `HandleWebhook()` - Process webhook notifications
  - Helper functions: `IsPaymentSuccess()`, `IsPaymentPending()`, `IsPaymentFailed()`

#### Business Logic
- ✅ `internal/usecase/billing_usecase.go`
  - `CreatePayment()` - Create Midtrans transaction
  - `HandlePaymentWebhook()` - Process webhook & update DB

#### Database
- ✅ `internal/repository/billing_repository.go`
  - `CreatePayment()` - Store payment record
  - `GetPaymentByOrderID()` - Find payment
  - `UpdatePaymentStatus()` - Update status from webhook

#### Configuration
- ✅ `.env.example` updated dengan sandbox credentials template:
  ```bash
  MIDTRANS_SERVER_KEY=SB-Mid-server-YOUR_SANDBOX_SERVER_KEY
  MIDTRANS_CLIENT_KEY=SB-Mid-client-YOUR_SANDBOX_CLIENT_KEY
  MIDTRANS_IS_PRODUCTION=false
  ```
- ✅ Test cards documented

### Frontend (Next.js Customer App)
**Lokasi**: `/home/bayu/insurance-policy-app`

#### Library
- ✅ `src/lib/midtrans.ts` - Midtrans Snap integration helper
  - `loadMidtransScript()` - Load Snap SDK dynamically
  - `openSnapPayment()` - Open payment popup
  - TypeScript types untuk Midtrans responses

#### UI Integration
- ✅ `src/app/billing/page.tsx` updated
  - `handlePayment()` calls backend `/billing/pay`
  - Opens Midtrans Snap popup dengan `snap_token`
  - Handles callbacks: `onSuccess`, `onPending`, `onError`, `onClose`
  - Auto-refresh invoices after payment

#### Configuration
- ✅ `.env.local.example` created:
  ```bash
  NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=SB-Mid-client-YOUR_SANDBOX_CLIENT_KEY
  NEXT_PUBLIC_MIDTRANS_IS_PRODUCTION=false
  ```
- ✅ `.env.local` updated dengan Midtrans config

#### Build
- ✅ Next.js build successful (verified)

## Payment Flow

```
┌─────────────┐
│   Customer  │
└──────┬──────┘
       │ 1. Click "Bayar"
       ▼
┌─────────────────┐
│  Billing Page   │
│  (Next.js UI)   │
└──────┬──────────┘
       │ 2. POST /billing/pay
       │    {invoice_id: "..."}
       ▼
┌──────────────────┐
│   Go Backend     │
│  - Create Snap   │
│    transaction   │
└──────┬───────────┘
       │ 3. Return snap_token
       ▼
┌──────────────────┐
│ Midtrans Snap    │
│   Popup Opens    │
│  (Credit Card,   │
│   GoPay, QRIS)   │
└──────┬───────────┘
       │ 4. Customer completes payment
       ▼
┌──────────────────┐
│    Midtrans      │
│  sends webhook   │
└──────┬───────────┘
       │ 5. POST /webhooks/payment
       │    {transaction_status: "settlement"}
       ▼
┌──────────────────┐
│   Go Backend     │
│  - Update DB     │
│  - Update status │
└──────────────────┘
```

## Testing

### Setup Sandbox
1. Register: https://dashboard.sandbox.midtrans.com/
2. Get keys dari Settings → Access Keys
3. Update `.env` (backend) & `.env.local` (frontend)

### Test Cards (Sandbox)
- **Success**: `4811 1111 1111 1114` | CVV: `123` | Exp: `12/28` | OTP: `112233`
- **Failed**: `4911 1111 1111 1113` | CVV: `123` | Exp: `12/28`

### Run Locally
```bash
# Backend
cd /home/bayu/Project/insurance-policy-core-api
# Set MIDTRANS_SERVER_KEY in .env
go run cmd/api/main.go  # Port 8080

# Frontend
cd /home/bayu/insurance-policy-app
# Set NEXT_PUBLIC_MIDTRANS_CLIENT_KEY in .env.local
npm run dev  # Port 3000
```

### Test Flow
1. Login ke customer app
2. Browse ke `/billing`
3. Klik "Bayar" pada invoice pending
4. Snap popup terbuka
5. Gunakan test card untuk complete payment
6. Webhook updates payment status
7. UI refresh menampilkan status "Lunas"

### Webhook Testing (Local)
Use **ngrok** untuk expose local backend:
```bash
ngrok http 8080
# Set webhook URL di Midtrans Dashboard:
# https://abc123.ngrok.io/api/v1/webhooks/payment
```

## Files Created/Modified

### Backend
- ✅ `internal/infrastructure/payment/midtrans_client.go` (exists)
- ✅ `internal/usecase/billing_usecase.go` (modified)
- ✅ `internal/delivery/http/billing_handler.go` (modified)
- ✅ `internal/repository/billing_repository.go` (modified)
- ✅ `cmd/api/main.go` (Midtrans client wired, routes registered)
- ✅ `.env.example` (updated dengan sandbox credentials docs)
- ✅ `MIDTRANS_INTEGRATION.md` (exists)
- ✅ `MIDTRANS_TESTING.md` (created - comprehensive testing guide)
- ✅ `test_midtrans_payment.sh` (exists)

### Frontend
- ✅ `src/lib/midtrans.ts` (created - Snap integration helper)
- ✅ `src/app/billing/page.tsx` (updated - integrated Snap payment)
- ✅ `.env.local.example` (created - with Midtrans config)
- ✅ `.env.local` (updated - Midtrans keys added)

## Production Checklist

Before going live:
- [ ] Get production keys from https://dashboard.midtrans.com/
- [ ] Update backend `.env`: `MIDTRANS_SERVER_KEY=Mid-server-...`
- [ ] Update frontend `.env.local`: `NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=Mid-client-...`
- [ ] Set `MIDTRANS_IS_PRODUCTION=true` (both backend & frontend)
- [ ] Configure webhook URL di Midtrans Dashboard (production)
- [ ] Enable signature verification di webhook handler
- [ ] Test all payment methods (CC, GoPay, QRIS, Bank Transfer)
- [ ] Setup monitoring untuk failed payments
- [ ] Use HTTPS untuk webhook endpoint

## Documentation

Lengkap di:
- **`MIDTRANS_INTEGRATION.md`** - Overview & API specs
- **`MIDTRANS_TESTING.md`** - Comprehensive testing guide (NEW)
- **`.env.example`** - Sandbox credentials template

## Support

- Midtrans Docs: https://docs.midtrans.com/
- Snap Integration: https://docs.midtrans.com/docs/snap-integration
- Sandbox Dashboard: https://dashboard.sandbox.midtrans.com/
