#!/bin/bash
# KANBAN Monitor - Runs every 10 minutes
# Reports current sprint status

echo "=== KANBAN Status Check $(date) ==="
echo ""
echo "🔴 BACKLOG:"
grep -A1 "^## 🔴 BACKLOG" /root/.openclaw/workspace/KANBAN.md | grep "^- \[ \]" | head -5

echo ""
echo "🟡 IN PROGRESS:"
grep -A1 "^## 🟡 IN PROGRESS" /root/.openclaw/workspace/KANBAN.md | grep "^- \[ \]"

echo ""
echo "🟢 DONE:"
DONE_COUNT=$(grep -c "^- \[x\]" /root/.openclaw/workspace/KANBAN.md || echo "0")
echo "Total completed: $DONE_COUNT"

echo ""
echo "=== Next Actions ==="
# Check if IN PROGRESS is empty
IN_PROGRESS=$(awk '/^## 🟡 IN PROGRESS/{found=1} found && /^- \[ \]/{print; exit}' /root/.openclaw/workspace/KANBAN.md)
if [ -z "$IN_PROGRESS" ]; then
    echo "⚠️ No active task! Architect may be idle."
else
    echo "✅ Architect has active work."
fi
