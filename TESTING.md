# SMTP Email Integration - Test Guide

## Quick Test (After Go Installation)

### 1. Install Dependencies
```bash
cd /home/bayu/Project/insurance-policy-core-api
go mod tidy
```

### 2. Configure SMTP in .env

Add these lines to `.env`:
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM=Insurance Platform <your-email@gmail.com>
```

**For Gmail:**
- Visit: https://myaccount.google.com/apppasswords
- Generate App Password (need 2FA enabled)
- Use that 16-char password, NOT your regular Gmail password

**For Mailtrap (Testing - emails won't actually send):**
```bash
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=<from mailtrap.io inbox>
SMTP_PASSWORD=<from mailtrap.io inbox>
SMTP_FROM=noreply@insurance.test
```

### 3. Run Test Script
```bash
go run test_smtp.go
```

Expected output:
```
📧 Testing SMTP Email Service
Host: smtp.gmail.com:587
User: your-email@gmail.com

Test 1: Sending Welcome Email...
✅ Welcome email sent successfully!

Test 2: Sending Policy Issued Email...
✅ Policy issued email sent successfully!

Test 3: Sending Claim Status Update Email...
✅ Claim status email sent successfully!

✅ All email tests completed! Check your inbox: your-email@gmail.com
```

### 4. Test via API

**Test 1: Register User (Welcome Email)**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "phone": "081234567890"
  }'
```

**Test 2: Approve Application (Policy Issued Email)**
```bash
# First create an application, then approve it
curl -X PUT http://localhost:8080/api/v1/admin/applications/{app_id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "underwriter_id": "underwriter-123",
    "underwriter_notes": "Application approved"
  }'
```

**Test 3: Update Claim Status (Claim Status Email)**
```bash
# This requires adding endpoint in claim_handler.go
# POST /api/v1/admin/claims/:id/status
```

## Verification

### Email Content Check

1. **Welcome Email** contains:
   - Greeting with user's full name
   - Platform features list
   - Professional styling

2. **Policy Issued Email** contains:
   - Policy number
   - Product name
   - Sum assured amount
   - Congratulations message

3. **Claim Status Update Email** contains:
   - Claim number
   - Status with emoji (✅ approved, ❌ rejected, 🔍 under review)
   - Notes/reason
   - Color coding by status

## Troubleshooting

### "535 Authentication failed"
- **Gmail**: Use App Password, not regular password
- **Mailtrap**: Copy credentials exactly from dashboard

### "dial tcp: i/o timeout"
- Port 587 might be blocked
- Try port 465 (SSL) or 2525 (Mailtrap)
- Check firewall rules

### "certificate verify failed"
- TLS issue with SMTP server
- Only for dev: can temporarily set `InsecureSkipVerify: true`

### Email not received
- Check spam folder
- Verify SMTP_FROM matches SMTP_USER domain
- Check provider's sending limits

## Files Modified/Created

### Core Implementation
- ✅ `internal/infrastructure/email/smtp_client.go` - SMTP client with 3 email methods
- ✅ `go.mod` - Added `github.com/wneessen/go-mail v0.4.1`
- ✅ `config/config.go` - Added SMTP config fields

### Integration Points
- ✅ `internal/usecase/auth_usecase.go` - Welcome email on register
- ✅ `internal/usecase/application_usecase.go` - Policy issued on approve
- ✅ `internal/usecase/claim_usecase.go` - Claim status update email
- ✅ `cmd/api/main.go` - SMTP client initialization and wiring

### Documentation
- ✅ `.env.example` - SMTP config template
- ✅ `EMAIL_SETUP.md` - Full setup guide
- ✅ `test_smtp.go` - Standalone test script
- ✅ `TESTING.md` - This file

## Production Deployment

1. Use dedicated SMTP service (Mailgun/SendGrid)
2. Set proper SPF/DKIM records
3. Monitor delivery rates
4. Consider email queue (Redis) for reliability
5. Implement retry logic
6. Track bounces and unsubscribes

## Status: READY FOR TESTING ✅

All code is in place. Once Go is installed:
```bash
go mod tidy
go run test_smtp.go
```

Then start the API and test via registration endpoint.
