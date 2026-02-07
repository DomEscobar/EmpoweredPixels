#!/bin/bash
set -euo pipefail

WORKSPACE="/root/EmpoweredPixels"
KANBAN_JSON="$WORKSPACE/tools/kanban-ui/kanban.json"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] [ORCHESTRATOR] $1"; }

# 1. Get first TODO task
TASK_ID=$(grep -B 3 '"status": "todo"' "$KANBAN_JSON" | grep '"id":' | head -n 1 | cut -d'"' -f4 || true)

if [ -z "$TASK_ID" ]; then
    log "✅ No pending tasks."
    exit 0
fi

TASK_TITLE=$(grep -A 5 "\"id\": \"$TASK_ID\"" "$KANBAN_JSON" | grep '"title":' | head -n 1 | cut -d'"' -f4)
TASK_ASSIGNEE=$(grep -A 5 "\"id\": \"$TASK_ID\"" "$KANBAN_JSON" | grep '"assignee":' | head -n 1 | cut -d'"' -f4)

log "🚀 Starting agent turn for $TASK_ID via @$TASK_ASSIGNEE"

# 2. Update status to in_progress
sed -i "/\"id\": \"$TASK_ID\"/,/\"status\": \"todo\"/s/\"status\": \"todo\"/\"status\": \"in_progress\"/" "$KANBAN_JSON"

# 3. Use 'openclaw agent' with --agent flag to trigger the specialized agent
# We use & to run in background so orchestrator doesn't hang
openclaw agent \
    --agent "$TASK_ASSIGNEE" \
    --message "MISSION START: Work on '$TASK_TITLE'. Context: in $WORKSPACE. Follow AGENTS.md. Commit and push everything." \
    --deliver \
    --session-id "task-$TASK_ID" > /tmp/agent-$TASK_ID.log 2>&1 &

# 4. Notify Telegram
openclaw message send \
    --target "-1003830885315" \
    --message "⚡️ **AGENT DEPLOYED**
👤 **Role:** @$TASK_ASSIGNEE
🎯 **Mission:** $TASK_TITLE
📍 **Status:** autonomous_turn_started"
