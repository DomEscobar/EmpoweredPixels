#!/bin/bash
set -euo pipefail

WORKSPACE="/root/EmpoweredPixels"
KANBAN_JSON="$WORKSPACE/tools/kanban-ui/kanban.json"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] [ORCHESTRATOR] $1"; }

TASK_ID=$(grep -B 5 '"status": "todo"' "$KANBAN_JSON" | grep '"id":' | head -n 1 | cut -d'"' -f4 || true)

if [ -z "$TASK_ID" ]; then
    log "✅ Keine offenen Tasks."
    exit 0
fi

TASK_TITLE=$(grep -A 5 "\"id\": \"$TASK_ID\"" "$KANBAN_JSON" | grep '"title":' | head -n 1 | cut -d'"' -f4)
TASK_ASSIGNEE=$(grep -A 5 "\"id\": \"$TASK_ID\"" "$KANBAN_JSON" | grep '"assignee":' | head -n 1 | cut -d'"' -f4)

log "🚀 Starte Task: $TASK_ID"

sed -i "/\"id\": \"$TASK_ID\"/,/\"status\": \"todo\"/s/\"status\": \"todo\"/\"status\": \"in_progress\"/" "$KANBAN_JSON"

# Da sub-agents im Test-Setup noch konfiguriert werden müssen, 
# markieren wir es im Board und senden die Notification.

openclaw message send \
    --target "-1003830885315" \
    --message "🚀 **Delegation-Flow Aktiv**
🎯 **Task:** $TASK_TITLE
👤 **Assignee:** $TASK_ASSIGNEE
📍 **Status:** in_progress
---
Board: http://v2202502215330313077.supersrv.de:8666"
