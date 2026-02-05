#!/bin/bash
# Enhanced Orchestrator with Guaranteed Git Safety
# Runs every 2 minutes, delegates tasks to agents

set -euo pipefail

CONFIG_FILE="/root/.clawdbot/clawdbot.json"
WORKSPACE="/root/EmpoweredPixels"
KANBAN_FILE="$WORKSPACE/KANBAN.md"
LOG_FILE="/var/log/agency/orchestrator.log"
PID_FILE="/var/run/orchestrator.pid"
MAX_AGENTS=4

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [ORCHESTRATOR] $1" | tee -a "$LOG_FILE"
}

# Prevent multiple instances
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE" 2>/dev/null || echo "")
    if [ -n "$OLD_PID" ] && kill -0 "$OLD_PID" 2>/dev/null; then
        log "⚠️ Orchestrator already running (PID: $OLD_PID), exiting"
        exit 0
    fi
fi
echo $$ > "$PID_FILE"

# Cleanup on exit
trap 'rm -f "$PID_FILE"' EXIT

log "🎯 Starting orchestration cycle..."

# Check active agents
active_agents=$(pgrep -f "openclaw.*agent" | wc -l)
log "📊 Active agents: $active_agents/$MAX_AGENTS"

if [ "$active_agents" -ge "$MAX_AGENTS" ]; then
    log "⏸️ Max agents reached, waiting..."
    exit 0
fi

# Parse KANBAN for next task
parse_kanban() {
    local section="$1"
    local task_pattern="$2"
    
    awk -v section="$section" '
        BEGIN { in_section = 0 }
        $0 ~ section { in_section = 1; next }
        /^###/ { in_section = 0 }
        in_section && /^- \[ \]/ { print; exit }
    ' "$KANBAN_FILE"
}

# Find next task (P0 first, then P1)
next_task=""
priority=""

for prio in "P0" "P1" "P2"; do
    next_task=$(parse_kanban "🟢 TO DO" "$prio")
    if [ -n "$next_task" ]; then
        priority="$prio"
        break
    fi
done

if [ -z "$next_task" ]; then
    log "✅ No pending tasks in KANBAN"
    exit 0
fi

# Extract task details
task_name=$(echo "$next_task" | sed 's/- \[ \] //' | sed 's/\*\*//g' | cut -d'-' -f1 | xargs)
task_id=$(echo "$task_name" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | sed 's/[^a-z0-9-]//g' | cut -c1-40)

log "🎯 Found task: $task_name (ID: $task_id, Priority: $priority)"

# Determine agent type based on task
determine_agent() {
    local task="$1"
    if echo "$task" | grep -qi "test\|qa\|coverage"; then
        echo "qa"
    elif echo "$task" | grep -qi "skill\|pattern\|refactor"; then
        echo "foundry"
    else
        echo "coder"
    fi
}

agent_type=$(determine_agent "$task_name")
log "🤖 Assigning to agent: $agent_type"

# Create feature branch
cd "$WORKSPACE"

if [ ! -d ".git" ]; then
    log "❌ Not a git repository"
    exit 1
fi

# Stash any uncommitted changes
git stash push -m "orchestrator-auto-stash-$(date +%s)" 2>/dev/null || true

# Checkout main and pull
git checkout main 2>/dev/null || git checkout master
git pull origin main 2>/dev/null || git pull origin master

# Create feature branch
branch_name="feature/$task_id"
if git branch | grep -q "$branch_name"; then
    log "🌿 Branch exists, checking out: $branch_name"
    git checkout "$branch_name"
    git pull origin "$branch_name" 2>/dev/null || true
else
    log "🌿 Creating branch: $branch_name"
    git checkout -b "$branch_name"
fi

# IMMEDIATELY push branch to origin (guarantee remote tracking)
if git push -u origin "$branch_name" 2>/dev/null; then
    log "✅ Branch $branch_name pushed to origin"
else
    log "⚠️ Branch may already exist on origin"
fi

# Create agent instruction file with MANDATORY git commands
agent_instruction_file="$WORKSPACE/.agent-instructions-$task_id.sh"
cat > "$agent_instruction_file" << 'EOF'
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
EOF

