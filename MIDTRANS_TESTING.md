# Midtrans Payment Testing Guide

## Quick Start

### 1. Setup Backend
```bash
cd /home/bayu/Project/insurance-policy-core-api

# Configure Midtrans credentials in .env
cat >> .env << EOF
MIDTRANS_SERVER_KEY=SB-Mid-server-YOUR_SANDBOX_SERVER_KEY
MIDTRANS_CLIENT_KEY=SB-Mid-client-YOUR_SANDBOX_CLIENT_KEY
MIDTRANS_IS_PRODUCTION=false
EOF

# Start backend
go run cmd/api/main.go
```

### 2. Setup Frontend
```bash
cd /home/bayu/insurance-policy-app

# Configure Midtrans client key in .env.local
cat >> .env.local << EOF
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=SB-Mid-client-YOUR_SANDBOX_CLIENT_KEY
NEXT_PUBLIC_MIDTRANS_IS_PRODUCTION=false
EOF

# Start frontend
npm run dev
```

### 3. Get Sandbox Credentials

1. Register at [Midtrans Sandbox](https://dashboard.sandbox.midtrans.com/)
2. Login and go to **Settings → Access Keys**
3. Copy:
   - **Server Key** → Backend `.env` (`MIDTRANS_SERVER_KEY`)
   - **Client Key** → Frontend `.env.local` (`NEXT_PUBLIC_MIDTRANS_CLIENT_KEY`)

## Payment Flow Testing

### End-to-End Flow

1. **Create Invoice** (via policy purchase or manual billing)
2. **Navigate to Billing Page** (`http://localhost:3000/billing`)
3. **Click "Bayar"** on pending invoice
4. **Midtrans Snap popup opens** with payment options
5. **Complete payment** using test card
6. **Webhook notification** updates payment status in backend
7. **Frontend refreshes** showing updated payment status

### Test Cards (Sandbox)

#### Success Scenario
```
Card Number: 4811 1111 1111 1114
CVV: 123
Expiry: 12/28 (any future date)
OTP: 112233
```

#### Failed Scenario
```
Card Number: 4911 1111 1111 1113
CVV: 123
Expiry: 12/28
```

### Other Payment Methods

- **GoPay**: Click GoPay → Use simulator app
- **QRIS**: Click QRIS → Scan with simulator
- **Bank Transfer**: Get virtual account number → Use simulator

## API Testing

### 1. Login
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }' | jq -r '.token')

echo "Token: $TOKEN"
```

### 2. Get Invoices
```bash
curl -X GET http://localhost:8080/api/v1/billing/invoices \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 3. Create Payment Transaction
```bash
INVOICE_ID="paste-invoice-id-here"

curl -X POST http://localhost:8080/api/v1/billing/pay \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"invoice_id\": \"$INVOICE_ID\"}" | jq
```

Response:
```json
{
  "message": "payment transaction created",
  "snap_token": "66e4fa55-fdac-4ef9-91b5-733b97d1b862",
  "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/..."
}
```

### 4. Test Webhook (Manual)
```bash
curl -X POST http://localhost:8080/api/v1/webhooks/payment \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_time": "2026-09-05 05:30:00",
    "transaction_status": "settlement",
    "transaction_id": "test-123",
    "status_message": "Success",
    "status_code": "200",
    "signature_key": "test",
    "payment_type": "credit_card",
    "order_id": "INV-2024-001-1725501000",
    "merchant_id": "G123456789",
    "gross_amount": "100000",
    "fraud_status": "accept",
    "currency": "IDR"
  }' | jq
```

### 5. Check Transaction Status
```bash
ORDER_ID="paste-order-id-here"

curl -X GET "http://localhost:8080/api/v1/billing/status/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

## Frontend Integration

### Midtrans Snap Library
Located at: `/home/bayu/insurance-policy-app/src/lib/midtrans.ts`

Usage in component:
```typescript
import { openSnapPayment } from '@/lib/midtrans'

