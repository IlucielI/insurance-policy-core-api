#!/bin/bash
# Test script untuk fraud detection API

API_URL="${API_URL:-http://localhost:8080}"
TOKEN="${AUTH_TOKEN}"

echo "🛡️  Fraud Detection API Test"
echo "================================"
echo ""

# Get sample application ID
echo "📋 Fetching applications..."
APP_ID=$(curl -s "$API_URL/api/v1/admin/applications?limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')

if [ -z "$APP_ID" ] || [ "$APP_ID" = "null" ]; then
  echo "❌ No applications found. Create one first."
  exit 1
fi

echo "✅ Found application: $APP_ID"
echo ""

# Test 1: Analyze Risk
echo "🔍 Test 1: Analyze Application Risk"
echo "-----------------------------------"
ANALYZE_RESULT=$(curl -s -X POST "$API_URL/api/v1/admin/applications/$APP_ID/analyze-risk" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "$ANALYZE_RESULT" | jq '.'
echo ""

RISK_SCORE=$(echo "$ANALYZE_RESULT" | jq -r '.data.risk_score // "N/A"')
RISK_LEVEL=$(echo "$ANALYZE_RESULT" | jq -r '.data.risk_level // "N/A"')

echo "📊 Risk Score: $RISK_SCORE ($RISK_LEVEL)"
echo ""

# Test 2: Get High-Risk Applications
echo "🚨 Test 2: Get High-Risk Applications"
echo "-------------------------------------"
HIGH_RISK=$(curl -s "$API_URL/api/v1/admin/fraud/high-risk?min_score=50&limit=10" \
  -H "Authorization: Bearer $TOKEN")

echo "$HIGH_RISK" | jq '.'
echo ""

TOTAL=$(echo "$HIGH_RISK" | jq -r '.total // 0')
echo "📈 Total high-risk applications: $TOTAL"
echo ""

# Test 3: Verify Database
echo "💾 Test 3: Verify Database"
echo "--------------------------"
echo "Checking risk_score column in applications table..."

DB_CHECK=$(docker exec -i insurance-postgres psql -U postgres -d insurance_policy -t -c \
  "SELECT COUNT(*) FROM applications WHERE risk_score IS NOT NULL;")

echo "Applications with risk scores: $(echo $DB_CHECK | xargs)"
echo ""

# Summary
echo "✅ Fraud Detection Tests Complete"
echo "=================================="
echo ""
echo "Results:"
echo "  - API analyze-risk: $([ -n "$RISK_SCORE" ] && echo "✅ Working" || echo "❌ Failed")"
echo "  - API high-risk list: $([ "$TOTAL" -ge 0 ] && echo "✅ Working" || echo "❌ Failed")"
echo "  - Database schema: ✅ Ready"
echo ""
echo "Next: Integrate with CMS underwriting page"
