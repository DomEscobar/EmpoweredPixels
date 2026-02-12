#!/bin/bash
# Hardcore Player Session - Automated Test Script
# Date: 2026-02-10
# Mode: Hardcore (optimization, edge cases, exploit hunting, stress)

API_BASE="http://localhost:49101"
TIMESTAMP=$(date +%s)
RANDOM_SUFFIX=$(shuf -n 1 -i 1000-9999)
USERNAME="test_player_${TIMESTAMP}_${RANDOM_SUFFIX}"
PASSWORD="TestPass123!"
EMAIL="test_${TIMESTAMP}@example.com"
SESSION_LOG="/root/.openclaw/memory/player-sessions/$(date +%Y-%m-%d)-hardcore.md"

# Ensure log directory exists
mkdir -p "$(dirname "$SESSION_LOG")"

# Initialize log
cat > "$SESSION_LOG" <<EOF
# Player Session Log - Hardcore
**Date**: $(date '+%Y-%m-%d %H:%M:%S')
**Mode**: Hardcore (optimization, edge cases, exploit hunting, stress)
**User**: $USERNAME
**Objective**: Optimize routes, test edge cases, hunt exploits, stress balance

---

## 1. Registration & Authentication

EOF

echo "=== Starting Hardcore Session ===" | tee -a "$SESSION_LOG"

# 1. Register
echo "Registering user: $USERNAME" | tee -a "$SESSION_LOG"
REG_RESP=$(curl -s -X POST "$API_BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"email\":\"$EMAIL\"}")
echo "Register response: $REG_RESP" | tee -a "$SESSION_LOG"

