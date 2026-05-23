#!/usr/bin/env bash
# Демонстрация API TransitOS. Запуск: ./curl-demo.sh [BASE_URL]
set -euo pipefail
BASE="${1:-http://localhost:8080}"

echo "== Health =="
curl -s "$BASE/api/health" | jq .

echo "== Login dispatcher =="
RESP=$(curl -s -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d '{"login":"dispatcher","password":"disp123"}')
TOKEN=$(echo "$RESP" | jq -r .token)
echo "Token: ${TOKEN:0:16}..."

AUTH="Authorization: Bearer $TOKEN"

echo "== Me =="
curl -s "$BASE/api/me" -H "$AUTH" | jq .

echo "== Dashboard =="
curl -s "$BASE/api/dashboard" -H "$AUTH" | jq .

echo "== Schedule (first vehicle) =="
curl -s "$BASE/api/schedule" -H "$AUTH" | jq '.[0]'

echo "== Orders count =="
curl -s "$BASE/api/orders" -H "$AUTH" | jq 'length'

echo "Done."
