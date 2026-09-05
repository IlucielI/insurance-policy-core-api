# Midtrans Payment Integration

## Overview

Midtrans payment gateway integration untuk Insurance Policy Core API. Mendukung berbagai metode pembayaran: credit card, GoPay, QRIS, bank transfer.

## Environment Variables

Tambahkan ke `.env`:

```bash
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx  # Sandbox: SB-Mid-server-xxx, Production: Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx  # Untuk frontend (opsional)
MIDTRANS_IS_PRODUCTION=false             # false = sandbox, true = production
```

## API Endpoints

### 1. Create Payment Transaction
**POST** `/api/v1/billing/pay`

Request:
```json
{
  "invoice_id": "invoice-uuid-here"
}
```

Response:
```json
{
  "message": "payment transaction created",
  "snap_token": "66e4fa55-fdac-4ef9-91b5-733b97d1b862",
  "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/66e4fa55-fdac-4ef9-91b5-733b97d1b862"
}
```

### 2. Payment Webhook (Midtrans Callback)
**POST** `/api/v1/webhooks/payment`

Midtrans akan mengirim notification ke endpoint ini ketika status pembayaran berubah.

Request (dari Midtrans):
```json
{
  "transaction_time": "2026-09-05 05:30:00",
  "transaction_status": "settlement",
  "transaction_id": "abc123",
  "status_message": "Payment successful",
  "status_code": "200",
  "signature_key": "...",
  "payment_type": "credit_card",
  "order_id": "INV-2024-001-1725501000",
  "merchant_id": "G123456789",
  "gross_amount": "100000",
  "fraud_status": "accept",
  "currency": "IDR"
}
```

## Payment Flow

```
1. Customer → POST /billing/pay
   ↓
2. Backend → Midtrans Snap API (create transaction)
   ↓
3. Backend ← Snap Token & Redirect URL
   ↓
4. Customer redirected to Snap payment page
   ↓
5. Customer completes payment
   ↓
6. Midtrans → Webhook notification to /webhooks/payment
   ↓
7. Backend updates payment status
```

## Testing

### Sandbox Credentials

1. Daftar di [Midtrans Sandbox](https://dashboard.sandbox.midtrans.com/)
2. Dapatkan **Server Key** dari Settings → Access Keys
3. Set environment variable `MIDTRANS_SERVER_KEY`

### Test Cards (Sandbox)

**Successful Payment:**
- Card: `4811 1111 1111 1114`
- CVV: `123`
- Exp: Any future date (e.g., `12/28`)

**Failed Payment:**
- Card: `4911 1111 1111 1113`
- CVV: `123`
- Exp: Any future date

**Other Payment Methods:**
- GoPay: Use simulator in Sandbox
- QRIS: Use simulator QR code
- Bank Transfer: Virtual account akan di-generate

### Automated Test Script

```bash
cd /home/bayu/Project/insurance-policy-core-api
chmod +x test_midtrans_payment.sh
./test_midtrans_payment.sh
```

### Manual Testing

1. Start the API server:
```bash
cd /home/bayu/Project/insurance-policy-core-api
go run cmd/api/main.go
```

2. Login as customer:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

3. Get invoices:
```bash
curl -X GET http://localhost:8080/api/v1/billing/invoices \
  -H "Authorization: Bearer YOUR_TOKEN"
```

4. Create payment:
```bash
curl -X POST http://localhost:8080/api/v1/billing/pay \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_id": "INVOICE_ID_FROM_STEP_3"
  }'
```

5. Open `redirect_url` dari response di browser
6. Complete payment dengan test card
7. Verify webhook dipanggil (check server logs)

## Webhook Configuration

Di Midtrans Dashboard:

1. Go to **Settings → Configuration**
2. Set **Payment Notification URL**: `https://your-domain.com/api/v1/webhooks/payment`
3. Enable all notification types

## Production Checklist

- [ ] Update `MIDTRANS_SERVER_KEY` dengan production key
- [ ] Set `MIDTRANS_IS_PRODUCTION=true`
- [ ] Configure webhook URL di Midtrans production dashboard
- [ ] Enable signature verification (currently simplified for demo)
- [ ] Set up monitoring for payment failures
- [ ] Test semua payment methods (CC, GoPay, QRIS, Bank Transfer)
- [ ] Configure proper error handling dan retry logic
- [ ] Set up invoice-payment linking (currently simplified)

## Transaction Status Mapping

| Midtrans Status | Backend Status | Description |
|----------------|----------------|-------------|
| `pending` | `pending` | Menunggu pembayaran |
| `capture` | `success` | CC authorized |
| `settlement` | `success` | Payment settled |
| `deny` | `failed` | Payment ditolak |
| `cancel` | `failed` | Payment dibatalkan |
| `expire` | `failed` | Payment expired |
| `failure` | `failed` | Payment gagal |

## Files Modified/Created

### New Files
- `internal/infrastructure/payment/midtrans_client.go` - Midtrans SDK wrapper
- `test_midtrans_payment.sh` - Test script

### Modified Files
- `internal/usecase/billing_usecase.go` - Added `CreatePayment()`, `HandlePaymentWebhook()`
- `internal/delivery/http/billing_handler.go` - Updated `PayInvoice()`, added `HandlePaymentWebhook()`
- `internal/repository/billing_repository.go` - Added `CreatePayment()`, `GetPaymentByOrderID()`, `UpdatePaymentStatus()`
- `cmd/api/main.go` - Wire Midtrans client, register webhook route
- `go.mod` - Added `github.com/midtrans/midtrans-go v1.3.7`

## Troubleshooting

**Error: "MIDTRANS_SERVER_KEY not set"**
- Set environment variable di `.env` atau shell

**Error: "payment not found" saat webhook**
- Order ID mismatch. Pastikan order_id dari webhook match dengan database

**Payment stuck di pending**
- Check Midtrans dashboard untuk status sebenarnya
- Verify webhook URL bisa diakses dari internet (use ngrok untuk local testing)

**Webhook tidak dipanggil**
- Check firewall
- Use ngrok: `ngrok http 8080` dan set webhook URL ke ngrok URL
- Verify di Midtrans dashboard → Transactions → notification history

## Support

Midtrans Documentation: https://docs.midtrans.com/
Midtrans Sandbox: https://dashboard.sandbox.midtrans.com/
