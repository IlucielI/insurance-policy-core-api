# SMTP Email Integration Guide

## Setup Overview

Email service integrated untuk 3 notifikasi utama:
1. **Welcome Email** - saat user register
2. **Policy Issued Email** - saat application approved
3. **Claim Status Update** - saat status claim berubah

## Configuration

Add to `.env`:
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=Insurance Platform <your-email@gmail.com>
```

## Gmail Setup (Recommended for Testing)

1. **Enable 2FA** di Google Account
2. **Generate App Password**:
   - Go to: https://myaccount.google.com/apppasswords
   - Select "Mail" and device
   - Copy 16-character password
   - Use ini sebagai `SMTP_PASSWORD`

## Alternative SMTP Providers

### Mailtrap (Development/Testing)
```bash
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your-mailtrap-username
SMTP_PASSWORD=your-mailtrap-password
```
Free tier: unlimited emails (inbox only, tidak terkirim real)

### Mailgun (Production)
```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.mailgun.org
SMTP_PASSWORD=your-mailgun-smtp-password
```

### SendGrid (Production)
```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

## Testing SMTP

Run test script:
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

✅ All email tests completed!
```

## Troubleshooting

### Error: "535 Authentication failed"
- Gmail: pastikan App Password benar (bukan password biasa)
- Pastikan 2FA enabled

### Error: "dial tcp: i/o timeout"
- Check firewall/network blocking port 587
- Try port 465 (SSL) atau 25

### Error: "certificate verification failed"
- Production: set proper TLS config
- Development only: bisa set `InsecureSkipVerify: true` (NOT for prod!)

## Email Templates

Semua template HTML responsive dengan styling inline. Located di:
`internal/infrastructure/email/smtp_client.go`

### Customize Templates
Edit methods:
- `SendWelcomeEmail()` - welcome email template
- `SendPolicyIssuedEmail()` - policy issued template  
- `SendClaimStatusUpdateEmail()` - claim update template

## Production Checklist

- [ ] Use dedicated SMTP service (Mailgun/SendGrid)
- [ ] Set proper SPF/DKIM/DMARC records
- [ ] Monitor email delivery rates
- [ ] Implement retry logic for failed sends
- [ ] Add email queue (Redis/RabbitMQ) for high volume
- [ ] Track email open/click rates
- [ ] Unsubscribe mechanism untuk marketing emails

## Notes

- Email sending non-blocking (goroutine) - tidak fail registration jika email fail
- Email service optional - app works tanpa SMTP config
- User email di-fetch dari user repository (TODO: implement di production)
