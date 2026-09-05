# ✅ SMTP Email Integration - Complete

## Status: Ready for Testing

All 4 email templates implemented and integrated into backend API.

## 📧 Email Templates

### 1. Welcome Email
- **Trigger**: User registration (`POST /api/v1/auth/register`)
- **Template**: Green header, feature list, welcome message
- **Method**: `SendWelcomeEmail(to, fullName)`

### 2. Policy Issued Email
- **Trigger**: Application approved → policy created
- **Template**: Blue header, policy details box (number, product, sum assured)
- **Method**: `SendPolicyIssuedEmail(to, fullName, policyNumber, productName, sumAssured)`

### 3. Claim Status Update Email
- **Trigger**: Claim status changes (approved/rejected/under_review)
- **Template**: Dynamic color coding (green/red/orange), status emoji
- **Method**: `SendClaimStatusUpdateEmail(to, fullName, claimNumber, status, notes)`

### 4. Password Reset Email ✨ NEW
- **Trigger**: Password reset request
- **Template**: Red/orange header, reset button + link, 1-hour expiry warning
- **Method**: `SendPasswordResetEmail(to, fullName, resetToken)`
- **URL**: `http://localhost:3000/reset-password?token={token}`

## 🔧 SMTP Configuration

### Option A: Mailtrap (Recommended for Development)

1. Sign up at https://mailtrap.io (free)
2. Get credentials from your inbox
3. Add to `.env`:

```bash
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your_mailtrap_username
SMTP_PASSWORD=your_mailtrap_password
SMTP_FROM=Insurance Platform <noreply@insurance.com>
```

### Option B: Gmail

1. Enable 2FA on your Google account
2. Generate App Password at https://myaccount.google.com/apppasswords
3. Add to `.env`:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM=Insurance Platform <your-email@gmail.com>
```

## 🧪 Testing

### Quick Test (All 4 Emails)

```bash
cd /home/bayu/Project/insurance-policy-core-api

# 1. Configure SMTP in .env (see above)

# 2. Run test script
/home/bayu/go-local/go/bin/go run test_smtp.go
```

**Expected Output:**
```
📧 Testing SMTP Email Service
Host: sandbox.smtp.mailtrap.io:587
User: your_username

Test 1: Sending Welcome Email...
✅ Welcome email sent successfully!

Test 2: Sending Policy Issued Email...
✅ Policy issued email sent successfully!

Test 3: Sending Claim Status Update Email...
✅ Claim status email sent successfully!

Test 4: Sending Password Reset Email...
✅ Password reset email sent successfully!

✅ All email tests completed! Check your inbox: your@email.com
```

### API Integration Test

```bash
# Start backend
/home/bayu/go-local/go/bin/go run cmd/api/main.go

# Test welcome email via registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "phone": "081234567890"
  }'
```

Check email inbox - welcome email should arrive.

## 📁 Files

### Created/Modified
```
internal/infrastructure/email/smtp_client.go  - Email service (4 methods)
test_smtp.go                                  - Standalone test script
.env.example                                  - SMTP config template
EMAIL_INTEGRATION_COMPLETE.md                 - This file
```

### Integration Points
```
internal/usecase/auth_usecase.go              - Welcome email on register
internal/usecase/application_usecase.go        - Policy issued on approval
internal/usecase/claim_usecase.go             - Claim status updates
cmd/api/main.go                               - SMTP client initialization
config/config.go                              - SMTP config struct
```

## 🔒 Implementation Features

✅ **Non-blocking sends** - Emails sent in goroutines  
✅ **Graceful degradation** - App works without SMTP config  
✅ **HTML templates** - Professional styling, responsive  
✅ **TLS encryption** - Secure SMTP connection  
✅ **Error handling** - Failures logged, don't break user flow  
✅ **Timeout protection** - 10s per email send  

## 📊 Library

**github.com/wneessen/go-mail v0.4.1**
- Modern Go SMTP client
- Better than stdlib net/smtp
- Built-in TLS, authentication
- Active maintenance

## 🚀 Production Checklist

- [ ] Use production SMTP service (SendGrid, Mailgun, AWS SES)
- [ ] Configure proper FROM domain with SPF/DKIM/DMARC
- [ ] Add email queue (Redis/RabbitMQ) for reliability
- [ ] Implement retry logic for failed sends
- [ ] Add unsubscribe links (marketing emails only)
- [ ] Monitor delivery rates and bounces
- [ ] Rate limiting per user
- [ ] Email preferences management

## 🎯 Next Steps

1. **Get SMTP credentials** (Mailtrap or Gmail)
2. **Add to `.env`**
3. **Run `go run test_smtp.go`**
4. **Verify all 4 emails received**
5. **Test API registration** → welcome email

---

**Implementation Date**: 2026-09-05  
**Status**: ✅ Complete - 4 Email Templates Ready  
**Test Status**: Pending SMTP credentials
