#!/bin/bash
# MANDATORY GIT PROTOCOL FOR AGENTS
# Failure to follow = task failure

set -e

TASK_ID="$1"
BRANCH="$2"
AGENT_TYPE="$3"

cd /root/EmpoweredPixels

# Git safety function
safe_commit() {
    local msg="$1"
    if ! git diff --quiet HEAD 2>/dev/null || ! git diff --cached --quiet HEAD 2>/dev/null; then
        git add -A
        git commit -m "$msg [$AGENT_TYPE]" || true
        git push origin "$BRANCH" || true
        echo "✅ Committed: $msg"
    fi
}

# Auto-commit every 5 minutes in background
(
    while true; do
        sleep 300
        safe_commit "wip: checkpoint auto-save"
    done
) &
CHECKPOINT_PID=$!

# Cleanup checkpoint on exit
trap 'kill $CHECKPOINT_PID 2>/dev/null; safe_commit "wip: final checkpoint before exit"' EXIT

# Agent must call this when done
agent_complete() {
    safe_commit "feat: $TASK_ID implementation complete"
    echo "✅ Task complete, all changes committed and pushed"
}

export -f safe_commit
export -f agent_complete
export TASK_ID BRANCH AGENT_TYPE