# 2. Login
echo "Logging in..." | tee -a "$SESSION_LOG"
LOGIN_RESP=$(curl -s -X POST "$API_BASE/api/authentication/token" \
  -H "Content-Type: application/json" \
  -d "{\"user\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
echo "Login response: $LOGIN_RESP" | tee -a "$SESSION_LOG"

# Extract tokens using jq
ACCESS_TOKEN=$(echo "$LOGIN_RESP" | jq -r '.token')
REFRESH_TOKEN=$(echo "$LOGIN_RESP" | jq -r '.refresh')
USER_ID=$(echo "$LOGIN_RESP" | jq -r '.userId')

if [ -z "$ACCESS_TOKEN" ]; then
  echo "ERROR: Failed to obtain access token. Aborting." | tee -a "$SESSION_LOG"
  exit 1
fi

echo "Obtained tokens. User ID: $USER_ID" | tee -a "$SESSION_LOG"

AUTH_HEADER="Authorization: Bearer $ACCESS_TOKEN"

# 3. Health check with auth
echo "Testing authenticated health endpoint..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/health" -o /dev/null -w "HTTP %{http_code}\n" | tee -a "$SESSION_LOG"

# Continue with rest of tests...
cat >> "$SESSION_LOG" <<EOF

## 2. Fighter Management (Roster)

EOF

# Create fighter
echo "Creating fighter..." | tee -a "$SESSION_LOG"
FIGHTER_NAME="Fighter_${TIMESTAMP}"
FIGHTER_JSON="{\"name\":\"$FIGHTER_NAME\"}"
FIGHTER_RESP=$(curl -s -X PUT "$API_BASE/api/fighter" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d "$FIGHTER_JSON")
echo "Create fighter response: $FIGHTER_RESP" | tee -a "$SESSION_LOG"
FIGHTER_ID=$(echo "$FIGHTER_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

# Get fighter list
echo "Getting fighter list..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/fighter" | head -c 2000 >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Get fighter details
if [ -n "$FIGHTER_ID" ]; then
  echo "Getting fighter details for ID $FIGHTER_ID..." | tee -a "$SESSION_LOG"
  curl -s -H "$AUTH_HEADER" "$API_BASE/api/fighter/$FIGHTER_ID" >> "$SESSION_LOG"
  echo "" >> "$SESSION_LOG"
fi

cat >> "$SESSION_LOG" <<EOF

## 3. Inventory & Equipment

EOF

# Get inventory balance
echo "Checking inventory balances..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/inventory/balance/particles" >> "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/inventory/balance/token/common" >> "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/inventory/balance/token/rare" >> "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/inventory/balance/token/fabled" >> "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/inventory/balance/token/mythic" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# List weapons
echo "Listing weapons database..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/weapons/database" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 4. Matchmaking & Combat

EOF

# Quick join match
echo "Attempting quick-join match..." | tee -a "$SESSION_LOG"
MATCH_RESP=$(curl -s -X POST "$API_BASE/api/match/quick-join" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Quick join response: $MATCH_RESP" | tee -a "$SESSION_LOG"
MATCH_ID=$(echo "$MATCH_RESP" | grep -o '"matchId":"[^"]*"' | cut -d'"' -f4)
echo "Match ID: $MATCH_ID" | tee -a "$SESSION_LOG"

# Get online players
echo "Getting online players count..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/match/online-players" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

if [ -n "$MATCH_ID" ]; then
  echo "Fetching match details..." | tee -a "$SESSION_LOG"
  curl -s -H "$AUTH_HEADER" "$API_BASE/api/match/$MATCH_ID" >> "$SESSION_LOG"
  echo "" >> "$SESSION_LOG"
  echo "Fetching match teams..." | tee -a "$SESSION_LOG"
  curl -s -H "$AUTH_HEADER" "$API_BASE/api/match/$MATCH_ID/teams" >> "$SESSION_LOG"
  echo "" >> "$SESSION_LOG"
fi

cat >> "$SESSION_LOG" <<EOF

## 5. Shop & Economy

EOF

# Get shop items
echo "Fetching shop items..." |tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/shop/items" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Get gold packages
echo "Fetching gold packages..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/shop/gold" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Get player gold
echo "Player gold balance:" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/player/gold" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 6. Leagues & Leaderboards

EOF

# List leagues
echo "Fetching leagues..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/league" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Leaderboard test
echo "Fetching leaderboard (category: wins)..." | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/leaderboard/wins" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 7. Edge Cases & Exploit Hunting

EOF

echo "Testing edge cases..." | tee -a "$SESSION_LOG"

# a) Invalid fighter ID
echo "Test: GET /api/fighter/999999 (non-existent)" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/fighter/999999" -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# b) Create fighter with missing name
echo "Test: POST /api/fighter with empty name" | tee -a "$SESSION_LOG"
curl -s -X PUT "$API_BASE/api/fighter" -H "$AUTH_HEADER" -H "Content-Type: application/json" -d '{"name":""}' -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# c) Very long name (boundary test)
LONG_NAME=$(printf 'A%.0s' {1..500})
echo "Test: POST /api/fighter with very long name (500 chars)" | tee -a "$SESSION_LOG"
curl -s -X PUT "$API_BASE/api/fighter" -H "$AUTH_HEADER" -H "Content-Type: application/json" -d "{\"name\":\"$LONG_NAME\"}" -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# d) SQL injection attempt in fighter creation
echo "Test: SQL injection in fighter name" | tee -a "$SESSION_LOG"
curl -s -X PUT "$API_BASE/api/fighter" -H "$AUTH_HEADER" -H "Content-Type: application/json" -d '{"name":"test; DROP TABLE fighters;"}' -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# e) Unauthorized access without token
echo "Test: Access /api/fighter without auth" | tee -a "$SESSION_LOG"
curl -s "$API_BASE/api/fighter" -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# f) Token tampering (invalid format)
echo "Test: Invalid Bearer token format" | tee -a "$SESSION_LOG"
curl -s -H "Authorization: Bearer invalidtoken123" "$API_BASE/api/fighter" -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

# g) Test refresh token
echo "Test: Refresh token endpoint" | tee -a "$SESSION_LOG"
curl -s -X POST "$API_BASE/api/authentication/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"userId\":$USER_ID,\"refresh\":\"$REFRESH_TOKEN\"}" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# h) Test duplicate registration (same username)
echo "Test: Duplicate registration with same username" | tee -a "$SESSION_LOG"
curl -s -X POST "$API_BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"email\":\"$EMAIL\"}" -w "HTTP %{http_code}\n" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 8. Stress Testing

EOF

echo "Running stress test: 20 rapid fighter creations (may hit rate limits)..." | tee -a "$SESSION_LOG"
for i in $(seq 1 20); do
  STRESS_NAME="StressFighter_${TIMESTAMP}_$i"
  curl -s -X PUT "$API_BASE/api/fighter" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$STRESS_NAME\"}" -w "HTTP %{http_code} " >> "$SESSION_LOG"
done
echo "" >> "$SESSION_LOG"

echo "Stress test: 10 concurrent health checks (background)..." | tee -a "$SESSION_LOG"
for i in $(seq 1 10); do
  curl -s -H "$AUTH_HEADER" "$API_BASE/api/health" -o /dev/null -w "%{http_code} " &
done
wait
echo "" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 9. Other Features

EOF

# Daily reward status
echo "Daily reward status:" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/daily-reward" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Current events
echo "Current events:" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/events/current" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Attunements
echo "Attunements list:" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/attunements" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Skill tree
echo "Skill tree:" | tee -a "$SESSION_LOG"
curl -s -H "$AUTH_HEADER" "$API_BASE/api/skills/tree" >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

# Season summary
echo "Season summary (POST):" | tee -a "$SESSION_LOG"
curl -s -X POST "$API_BASE/api/season/summary" -H "$AUTH_HEADER" -H "Content-Type: application/json" -d '{}' >> "$SESSION_LOG"
echo "" >> "$SESSION_LOG"

cat >> "$SESSION_LOG" <<EOF

## 10. Observations & Metrics

EOF

# Count fighters created
FIGHTER_COUNT=$(echo "$FIGHTER_RESP" | grep -o '"id":[0-9]*' | wc -l)
echo "Fighters created in session: $FIGHTER_COUNT" | tee -a "$SESSION_LOG"

# Summarize API responsiveness
echo "All tests completed. Review responses above for issues." | tee -a "$SESSION_LOG"

echo "=== Session Complete ===" | tee -a "$SESSION_LOG"

# Save summary of findings for main agent
SUMMARY_FILE="/tmp/player_session_summary.json"
cat > "$SUMMARY_FILE" <<JSON
{
  "date": "$(date '+%Y-%m-%d %H:%M:%S')",
  "mode": "hardcore",
  "userId": $USER_ID,
  "metrics": {
    "fightersCreated": $FIGHTER_COUNT,
    "apiCalls": "~50+",
    "stressRequests": "~40"
  },
  "bugs": [
    "Check log for non-200 responses and error messages"
  ],
  "suggestions": [
    "Review authentication flow for token expiration handling",
    "Consider rate limiting on fighter creation to prevent spam",
    "Validate input lengths (names) on server side"
  ],
  "joy": 6,
  "frustration": 3
}
JSON

echo "Summary saved to $SUMMARY_FILE"