// Create payment
const handlePayment = async (invoiceId: string) => {
  // 1. Call backend to create transaction
  const res = await fetch(`${API_URL}/billing/pay`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ invoice_id: invoiceId })
  })
  
  const data = await res.json()
  
  // 2. Open Midtrans Snap popup
  await openSnapPayment({
    token: data.snap_token,
    onSuccess: (result) => {
      console.log('Payment success:', result)
      // Refresh invoices
    },
    onPending: (result) => {
      console.log('Payment pending:', result)
    },
    onError: (result) => {
      console.error('Payment error:', result)
    },
    onClose: () => {
      console.log('Popup closed')
    }
  })
}
```

## Webhook Setup

### Local Development (ngrok)

```bash
# Install ngrok
brew install ngrok  # or download from ngrok.com

# Expose local backend
ngrok http 8080

# Copy HTTPS URL (e.g., https://abc123.ngrok.io)
# Go to Midtrans Dashboard → Settings → Configuration
# Set Payment Notification URL: https://abc123.ngrok.io/api/v1/webhooks/payment
```

### Production

Set webhook URL in Midtrans Dashboard:
```
https://your-api-domain.com/api/v1/webhooks/payment
```

## Transaction Status Mapping

| Midtrans Status | Backend Status | Description |
|----------------|----------------|-------------|
| `pending` | `pending` | Awaiting payment |
| `capture` | `success` | Card authorized |
| `settlement` | `success` | Payment settled |
| `deny` | `failed` | Payment denied |
| `cancel` | `failed` | Cancelled by user |
| `expire` | `failed` | Payment expired |
| `failure` | `failed` | Payment failed |

## Troubleshooting

### Error: "MIDTRANS_SERVER_KEY not set"
**Solution**: Add credentials to backend `.env`:
```bash
MIDTRANS_SERVER_KEY=SB-Mid-server-YOUR_KEY
```

### Error: "Failed to load Midtrans Snap script"
**Solution**: 
1. Check frontend `.env.local` has `NEXT_PUBLIC_MIDTRANS_CLIENT_KEY`
2. Verify network can reach `https://app.sandbox.midtrans.com`
3. Check browser console for errors

### Webhook Not Received
**Solution**:
1. For local dev, use ngrok to expose backend
2. Verify webhook URL in Midtrans Dashboard
3. Check backend logs for incoming requests
4. Test webhook manually with curl

### Payment Stuck at Pending
**Solution**:
1. Check Midtrans Dashboard → Transactions for actual status
2. Verify webhook URL is accessible from internet
3. Manually trigger webhook with correct `order_id`
4. Check backend logs for webhook processing errors

### "Transaction Not Found" Error
**Solution**:
1. Verify `order_id` format matches backend expectations
2. Check database for payment record
3. Ensure transaction was created before webhook call

## Monitoring

### Backend Logs
```bash
cd /home/bayu/Project/insurance-policy-core-api
go run cmd/api/main.go

# Look for:
# - ✅ Midtrans payment gateway initialized
# - Payment creation logs
# - Webhook notification logs
```

### Midtrans Dashboard
1. Go to **Transactions** → **All Transactions**
2. Filter by date/status
3. Click transaction for details
4. Check **Notification History** for webhook attempts

## Security Checklist (Production)

- [ ] Use production Server Key (starts with `Mid-server-`)
- [ ] Use production Client Key (starts with `Mid-client-`)
- [ ] Set `MIDTRANS_IS_PRODUCTION=true` in backend
- [ ] Set `NEXT_PUBLIC_MIDTRANS_IS_PRODUCTION=true` in frontend
- [ ] Enable signature verification in webhook handler
- [ ] Use HTTPS for webhook URL
- [ ] Store credentials in secrets manager (not .env files)
- [ ] Set up monitoring alerts for failed payments
- [ ] Test all payment methods (CC, GoPay, QRIS, Bank Transfer)
- [ ] Configure proper error handling and retry logic

## Support Resources

- [Midtrans Docs](https://docs.midtrans.com/)
- [Snap Integration Guide](https://docs.midtrans.com/docs/snap-integration)
- [API Reference](https://docs.midtrans.com/reference/api-reference)
- [Sandbox Dashboard](https://dashboard.sandbox.midtrans.com/)
- [Production Dashboard](https://dashboard.midtrans.com/)