chmod +x "$agent_instruction_file"

# Update KANBAN - mark as IN PROGRESS
log "📝 Updating KANBAN..."

# Use atomic sed replacement
sed -i "s|- \[ \] \*\*$task_name\*\*|- [x] **$task_name** (IN PROGRESS)|" "$KANBAN_FILE" 2>/dev/null || \
sed -i "s|- \[ \] $task_name|- [x] $task_name (IN PROGRESS)|" "$KANBAN_FILE"

# Add branch info to KANBAN
sed -i "/- \[x\].*$task_id/,/^$/s|Assignee:.*|Assignee: $agent_type @ $branch_name|" "$KANBAN_FILE" 2>/dev/null || true

# Commit and push KANBAN update immediately
git add "$KANBAN_FILE"
if git commit -m "chore: $task_id moved to IN PROGRESS [auto]" 2>/dev/null; then
    git push origin main 2>/dev/null || log "⚠️ Could not push KANBAN update"
    log "✅ KANBAN updated and pushed"
else
    log "ℹ️ No KANBAN changes to commit"
fi

# Create task manifest for tracking
manifest_file="$WORKSPACE/.agent-manifest-$task_id.json"
cat > "$manifest_file" << EOF
{
  "task_id": "$task_id",
  "task_name": "$task_name",
  "priority": "$priority",
  "agent_type": "$agent_type",
  "branch": "$branch_name",
  "created_at": "$(date -Iseconds)",
  "expires_at": "$(date -Iseconds -d '+24 hours')",
  "status": "in_progress",
  "manifest_version": "1.0"
}
EOF

git add "$manifest_file"
git commit -m "chore: manifest for $task_id [auto]" 2>/dev/null || true
git push origin "$branch_name" 2>/dev/null || true

# Determine task-specific prompt
case "$agent_type" in
    coder)
        PROMPT="CRITICAL GIT PROTOCOL:
1. Work ONLY in branch: $branch_name
2. Run 'source /root/EmpoweredPixels/.agent-instructions-$task_id.sh' first
3. Call safe_commit() after EVERY file change
4. Call agent_complete() when done
5. NEVER work on main branch

Task: $task_name
Implement the feature. Write tests. Ensure 80%+ coverage."
        ;;
    qa)
        PROMPT="CRITICAL GIT PROTOCOL:
1. Work ONLY in branch: $branch_name  
2. Run 'source /root/EmpoweredPixels/.agent-instructions-$task_id.sh' first
3. Call safe_commit() after test file creation
4. Call agent_complete() when done

Task: $task_name
Write comprehensive tests. Verify coverage."
        ;;
    foundry)
        PROMPT="CRITICAL GIT PROTOCOL:
1. Work ONLY in branch: $branch_name
2. Run 'source /root/EmpoweredPixels/.agent-instructions-$task_id.sh' first
3. Call safe_commit() after pattern extraction
4. Call agent_complete() when done

Task: $task_name
Analyze patterns. Create reusable skill."
        ;;
    *)
        PROMPT="Task: $task_name. Branch: $branch_name. Follow git protocol."
        ;;
esac

# Log the delegation
log "📝 Agent prompt prepared"
log "📝 Instruction file: $agent_instruction_file"

# Send Telegram notification
if command -v openclaw &> /dev/null; then
    openclaw message send \
        --target "-1003830885315" \
        --message "🎯 Task Delegated
Task: $task_name
Agent: $agent_type
Branch: $branch_name
Priority: $priority
Instructions: $agent_instruction_file" \
        2>/dev/null || true
fi

# In production, this would spawn the actual agent
# For now, we create a marker that the agent was delegated
marker_file="$WORKSPACE/.agent-delegated-$task_id"
echo "$agent_type:$branch_name:$(date +%s)" > "$marker_file"
git add "$marker_file" 2>/dev/null || true
git commit -m "chore: delegation marker for $task_id [auto]" 2>/dev/null || true
git push origin "$branch_name" 2>/dev/null || true

log "✅ Delegation complete for $task_id"
log "🎉 Orchestration cycle complete"

exit 0
