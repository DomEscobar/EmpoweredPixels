#!/usr/bin/env bash
set -euo pipefail

# Telegram group where everything reports
TG="-1003830885315"

echo "Setting up agency cron jobs..."

# Morning standup — 09:00 weekdays
openclaw cron add \
  --name "Morning Standup" \
  --cron "0 9 * * 1-5" \
  --tz "Europe/Berlin" \
  --session isolated \
  --agent main \
  --message "Read /root/EmpoweredPixels/kanban.json. Post a standup to Telegram: done yesterday, in progress today, blockers. Assign any unassigned backlog tasks to idle devs via agent-to-agent." \
  --announce \
  --channel telegram \
  --to "$TG"

# Evening wrap — 18:00 weekdays
openclaw cron add \
  --name "Evening Wrap" \
  --cron "0 18 * * 1-5" \
  --tz "Europe/Berlin" \
  --session isolated \
  --agent main \
  --message "Read /root/EmpoweredPixels/kanban.json. Post end-of-day summary to Telegram: completed, still in progress, any risks." \
  --announce \
  --channel telegram \
  --to "$TG"

# Stale check — every 3 hours
openclaw cron add \
  --name "Stale Check" \
  --cron "0 */3 * * *" \
  --tz "Europe/Berlin" \
  --session isolated \
  --agent main \
  --message "Read /root/EmpoweredPixels/kanban.json. Find in_progress tasks with updatedAt older than 3 hours. Ping the assigned dev via agent-to-agent. If unresponsive for 6+ hours, report to Telegram." \
  --announce \
  --channel telegram \
  --to "$TG"

echo "Done. Verify: openclaw cron list"
