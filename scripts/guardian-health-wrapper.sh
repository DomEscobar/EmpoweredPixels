#!/bin/bash
# Guardian Health Check Wrapper
# Runs pre-merge checks and reports status with alerts on failure

set -euo pipefail

REPO_ROOT="/root/EmpoweredPixels"
SCRIPT_PATH="$REPO_ROOT/scripts/pre-merge-check.sh"
REPORT_DIR="/root/.openclaw/reports"
REPORT_FILE="$REPORT_DIR/lastGuardianHealth.txt"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S %Z')
TELEGRAM_GROUP="-1003830885315"

# Ensure report directory exists
mkdir -p "$REPORT_DIR"

# Run the health check from repo root
cd "$REPO_ROOT"
OUTPUT=$(bash "$SCRIPT_PATH" 2>&1)
EXIT_CODE=$?

# Write summary report
{
  echo "=== Guardian Health Check ==="
  echo "Timestamp: $TIMESTAMP"
  echo "Exit Code: $EXIT_CODE"
  echo "=== Output ==="
  echo "$OUTPUT"
  echo "=== End of Report ==="
} > "$REPORT_FILE"

# If check failed, send alerts
if [ $EXIT_CODE -ne 0 ]; then {
  # Alert summary for main session (system event)
  ALERT_MSG="[Guardian Health] Pre-merge check FAILED (exit $EXIT_CODE). Timestamp: $TIMESTAMP. First lines:\n$(echo "$OUTPUT" | head -n 5)"

  # Send to main session via system event
  echo "$ALERT_MSG" | openclaw cron wake --mode next-heartbeat

  # Post to Telegram group
  TG_MSG="🚨 Guardian Health Check FAILED\nTimestamp: $TIMESTAMP\nExit Code: $EXIT_CODE\n\nOutput (first 500 chars):\n$(echo "$OUTPUT" | head -c 500)"
  openclaw message send --target "$TELEGRAM_GROUP" --message "$TG_MSG" --best-effort
} fi

# Return appropriate exit code for cron
exit $EXIT_CODE
