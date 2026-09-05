#!/bin/bash
# SMTP Email Testing Guide
# Run this after configuring SMTP credentials in .env

set -e

echo "🔍 SMTP Email Integration Test"
echo "================================"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo ""
    echo "Create .env with SMTP configuration:"
    echo ""
    echo "SMTP_HOST=sandbox.smtp.mailtrap.io"
    echo "SMTP_PORT=2525"
    echo "SMTP_USER=your_mailtrap_username"
    echo "SMTP_PASSWORD=your_mailtrap_password"
    echo "SMTP_FROM=Insurance Platform <noreply@insurance.com>"
    echo ""
    echo "Get free Mailtrap credentials at: https://mailtrap.io"
    exit 1
fi

# Check if SMTP_HOST is configured
SMTP_HOST=$(grep "^SMTP_HOST=" .env | cut -d'=' -f2)
if [ -z "$SMTP_HOST" ]; then
    echo "❌ SMTP_HOST not configured in .env"
    echo ""
    echo "Add SMTP configuration to .env (see .env.example)"
    exit 1
fi

echo "📧 SMTP Configuration Found"
echo "   Host: $SMTP_HOST"
echo ""

# Run Go test
echo "🧪 Running Email Test Script..."
echo "   This will send 4 test emails:"
echo "   1. Welcome Email"
echo "   2. Policy Issued Email"
echo "   3. Claim Status Update Email"
echo "   4. Password Reset Email"
echo ""

/home/bayu/go-local/go/bin/go run test_smtp.go

echo ""
echo "✅ Test completed!"
echo ""
echo "📬 Check your email inbox (SMTP_USER address)"
echo "   If using Mailtrap: https://mailtrap.io/inboxes"
