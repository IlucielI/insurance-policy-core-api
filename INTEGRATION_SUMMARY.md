# Midtrans Payment Integration - Summary

## ✅ Completed

Successfully integrated Midtrans payment gateway into Insurance Policy Core API.

## Files Created

1. **internal/infrastructure/payment/midtrans_client.go**
   - `NewMidtransClient()` - Initialize Snap & Core API clients
   - `CreateTransaction()` - Generate Snap token for payment
   - `GetStatus()` - Query payment status
   - `HandleWebhook()` - Process Midtrans notifications
   - Payment status helpers: `IsPaymentSuccess()`, `IsPaymentPending()`, `IsPaymentFailed()`

2. **test_midtrans_payment.sh**
   - Automated integration test script
   - Tests full payment flow: login → get invoices → create payment → webhook

3. **MIDTRANS_INTEGRATION.md**
   - Complete documentation
   - API endpoints, flow diagrams, testing guide
   - Sandbox credentials and test cards
   - Production checklist

## Files Modified

1. **go.mod**
   - Added `github.com/midtrans/midtrans-go v1.3.7`

2. **config/config.go**
   - Already had Midtrans config fields (no changes needed)

3. **internal/usecase/billing_usecase.go**
   - Added `CreatePayment()` - Creates Midtrans transaction
   - Added `HandlePaymentWebhook()` - Processes payment notifications
   - Updated constructor to accept `userRepo` and `midtransClient`

4. **internal/delivery/http/billing_handler.go**
   - Updated `PayInvoice()` - Now creates Midtrans Snap transaction
   - Added `HandlePaymentWebhook()` - Webhook endpoint handler

5. **internal/repository/billing_repository.go**
   - Added `CreatePayment()` - Insert payment record
   - Added `GetPaymentByOrderID()` - Query by Midtrans order ID
   - Added `UpdatePaymentStatus()` - Update payment after webhook

6. **cmd/api/main.go**
   - Initialize Midtrans client from config
   - Wire client to billing usecase
   - Register webhook route: `POST /api/v1/webhooks/payment`

7. **internal/usecase/document_usecase.go**
   - Removed unused storage import (build fix)

## API Endpoints

### POST /api/v1/billing/pay
Create payment transaction.

**Request:**
```json
{
  "invoice_id": "uuid"
}
```

**Response:**
```json
{
  "message": "payment transaction created",
  "snap_token": "66e4fa55-fdac-4ef9-91b5-733b97d1b862",
  "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/..."
}
```

### POST /api/v1/webhooks/payment
Midtrans notification callback (no auth).

## Payment Flow

```
Customer → POST /billing/pay
  ↓
Backend → Midtrans Snap API
  ↓
Customer → Redirected to Snap page
  ↓
Customer → Completes payment (CC/GoPay/QRIS/Bank)
  ↓
Midtrans → POST /webhooks/payment
  ↓
Backend → Updates payment status in DB
```

## Database Schema

**payments table** (already exists in migrations/001_init_schema.sql):
- id, application_id, order_id
- midtrans_transaction_id, payment_type
- gross_amount, status
- paid_at, expired_at
- created_at, updated_at

## Environment Variables

```bash
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx  # Sandbox key
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx  # Optional for frontend
MIDTRANS_IS_PRODUCTION=false             # false = sandbox
```

## Testing

### 1. Run Test Script
```bash
cd /home/bayu/Project/insurance-policy-core-api
chmod +x test_midtrans_payment.sh
./test_midtrans_payment.sh
```

### 2. Manual Testing
1. Get sandbox key from https://dashboard.sandbox.midtrans.com/
2. Set `MIDTRANS_SERVER_KEY` in `.env`
3. Start server: `docker-compose up` or `go run cmd/api/main.go`
4. Call `POST /billing/pay` with invoice_id
5. Open returned `redirect_url` in browser
6. Use test card: `4811 1111 1111 1114`, CVV: `123`, Exp: `12/28`
7. Complete payment → webhook will be called

### Supported Payment Methods
- Credit Card (Visa, Mastercard, JCB, Amex)
- GoPay
- QRIS
- Bank Transfer (BCA, BNI, BRI, Permata, etc.)
- Alfamart/Indomaret
- Akulaku, Kredivo

## Build Verification

✅ Docker build successful:
```bash
docker build -t insurance-api-test:latest .
```

Build output: No errors, image created successfully.

## Production Checklist

- [ ] Replace sandbox key with production key
- [ ] Set `MIDTRANS_IS_PRODUCTION=true`
- [ ] Configure webhook URL in Midtrans dashboard
- [ ] Enable signature verification in webhook handler
- [ ] Add logging for payment events
- [ ] Set up monitoring/alerts for payment failures
- [ ] Test all payment methods in production
- [ ] Link payment records to invoices (currently simplified)
- [ ] Handle payment expiry notifications
- [ ] Implement retry logic for failed webhooks

## Known Limitations

1. **Invoice-Payment Linking**: Simplified approach. In production, store `invoice_id` in payment record for direct linking.
2. **Webhook Signature**: Currently not validating Midtrans signature. Add verification for production.
3. **Payment Expiry**: Not handling expired payment notifications. Add handler for better UX.
4. **Refunds**: Refund flow not implemented. Add if needed.

## Next Steps

1. Test with real sandbox account
2. Configure webhook URL (use ngrok for local testing)
3. Verify webhook receives notifications
4. Test all payment methods
5. Monitor payment success rate
6. Implement invoice update after successful payment

## Support

- Midtrans Docs: https://docs.midtrans.com/
- Sandbox Dashboard: https://dashboard.sandbox.midtrans.com/
- Go SDK: https://github.com/midtrans/midtrans-go
